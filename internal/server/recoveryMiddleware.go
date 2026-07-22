package server

import (
	"log"
	"net/http"
)

// recoveryMiddleware catches panics and returns 500 errors.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Panic: %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
