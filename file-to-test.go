package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface
}

type CommandCenterInterface interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

type ActivityRecorderInterface interface {
	UserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error)
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

// Activity recorder implementation
type ActivityRecorder struct{}

func (ar *ActivityRecorder) UserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	return nil, nil
}

func UserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	return nil, nil
}

// Bank implementation of BankInterface
func (b *Bank) GetBUserByID(BUserID int) (*BUser, error) {
	return &BUser{}, nil
}

func (b *Bank) V2UnitGetById(unitID int, siteID int) (*Unit, error) {
	return &Unit{}, nil
}

func (b *Bank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return []GateAccessCode{}, nil
}

func (b *Bank) UpdateAccessCodes(codes []string, siteID int) error {
	return nil
}

func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface {
	return &CommandCenterClient{}
}

// GateAccessCodes validation method
func (gacs GateAccessCodes) Validate(bank BankInterface) error {
	if len(gacs) == 0 {
		return fmt.Errorf("no access codes provided")
	}

	for i := range gacs {
		gacs[i].IsValid = true
	}
	return nil
}

// CommandCenterClient implementation of CommandCenterInterface
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

// AccessCodeEditRequest represents the input data for access code editing
type AccessCodeEditRequest struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

// AccessCodeEditService handles the business logic for access code editing
type AccessCodeEditService struct {
	bank             BankInterface
	activityRecorder ActivityRecorderInterface
}

func NewAccessCodeEditService(bank BankInterface, activityRecorder ActivityRecorderInterface) *AccessCodeEditService {
	return &AccessCodeEditService{
		bank:             bank,
		activityRecorder: activityRecorder,
	}
}

// ConvertUUIDs converts UUIDs to IDs in the request
func (s *AccessCodeEditService) ConvertUUIDs(input *AccessCodeEditRequest) error {
	// Convert UUIDs to IDs or vice versa
	if input.UserUUID != "" {
		userID, err := strconv.Atoi(input.UserUUID)
		if err != nil {
			return fmt.Errorf("invalid UUID: %s", input.UserUUID)
		}
		input.UserID = userID
	} else {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	// loop through each unitUUID and convert to unitID because we only receive UUIDs from the front end
	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return fmt.Errorf("invalid UUID: %s", unitUUID)
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	return nil
}

// ValidateUserAccess validates that the user exists and has proper access
func (s *AccessCodeEditService) ValidateUserAccess(userID int, claims *Claims) (*BUser, error) {
	// Check if provided user exists (needs DB)
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

	// Checks to make sure the user is associated with the site
	// Input site is contained in list of user sites
	siteFound := false
	for _, siteID := range user.Sites {
		if siteID == strconv.Itoa(claims.CurrentSite) {
			siteFound = true
			break
		}
	}

	// If the user is not associated with the site, return an error
	if !siteFound {
		return nil, fmt.Errorf("invalid user, missing association with target site")
	}

	return user, nil
}

// ValidateUnitAccess validates that the unit exists and allows access code changes
func (s *AccessCodeEditService) ValidateUnitAccess(unitID int, claims *Claims) (*Unit, error) {
	// Get unit from DB
	unit, err := s.bank.V2UnitGetById(unitID, claims.CurrentSite)
	if err != nil {
		return nil, fmt.Errorf("unit not found")
	}

	// If unit not found, or not associated with the site, return error
	if !(unit.SiteID == claims.CurrentSite) {
		return nil, fmt.Errorf("invalid unit, missing association with target site")
	}

	// Overlock, Gatelock, Prelet units cannot have access codes changed
	if unit.RentalState == LOCK_STATE_OVERLOCK || unit.RentalState == LOCK_STATE_GATELOCK ||
		unit.RentalState == LOCK_STATE_PRELET {
		return nil, fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return unit, nil
}

// ProcessAccessCodeForUnit processes access code changes for a single unit
func (s *AccessCodeEditService) ProcessAccessCodeForUnit(ctx context.Context, unitID int, userID int, accessCode string, claims *Claims) error {
	if unitID == 0 {
		return nil // Skip unitID 0
	}

	// Validate unit access
	_, err := s.ValidateUnitAccess(unitID, claims)
	if err != nil {
		return err
	}

	// Create a GateAccessCode object from unit
	var gacs GateAccessCodes
	gacs = append(gacs, GateAccessCode{
		AccessCode: accessCode,
		UnitID:     unitID,
		UserID:     userID,
		SiteID:     claims.CurrentSite,
		State:      AccessCodeStateSetup,
	})

	// Checks validity of GateAccessCode object
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

	var units []int
	units = append(units, unitID)

	// Get existing access codes for the (unit, site) to be revoked
	revokeGacs, err := s.bank.GetCodesForUnits(units, claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	var revokeUnits []int
	var updateRemoveGacs GateAccessCodes

	// Loop through existing access codes to determine which need to be revoked
	for _, revokeGac := range revokeGacs {
		// if the existing code is the same as the new code and is active/setup/pending, no changes needed
		if (revokeGac.State == AccessCodeStateActive || revokeGac.State == AccessCodeStateSetup || revokeGac.State == AccessCodeStatePending) && revokeGac.AccessCode == accessCode {
			return nil // No changes needed - early return for this condition
		}

		// if the existing code is not already inactive/removed/removing/overlocking/overlocked, mark it for removal
		if revokeGac.State != AccessCodeStateInactive && revokeGac.State != AccessCodeStateRemoved && revokeGac.State != AccessCodeStateRemoving && revokeGac.State != AccessCodeStateOverlocking && revokeGac.State != AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
			revokeGac.State = AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, revokeGac)
		}
	}

	// Updates access code table - sets status to "remove" for old GACs
	err = s.bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	// Updates access code table - inserts new GACs, or updates their status to "setup"
	err = s.bank.UpdateAccessCodes(convertToStringSlice(gacs), claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	revokeUnits, _ = uniqueIntSlice(revokeUnits)
	// Uses command center to remove the old access code for the unit
	cc := s.bank.NewCommandCenterClient(claims.CurrentSite, ctx)
	err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0))
	if err != nil {
		return fmt.Errorf("failed to revoke previous access codes: %v", err)
	}

	// Uses command center to set the new access code for the unit
	err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
	if err != nil {
		return fmt.Errorf("failed to set access codes: %v", err)
	}

	return nil
}

