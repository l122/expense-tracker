package database

import (
	"context"
	"fmt"

	"github.com/l122/expense-tracker/internal/domain"
	"github.com/l122/expense-tracker/pkgs/token"
)

func (s *service) GetUsers(ctx context.Context) ([]domain.User, error) {
	if token, ok := token.FromContext(ctx); ok {
		fmt.Println(token)
		// request.Set("Bearer", token)
	}

	// TODO:get []User from db
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// rows, err := s.db.Query(ctx, "SELECT id, name, email, username, role FROM Users")
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	var result []domain.User
	// for rows.Next() {
	// 	user := domain.User{}
	// 	if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Role); err != nil {
	// 		return nil, err
	// 	}

	// 	result = append(result, user)
	// }
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	return result, nil
}
