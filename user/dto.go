package user

type UserDto struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
