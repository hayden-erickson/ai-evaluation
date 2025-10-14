package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BUser struct {
	CompanyUUID string
	Id          int
	Sites       []string
}
type Unit struct {
	SiteID      int
	RentalState string
}

type GateAccessCode struct {
	AccessCode         string
	UnitID             int
	UserID             int
	SiteID             int
	State              string
	IsValid            bool
	ValidationMessages []string
}

type GateAccessCodes []GateAccessCode

// Interfaces for dependency injection and mocking
type BankInterface interface {
	GetBUserByID(BUserID int) (*BUser, error)
	V2UnitGetById(unitID int, siteID int) (*Unit, error)
	GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodes(codes []string, siteID int) error
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientInterface
}

type CommandCenterClientInterface interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

type CommandCenterClient struct{}
type Bank struct{}

const (
	AccessCodeStateActive      = "active"
	AccessCodeStateSetup       = "setup"
	AccessCodeStatePending     = "pending"
	AccessCodeStateInactive    = "inactive"
	AccessCodeStateRemoved     = "removed"
	AccessCodeStateRemoving    = "removing"
	AccessCodeStateOverlocking = "overlocking"
	AccessCodeStateOverlocked  = "overlocked"
	AccessCodeStateRemove      = "remove"

	LOCK_STATE_OVERLOCK = "overlock"
	LOCK_STATE_GATELOCK = "gatelock"
	LOCK_STATE_PRELET   = "prelet"

	AccessCodeMsgDuplicate = "duplicate_code"
)

func UserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	return nil, nil
}

func (b *Bank) GetBUserByID(BUserID int) (*BUser, error) {
	return &BUser{}, nil
}

func (b *Bank) V2UnitGetById(unitID int, siteID int) (*Unit, error) {
	return &Unit{}, nil
}

func (b *Bank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return []GateAccessCode{}, nil
}

func (gacs GateAccessCodes) Validate(bank BankInterface) error {
	if len(gacs) == 0 {
		return fmt.Errorf("no access codes provided")
	}

	for i := range gacs {
		gacs[i].IsValid = true
	}
	return nil
}

func (b *Bank) UpdateAccessCodes(codes []string, siteID int) error {
	return nil
}

func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientInterface {
	return &CommandCenterClient{}
}

func (cc *CommandCenterClient) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	return nil
}

func (cc *CommandCenterClient) SetAccessCodes(units []int, options map[string]struct{}) error {
	return nil
}

type contextKey string

const (
	claimsKey contextKey = "claims"
	bankKey   contextKey = "db"
)

func NewBankContext(ctx context.Context, bank *Bank) context.Context {
	return context.WithValue(ctx, bankKey, bank)
}

// FromContext returns the BUser value stored in ctx, if any.
func BankFromContext(ctx context.Context) (*Bank, bool) {
	bank, ok := ctx.Value(bankKey).(*Bank)
	return bank, ok
}

type Claims struct {
	CompanyUUID     string
	CurrentSite     int
	CurrentSiteUUID string
	UserID          int
}

func uniqueIntSlice(input []int) ([]int, error) {
	keys := make(map[int]bool)
	result := []int{}
	for _, val := range input {
		if _, exists := keys[val]; !exists {
			keys[val] = true
			result = append(result, val)
		}
	}
	return result, nil
}

func convertToStringSlice(gacs GateAccessCodes) []string {
	result := make([]string, len(gacs))
	for i, gac := range gacs {
		result[i] = gac.AccessCode
	}
	return result
}

func getClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

type AccessCodeEditRequest struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

type AccessCodeEditService struct {
	bank BankInterface
}

func NewAccessCodeEditService(bank BankInterface) *AccessCodeEditService {
	return &AccessCodeEditService{bank: bank}
}

// validateUserAccess checks if user exists and has access to the site
func (s *AccessCodeEditService) validateUserAccess(userID int, claims *Claims) (*BUser, error) {
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

	// Check if user is associated with the site
	siteFound := false
	for _, siteID := range user.Sites {
		if siteID == strconv.Itoa(claims.CurrentSite) {
			siteFound = true
			break
		}
	}

	if !siteFound {
		return nil, fmt.Errorf("invalid user, missing association with target site")
	}

	return user, nil
}

