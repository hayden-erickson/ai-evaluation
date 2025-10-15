package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hayden-erickson/ai-evaluation/constants"
	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

// AccessCodeService handles business logic for access codes
type AccessCodeService struct {
	bank *repository.Bank
}

// NewAccessCodeService creates a new AccessCodeService
func NewAccessCodeService(bank *repository.Bank) *AccessCodeService {
	return &AccessCodeService{
		bank: bank,
	}
}

// ValidateUser checks if the user exists and belongs to the same company
func (s *AccessCodeService) ValidateUser(userID int, claims *models.Claims) (*models.BUser, error) {
	user, err := s.bank.GetBUserByID(userID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("internal server error")
	}

	// Check if user belongs to the same company as the requester
	if user.CompanyUUID != claims.CompanyUUID {
		return nil, fmt.Errorf("user not found in company")
	}

	return user, nil
}

// ValidateUserSiteAccess checks if the user has access to the site
func (s *AccessCodeService) ValidateUserSiteAccess(user *models.BUser, claims *models.Claims) error {
	siteFound := false
	for _, siteID := range user.Sites {
		if siteID == strconv.Itoa(claims.CurrentSite) {
			siteFound = true
			break
		}
	}

	if !siteFound {
		return fmt.Errorf("invalid user, missing association with target site")
	}

	return nil
}

// ValidateUnit checks if the unit exists and is in a valid state for access code changes
func (s *AccessCodeService) ValidateUnit(unitID int, claims *models.Claims) (*models.Unit, error) {
	unit, err := s.bank.V2UnitGetById(unitID, claims.CurrentSite)
	if err != nil {
		return nil, fmt.Errorf("unit not found")
	}

	// Check if unit is associated with the site
	if unit.SiteID != claims.CurrentSite {
		return nil, fmt.Errorf("invalid unit, missing association with target site")
	}

	// Check if unit is in a state that allows access code changes
	if unit.RentalState == constants.LockStateOverlock ||
		unit.RentalState == constants.LockStateGatelock ||
		unit.RentalState == constants.LockStatePrelet {
		return nil, fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return unit, nil
}

// ProcessAccessCodeUpdate handles the complete access code update process
func (s *AccessCodeService) ProcessAccessCodeUpdate(unitID int, userID int, accessCode string, claims *models.Claims) error {
	// Create new access code
	var gacs models.GateAccessCodes
	gacs = append(gacs, models.GateAccessCode{
		AccessCode: accessCode,
		UnitID:     unitID,
		UserID:     userID,
		SiteID:     claims.CurrentSite,
		State:      constants.AccessCodeStateSetup,
	})

	// Validate access codes
	err := gacs.Validate()
	if err != nil {
		return fmt.Errorf("internal server error during validation")
	}

	if !gacs[0].IsValid {
		for _, validationMsg := range gacs[0].ValidationMessages {
			if validationMsg == constants.AccessCodeMsgDuplicate {
				return fmt.Errorf("duplicate access code")
			}
		}
		return fmt.Errorf("invalid access code")
	}

	// Handle existing access codes
	err = s.handleExistingAccessCodes(unitID, accessCode, claims)
	if err != nil {
		return err
	}

	// Update database with new access codes
	err = s.bank.UpdateAccessCodes(utils.ConvertToStringSlice(gacs), claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	return nil
}

// handleExistingAccessCodes processes existing access codes for revocation
func (s *AccessCodeService) handleExistingAccessCodes(unitID int, newAccessCode string, claims *models.Claims) error {
	units := []int{unitID}

	// Get existing access codes
	revokeGacs, err := s.bank.GetCodesForUnits(units, claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	var revokeUnits []int
	var updateRemoveGacs models.GateAccessCodes

	// Process existing codes
	for _, revokeGac := range revokeGacs {
		// If existing code is the same as new code and is active/setup/pending, no changes needed
		if (revokeGac.State == constants.AccessCodeStateActive ||
			revokeGac.State == constants.AccessCodeStateSetup ||
			revokeGac.State == constants.AccessCodeStatePending) &&
			revokeGac.AccessCode == newAccessCode {
			return nil // No change needed
		}

		// Mark for removal if not already in a removal state
		if revokeGac.State != constants.AccessCodeStateInactive &&
			revokeGac.State != constants.AccessCodeStateRemoved &&
			revokeGac.State != constants.AccessCodeStateRemoving &&
			revokeGac.State != constants.AccessCodeStateOverlocking &&
			revokeGac.State != constants.AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
			revokeGac.State = constants.AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, revokeGac)
		}
	}

	// Update database to mark old codes for removal
	if len(updateRemoveGacs) > 0 {
		err = s.bank.UpdateAccessCodes(utils.ConvertToStringSlice(updateRemoveGacs), claims.CurrentSite)
		if err != nil {
			return fmt.Errorf("internal server error updating access codes")
		}
	}

	return nil
}

// UserAccessCodeEdit records the access code edit activity
func UserAccessCodeEdit(userID int, user *models.BUser, siteUUID string, timestamp time.Time) (*string, error) {
	// This would typically record the activity in an activity log
	return nil, nil
}
