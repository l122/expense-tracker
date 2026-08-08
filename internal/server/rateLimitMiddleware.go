package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"golang.org/x/time/rate"
)

func rateLimitMiddleware(next http.Handler) http.Handler {
	rps, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		log.Fatal("rateLimitMiddleware: error parsing RATE_LIMIT:", err)
	}

	limiter := rate.NewLimiter(rate.Limit(rps), rps)
	var mu sync.Mutex
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		err := limiter.Wait(r.Context())
		mu.Unlock()
		if err != nil {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