// ProcessAccessCodeEdit handles the complete access code editing process
func (s *AccessCodeEditService) ProcessAccessCodeEdit(ctx context.Context, input *AccessCodeEditRequest, claims *Claims) error {
	// Convert UUIDs to IDs
	err := s.ConvertUUIDs(input)
	if err != nil {
		return err
	}

	// Validate user access
	user, err := s.ValidateUserAccess(input.UserID, claims)
	if err != nil {
		return err
	}

	// loop through each unitID and process the access code change
	for _, unitID := range input.UnitIDs {
		err := s.ProcessAccessCodeForUnit(ctx, unitID, input.UserID, input.AccessCode, claims)
		if err != nil {
			return err
		}
	}

	// Record access code update activity using activityEvents package
	// Record activity using the activity package
	_, err = s.activityRecorder.UserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		return fmt.Errorf("recording activity failed: %v", err)
	}

	return nil
}

func getClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
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

	// Create service with dependencies
	activityRecorder := &ActivityRecorder{}
	service := NewAccessCodeEditService(bank, activityRecorder)

	// decode input struct
	var input AccessCodeEditRequest
	// Note: In a real implementation, you'd decode JSON from request body here

	// Process the access code edit
	err = service.ProcessAccessCodeEdit(r.Context(), &input, claims)
	if err != nil {
		// Map service errors to appropriate HTTP status codes
		switch {
		case err.Error() == "user not found":
			http.Error(w, "User not found", http.StatusNotFound)
		case err.Error() == "user not found in company":
			http.Error(w, "User not found in company", http.StatusForbidden)
		case err.Error() == "invalid user, missing association with target site":
			http.Error(w, "invalid user, missing association with target site", http.StatusForbidden)
		case err.Error() == "unit not found":
			http.Error(w, "Unit not found", http.StatusNotFound)
		case err.Error() == "invalid unit, missing association with target site":
			http.Error(w, "invalid unit, missing association with target site", http.StatusForbidden)
		case err.Error() == "duplicate access code":
			http.Error(w, "Duplicate access code", http.StatusConflict)
		case err.Error() == "invalid access code":
			http.Error(w, "Invalid access code", http.StatusBadRequest)
		case strings.Contains(err.Error(), "invalid UUID"):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case strings.Contains(err.Error(), "access code changes not allowed"):
			http.Error(w, err.Error(), http.StatusForbidden)
		case strings.Contains(err.Error(), "failed to revoke") || strings.Contains(err.Error(), "failed to set"):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case strings.Contains(err.Error(), "recording activity failed"):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
