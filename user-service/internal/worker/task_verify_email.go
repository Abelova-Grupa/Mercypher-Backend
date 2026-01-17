package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Abelova-Grupa/Mercypher/user-service/internal/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"


func (dist *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context,payload *email.EmailPayload,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload);
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %v", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	_, err = dist.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}
	return nil
}

func (proc *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload email.EmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("Failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	//TODO: This code makes some much fucking problems for me, even though it was transactioned
	//TODO: Find a solution to this problem
	// var user models.User
	// res := proc.db.WithContext(ctx).Where("username = ?", payload.Username).First(&user)
	// if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	// 	return fmt.Errorf("Task couldn't be processed, username doesn't exist %v %w", payload.Username, asynq.SkipRetry)
	// }
	
	if err := proc.emailSender.SendEmail(payload); err != nil {
		return fmt.Errorf("failed to send verify email %w", asynq.SkipRetry)
	}
	log.Info().Str("type",task.Type()).Bytes("payload", task.Payload()).Str("email", payload.ToEmail).Msg("processed task")

	return nil
}