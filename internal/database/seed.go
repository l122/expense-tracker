package database

import (
	"context"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Create the schema if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT UNIQUE,
		username TEXT UNIQUE,
		role TEXT
	);`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return err
	}

	// 2. Double check if data already exists to avoid duplicate seeding
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Users").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// 3. Insert seed users
	for _, u := range users {
		_, err := s.db.ExecContext(ctx, "INSERT INTO Users (name, email, username, role) VALUES (?, ?, ?, ?)", u.Name, u.Email, u.Username, u.Role)
		if err != nil {
			return err
		}
	}
	return nil
}
