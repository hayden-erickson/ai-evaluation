package main

import (
	"context"
	"encoding/json"
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

// BankStore abstracts DB operations needed by access code editing logic for easier testing.
type BankStore interface {
	GetBUserByID(BUserID int) (*BUser, error)
	V2UnitGetById(unitID int, siteID int) (*Unit, error)
	GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodes(codes []string, siteID int) error
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientAPI
}

// CommandCenterClientAPI abstracts external command center actions.
type CommandCenterClientAPI interface {
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

func (gacs GateAccessCodes) Validate(bank *Bank) error {
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

func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientAPI {
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

// accessCodeEditInput holds decoded request body values.
type accessCodeEditInput struct {
	UserID     int      `json:"userId"`
	UserUUID   string   `json:"userUuid"`
	UnitIDs    []int    `json:"unitIDs"`
	UnitUUIDs  []string `json:"unitUUIDs"`
	AccessCode string   `json:"accessCode"`
}

// userAccessCodeEdit is a package level variable wrapping UserAccessCodeEdit for test injection.
var userAccessCodeEdit = UserAccessCodeEdit

// validateAccessCodes wraps GateAccessCodes.Validate for test overrides.
var validateAccessCodes = func(gacs GateAccessCodes, bank *Bank) error { return gacs.Validate(bank) }

// AccessCodeEdit contains the business logic so it can be unit tested without HTTP concerns.
func AccessCodeEdit(ctx context.Context, bank BankStore, claims *Claims, input accessCodeEditInput) (int, string) {
	if claims == nil {
		return http.StatusUnauthorized, "Unauthorized"
	}
	if bank == nil {
		return http.StatusInternalServerError, "Failed to retrieve bank from context"
	}

	// Convert / normalize user IDs
	if input.UserUUID != "" {
		userID, err := strconv.Atoi(input.UserUUID)
		if err != nil {
			return http.StatusBadRequest, fmt.Sprintf("Invalid UUID: %s", input.UserUUID)
		}
		input.UserID = userID
	} else if input.UserID != 0 {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	// Convert unit UUIDs
	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return http.StatusBadRequest, fmt.Sprintf("Invalid UUID: %s", unitUUID)
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	user, err := bank.GetBUserByID(input.UserID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			return http.StatusNotFound, "User not found"
		}
		return http.StatusInternalServerError, "Internal server error"
	}

	if user.CompanyUUID != claims.CompanyUUID {
		return http.StatusForbidden, "User not found in company"
	}

	siteFound := false
	for _, siteID := range user.Sites {
		if siteID == strconv.Itoa(claims.CurrentSite) {
			siteFound = true
			break
		}
	}
	if !siteFound {
		return http.StatusForbidden, "invalid user, missing association with target site"
	}

	for _, unitID := range input.UnitIDs {
		if unitID == 0 {
			continue
		}
		unit, err := bank.V2UnitGetById(unitID, claims.CurrentSite)
		if err != nil {
			return http.StatusNotFound, "Unit not found"
		}
		if unit.SiteID != claims.CurrentSite {
			return http.StatusForbidden, "invalid unit, missing association with target site"
		}
		if unit.RentalState == LOCK_STATE_OVERLOCK || unit.RentalState == LOCK_STATE_GATELOCK || unit.RentalState == LOCK_STATE_PRELET {
			return http.StatusForbidden, fmt.Sprintf("access code changes not allowed - unit in %s", unit.RentalState)
		}

		gacs := GateAccessCodes{GateAccessCode{AccessCode: input.AccessCode, UnitID: unitID, UserID: input.UserID, SiteID: claims.CurrentSite, State: AccessCodeStateSetup}}
		// Attempt to retrieve concrete *Bank for legacy Validate signature; fall back to empty bank for mocks
		realBank, _ := bank.(*Bank)
		if realBank == nil {
			realBank = &Bank{}
		}
		if err := validateAccessCodes(gacs, realBank); err != nil { // original Validate signature expects *Bank
			return http.StatusInternalServerError, "Internal server error during validation"
		}
		if !gacs[0].IsValid {
			for _, msg := range gacs[0].ValidationMessages {
				if msg == AccessCodeMsgDuplicate {
					return http.StatusConflict, "Duplicate access code"
				}
			}
			return http.StatusBadRequest, "Invalid access code"
		}

		units := []int{unitID}
		revokeGacs, err := bank.GetCodesForUnits(units, claims.CurrentSite)
		if err != nil {
			return http.StatusInternalServerError, "Internal server error updating access codes"
		}
		var revokeUnits []int
		var updateRemoveGacs GateAccessCodes
		for _, revokeGac := range revokeGacs {
			if (revokeGac.State == AccessCodeStateActive || revokeGac.State == AccessCodeStateSetup || revokeGac.State == AccessCodeStatePending) && revokeGac.AccessCode == input.AccessCode {
				// early success (no change)
				if _, err := userAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now()); err != nil {
					return http.StatusInternalServerError, fmt.Sprintf("Recording activity failed: %v", err)
				}
				return http.StatusOK, ""
			}
			if revokeGac.State != AccessCodeStateInactive && revokeGac.State != AccessCodeStateRemoved && revokeGac.State != AccessCodeStateRemoving && revokeGac.State != AccessCodeStateOverlocking && revokeGac.State != AccessCodeStateOverlocked {
				revokeUnits = append(revokeUnits, revokeGac.UnitID)
				revokeGac.State = AccessCodeStateRemove
				updateRemoveGacs = append(updateRemoveGacs, revokeGac)
			}
		}

		if err := bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), claims.CurrentSite); err != nil {
			return http.StatusInternalServerError, "Internal server error updating access codes"
		}
		if err := bank.UpdateAccessCodes(convertToStringSlice(gacs), claims.CurrentSite); err != nil {
			return http.StatusInternalServerError, "Internal server error updating access codes"
		}
		revokeUnits, _ = uniqueIntSlice(revokeUnits)
		cc := bank.NewCommandCenterClient(claims.CurrentSite, ctx)
		if err := cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0)); err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("failed to revoke previous access codes: %v", err)
		}
		if err := cc.SetAccessCodes(units, make(map[string]struct{}, 0)); err != nil {
			return http.StatusInternalServerError, fmt.Sprintf("failed to set access codes: %v", err)
		}
	}

	if _, err := userAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now()); err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Recording activity failed: %v", err)
	}
	return http.StatusOK, ""
}

// AccessCodeEditHandler now focuses only on HTTP concerns and delegates business logic.
func AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := getClaimsFromContext(r.Context())
	bank, _ := BankFromContext(r.Context())

	var input accessCodeEditInput
	_ = json.NewDecoder(r.Body).Decode(&input) // ignore decode error for brevity; could be validated

	status, msg := AccessCodeEdit(r.Context(), bank, claims, input)
	if msg != "" && status != http.StatusOK {
		http.Error(w, msg, status)
		return
	}
	w.WriteHeader(status)
}
