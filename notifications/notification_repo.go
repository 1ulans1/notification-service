package notifications

import (
	"gorm.io/gorm"
)

type NotificationRepo interface {
	SaveNotification(notification Notification) error
	GetNotification(id string) (Notification, error)
	GetAllNotifications() ([]Notification, error)
	DeleteNotification(id string) error
	GetAllNotificationsByUserId(userId uint) ([]Notification, error)
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) SaveNotification(notification Notification) error {
	return r.db.Create(&notification).Error
}

func (r *notificationRepo) GetNotification(id string) (Notification, error) {
	var notification Notification
	err := r.db.Where("id = ?", id).First(&notification).Error
	return notification, err
}

func (r *notificationRepo) GetAllNotifications() ([]Notification, error) {
	var notifications []Notification
	err := r.db.Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepo) DeleteNotification(id string) error {
	return r.db.Where("id = ?", id).Delete(&Notification{}).Error
}

func (r *notificationRepo) GetAllNotificationsByUserId(userId uint) ([]Notification, error) {
	var notifications []Notification
	err := r.db.Where("user_id = ?", userId).Find(&notifications).Error
	return notifications, err
}
