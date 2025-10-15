package middleware

import (
	"net/http"

	appContext "github.com/hayden-erickson/ai-evaluation/context"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// BankMiddleware creates middleware that adds a Bank instance to the request context
func BankMiddleware(bank *repository.Bank) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add bank to request context
			ctx := appContext.NewBankContext(r.Context(), bank)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// BankHandlerFunc wraps a HandlerFunc with bank context middleware
func BankHandlerFunc(bank *repository.Bank, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add bank to request context
		ctx := appContext.NewBankContext(r.Context(), bank)
		r = r.WithContext(ctx)

		// Call the handler function
		handlerFunc(w, r)
	}
}
