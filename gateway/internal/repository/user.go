package repository

import (
	"fmt"
)

type User struct {
	Login    string
	Password string
}

type UserRepo interface {
	GetPassword(login string) (string, error)
}

type InMemoryUserRepo struct {
	users map[string]string
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	users := map[string]string{
		"alice": "password123",
		"bob":   "qwerty",
	}
	return &InMemoryUserRepo{users: users}
}

func (r *InMemoryUserRepo) GetPassword(login string) (string, error) {
	pw, ok := r.users[login]
	if !ok {
		return "", fmt.Errorf("user %q not found", login)
	}
	return pw, nil
}
