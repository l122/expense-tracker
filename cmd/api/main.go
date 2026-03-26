package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	AppVersion = "dev" // Default to 'dev' if not set
	CommitSHA  = "commit-sha-not-set"
	BuildDate  = "build-date"
)

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
	fmt.Fprintf(w, "AppVersion: %v\n", AppVersion)
	fmt.Fprintf(w, "CommitSHA: %v\n", CommitSHA)
	fmt.Fprintf(w, "BuildDate: %v\n", BuildDate)
}

func main() {

	// Debug: figure out how to display it on the main page
	log.Printf("AppVersion: %v", AppVersion)
	log.Printf("CommitSHA: %v", CommitSHA)
	log.Printf("BuildDate: %v", BuildDate)

	log.Println("Creating server")

	http.HandleFunc("/", version)

	http.ListenAndServe(":8080", nil)
}
