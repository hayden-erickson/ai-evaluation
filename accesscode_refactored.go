package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Input data structure for access code edit requests
type AccessCodeEditInput struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

// BUserService defines the interface for user-related operations
type BUserService interface {
	GetBUserByID(userID int) (*BUser, error)
}

// UnitService defines the interface for unit-related operations
type UnitService interface {
	GetUnitByID(unitID, siteID int) (*Unit, error)
}

// AccessCodeService defines the interface for access code operations
type AccessCodeService interface {
	GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodes(codes []string, siteID int) error
	ValidateAccessCodes(codes GateAccessCodes) error
}

// CommandCenterService defines the interface for command center operations
type CommandCenterService interface {
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClient
}

// CommandCenterClient interface for abstracting command center operations
type CommandCenterClient interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

// ActivityService defines the interface for activity logging
type ActivityService interface {
	LogUserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error)
}

// AccessCodeHandler encapsulates dependencies for the handler
type AccessCodeHandler struct {
	userService          BUserService
	unitService          UnitService
	accessCodeService    AccessCodeService
	commandCenterService CommandCenterService
	activityService      ActivityService
}

// NewAccessCodeHandler creates a new instance of AccessCodeHandler
func NewAccessCodeHandler(
	userService BUserService,
	unitService UnitService,
	accessCodeService AccessCodeService,
	commandCenterService CommandCenterService,
	activityService ActivityService,
) *AccessCodeHandler {
	return &AccessCodeHandler{
		userService:          userService,
		unitService:          unitService,
		accessCodeService:    accessCodeService,
		commandCenterService: commandCenterService,
		activityService:      activityService,
	}
}

// parseInput parses and validates the request input
func (h *AccessCodeHandler) parseInput(r *http.Request) (*AccessCodeEditInput, error) {
	var input AccessCodeEditInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	// Convert UUIDs to IDs if necessary
	if input.UserUUID != "" {
		userID, err := strconv.Atoi(input.UserUUID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %s", input.UserUUID)
		}
		input.UserID = userID
	} else {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	// Convert unit UUIDs to IDs
	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %s", unitUUID)
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	return &input, nil
}

// validateUser checks if the user exists and belongs to the right company/site
func (h *AccessCodeHandler) validateUser(userID int, companyUUID string, siteID int) (*BUser, error) {
	user, err := h.userService.GetBUserByID(userID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	// Check if user belongs to the same company
	if user.CompanyUUID != companyUUID {
		return nil, fmt.Errorf("user not found in company")
	}

	// Check if user is associated with the site
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

// validateUnit checks if the unit exists, belongs to the site, and is in a valid state
func (h *AccessCodeHandler) validateUnit(unitID int, siteID int) error {
	unit, err := h.unitService.GetUnitByID(unitID, siteID)
	if err != nil {
		return fmt.Errorf("unit not found")
	}

	if unit.SiteID != siteID {
		return fmt.Errorf("invalid unit, missing association with target site")
	}

	// Check if unit is in a state that allows access code changes
	if unit.RentalState == LOCK_STATE_OVERLOCK ||
		unit.RentalState == LOCK_STATE_GATELOCK ||
		unit.RentalState == LOCK_STATE_PRELET {
		return fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return nil
}

// processUnitAccessCode processes the access code edit for a single unit
func (h *AccessCodeHandler) processUnitAccessCode(ctx context.Context, unitID int, userID int, accessCode string, siteID int) error {
	// Validate unit
	if err := h.validateUnit(unitID, siteID); err != nil {
		return err
	}

	// Create new access code
	gacs := GateAccessCodes{
		{
			AccessCode: accessCode,
			UnitID:     unitID,
			UserID:     userID,
			SiteID:     siteID,
			State:      AccessCodeStateSetup,
		},
	}

	// Validate access code
	if err := h.accessCodeService.ValidateAccessCodes(gacs); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Check validation results
	if !gacs[0].IsValid {
		for _, validationMsg := range gacs[0].ValidationMessages {
			if validationMsg == AccessCodeMsgDuplicate {
				return fmt.Errorf("duplicate access code")
			}
		}
		return fmt.Errorf("invalid access code")
	}

	units := []int{unitID}

	// Get existing access codes for the unit
	existingCodes, err := h.accessCodeService.GetCodesForUnits(units, siteID)
	if err != nil {
		return fmt.Errorf("error retrieving existing access codes: %w", err)
	}

	// Process existing codes
	revokeUnits, updateRemoveGacs, skipUpdate := h.processExistingCodes(existingCodes, accessCode)

	if skipUpdate {
		return nil // Code already exists and is active/setup/pending
	}

	// Update codes marked for removal
	if len(updateRemoveGacs) > 0 {
		err = h.accessCodeService.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), siteID)
		if err != nil {
			return fmt.Errorf("error updating access codes: %w", err)
		}
	}

	// Update new access code
	err = h.accessCodeService.UpdateAccessCodes(convertToStringSlice(gacs), siteID)
	if err != nil {
		return fmt.Errorf("error updating access codes: %w", err)
	}

	// Revoke old codes via command center
	if len(revokeUnits) > 0 {
		uniqueRevokeUnits, _ := uniqueIntSlice(revokeUnits)
		cc := h.commandCenterService.NewCommandCenterClient(siteID, ctx)
		err = cc.RevokeAccessCodes(uniqueRevokeUnits, make(map[string]struct{}, 0))
		if err != nil {
			return fmt.Errorf("failed to revoke previous access codes: %w", err)
		}
	}

	// Set new code via command center
	cc := h.commandCenterService.NewCommandCenterClient(siteID, ctx)
	err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
	if err != nil {
		return fmt.Errorf("failed to set access codes: %w", err)
	}

	return nil
}

// processExistingCodes processes the existing access codes for a unit
// returns revokeUnits, updateRemoveGacs, and a boolean indicating if we should skip further updates
func (h *AccessCodeHandler) processExistingCodes(existingCodes []GateAccessCode, newAccessCode string) ([]int, GateAccessCodes, bool) {
	var revokeUnits []int
	var updateRemoveGacs GateAccessCodes

	for _, code := range existingCodes {
		// If code already exists and is active/setup/pending, skip further processing
		if (code.State == AccessCodeStateActive ||
			code.State == AccessCodeStateSetup ||
			code.State == AccessCodeStatePending) &&
			code.AccessCode == newAccessCode {
			return nil, nil, true
		}

		// If code is not already marked as inactive/removed/etc, mark for removal
		if code.State != AccessCodeStateInactive &&
			code.State != AccessCodeStateRemoved &&
			code.State != AccessCodeStateRemoving &&
			code.State != AccessCodeStateOverlocking &&
			code.State != AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, code.UnitID)
			code.State = AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, code)
		}
	}

	return revokeUnits, updateRemoveGacs, false
}

