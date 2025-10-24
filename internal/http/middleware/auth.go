package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	ctxUserID ctxKey = "user_id"
	ctxUserRole ctxKey = "user_role"
)

// Authenticate validates the Authorization: Bearer <token> header and injects user info into context.
func Authenticate(secret, issuer string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		tok := strings.TrimSpace(auth[len("Bearer "):])
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tok, claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 { return nil, errors.New("invalid alg") }
			return []byte(secret), nil
		})
		if err != nil { http.Error(w, "invalid token", http.StatusUnauthorized); return }
		if iss, _ := claims["iss"].(string); issuer != "" && iss != issuer { http.Error(w, "invalid token issuer", http.StatusUnauthorized); return }
		uidFloat, ok := claims["sub"].(float64)
		if !ok { http.Error(w, "invalid subject", http.StatusUnauthorized); return }
		role, _ := claims["role"].(string)
		ctx := context.WithValue(r.Context(), ctxUserID, int64(uidFloat))
		ctx = context.WithValue(ctx, ctxUserRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the authenticated user id from context.
func GetUserID(r *http.Request) int64 {
	v := r.Context().Value(ctxUserID)
	if v == nil { return 0 }
	if id, ok := v.(int64); ok { return id }
	return 0
}

// GetUserRole returns role from context.
func GetUserRole(r *http.Request) string {
	v := r.Context().Value(ctxUserRole)
	if v == nil { return "" }
	if s, ok := v.(string); ok { return s }
	return ""
}
