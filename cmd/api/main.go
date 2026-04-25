package main

import (
	"log"

	"github.com/l122/expense-tracker/internal/server"
)

var (
	AppVersion = "dev" // Default to 'dev' if not set
	CommitSHA  = "commit-sha-not-set"
	BuildDate  = "build-date"
)

func main() {

	// Debug: figure out how to display it on the main page
	log.Printf("AppVersion: %v", AppVersion)
	log.Printf("CommitSHA: %v", CommitSHA)
	log.Printf("BuildDate: %v", BuildDate)

	server := server.NewServer()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
