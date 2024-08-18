package user

type User struct {
	Id           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"not null"`
	Email        string         `gorm:"not null"`
	UserServices []UserServices `gorm:"foreignKey:UserId"`
}

type UserServices struct {
	UserId        uint `gorm:"primaryKey"`
	Service       OwnerService
	ServiceUserId string
}

type OwnerService string

const (
	Tg   OwnerService = "tg"
	Site OwnerService = "site"
)
