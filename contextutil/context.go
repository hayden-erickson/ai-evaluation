package contextutil

import (
	"context"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/database"
	"github.com/hayden-erickson/ai-evaluation/models"
)

type contextKey string

const (
	claimsKey contextKey = "claims"
	bankKey   contextKey = "db"
)

// NewBankContext creates a new context with the bank instance
func NewBankContext(ctx context.Context, bank *database.Bank) context.Context {
	return context.WithValue(ctx, bankKey, bank)
}

// BankFromContext returns the Bank value stored in ctx, if any
func BankFromContext(ctx context.Context) (*database.Bank, bool) {
	bank, ok := ctx.Value(bankKey).(*database.Bank)
	return bank, ok
}

// GetClaimsFromContext retrieves claims from the context
func GetClaimsFromContext(ctx context.Context) (*models.Claims, error) {
	claims, ok := ctx.Value(claimsKey).(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}
