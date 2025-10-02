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

type CommandCenterClient struct{}
type Bank struct{}

// BankInterface defines the interface for database operations
type BankInterface interface {
	GetBUserByID(BUserID int) (*BUser, error)
	V2UnitGetById(unitID int, siteID int) (*Unit, error)
	GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodes(codes []string, siteID int) error
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface
}

// CommandCenterInterface defines the interface for command center operations
type CommandCenterInterface interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

// Ensure Bank implements BankInterface
var _ BankInterface = (*Bank)(nil)

// Ensure CommandCenterClient implements CommandCenterInterface
var _ CommandCenterInterface = (*CommandCenterClient)(nil)

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

func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface {
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

func NewBankContext(ctx context.Context, bank BankInterface) context.Context {
	return context.WithValue(ctx, bankKey, bank)
}

// FromContext returns the BUser value stored in ctx, if any.
func BankFromContext(ctx context.Context) (BankInterface, bool) {
	bank, ok := ctx.Value(bankKey).(BankInterface)
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

// inputData represents the request body structure
type inputData struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

// validateAndAuthorizeUser validates the user exists and has access to the site
func validateAndAuthorizeUser(bank BankInterface, userID int, companyUUID string, currentSite int) (*BUser, error) {
	user, err := bank.GetBUserByID(userID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("internal server error")
	}

	if user.CompanyUUID != companyUUID {
		return nil, fmt.Errorf("user not found in company")
	}

	// Check if user is associated with the site
	siteFound := false
	for _, siteID := range user.Sites {
		if siteID == strconv.Itoa(currentSite) {
			siteFound = true
			break
		}
	}

	if !siteFound {
		return nil, fmt.Errorf("invalid user, missing association with target site")
	}

	return user, nil
}

// validateUnit checks if unit exists and is in valid state for access code changes
func validateUnit(bank BankInterface, unitID int, siteID int) (*Unit, error) {
	unit, err := bank.V2UnitGetById(unitID, siteID)
	if err != nil {
		return nil, fmt.Errorf("unit not found")
	}

	if unit.SiteID != siteID {
		return nil, fmt.Errorf("invalid unit, missing association with target site")
	}

	if unit.RentalState == LOCK_STATE_OVERLOCK || unit.RentalState == LOCK_STATE_GATELOCK ||
		unit.RentalState == LOCK_STATE_PRELET {
		return nil, fmt.Errorf("access code changes not allowed - unit in %s", unit.RentalState)
	}

	return unit, nil
}

// processAccessCodeUpdate handles the logic for updating access codes for a single unit
func processAccessCodeUpdate(ctx context.Context, bank BankInterface, unitID int, userID int, siteID int, accessCode string) error {
	// Create a GateAccessCode object
	gacs := GateAccessCodes{
		{
			AccessCode: accessCode,
			UnitID:     unitID,
			UserID:     userID,
			SiteID:     siteID,
			State:      AccessCodeStateSetup,
		},
	}

	// Validate the new access code
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

	units := []int{unitID}

	// Get existing access codes for the unit to be revoked
	revokeGacs, err := bank.GetCodesForUnits(units, siteID)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	var revokeUnits []int
	var updateRemoveGacs GateAccessCodes

	// Loop through existing access codes to determine which need to be revoked
	for _, revokeGac := range revokeGacs {
		// If the existing code is the same as the new code and is active/setup/pending, no changes needed
		if (revokeGac.State == AccessCodeStateActive || revokeGac.State == AccessCodeStateSetup || revokeGac.State == AccessCodeStatePending) && revokeGac.AccessCode == accessCode {
			return nil // No update needed
		}

		// If the existing code is not already inactive/removed/removing/overlocking/overlocked, mark it for removal
		if revokeGac.State != AccessCodeStateInactive && revokeGac.State != AccessCodeStateRemoved && revokeGac.State != AccessCodeStateRemoving && revokeGac.State != AccessCodeStateOverlocking && revokeGac.State != AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
			revokeGac.State = AccessCodeStateRemove
			updateRemoveGacs = append(updateRemoveGacs, revokeGac)
		}
	}

	// Update access code table - set status to "remove" for old GACs
	if len(updateRemoveGacs) > 0 {
		err = bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), siteID)
		if err != nil {
			return fmt.Errorf("internal server error updating access codes")
		}
	}

	// Update access code table - insert new GACs, or update their status to "setup"
	err = bank.UpdateAccessCodes(convertToStringSlice(gacs), siteID)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	revokeUnits, _ = uniqueIntSlice(revokeUnits)

	// Use command center to remove the old access code for the unit
	if len(revokeUnits) > 0 {
		cc := bank.NewCommandCenterClient(siteID, ctx)
		err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0))
		if err != nil {
			return fmt.Errorf("failed to revoke previous access codes: %v", err)
		}
	}

	// Use command center to set the new access code for the unit
	cc := bank.NewCommandCenterClient(siteID, ctx)
	err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
	if err != nil {
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

	// decode input struct
	var input inputData
	// TODO: actual JSON decoding should be here

	// Convert UUIDs to IDs or vice versa
	if input.UserUUID != "" {
		input.UserID, err = strconv.Atoi(input.UserUUID)
		if err != nil {
			errStr := fmt.Sprintf("Invalid UUID: %s", input.UserUUID)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}
	} else {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	// loop through each unitUUID and convert to unitID because we only receive UUIDs from the front end
	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			errStr := fmt.Sprintf("Invalid UUID: %s", unitUUID)
			http.Error(w, errStr, http.StatusBadRequest)
			return
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	// Validate and authorize user
	user, err := validateAndAuthorizeUser(bank, input.UserID, claims.CompanyUUID, claims.CurrentSite)
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

	// loop through each unitID and process the access code change
	for _, unitID := range input.UnitIDs {
		if unitID != 0 {
			// Validate unit
			_, err := validateUnit(bank, unitID, claims.CurrentSite)
			if err != nil {
				switch {
				case err.Error() == "unit not found":
					http.Error(w, "Unit not found", http.StatusNotFound)
				case err.Error() == "invalid unit, missing association with target site":
					http.Error(w, "invalid unit, missing association with target site", http.StatusForbidden)
				case err.Error() == fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_OVERLOCK),
					err.Error() == fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_GATELOCK),
					err.Error() == fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_PRELET):
					http.Error(w, err.Error(), http.StatusForbidden)
				default:
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
				return
			}

			// Process the access code update
			err = processAccessCodeUpdate(r.Context(), bank, unitID, input.UserID, claims.CurrentSite, input.AccessCode)
			if err != nil {
				switch err.Error() {
				case "duplicate access code":
					http.Error(w, "Duplicate access code", http.StatusConflict)
				case "invalid access code":
					http.Error(w, "Invalid access code", http.StatusBadRequest)
				default:
					if err.Error() == "internal server error during validation" ||
						err.Error() == "internal server error updating access codes" {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					} else {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
				return
			}
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
