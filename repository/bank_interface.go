package repository

import (
	"context"

	"github.com/hayden-erickson/ai-evaluation/models"
)

// Bank represents the database layer for business operations
type Bank struct{}

// NewBank creates a new Bank instance
func NewBank() *Bank {
	return &Bank{}
}

// GetBUserByID retrieves a business user by ID
func (b *Bank) GetBUserByID(BUserID int) (*models.BUser, error) {
	return &models.BUser{}, nil
}

// V2UnitGetById retrieves a unit by ID and site ID
func (b *Bank) V2UnitGetById(unitID int, siteID int) (*models.Unit, error) {
	return &models.Unit{}, nil
}

// GetCodesForUnits retrieves access codes for given units and site
func (b *Bank) GetCodesForUnits(units []int, siteID int) ([]models.GateAccessCode, error) {
	return []models.GateAccessCode{}, nil
}

// UpdateAccessCodes updates access codes in the database
func (b *Bank) UpdateAccessCodes(codes []string, siteID int) error {
	return nil
}

// NewCommandCenterClient creates a new command center client for the given site
func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClient {
	return &CommandCenter{}
}
