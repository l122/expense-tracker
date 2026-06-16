package database

import (
	"github.com/l122/expense-tracker/internal/domain"
)

var users = []domain.User{
	{
		Id:       1,
		Name:     "Alice Johnson",
		Email:    "alice.johnson@example.com",
		Username: "usr_982341",
		Role:     "admin",
	},
	{
		Id:       2,
		Name:     "Bob Smith",
		Email:    "bob.smith@provider.net",
		Username: "usr_441209",
		Role:     "user",
	},
	{
		Id:       3,
		Name:     "Charlie Davis",
		Email:    "charlie.d@service.org",
		Username: "usr_112233",
		Role:     "user",
	},
}

func (s *service) SeedDb() error {
	// TODO: seed users
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()

	// rows, err := s.db.QueryContext(ctx, "SELECT id, name, email, username, role FROM Users")
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var result []domain.User
	// for rows.Next() {
	// 	user := domain.User{}
	// 	if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Name, &user.Role); err != nil {
	// 		return nil, err
	// 	}

	// 	result = append(result, user)
	// }
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }

	// return result, nil
	return nil
}
