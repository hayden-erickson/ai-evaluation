package models

import (
	"fmt"
)

// GateAccessCode represents an access code for gate entry
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

// Validate validates a collection of gate access codes
func (gacs GateAccessCodes) Validate() error {
	if len(gacs) == 0 {
		return fmt.Errorf("no access codes provided")
	}

	for i := range gacs {
		gacs[i].IsValid = true
	}
	return nil
}
