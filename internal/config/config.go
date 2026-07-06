package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion string
	CommitSHA  string
	BuildDate  string
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("No .env file found")
		}
	}

	if err := godotenv.Overload("../../.env.local"); err != nil {
		err = godotenv.Overload(".env.local")
		if err != nil {
			log.Println("No .env.local file found")
		}
	}

	if err != nil {
		log.Fatalf("No .env or .env.local files found")
	}
}

func GetConfigValue(key string) int {
	value := os.Getenv(key)
	result, _ := strconv.Atoi(value)

	return result
}
