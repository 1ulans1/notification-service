package user

type UserService interface {
	CreateUser(name string, email string) error
	CreateUserByServiceAndServiceUserID(name string, email string, service OwnerService, serviceUserID string) error
	GetUserByServiceAndServiceUserID(service OwnerService, serviceID string) (UserDto, error)
	GetAllUsers() ([]UserDto, error)
	DeleteUser(id string) error
	AddServiceToUser(id string, service OwnerService, serviceUserId string) error
	GetUserById(id string) (UserDto, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(name string, email string) error {
	return s.userRepo.SaveUser(User{
		Name:  name,
		Email: email,
	})
}

func (s *userService) CreateUserByServiceAndServiceUserID(name string, email string, service OwnerService, serviceUserID string) error {
	user := User{
		Name:  name,
		Email: email,
		UserServices: []UserServices{
			{
				Service:       service,
				ServiceUserId: serviceUserID,
			},
		},
	}

	return s.userRepo.SaveUser(user)
}

func (s *userService) GetUserByServiceAndServiceUserID(service OwnerService, serviceID string) (UserDto, error) {
	user, err := s.userRepo.GetUserByServiceAndServiceID(service, serviceID)
	if err != nil {
		return UserDto{}, err
	}

	return mapUserToDto(user), nil
}

func mapUserToDto(user User) UserDto {
	return UserDto{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}

func (s *userService) GetAllUsers() ([]UserDto, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userDtos []UserDto
	for _, user := range users {
		userDtos = append(userDtos, mapUserToDto(user))
	}

	return userDtos, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.DeleteUser(id)
}

func (s *userService) AddServiceToUser(id string, service OwnerService, serviceUserId string) error {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return err
	}

	user.UserServices = append(user.UserServices, UserServices{
		Service:       service,
		ServiceUserId: serviceUserId,
	})

	return s.userRepo.SaveUser(user)
}

func (s *userService) GetUserById(id string) (UserDto, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return UserDto{}, err
	}

	return mapUserToDto(user), nil
}
