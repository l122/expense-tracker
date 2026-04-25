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

func GetConfigValue(key string) int {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	value := os.Getenv(key)
	result, _ := strconv.Atoi(value)

	return result
}
