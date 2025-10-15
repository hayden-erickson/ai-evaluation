package context

import (
	"context"
	"fmt"

	"github.com/hayden-erickson/ai-evaluation/constants"
	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// NewBankContext creates a new context with a Bank instance
func NewBankContext(ctx context.Context, bank *repository.Bank) context.Context {
	return context.WithValue(ctx, constants.BankKey, bank)
}

// BankFromContext returns the Bank value stored in ctx, if any
func BankFromContext(ctx context.Context) (*repository.Bank, bool) {
	bank, ok := ctx.Value(constants.BankKey).(*repository.Bank)
	return bank, ok
}

// GetClaimsFromContext retrieves claims from the request context
func GetClaimsFromContext(ctx context.Context) (*models.Claims, error) {
	claims, ok := ctx.Value(constants.ClaimsKey).(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}
