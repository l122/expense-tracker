package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

var (
	AppVersion = "dev" // Default to 'dev' if not set
	CommitSHA  = "commit-sha-not-set"
	BuildDate  = "build-date"

	//go:embed *.gohtml
	templates embed.FS
)

// func version(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "hello\n")
// 	fmt.Fprintf(w, "AppVersion: %v\n", AppVersion)
// 	fmt.Fprintf(w, "CommitSHA: %v\n", CommitSHA)
// 	fmt.Fprintf(w, "BuildDate: %v\n", BuildDate)
// }

func main() {

	// Debug: figure out how to display it on the main page
	log.Printf("AppVersion: %v", AppVersion)
	log.Printf("CommitSHA: %v", CommitSHA)
	log.Printf("BuildDate: %v", BuildDate)

	tmpl, err := template.ParseFS(templates, "*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	indexView := NewIndexView(tmpl)

	hander, err := NewHandler(indexView)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Creating server")

	// http.HandleFunc("/", version)

	// http.ListenAndServe(":8080", nil)

	if err := http.ListenAndServe(":8080", hander); err != nil {
		log.Fatal(err)
	}
}
