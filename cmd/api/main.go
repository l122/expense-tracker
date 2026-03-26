package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	AppVersion = "dev" // Default to 'dev' if not set
	CommitSHA  = "n/a"
	BuildDate  = "n/a"
)

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
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