// AccessCodeEditHandler handles HTTP requests for editing access codes
func (h *AccessCodeHandler) AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, err := getClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse input
	input, err := h.parseInput(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user
	user, err := h.validateUser(input.UserID, claims.CompanyUUID, claims.CurrentSite)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "user not found in company" ||
			err.Error() == "invalid user, missing association with target site" {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	// Process each unit
	for _, unitID := range input.UnitIDs {
		if unitID == 0 {
			continue
		}

		err := h.processUnitAccessCode(r.Context(), unitID, input.UserID, input.AccessCode, claims.CurrentSite)
		if err != nil {
			status := http.StatusInternalServerError

			// Handle different error cases with appropriate status codes
			errMsg := err.Error()
			
			switch {
			case errMsg == "unit not found":
				status = http.StatusNotFound
			case errMsg == "invalid unit, missing association with target site" || strings.Contains(errMsg, "access code changes not allowed"):
				status = http.StatusForbidden
			case errMsg == "duplicate access code":
				status = http.StatusConflict
			case errMsg == "invalid access code":
				status = http.StatusBadRequest
			case strings.Contains(errMsg, "error retrieving existing access codes") || 
				 strings.Contains(errMsg, "failed to revoke previous access codes") ||
				 strings.Contains(errMsg, "failed to set access codes") ||
				 strings.Contains(errMsg, "error updating access codes"):
				status = http.StatusInternalServerError
			}

			http.Error(w, err.Error(), status)
			return
		}
	}

	// Log activity
	_, err = h.activityService.LogUserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// BankAdapter implements all the interfaces using the Bank struct
type BankAdapter struct {
	bank *Bank
}

func NewBankAdapter(bank *Bank) *BankAdapter {
	return &BankAdapter{bank: bank}
}

// BUserService implementation
func (a *BankAdapter) GetBUserByID(userID int) (*BUser, error) {
	return a.bank.GetBUserByID(userID)
}

// UnitService implementation
func (a *BankAdapter) GetUnitByID(unitID, siteID int) (*Unit, error) {
	return a.bank.V2UnitGetById(unitID, siteID)
}

// AccessCodeService implementation
func (a *BankAdapter) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return a.bank.GetCodesForUnits(units, siteID)
}

func (a *BankAdapter) UpdateAccessCodes(codes []string, siteID int) error {
	return a.bank.UpdateAccessCodes(codes, siteID)
}

func (a *BankAdapter) ValidateAccessCodes(codes GateAccessCodes) error {
	return codes.Validate(a.bank)
}

// CommandCenterService implementation
func (a *BankAdapter) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClient {
	return a.bank.NewCommandCenterClient(siteID, ctx)
}

// ActivityAdapter implements ActivityService
type ActivityAdapter struct{}

func NewActivityAdapter() *ActivityAdapter {
	return &ActivityAdapter{}
}

// ActivityService implementation
func (a *ActivityAdapter) LogUserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	return UserAccessCodeEdit(userID, user, siteUUID, timestamp)
}

// CreateAccessCodeHandlerFromBank creates a handler using Bank as the backend
func CreateAccessCodeHandlerFromBank(bank *Bank) *AccessCodeHandler {
	adapter := NewBankAdapter(bank)
	activityAdapter := NewActivityAdapter()

	return NewAccessCodeHandler(
		adapter,         // BUserService
		adapter,         // UnitService
		adapter,         // AccessCodeService
		adapter,         // CommandCenterService
		activityAdapter, // ActivityService
	)
}

// AccessCodeEditHandlerFunc creates an http.HandlerFunc that uses the bank from context
func AccessCodeEditHandlerFunc(w http.ResponseWriter, r *http.Request) {
	bank, ok := BankFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve bank from context", http.StatusInternalServerError)
		return
	}

	handler := CreateAccessCodeHandlerFromBank(bank)
	handler.AccessCodeEditHandler(w, r)
}
