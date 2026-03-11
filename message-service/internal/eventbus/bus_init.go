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
	SubscriptionNames []string
	SubscriptionProperties *admin.SubscriptionProperties
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
		log.Info().Msg(fmt.Sprintf("unable to retrieve %s topic", e.TopicName))
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

		for _, sub := range e.SubscriptionNames {
			err := createSubscription(e, sub,adminCli)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func createSubscription(e *EventBusArgs,subscription string, adminCli *admin.Client) error {

	res, err := adminCli.GetSubscription(e.Ctx,e.TopicName,subscription,&admin.GetSubscriptionOptions{})
	if res != nil {
		log.Info().Msg("Subscription with specified name already exists, ignoring subscription creation")
		return nil

	} else if err != nil {
		log.Info().Msg(fmt.Sprintf("unable to retrieve %s subscription", subscription))
		return err

	} else {
		subOpt := &admin.CreateSubscriptionOptions{
			Properties: e.SubscriptionProperties,
		}
		_, err := adminCli.CreateSubscription(e.Ctx,e.TopicName,subscription,subOpt)
		if err != nil {
			return err
		} else {
			log.Info().Msg(fmt.Sprintf("Created a new %s subscription", subscription))
		}
		return nil
	}

}
