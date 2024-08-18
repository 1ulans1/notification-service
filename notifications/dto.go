package notifications

import (
	"notification-service/user"
	"time"
)

type NotificationInDto struct {
	Id        string    `json:"id"`
	Time      time.Time `json:"time"`
	Payload   string    `json:"payload"`
	OwnerInfo OwnerInfo `json:"owner_info"`
}

type OwnerInfo struct {
	OwnerID      string            `json:"owner_id"`
	OwnerService user.OwnerService `json:"owner_service"`
}

type NotificationDto struct {
	Id      int       `json:"id"`
	Time    time.Time `json:"time"`
	Payload string    `json:"payload"`
}
