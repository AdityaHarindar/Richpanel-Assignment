package transport

import (
	"net/http"

	"github.com/gorilla/mux"
)

// AuthMiddleware is an authentication middleware that checks if the API key in the header is authorized
func AuthMiddleware(expectedKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-API-Key") != expectedKey {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
