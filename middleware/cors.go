package middleware

import (
	"net/http"
	"os"
)

// CorsMiddleware adds CORS headers to allow the frontend to call the API from a different origin.
// Allowed origins can be configured via the ALLOWED_ORIGINS environment variable (comma-separated).
// For local development, it defaults to http://localhost:5173 (Vite default port).
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:5173"
		}

		origin := r.Header.Get("Origin")
		// In a simple implementation, reflect the origin if it's contained in allowedOrigins; otherwise, do not set.
		// For a small project, a contains check is sufficient. In production, parse and validate.
		if origin != "" && (origin == allowedOrigins || containsCommaSeparated(allowedOrigins, origin)) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// containsCommaSeparated returns true if value is present in a comma-separated list (trim-aware).
func containsCommaSeparated(list string, value string) bool {
	start := 0
	for i := 0; i <= len(list); i++ {
		if i == len(list) || list[i] == ',' {
			segment := trimSpace(list[start:i])
			if segment == value {
				return true
			}
			start = i + 1
		}
	}
	return false
}

// trimSpace is a lightweight ASCII-only space trimmer to avoid importing strings for a small helper.
func trimSpace(s string) string {
	// Left trim
	start := 0
	for start < len(s) {
		c := s[start]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			start++
			continue
		}
		break
	}
	// Right trim
	end := len(s)
	for end > start {
		c := s[end-1]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v' {
			end--
			continue
		}
		break
	}
	return s[start:end]
}
