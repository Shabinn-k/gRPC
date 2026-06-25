package repository

import (
	"errors"
	"sync"
)

type User struct {
	ID       int32
	Email    string
	Password string
}

type UserRepository struct {
	users map[int32]*User
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	repo := &UserRepository{
		users: make(map[int32]*User),
	}

	repo.users[1] = &User{
		ID:       1,
		Email:    "user1@example.com",
		Password: "password123",
	}
	repo.users[2] = &User{
		ID:       2,
		Email:    "user2@example.com",
		Password: "password456",
	}

	return repo
}

func (repo *UserRepository) GetUserByEmail(email string) (*User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	for _, user := range repo.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (repo *UserRepository) GetUserByID(id int32) (*User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, ok := repo.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (repo *UserRepository) ValidateUser(email, password string) (int32, error) {
	user, err := repo.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}

	if user.Password != password {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}