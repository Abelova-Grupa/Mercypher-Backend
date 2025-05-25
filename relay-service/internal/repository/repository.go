package repository

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Abelova-Grupa/Mercypher-Backend/relay-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher-Backend/relay-service/internal/config"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

var RedisRepo = *redis.NewClient(&redis.Options{
	Addr:     config.GetEnv("REDIS_REPO_ADDR", "localhost:6379"),
	Password: "",
	DB:       0,
})

var Ctx = context.Background()

func SaveMessage(msg *pb.ChatMessage) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return RedisRepo.RPush(Ctx, "userid:"+msg.RecipientId, data).Err()
}

func GetMessages(id *pb.UserId) ([]*pb.ChatMessage, error) {
	rawMessages, err := RedisRepo.LRange(Ctx, "userid:"+id.Id, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*pb.ChatMessage
	for _, raw := range rawMessages {
		var msg pb.ChatMessage
		if err := proto.Unmarshal([]byte(raw), &msg); err != nil {
			log.Fatalf("failed to deserialize message: %v", err)
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func DeleteMessages(id *pb.UserId) {
	fmt.Println("To be implemented...")
}
