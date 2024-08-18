package notifications

import (
	"notification-service/user"
)

type NotificationService interface {
	SaveNotification(notification NotificationInDto) error
	GetNotification(id string) (NotificationDto, error)
	GetAllNotifications() ([]NotificationDto, error)
	DeleteNotification(id string) error
	GetAllNotificationsByUserId(userId uint) ([]NotificationDto, error)
}

type notificationService struct {
	notificationRepo NotificationRepo
	userService      user.UserService
}

func NewNotificationService(notificationRepo NotificationRepo) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *notificationService) SaveNotification(notification NotificationInDto) error {
	user, err := s.userService.GetUserByServiceAndServiceUserID(notification.OwnerInfo.OwnerService, notification.OwnerInfo.OwnerID)
	if err != nil {
		return err
	}

	return s.notificationRepo.SaveNotification(Notification{
		Time:    notification.Time,
		Payload: notification.Payload,
		UserId:  user.Id,
	})
}

func (s *notificationService) GetNotification(id string) (NotificationDto, error) {
	notification, err := s.notificationRepo.GetNotification(id)
	if err != nil {
		return NotificationDto{}, err
	}

	return mapNotificationToDto(notification), nil
}

func mapNotificationToDto(notification Notification) NotificationDto {
	return NotificationDto{
		Id:      int(notification.Id),
		Time:    notification.Time,
		Payload: notification.Payload,
	}
}

func (s *notificationService) GetAllNotifications() ([]NotificationDto, error) {
	notifications, err := s.notificationRepo.GetAllNotifications()
	if err != nil {
		return nil, err
	}

	return mapNotificationsToDto(notifications), nil
}

func mapNotificationsToDto(notifications []Notification) []NotificationDto {
	var dtos []NotificationDto

	for _, n := range notifications {
		dtos = append(dtos, mapNotificationToDto(n))
	}

	return dtos
}

func (s *notificationService) DeleteNotification(id string) error {
	return s.notificationRepo.DeleteNotification(id)
}

func (s *notificationService) GetAllNotificationsByUserId(userId uint) ([]NotificationDto, error) {
	notifications, err := s.notificationRepo.GetAllNotificationsByUserId(userId)
	if err != nil {
		return nil, err
	}

	return mapNotificationsToDto(notifications), nil
}
