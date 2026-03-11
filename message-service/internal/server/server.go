package server

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/Abelova-Grupa/Mercypher/message-service/internal/eventbus"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/kafka"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/model"
	"github.com/Abelova-Grupa/Mercypher/message-service/internal/repository"
	pb "github.com/Abelova-Grupa/Mercypher/proto/message"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MessageServer handles incoming gRPC requests implementing protobuf

type KafkaMessageServer struct {
	pb.UnimplementedMessageServiceServer
	brokers []string
	repo    repository.MessageRepository
}

func NewKafkaMessageServer(brokers []string, repo repository.MessageRepository) *KafkaMessageServer {
	// return &MessageServer{
	// 	brokers: brokers,
	// 	repo:    repo,
	// }
	return &KafkaMessageServer{
		brokers: brokers,
		repo:    repo,
	}
}

func (k *KafkaMessageServer) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.MessageAck, error) {
	// bare minimum checks
	if req.Body == "" || req.RecieverId == "" || req.SenderId == "" {
		return nil, status.Error(codes.InvalidArgument, "body and recipient_id and sender_id are required")
	}

	req.Id = uuid.New().String()
	// when is timestamp added? here maybe?
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().Unix()
	}

	log.Printf("Queueing message from %s to %s", req.SenderId, req.RecieverId)

	generatedID, err := kafka.PublishMessage(ctx, k.brokers, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to queue: %v", err)
	}

	return &pb.MessageAck{
		MessageId: generatedID,
	}, nil
}

func (k *KafkaMessageServer) GetMessages(ctx context.Context, req *pb.MessageRange) (*pb.MessageList, error) {
	lastSeen := time.Unix(req.LastSeen, 0)
	if req.Limit < 1 {
		req.Limit = 20
	}

	var messages []model.ChatMessage // Assuming this is your repo model type
	var err error

	// Check if Participant2 is a Group UUID
	_, uuidErr := uuid.Parse(req.Participant2)

	if uuidErr == nil {
		// --- GROUP HISTORY ---
		// We only care about messages where Receiver_id == GroupUUID
		messages, err = k.repo.GetGroupHistory(ctx, req.Participant2, lastSeen, int(req.Limit))
	} else {
		// --- 1-on-1 HISTORY ---
		// Standard logic: messages between P1 and P2
		messages, err = k.repo.GetChatHistory(ctx, req.Participant1, req.Participant2, lastSeen, int(req.Limit))
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch history: %v", err)
	}

	var pbMessages []*pb.ChatMessage
	for _, m := range messages {
		pbMessages = append(pbMessages, &pb.ChatMessage{
			Id:         m.Message_id,
			SenderId:   m.Sender_id,
			RecieverId: m.Receiver_id,
			Body:       m.Body,
			Timestamp:  m.Timestamp.Unix(),
		})
	}

	return &pb.MessageList{Messages: pbMessages}, nil
}

type AzureMessageServer struct {
	pb.UnimplementedMessageServiceServer
	repo          *repository.MessageRepo
	AzureCli      *azservicebus.Client
	QueueSender   *eventbus.BusSender
	TopicSender   *eventbus.BusSender
	QueueReceiver *eventbus.BusReceiver
	TopicReceiver *eventbus.BusReceiver
}

func NewAzureMessageServer(ctx context.Context, repo *repository.MessageRepo) *AzureMessageServer {
	namespace := os.Getenv("AZURE_SERVICE_BUS_CONN_STR")
	if namespace == "" {
		panic("azure service bus connection string not loaded")
	}
	clientOpt := &azservicebus.ClientOptions{
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		RetryOptions: azservicebus.RetryOptions{
			MaxRetries: 5,
		},
	}
	client, err := azservicebus.NewClientFromConnectionString(namespace, clientOpt)
	if err != nil {
		panic("unable to create azure message server")
	}

	// Initializing azure service bus topics and queues if they don't exist
	busArgs := eventbus.EventBusArgs{
		Ctx:            ctx,
		QueueName:      "contact-queue",
		QueueProerties: &admin.QueueProperties{},

		TopicName:       "azure-topic",
		TopicProperites: &admin.TopicProperties{},

		SubscriptionNames:      []string{"gateway-sub", "message-service-sub"},
		SubscriptionProperties: &admin.SubscriptionProperties{},
	}
	busArgs.InitEventBus()

	// Initializing azure service bus sender, variable names could be better
	queueSender := eventbus.NewBusSender("contact-queue")
	queueSender.CreateSender(client)

	topicSender := eventbus.NewBusSender("azure-topic")
	topicSender.CreateSender(client)

	busConsumer, _ := eventbus.CreateQueueReceiver(client, "contact-queue", nil)
	queueConsumer := eventbus.NewBusReceiver(repo, busConsumer)

	topicConsumer, _ := eventbus.CreateTopicReceiver(client, "azure-topic", "message-service-sub", nil)
	azTopicConsumer := eventbus.NewBusReceiver(repo, topicConsumer)

	return &AzureMessageServer{
		repo:          repo,
		AzureCli:      client,
		QueueSender:   queueSender,
		TopicSender:   topicSender,
		QueueReceiver: queueConsumer,
		TopicReceiver: azTopicConsumer,
	}
}

func (a *AzureMessageServer) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.MessageAck, error) {
	if req.Body == "" || req.RecieverId == "" || req.SenderId == "" {
		return nil, status.Error(codes.InvalidArgument, "body and recipient_id and sender_id are required")
	}

	req.Id = uuid.New().String()
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().Unix()
	}
	log.Printf("Queueing message from %s to %s", req.SenderId, req.RecieverId)

	// For now it only works for contact messaging, not with groups
	generatedID, err := a.TopicSender.SendMessage(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to queue: %v", err)
	}

	return &pb.MessageAck{
		MessageId: generatedID,
	}, nil
}

func (a *AzureMessageServer) GetMessages(ctx context.Context, req *pb.MessageRange) (*pb.MessageList, error) {
	lastSeen := time.Unix(req.LastSeen, 0)
	if req.Limit < 1 {
		req.Limit = 10
	}
	var messages []model.ChatMessage // Assuming this is your repo model type
	var err error
	_, uuidErr := uuid.Parse(req.Participant2)

	if uuidErr == nil {
		
		messages, err = a.repo.GetGroupHistory(ctx, req.Participant2, lastSeen, int(req.Limit))
	} else {
		
		messages, err = a.repo.GetChatHistory(ctx, req.Participant1, req.Participant2, lastSeen, int(req.Limit))
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch history: %v", err)
	}

	var pbMessages []*pb.ChatMessage
	for _, m := range messages {
		pbMessages = append(pbMessages, &pb.ChatMessage{
			Id:         m.Message_id,
			SenderId:   m.Sender_id,
			RecieverId: m.Receiver_id,
			Body:       m.Body,
			Timestamp:  m.Timestamp.Unix(),
		})
	}

	return &pb.MessageList{
		Messages: pbMessages,
	}, nil
}
