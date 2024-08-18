package notifications

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification-service/clients"
)

type NotificationConsumer interface {
	StartConsuming()
}

type notificationConsumerImpl struct {
	rmqClient           clients.RmqClient
	notificationService NotificationService
}

func NewNotificationConsumer(rmqClient clients.RmqClient, notificationService NotificationService) NotificationConsumer {
	return &notificationConsumerImpl{
		rmqClient:           rmqClient,
		notificationService: notificationService,
	}
}

const (
	exchangeName = "backend.notifications"
	queueName    = "backend.notifications"
)

func (c *notificationConsumerImpl) StartConsuming() {
	go c.rmqClient.ConsumeWithRetry(exchangeName, queueName, func(msg amqp.Delivery) error {
		var notification NotificationInDto
		err := json.Unmarshal(msg.Body, &notification)
		if err != nil {
			return fmt.Errorf("failed to unmarshal notification: %w", err)
		}
		err = c.notificationService.SaveNotification(notification)
		if err != nil {
			return fmt.Errorf("failed to save notification: %w", err)
		}
		return nil
	}, 5)
}
