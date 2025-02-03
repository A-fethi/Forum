package handlers

import (
	"database/sql"
	"forum/app/utils"
	"net/http"
	"time"
)

func RegisterRoutes(DB *sql.DB) {
	http.HandleFunc("/static/", Static)

	http.HandleFunc("/", utils.RateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, DB)
	}, 5, 1*time.Second))
	http.HandleFunc("/api/", utils.RateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		Router(w, r, DB)
	}, 5, 1*time.Second))

}
