package service

import (
	"gorepository/repository" // Adjust the import path as needed
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
    // Setup mock repository
    mockRepo := &repository.MockUserRepository{
        Users: map[int]*repository.User{
            1: {ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
        },
    }
    
    service := &UserService{Repo: mockRepo}
    
    // Test getting an existing user
    user, err := service.GetUser(1)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
    
    // Test getting a non-existing user
    user, err = service.GetUser(2)
    assert.Error(t, err)
    assert.Nil(t, user)
}

func TestCreateUser(t *testing.T) {
    // Setup mock repository
    mockRepo := &repository.MockUserRepository{
        Users: map[int]*repository.User{},
    }
    
    service := &UserService{Repo: mockRepo}
    
    // Test creating a user
    user := &repository.User{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"}
    err := service.CreateUser(user)
    assert.NoError(t, err)
    
    // Verify that the user was saved
    savedUser, err := mockRepo.FindUserByID(2)
    assert.NoError(t, err)
    assert.NotNil(t, savedUser)
    assert.Equal(t, "Jane Doe", savedUser.Name)
}
