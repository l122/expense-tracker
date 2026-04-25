package main

import (
	"log"

	"github.com/l122/expense-tracker/internal/config"
	"github.com/l122/expense-tracker/internal/server"
)

var (
	AppVersion = "dev" // Default to 'dev' if not set
	CommitSHA  = "commit-sha-not-set"
	BuildDate  = "build-date"
)

func main() {
	cfg := &config.Config{
		AppVersion: AppVersion,
		CommitSHA:  CommitSHA,
		BuildDate:  BuildDate,
	}

	server := server.NewServer(cfg)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
