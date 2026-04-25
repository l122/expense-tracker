package main

import (
	"html/template"
	"log"
	"net/http"
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

	templates, err := template.ParseFS(Files, "*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	hander, err := NewHandler(NewIndexView(templates))
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", hander); err != nil {
		log.Fatal(err)
	}
}
