package middleware

import (
	"net/http"

	"github.com/hayden-erickson/ai-evaluation/internal/clients"
	appContext "github.com/hayden-erickson/ai-evaluation/internal/context"
)

// BankMiddleware creates middleware that injects the bank into the request context
func BankMiddleware(bank clients.BankInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add bank to request context
			ctx := appContext.NewBankContext(r.Context(), bank)
			r = r.WithContext(ctx)

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// WithBankMiddleware wraps an http.HandlerFunc with the bank middleware
func WithBankMiddleware(bank clients.BankInterface, handler http.HandlerFunc) http.HandlerFunc {
	return BankMiddleware(bank)(handler).ServeHTTP
}
