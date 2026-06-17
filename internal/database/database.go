package database

import (
	"database/sql"

	"log"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
	"github.com/l122/expense-tracker/internal/domain"
	_ "modernc.org/sqlite"
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
	db *sql.DB
}

var dbInstance *service
var dburl string

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	dburl = os.Getenv("BLUEPRINT_DB_URL")
	if dburl == "" {
		dburl = "expenses.db"
	}

	// Check if the database file already exists to determine if we should seed
	_, err := os.Stat(dburl)
	isFirstRun := os.IsNotExist(err)

	// Ensure the directory for the database file exists
	dir := filepath.Dir(dburl)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("failed to create database directory: %v", err)
		}
	}

	db, err := sql.Open("sqlite", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}

	// Initialize schema and seed if it's the first time the DB is created
	if isFirstRun {
		if err := dbInstance.SeedDb(); err != nil {
			log.Printf("Warning: failed to seed database: %v", err)
		}
	}

	return dbInstance
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}