// validateUnitAccess checks if unit exists and allows access code changes
func (s *AccessCodeEditService) validateUnitAccess(unitID int, claims *Claims) (*Unit, error) {
	if unitID == 0 {
		return nil, fmt.Errorf("invalid unit ID")
	}

	unit, err := s.bank.V2UnitGetById(unitID, claims.CurrentSite)
	if err != nil {
		return nil, fmt.Errorf("unit not found")
	}

	if unit.SiteID != claims.CurrentSite {
		return nil, fmt.Errorf("invalid unit, missing association with target site")
	}

	// Check if unit is in a state that allows access code changes
	if unit.RentalState == LOCK_STATE_OVERLOCK || unit.RentalState == LOCK_STATE_GATELOCK || unit.RentalState == LOCK_STATE_PRELET {
		return nil, fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return unit, nil
}

// processAccessCodeForUnit handles access code processing for a single unit
func (s *AccessCodeEditService) processAccessCodeForUnit(unitID int, userID int, accessCode string, claims *Claims, ctx context.Context) error {
	// Validate unit access
	_, err := s.validateUnitAccess(unitID, claims)
	if err != nil {
		return err
	}

	// Create and validate access code
	gacs := GateAccessCodes{
		GateAccessCode{
			AccessCode: accessCode,
			UnitID:     unitID,
			UserID:     userID,
			SiteID:     claims.CurrentSite,
			State:      AccessCodeStateSetup,
		},
	}

	err = gacs.Validate(s.bank)
	if err != nil {
		return fmt.Errorf("internal server error during validation")
	}

	if !gacs[0].IsValid {
		for _, validationMsg := range gacs[0].ValidationMessages {
			if validationMsg == AccessCodeMsgDuplicate {
				return fmt.Errorf("duplicate access code")
			}
		}
		return fmt.Errorf("invalid access code")
	}

	// Handle existing access codes
	err = s.handleExistingAccessCodes(unitID, accessCode, claims, ctx)
	if err != nil {
		return err
	}

	// Update access codes in database
	err = s.bank.UpdateAccessCodes(convertToStringSlice(gacs), claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	// Set new access codes via command center
	cc := s.bank.NewCommandCenterClient(claims.CurrentSite, ctx)
	err = cc.SetAccessCodes([]int{unitID}, make(map[string]struct{}))
	if err != nil {
		return fmt.Errorf("failed to set access codes: %v", err)
	}

	return nil
}

// handleExistingAccessCodes manages revocation of existing access codes
func (s *AccessCodeEditService) handleExistingAccessCodes(unitID int, newAccessCode string, claims *Claims, ctx context.Context) error {
	units := []int{unitID}
	revokeGacs, err := s.bank.GetCodesForUnits(units, claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	var revokeUnits []int
	var updateRemoveGacs GateAccessCodes

	// Check existing access codes
	for _, revokeGac := range revokeGacs {
		// If code is already active/setup/pending and matches new code, no change needed
		if (revokeGac.State == AccessCodeStateActive || revokeGac.State == AccessCodeStateSetup || revokeGac.State == AccessCodeStatePending) && revokeGac.AccessCode == newAccessCode {
			return fmt.Errorf("no change needed") // Special error to indicate success with no action
		}

		// Mark codes for removal if not already inactive/removed/removing/overlocking/overlocked
		if revokeGac.State != AccessCodeStateInactive && revokeGac.State != AccessCodeStateRemoved && revokeGac.State != AccessCodeStateRemoving && revokeGac.State != AccessCodeStateOverlocking && revokeGac.State != AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
			revokeGac.State = AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, revokeGac)
		}
	}

	// Update old access codes to "remove" status
	if len(updateRemoveGacs) > 0 {
		err = s.bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), claims.CurrentSite)
		if err != nil {
			return fmt.Errorf("internal server error updating access codes")
		}
	}

	// Revoke old access codes via command center
	if len(revokeUnits) > 0 {
		revokeUnits, _ = uniqueIntSlice(revokeUnits)
		cc := s.bank.NewCommandCenterClient(claims.CurrentSite, ctx)
		err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}))
		if err != nil {
			return fmt.Errorf("failed to revoke previous access codes: %v", err)
		}
	}

	return nil
}

// convertUUIDsToIDs converts UUID strings to integer IDs
func convertUUIDsToIDs(input *AccessCodeEditRequest) error {
	var err error
	if input.UserUUID != "" {
		input.UserID, err = strconv.Atoi(input.UserUUID)
		if err != nil {
			return fmt.Errorf("invalid UUID: %s", input.UserUUID)
		}
	} else {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	// Convert unit UUIDs to IDs
	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return fmt.Errorf("invalid UUID: %s", unitUUID)
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	return nil
}

func AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Use the extracted struct
	var input AccessCodeEditRequest

	bank, ok := BankFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve bank from context", http.StatusInternalServerError)
		return
	}

	// decode input struct
	// TODO: Add actual JSON decoding here

	// Convert UUIDs to IDs
	err = convertUUIDsToIDs(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create service
	service := NewAccessCodeEditService(bank)

	// Validate user access
	user, err := service.validateUserAccess(input.UserID, claims)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "User not found", http.StatusNotFound)
		case "user not found in company":
			http.Error(w, "User not found in company", http.StatusForbidden)
		case "invalid user, missing association with target site":
			http.Error(w, "invalid user, missing association with target site", http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Process each unit
	for _, unitID := range input.UnitIDs {
		err = service.processAccessCodeForUnit(unitID, input.UserID, input.AccessCode, claims, r.Context())
		if err != nil {
			if err.Error() == "no change needed" {
				w.WriteHeader(http.StatusOK)
				return
			}
			// Map specific errors to HTTP status codes
			switch {
			case err.Error() == "invalid unit ID" || err.Error() == "unit not found":
				http.Error(w, err.Error(), http.StatusNotFound)
			case err.Error() == "invalid unit, missing association with target site":
				http.Error(w, err.Error(), http.StatusForbidden)
			case err.Error() == "duplicate access code":
				http.Error(w, err.Error(), http.StatusConflict)
			case err.Error() == "invalid access code":
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				if err.Error() == "access code changes not allowed - unit in overlock" || err.Error() == "access code changes not allowed - unit in gatelock" || err.Error() == "access code changes not allowed - unit in prelet" {
					http.Error(w, err.Error(), http.StatusForbidden)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			return
		}
	}

	// Record access code update activity using activityEvents package
	// Record activity using the activity package
	_, err = UserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Record access code update activity using activityEvents package
	w.WriteHeader(http.StatusOK)
	return
}
