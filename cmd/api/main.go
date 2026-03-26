package main

import (
	"log"
)

var (
	AppVersion = "dev" // Default to 'dev' if not set
	CommitSHA  = "n/a"
	BuildDate  = "n/a"
)

func main() {

	// Debug: figure out how to display it on the main page
	log.Printf("AppVersion: %v", AppVersion)
	log.Printf("CommitSHA: %v", CommitSHA)
	log.Printf("BuildDate: %v", BuildDate)

	log.Println("Creating server")

}
