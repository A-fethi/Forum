package handlers

import (
	"database/sql"
	"forum/app/utils"
	"net/http"
	"time"
)

func RegisterRoutes(DB *sql.DB) {
	http.HandleFunc("/static/", Static)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, DB)
	})
	http.HandleFunc("/api/", utils.RateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		Router(w, r, DB)
	}, 5, 10*time.Second))

}
