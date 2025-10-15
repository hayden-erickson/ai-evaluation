package context

import (
	"context"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

type contextKey string

const (
	claimsKey contextKey = "claims"
	bankKey   contextKey = "db"
)

// NewBankContext creates a new context with the given bank
func NewBankContext(ctx context.Context, bank interface{}) context.Context {
	return context.WithValue(ctx, bankKey, bank)
}

// BankFromContext returns the Bank value stored in ctx, if any
func BankFromContext(ctx context.Context) (interface{}, bool) {
	bank := ctx.Value(bankKey)
	if bank == nil {
		return nil, false
	}
	return bank, true
}

// GetClaimsFromContext retrieves the claims from the context
func GetClaimsFromContext(ctx context.Context) (*models.Claims, error) {
	claims, ok := ctx.Value(claimsKey).(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}
