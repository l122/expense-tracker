package database

import (
	"context"

	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/l122/expense-tracker/internal/domain"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	GetUsers() ([]domain.User, error)
	DeleteUsers(id int) error
	SeedDb() error
}

type service struct {
	db *pgxpool.Pool
}

var dbInstance *service
var connStr string

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	dbInstance = &service{
		db: db,
	}

	// Initialize schema and seed if it's the first time the DB is created
	if err := dbInstance.SeedDb(); err != nil {
		log.Printf("Warning: failed to seed database: %v", err)
	}

	return dbInstance
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", connStr)
	return nil
	// return s.db.Close()
}
