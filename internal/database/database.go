package database

import (
	"context"
	"os"

	"log"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/l122/expense-tracker/internal/domain"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/supabase-go"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	// Health() map[string]string
	Health() string

	GetUsers(ctx context.Context) ([]domain.User, error)
	GetUserByAuthId(ctx context.Context, userId string) (domain.User, error)
	GetUserById(ctx context.Context, userId int) (domain.User, error)

	EnableUser(ctx context.Context, userId int) (domain.User, error)
	DisableUser(ctx context.Context, userId int) (domain.User, error)
	UpdateAvatar(ctx context.Context, userId int, avatarUrl string) (domain.User, error)
	// DeleteUsers(id int) error

	GetAuthClient() gotrue.Client
}

type service struct {
	db *supabase.Client

	dbUrl string
}

var dbInstance *service

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// 1. Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 2. Fetch the variables using standard os package
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_PUBLISHABLE_KEY")
	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("Supabase credentials are missing from the environment")
	}

	// 3. Initialize your client
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	dbInstance = &service{
		db:    client,
		dbUrl: supabaseURL,
	}

	return dbInstance
}

func (s *service) GetAuthClient() gotrue.Client {
	return s.db.Auth
}
