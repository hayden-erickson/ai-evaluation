package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AccessCodeService struct {
	bank BankService
}

type accessCodeInput struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

func NewAccessCodeService(bank BankService) *AccessCodeService {
	return &AccessCodeService{bank: bank}
}

// validateUser checks if the user exists and has proper permissions
func (s *AccessCodeService) validateUser(userID int, companyUUID string, siteID int) (*BUser, error) {
	user, err := s.bank.GetBUserByID(userID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("internal error: %v", err)
	}

	if user.CompanyUUID != companyUUID {
		return nil, fmt.Errorf("user not found in company")
	}

	siteFound := false
	for _, userSite := range user.Sites {
		if userSite == strconv.Itoa(siteID) {
			siteFound = true
			break
		}
	}

	if !siteFound {
		return nil, fmt.Errorf("invalid user, missing association with target site")
	}

	return user, nil
}

// validateUnit checks if the unit exists and can have its access code modified
func (s *AccessCodeService) validateUnit(unitID int, siteID int) error {
	unit, err := s.bank.V2UnitGetById(unitID, siteID)
	if err != nil {
		return fmt.Errorf("unit not found: %v", err)
	}

	if unit.SiteID != siteID {
		return fmt.Errorf("invalid unit, missing association with target site")
	}

	if unit.RentalState == LOCK_STATE_OVERLOCK ||
		unit.RentalState == LOCK_STATE_GATELOCK ||
		unit.RentalState == LOCK_STATE_PRELET {
		return fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return nil
}

// updateUnitAccessCode handles the access code update process for a single unit
func (s *AccessCodeService) updateUnitAccessCode(ctx context.Context, unitID int, userID int, siteID int, accessCode string) error {
	// Create and validate new access code
	gacs := GateAccessCodes{
		{
			AccessCode: accessCode,
			UnitID:     unitID,
			UserID:     userID,
			SiteID:     siteID,
			State:      AccessCodeStateSetup,
		},
	}

	if err := gacs.Validate(s.bank); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if !gacs[0].IsValid {
		for _, msg := range gacs[0].ValidationMessages {
			if msg == AccessCodeMsgDuplicate {
				return fmt.Errorf("duplicate access code")
			}
		}
		return fmt.Errorf("invalid access code")
	}

	units := []int{unitID}

	// Handle existing access codes
	revokeGacs, err := s.bank.GetCodesForUnits(units, siteID)
	if err != nil {
		return fmt.Errorf("error getting existing codes: %v", err)
	}

	var revokeUnits []int
	var updateRemoveGacs GateAccessCodes

	for _, revokeGac := range revokeGacs {
		if (revokeGac.State == AccessCodeStateActive ||
			revokeGac.State == AccessCodeStateSetup ||
			revokeGac.State == AccessCodeStatePending) &&
			revokeGac.AccessCode == accessCode {
			return nil // No changes needed
		}

		if revokeGac.State != AccessCodeStateInactive &&
			revokeGac.State != AccessCodeStateRemoved &&
			revokeGac.State != AccessCodeStateRemoving &&
			revokeGac.State != AccessCodeStateOverlocking &&
			revokeGac.State != AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
			revokeGac.State = AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, revokeGac)
		}
	}

	// Update access codes in database
	if err := s.bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), siteID); err != nil {
		return fmt.Errorf("error updating old access codes: %v", err)
	}

	if err := s.bank.UpdateAccessCodes(convertToStringSlice(gacs), siteID); err != nil {
		return fmt.Errorf("error updating new access codes: %v", err)
	}

	// Update access codes in command center
	cc := s.bank.NewCommandCenterClient(siteID, ctx)

	revokeUnits, _ = uniqueIntSlice(revokeUnits)
	if err := cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{})); err != nil {
		return fmt.Errorf("failed to revoke previous access codes: %v", err)
	}

	if err := cc.SetAccessCodes(units, make(map[string]struct{})); err != nil {
		return fmt.Errorf("failed to set access codes: %v", err)
	}

	return nil
}

func AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bank, ok := BankFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve bank from context", http.StatusInternalServerError)
		return
	}

	// Convert Bank to BankService
	var bankService BankService = bank
	service := NewAccessCodeService(bankService)

	var input accessCodeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Handle UUID/ID conversions
	if input.UserUUID != "" {
		input.UserID, err = strconv.Atoi(input.UserUUID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid UUID: %s", input.UserUUID), http.StatusBadRequest)
			return
		}
	}

	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid UUID: %s", unitUUID), http.StatusBadRequest)
			return
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	// Validate user
	user, err := service.validateUser(input.UserID, claims.CompanyUUID, claims.CurrentSite)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "user not found in company", "invalid user, missing association with target site":
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Process each unit
	for _, unitID := range input.UnitIDs {
		if unitID == 0 {
			continue
		}

		if err := service.validateUnit(unitID, claims.CurrentSite); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if err := service.updateUnitAccessCode(r.Context(), unitID, input.UserID, claims.CurrentSite, input.AccessCode); err != nil {
			switch err.Error() {
			case "duplicate access code":
				http.Error(w, err.Error(), http.StatusConflict)
			case "invalid access code":
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			}
			return
		}
	}

	// Record activity
	_, err = UserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
