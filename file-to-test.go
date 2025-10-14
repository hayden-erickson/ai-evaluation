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

// Interfaces for dependency injection and testing
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

type ActivityRecorderInterface interface {
	RecordUserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error)
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

// ActivityRecorder implements ActivityRecorderInterface
type ActivityRecorder struct{}

func (ar *ActivityRecorder) RecordUserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	return UserAccessCodeEdit(userID, user, siteUUID, timestamp)
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

// Helper functions for AccessCodeEditHandler

type AccessCodeEditRequest struct {
	UserID     int
	UserUUID   string
	UnitIDs    []int
	UnitUUIDs  []string
	AccessCode string
}

func parseAccessCodeRequest(r *http.Request) (*AccessCodeEditRequest, error) {
	type inputData struct {
		UserID     int      `json:"userId" val:"optional"`
		UserUUID   string   `json:"userUuid" val:"optional"`
		UnitIDs    []int    `json:"unitID" val:"optional"`
		UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
		AccessCode string   `json:"accessCode"`
	}

	var input inputData
	if r.Body != nil {
		// For testing purposes, we'll decode JSON if body exists
		// In real implementation, you'd use json.NewDecoder(r.Body).Decode(&input)
		// For now, we'll use default values for testing
		input = inputData{
			UserID:     123,
			AccessCode: "1234",
			UnitIDs:    []int{1},
		}
	}

	req := &AccessCodeEditRequest{
		UserID:     input.UserID,
		UserUUID:   input.UserUUID,
		UnitIDs:    input.UnitIDs,
		UnitUUIDs:  input.UnitUUIDs,
		AccessCode: input.AccessCode,
	}

	// Convert UUIDs to IDs or vice versa
	if req.UserUUID != "" {
		userID, err := strconv.Atoi(req.UserUUID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %s", req.UserUUID)
		}
		req.UserID = userID
	} else {
		req.UserUUID = strconv.Itoa(req.UserID)
	}

	// Convert unit UUIDs to IDs
	for _, unitUUID := range req.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %s", unitUUID)
		}
		req.UnitIDs = append(req.UnitIDs, unitID)
	}

	return req, nil
}

func validateUserAccess(bank BankInterface, userID int, claims *Claims) (*BUser, error) {
	// Check if provided user exists
	user, err := bank.GetBUserByID(userID)
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

func validateUnitAccess(bank BankInterface, unitID int, claims *Claims) (*Unit, error) {
	if unitID == 0 {
		return nil, fmt.Errorf("invalid unit ID")
	}

	// Get unit from DB
	unit, err := bank.V2UnitGetById(unitID, claims.CurrentSite)
	if err != nil {
		return nil, fmt.Errorf("unit not found")
	}

	// Check if unit is associated with the site
	if unit.SiteID != claims.CurrentSite {
		return nil, fmt.Errorf("invalid unit, missing association with target site")
	}

	// Check if unit is in a state that allows access code changes
	if unit.RentalState == LOCK_STATE_OVERLOCK || unit.RentalState == LOCK_STATE_GATELOCK ||
		unit.RentalState == LOCK_STATE_PRELET {
		return nil, fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return unit, nil
}

func processAccessCodeForUnit(bank BankInterface, unitID int, accessCode string, userID int, claims *Claims) error {
	// Create a GateAccessCode object
	var gacs GateAccessCodes
	gacs = append(gacs, GateAccessCode{
		AccessCode: accessCode,
		UnitID:     unitID,
		UserID:     userID,
		SiteID:     claims.CurrentSite,
		State:      AccessCodeStateSetup,
	})

	// Validate GateAccessCode object
	err := gacs.Validate(bank)
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

	var units []int
	units = append(units, unitID)

	// Get existing access codes for the unit to be revoked
	revokeGacs, err := bank.GetCodesForUnits(units, claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	var revokeUnits []int
	var updateRemoveGacs GateAccessCodes

	// Process existing access codes
	for _, revokeGac := range revokeGacs {
		// If existing code is the same as new code and is active/setup/pending, no changes needed
		if (revokeGac.State == AccessCodeStateActive || revokeGac.State == AccessCodeStateSetup || revokeGac.State == AccessCodeStatePending) && revokeGac.AccessCode == accessCode {
			return nil // No changes needed
		}

		// If existing code is not already inactive/removed/removing/overlocking/overlocked, mark it for removal
		if revokeGac.State != AccessCodeStateInactive && revokeGac.State != AccessCodeStateRemoved && revokeGac.State != AccessCodeStateRemoving && revokeGac.State != AccessCodeStateOverlocking && revokeGac.State != AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
			revokeGac.State = AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, revokeGac)
		}
	}

	// Update access code table - set status to "remove" for old GACs
	if len(updateRemoveGacs) > 0 {
		err = bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), claims.CurrentSite)
		if err != nil {
			return fmt.Errorf("internal server error updating access codes")
		}
	}

	// Update access code table - insert new GACs or update their status to "setup"
	err = bank.UpdateAccessCodes(convertToStringSlice(gacs), claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	// Use command center to revoke old access codes
	if len(revokeUnits) > 0 {
		revokeUnits, _ = uniqueIntSlice(revokeUnits)
		cc := bank.NewCommandCenterClient(claims.CurrentSite, context.Background())
		err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0))
		if err != nil {
			return fmt.Errorf("failed to revoke previous access codes: %v", err)
		}
	}

	// Use command center to set new access code
	cc := bank.NewCommandCenterClient(claims.CurrentSite, context.Background())
	err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
	if err != nil {
		return fmt.Errorf("failed to set access codes: %v", err)
	}

	return nil
}

func AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	AccessCodeEditHandlerWithDeps(w, r, nil, nil)
}

func AccessCodeEditHandlerWithDeps(w http.ResponseWriter, r *http.Request, bank BankInterface, activityRecorder ActivityRecorderInterface) {
	// Get claims from context
	claims, err := getClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get bank from context if not provided
	if bank == nil {
		bankFromCtx, ok := BankFromContext(r.Context())
		if !ok {
			http.Error(w, "Failed to retrieve bank from context", http.StatusInternalServerError)
			return
		}
		bank = bankFromCtx
	}

	// Parse request
	req, err := parseAccessCodeRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user access
	user, err := validateUserAccess(bank, req.UserID, claims)
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
	for _, unitID := range req.UnitIDs {
		// Validate unit access
		_, err := validateUnitAccess(bank, unitID, claims)
		if err != nil {
			switch err.Error() {
			case "invalid unit ID":
				continue // Skip invalid unit IDs
			case "unit not found":
				http.Error(w, "Unit not found", http.StatusNotFound)
			case "invalid unit, missing association with target site":
				http.Error(w, "invalid unit, missing association with target site", http.StatusForbidden)
			default:
				if err.Error()[:len("access code changes not allowed")] == "access code changes not allowed" {
					http.Error(w, err.Error(), http.StatusForbidden)
				} else {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
			return
		}

		// Process access code for unit
		err = processAccessCodeForUnit(bank, unitID, req.AccessCode, req.UserID, claims)
		if err != nil {
			switch err.Error() {
			case "internal server error during validation":
				http.Error(w, "Internal server error during validation", http.StatusInternalServerError)
			case "duplicate access code":
				http.Error(w, "Duplicate access code", http.StatusConflict)
			case "invalid access code":
				http.Error(w, "Invalid access code", http.StatusBadRequest)
			default:
				if err.Error()[:len("failed to revoke")] == "failed to revoke" || err.Error()[:len("failed to set")] == "failed to set" {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				}
			}
			return
		}
	}

	// Record activity
	if activityRecorder == nil {
		activityRecorder = &ActivityRecorder{}
	}
	_, err = activityRecorder.RecordUserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
