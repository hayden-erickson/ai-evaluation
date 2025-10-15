package models

import "fmt"

// BUser represents a business user
type BUser struct {
	CompanyUUID string
	Id          int
	Sites       []string
}

// Unit represents a rental unit
type Unit struct {
	SiteID      int
	RentalState string
}

// GateAccessCode represents an access code for a gate
type GateAccessCode struct {
	AccessCode         string
	UnitID             int
	UserID             int
	SiteID             int
	State              string
	IsValid            bool
	ValidationMessages []string
}

// GateAccessCodes is a slice of GateAccessCode
type GateAccessCodes []GateAccessCode

// Validate validates all gate access codes in the collection
func (gacs GateAccessCodes) Validate(validator GateAccessCodeValidator) error {
	if len(gacs) == 0 {
		return fmt.Errorf("no access codes provided")
	}

	for i := range gacs {
		gacs[i].IsValid = true
	}
	return nil
}

// GateAccessCodeValidator interface for validating gate access codes
type GateAccessCodeValidator interface {
	// Add validation methods as needed
}

// Claims represents JWT claims or session claims
type Claims struct {
	CompanyUUID     string
	CurrentSite     int
	CurrentSiteUUID string
	UserID          int
}
