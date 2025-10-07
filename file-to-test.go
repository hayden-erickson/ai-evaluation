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

type CommandCenterClientInterface interface {
	RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodes(units []int, options map[string]struct{}) error
}

type BankInterface interface {
	GetBUserByID(BUserID int) (*BUser, error)
	V2UnitGetById(unitID int, siteID int) (*Unit, error)
	GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodes(codes []string, siteID int) error
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientInterface
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

func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) *CommandCenterClient {
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

type inputData struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

var decodeAndValidateInput = func(r *http.Request) (inputData, string) {
	var input inputData
	// For testing, we skip actual decoding and assume input is set by test
	// In production, decode from JSON body
	// _ = json.NewDecoder(r.Body).Decode(&input)

	// Convert UUIDs to IDs or vice versa
	if input.UserUUID != "" {
		id, err := strconv.Atoi(input.UserUUID)
		if err != nil {
			return input, fmt.Sprintf("Invalid UUID: %s", input.UserUUID)
		}
		input.UserID = id
	} else {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return input, fmt.Sprintf("Invalid UUID: %s", unitUUID)
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}
	return input, ""
}

func AccessCodeEditHandlerWithDeps(w http.ResponseWriter, r *http.Request, bank BankInterface) {
	claims, err := getClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	input, inputErr := decodeAndValidateInput(r)
	if inputErr != "" {
		http.Error(w, inputErr, http.StatusBadRequest)
		return
	}

	user, err := bank.GetBUserByID(input.UserID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user.CompanyUUID != claims.CompanyUUID {
		http.Error(w, "User not found in company", http.StatusForbidden)
		return
	}

	siteFound := false
	for _, siteID := range user.Sites {
		if siteID == strconv.Itoa(claims.CurrentSite) {
			siteFound = true
			break
		}
	}
	if !siteFound {
		http.Error(w, "invalid user, missing association with target site", http.StatusForbidden)
		return
	}

	for _, unitID := range input.UnitIDs {
		if unitID != 0 {
			unit, err := bank.V2UnitGetById(unitID, claims.CurrentSite)
			if err != nil {
				http.Error(w, "Unit not found", http.StatusNotFound)
				return
			}
			if !(unit.SiteID == claims.CurrentSite) {
				http.Error(w, "invalid unit, missing association with target site", http.StatusForbidden)
				return
			}
			if unit.RentalState == LOCK_STATE_OVERLOCK || unit.RentalState == LOCK_STATE_GATELOCK || unit.RentalState == LOCK_STATE_PRELET {
				http.Error(w, fmt.Sprintf("access code changes not allowed - unit in %s", unit.RentalState), http.StatusForbidden)
				return
			}
			var gacs GateAccessCodes
			gacs = append(gacs, GateAccessCode{
				AccessCode: input.AccessCode,
				UnitID:     unitID,
				UserID:     input.UserID,
				SiteID:     claims.CurrentSite,
				State:      AccessCodeStateSetup,
			})
			err = gacs.Validate(bank)
			if err != nil {
				http.Error(w, "Internal server error during validation", http.StatusInternalServerError)
				return
			}
			if !gacs[0].IsValid {
				for _, validationMsg := range gacs[0].ValidationMessages {
					if validationMsg == AccessCodeMsgDuplicate {
						http.Error(w, "Duplicate access code", http.StatusConflict)
						return
					}
				}
				http.Error(w, "Invalid access code", http.StatusBadRequest)
				return
			}
			var units []int
			units = append(units, unitID)
			revokeGacs, err := bank.GetCodesForUnits(units, claims.CurrentSite)
			if err != nil {
				http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				return
			}
			var revokeUnits []int
			var updateRemoveGacs GateAccessCodes
			for _, revokeGac := range revokeGacs {
				if (revokeGac.State == AccessCodeStateActive || revokeGac.State == AccessCodeStateSetup || revokeGac.State == AccessCodeStatePending) && revokeGac.AccessCode == input.AccessCode {
					w.WriteHeader(http.StatusOK)
					return
				}
				if revokeGac.State != AccessCodeStateInactive && revokeGac.State != AccessCodeStateRemoved && revokeGac.State != AccessCodeStateRemoving && revokeGac.State != AccessCodeStateOverlocking && revokeGac.State != AccessCodeStateOverlocked {
					revokeUnits = append(revokeUnits, revokeGac.UnitID)
					revokeGac.State = AccessCodeStateRemove
					updateRemoveGacs = append(updateRemoveGacs, revokeGac)
				}
			}
			err = bank.UpdateAccessCodes(convertToStringSlice(updateRemoveGacs), claims.CurrentSite)
			if err != nil {
				http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				return
			}
			err = bank.UpdateAccessCodes(convertToStringSlice(gacs), claims.CurrentSite)
			if err != nil {
				http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				return
			}
			revokeUnits, _ = uniqueIntSlice(revokeUnits)
			cc := bank.NewCommandCenterClient(claims.CurrentSite, r.Context())
			err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0))
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to revoke previous access codes: %v", err), http.StatusInternalServerError)
				return
			}
			err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to set access codes: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
	_, err = UserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
