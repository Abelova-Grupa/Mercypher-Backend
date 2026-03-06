package eventbus

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/rs/zerolog/log"
)

type EventBusArgs struct {
	Ctx             context.Context
	QueueName       string
	QueueProerties  *admin.QueueProperties
	TopicName       string
	TopicProperites *admin.TopicProperties
}

func (e *EventBusArgs) InitEventBus() error {
	connectionStr := os.Getenv("AZURE_SERVICE_BUS_CONN_STR")
	if connectionStr == "" {
		panic("failed to retrieve azure service bus connection string")
	}
	azureOpt := &admin.ClientOptions{}
	adminCli, err := admin.NewClientFromConnectionString(connectionStr, azureOpt)
	if err != nil {
		return err
	}
	if err := createQueue(e, adminCli); err != nil {
		return err
	}

	if err := createTopic(e, adminCli); err != nil {
		return err
	}
	return nil
}

func createQueue(e *EventBusArgs, adminCli *admin.Client) error {
	res, err := adminCli.GetQueue(e.Ctx, e.QueueName, &admin.GetQueueOptions{})
	if res != nil {
		log.Info().Msg("queue with specified name already exists, ignoring queue creation")
		return nil
	} else if err != nil {
		log.Info().Msg(fmt.Sprintf("unable to retrieve %s queue", e.QueueName))
		return err
	} else {
		queueOpt := &admin.CreateQueueOptions{
			Properties: e.QueueProerties,
		}
		_, err := adminCli.CreateQueue(e.Ctx, e.QueueName, queueOpt)
		if err != nil {
			return err
		}
		log.Info().Msg(fmt.Sprintf("Created a new %s queue", e.QueueName))
		return nil
	}
}

func createTopic(e *EventBusArgs, adminCli *admin.Client) error {
	res, err := adminCli.GetTopic(e.Ctx, e.TopicName, &admin.GetTopicOptions{})
	if res != nil {
		log.Info().Msg("Topic with specified name already exists, ignoring topic creation")
		return nil
	} else if err != nil {
		log.Info().Msg(fmt.Sprintf("unable to retrieve %s queue", e.TopicName))
		return err
	} else {
		topicOpt := &admin.CreateTopicOptions{
			Properties: e.TopicProperites,
		}
		_, err := adminCli.CreateTopic(e.Ctx, e.TopicName, topicOpt)
		if err != nil {
			return err
		} else {
			log.Info().Msg(fmt.Sprintf("Created a new %s topic", e.TopicName))
		}
		return nil
	}
}
