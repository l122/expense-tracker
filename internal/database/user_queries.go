package database

// import (
// 	"context"
// 	"time"

// 	"github.com/l122/expense-tracker/internal/domain"
// )

// func (s *service) GetUserByEmail(email string) (domain.User, error) {
// 	// TODO:implement
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var user domain.User
// 	query := "SELECT id, email, role FROM users WHERE email = $1"
// 	err := s.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.Role)
// 	if err != nil {
// 		return domain.User{}, err
// 	}

// 	return user, nil
// }

// func (s *service) CreateUser(name, email, username, role string) (domain.User, error) {
// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()

// 	var user domain.User
// 	// query := "INSERT INTO users"

// 	return user, nil
// }
