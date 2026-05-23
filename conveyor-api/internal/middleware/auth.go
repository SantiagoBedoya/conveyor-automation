package middleware

import (
	"net/http"
	"os"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" {
			http.Error(w, `{"error":"missing api key"}`, http.StatusUnauthorized)
			return
		}
		if key != os.Getenv("API_KEY") {
			http.Error(w, `{"error":"invalid api key"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
