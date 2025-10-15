package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hayden-erickson/ai-evaluation/constants"
	appContext "github.com/hayden-erickson/ai-evaluation/context"
	"github.com/hayden-erickson/ai-evaluation/models"
	"github.com/hayden-erickson/ai-evaluation/repository"
	"github.com/hayden-erickson/ai-evaluation/services"
	"github.com/hayden-erickson/ai-evaluation/utils"
)

// AccessCodeHandler handles access code related HTTP requests
type AccessCodeHandler struct {
	accessCodeService *services.AccessCodeService
}

// NewAccessCodeHandler creates a new AccessCodeHandler
func NewAccessCodeHandler(accessCodeService *services.AccessCodeService) *AccessCodeHandler {
	return &AccessCodeHandler{
		accessCodeService: accessCodeService,
	}
}

// inputData represents the request payload for access code operations
type inputData struct {
	UserID     int      `json:"userId" val:"optional"`
	UserUUID   string   `json:"userUuid" val:"optional"`
	UnitIDs    []int    `json:"unitID" val:"optional"`
	UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
	AccessCode string   `json:"accessCode"`
}

// AccessCodeEditHandler handles HTTP requests for editing access codes
func (h *AccessCodeHandler) AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, err := appContext.GetClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get bank from context
	bank, ok := appContext.BankFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve bank from context", http.StatusInternalServerError)
		return
	}

	// Parse input data (in a real implementation, you'd decode JSON here)
	var input inputData

	// Convert UUIDs to IDs
	err = h.processInputData(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user
	user, err := h.accessCodeService.ValidateUser(input.UserID, claims)
	if err != nil {
		switch err.Error() {
		case "user not found":
			http.Error(w, "User not found", http.StatusNotFound)
		case "user not found in company":
			http.Error(w, "User not found in company", http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Validate user site access
	err = h.accessCodeService.ValidateUserSiteAccess(user, claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Process each unit
	for _, unitID := range input.UnitIDs {
		if unitID != 0 {
			err = h.processUnitAccessCode(unitID, input, claims, bank, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Record activity
	_, err = services.UserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// processInputData converts UUIDs to IDs in the input data
func (h *AccessCodeHandler) processInputData(input *inputData) error {
	// Convert UUIDs to IDs
	if input.UserUUID != "" {
		userID, err := strconv.Atoi(input.UserUUID)
		if err != nil {
			return fmt.Errorf("Invalid UUID: %s", input.UserUUID)
		}
		input.UserID = userID
	} else {
		input.UserUUID = strconv.Itoa(input.UserID)
	}

	// Convert unit UUIDs to unit IDs
	for _, unitUUID := range input.UnitUUIDs {
		unitID, err := strconv.Atoi(unitUUID)
		if err != nil {
			return fmt.Errorf("Invalid UUID: %s", unitUUID)
		}
		input.UnitIDs = append(input.UnitIDs, unitID)
	}

	return nil
}

// processUnitAccessCode handles access code processing for a single unit
func (h *AccessCodeHandler) processUnitAccessCode(unitID int, input inputData, claims *models.Claims, bank *repository.Bank, r *http.Request) error {
	// Validate unit
	_, err := h.accessCodeService.ValidateUnit(unitID, claims)
	if err != nil {
		return err
	}

	// Process access code update
	err = h.accessCodeService.ProcessAccessCodeUpdate(unitID, input.UserID, input.AccessCode, claims)
	if err != nil {
		return err
	}

	// Handle command center operations
	err = h.handleCommandCenterOperations(unitID, claims, bank, r)
	if err != nil {
		return err
	}

	return nil
}

// handleCommandCenterOperations manages command center interactions
func (h *AccessCodeHandler) handleCommandCenterOperations(unitID int, claims *models.Claims, bank *repository.Bank, r *http.Request) error {
	units := []int{unitID}

	// Get existing codes for revocation
	revokeGacs, err := bank.GetCodesForUnits(units, claims.CurrentSite)
	if err != nil {
		return fmt.Errorf("internal server error updating access codes")
	}

	var revokeUnits []int
	for _, revokeGac := range revokeGacs {
		if revokeGac.State != constants.AccessCodeStateInactive &&
			revokeGac.State != constants.AccessCodeStateRemoved &&
			revokeGac.State != constants.AccessCodeStateRemoving &&
			revokeGac.State != constants.AccessCodeStateOverlocking &&
			revokeGac.State != constants.AccessCodeStateOverlocked {
			revokeUnits = append(revokeUnits, revokeGac.UnitID)
		}
	}

	// Create command center client
	cc := bank.NewCommandCenterClient(claims.CurrentSite, r.Context())

	// Revoke old access codes
	if len(revokeUnits) > 0 {
		revokeUnits, _ = utils.UniqueIntSlice(revokeUnits)
		err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0))
		if err != nil {
			return fmt.Errorf("failed to revoke previous access codes: %v", err)
		}
	}

	// Set new access codes
	err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
	if err != nil {
		return fmt.Errorf("failed to set access codes: %v", err)
	}

	return nil
}
