package repository

import "errors"

type MockUserRepository struct {
    Users map[int]*User
    Err   error
}

func (m *MockUserRepository) FindUserByID(id int) (*User, error) {
    if m.Err != nil {
        return nil, m.Err
    }
    user, exists := m.Users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}

func (m *MockUserRepository) SaveUser(user *User) error {
    if m.Err != nil {
        return m.Err
    }
    m.Users[user.ID] = user
    return nil
}
