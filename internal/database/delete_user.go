package database

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// func (s *service) DeleteUsers(id int) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := "DELETE FROM Users WHERE id = ?"
// 	result, err := s.db.Exec(ctx, query, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute delete: %w", err)
// 	}

// 	rowsAffected := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("No user found with ID %d\n", id)
// 	}
// 	return nil
// }
