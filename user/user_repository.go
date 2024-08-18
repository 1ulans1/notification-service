package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user User) error
	GetUserByServiceAndServiceID(service OwnerService, serviceID string) (User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(id string) error
	GetUserById(id string) (User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveUser(user User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) GetUserByServiceAndServiceID(service OwnerService, serviceID string) (User, error) {
	var user User
	err := r.db.Where("service = ? AND service_user_id = ?", service, serviceID).First(&user).Error
	return user, err
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) DeleteUser(id string) error {
	return r.db.Where("id = ?", id).Delete(&User{}).Error
}

func (r *userRepository) GetUserById(id string) (User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}
