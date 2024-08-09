package service

import "gorepository/repository"

// UserService handles user-related operations.
type UserService struct {
    Repo repository.UserRepository
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(id int) (*repository.User, error) {
    return s.Repo.FindUserByID(id)
}

// CreateUser saves a new user to the repository.
func (s *UserService) CreateUser(user *repository.User) error {
    return s.Repo.SaveUser(user)
}
