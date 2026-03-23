package auth

import (
	"net/http"
)

// WithHeaderAuth wraps a handler and validates the X-Proxy-Auth header
func WithHeaderAuth(next http.Handler, secretToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("X-Proxy-Auth")

		if clientToken == "" || clientToken != secretToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
