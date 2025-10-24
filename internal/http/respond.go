package http

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON writes the response as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil { return }
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error encoding json: %v", err)
	}
}

// writeError writes a JSON error with code and message.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// methodNotAllowed writes a 405 error.
func methodNotAllowed(w http.ResponseWriter) { writeError(w, http.StatusMethodNotAllowed, "method not allowed") }
