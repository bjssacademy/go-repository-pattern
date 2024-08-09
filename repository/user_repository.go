package repository

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	FindUserByID(id int) (*User, error)
	SaveUser(user *User) error
}