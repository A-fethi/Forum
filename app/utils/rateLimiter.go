package utils

import (
	"forum/app/config"
	"forum/app/models"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	requestCounts = make(map[string]int)
	mu            sync.Mutex
)

func RateLimitMiddleware(next http.HandlerFunc, limit int, window time.Duration) http.HandlerFunc {
	go func() {
		for {
			time.Sleep(window) // reset the counts periodically
			mu.Lock()
			requestCounts = make(map[string]int)
			mu.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr) // ip adress of logger

		mu.Lock()
		requestCounts[ip]++
		count := requestCounts[ip]
		mu.Unlock()

		if count > limit {
			config.Logger.Println("IP: ", ip, " Exceed the request limit!!")
			models.SendErrorResponse(w, http.StatusTooManyRequests, "Rate Limit Exceeded. Try Again Later.")
			// http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
