package auth

import (
	"fmt"
	"net/http"

	"github.com/triviajon/liteproxy/processor/internal/logging"
)

// WithHeaderAuth wraps a handler and validates the X-Proxy-Auth header.
// Requires that next is not nil and secretToken is not empty.
// Returns an HTTP handler that enforces X-Proxy-Auth header validation, otherwise an error describing which constraint was violated.
func WithHeaderAuth(next http.Handler, secretToken string) (http.Handler, error) {
	if next == nil {
		return nil, fmt.Errorf("next handler must not be nil")
	}
	if secretToken == "" {
		return nil, fmt.Errorf("secretToken must not be empty")
	}

	logging.Infof("Middleware configured")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("X-Proxy-Auth")

		if clientToken == "" || clientToken != secretToken {
			logging.Warnf("Unauthorized request - method=%s path=%s remote_addr=%s", r.Method, r.RequestURI, r.RemoteAddr)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logging.Debugf("Authorized request - method=%s path=%s remote_addr=%s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	}), nil
}
