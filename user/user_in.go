package user

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification-service/clients"
	"notification-service/notifications"
)

type UserInService interface {
	StartConsuming()
}

type userInService struct {
	rmqClient   clients.RmqClient
	userService UserService
}

func NewUserInService(rmqClient clients.RmqClient, userService UserService) UserInService {
	return &userInService{
		rmqClient:   rmqClient,
		userService: userService,
	}
}

const (
	exchangeName = "backend.users"
	queueName    = "backend.users"
)

type UserInDto struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	OwnerInfo notifications.OwnerInfo
}

func (s *userInService) StartConsuming() {
	go s.rmqClient.ConsumeWithRetry(exchangeName, queueName, func(msg amqp.Delivery) error {
		var user UserInDto
		err := json.Unmarshal(msg.Body, &user)
		if err != nil {
			return fmt.Errorf("failed to unmarshal user: %w", err)
		}
		err = s.userService.SaveUser(user)
		if err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}
		return nil
	}, 5)
}
