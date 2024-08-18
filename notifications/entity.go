package notifications

import "time"

type Notification struct {
	Id        uint      `gorm:"primaryKey"`
	Time      time.Time `gorm:"not null"`
	Payload   string    `gorm:"not null"`
	UserId    uint      `gorm:"not null"`
	ownerInfo OwnerInfo `gorm:"embedded"`
}
