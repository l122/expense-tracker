package database

import "github.com/l122/expense-tracker/internal/domain"

func (s *service) GetUserByEmail(email string) (domain.User, error) {
	// TODO:implement
	return domain.User{}, nil
}

func (s *service) CreateUser(name, email, username, role string) (domain.User, error) {
	// TODO:implement
	return domain.User{}, nil
}
