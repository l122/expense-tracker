package database

// import (
// 	"context"
// 	"time"

// 	"github.com/l122/expense-tracker/internal/domain"
// )

// func (s *service) GetUsers() ([]domain.User, error) {
// 	// TODO:get []User from db
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
// 	defer cancel()

// 	rows, err := s.db.Query(ctx, "SELECT id, name, email, username, role FROM Users")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var result []domain.User
// 	for rows.Next() {
// 		user := domain.User{}
// 		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Role); err != nil {
// 			return nil, err
// 		}

// 		result = append(result, user)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
