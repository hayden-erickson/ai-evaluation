package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hayden-erickson/ai-evaluation/internal/clients"
	"github.com/hayden-erickson/ai-evaluation/internal/constants"
	appContext "github.com/hayden-erickson/ai-evaluation/internal/context"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/utils"
)

// AccessCodeEditHandler handles access code edit requests
func AccessCodeEditHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := appContext.GetClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type inputData struct {
		UserID     int      `json:"userId" val:"optional"`
		UserUUID   string   `json:"userUuid" val:"optional"`
		UnitIDs    []int    `json:"unitID" val:"optional"`
		UnitUUIDs  []string `json:"unitUUIDs" val:"optional"`
		AccessCode string   `json:"accessCode"`
	}

	bankInterface, ok := appContext.BankFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve bank from context", http.StatusInternalServerError)
		return
	}

	bank, ok := bankInterface.(clients.BankInterface)
	if !ok {
		http.Error(w, "Invalid bank interface", http.StatusInternalServerError)
		return
	}

	// decode input struct
	var input inputData

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

	// Check if provided user exists (needs DB)
	user, err := bank.GetBUserByID(input.UserID)
	if err != nil {
		if err.Error() == "no_ob_found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if user belongs to the same company as the requester
	if user.CompanyUUID != claims.CompanyUUID {
		http.Error(w, "User not found in company", http.StatusForbidden)
		return
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
		http.Error(w, "invalid user, missing association with target site", http.StatusForbidden)
		return
	}

	// loop through each unitID and process the access code change
	for _, unitID := range input.UnitIDs {

		if unitID != 0 {
			// Get unit from DB
			unit, err := bank.V2UnitGetById(unitID, claims.CurrentSite)
			if err != nil {
				http.Error(w, "Unit not found", http.StatusNotFound)
				return
			}

			// If unit not found, or not associated with the site, or in a state that doesn't allow access code changes, return error
			if !(unit.SiteID == claims.CurrentSite) {
				http.Error(w, "invalid unit, missing association with target site", http.StatusForbidden)
				return
			}

			// Overlock, Gatelock, Prelet units cannot have access codes changed
			if unit.RentalState == constants.LOCK_STATE_OVERLOCK || unit.RentalState == constants.LOCK_STATE_GATELOCK ||
				unit.RentalState == constants.LOCK_STATE_PRELET {
				http.Error(w, fmt.Sprintf("access code changes not allowed - unit in %s", unit.RentalState), http.StatusForbidden)
				return
			}

			// Create a GateAccessCode object from unit
			var gacs models.GateAccessCodes
			gacs = append(gacs, models.GateAccessCode{
				AccessCode: input.AccessCode,
				UnitID:     unitID,
				UserID:     input.UserID,
				SiteID:     claims.CurrentSite,
				State:      constants.AccessCodeStateSetup,
			})

			// Checks validity of GateAccessCode object
			err = gacs.Validate(bank)
			if err != nil {
				http.Error(w, "Internal server error during validation", http.StatusInternalServerError)
				return
			}
			if !gacs[0].IsValid {
				for _, validationMsg := range gacs[0].ValidationMessages {
					if validationMsg == constants.AccessCodeMsgDuplicate {
						http.Error(w, "Duplicate access code", http.StatusConflict)
						return
					}
				}
				http.Error(w, "Invalid access code", http.StatusBadRequest)
				return
			}

			var units []int
			units = append(units, unitID)

			// Get existing access codes for the (unit, site) to be revoked
			revokeGacs, err := bank.GetCodesForUnits(units, claims.CurrentSite)
			if err != nil {
				http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				return
			}
			var revokeUnits []int
			var updateRemoveGacs models.GateAccessCodes

			// Loop through existing access codes to determine which need to be revoked
			for _, revokeGac := range revokeGacs {
				// if the existing code is the same as the new code and is active/setup/pending, no changes needed
				if (revokeGac.State == constants.AccessCodeStateActive || revokeGac.State == constants.AccessCodeStateSetup || revokeGac.State == constants.AccessCodeStatePending) && revokeGac.AccessCode == input.AccessCode {
					w.WriteHeader(http.StatusOK)
					return
				}

				// if the existing code is not already inactive/removed/removing/overlocking/overlocked, mark it for removal
				if revokeGac.State != constants.AccessCodeStateInactive && revokeGac.State != constants.AccessCodeStateRemoved && revokeGac.State != constants.AccessCodeStateRemoving && revokeGac.State != constants.AccessCodeStateOverlocking && revokeGac.State != constants.AccessCodeStateOverlocked {
					revokeUnits = append(revokeUnits, revokeGac.UnitID)
					revokeGac.State = constants.AccessCodeStateRemove
					updateRemoveGacs = append(updateRemoveGacs, revokeGac)
				}
			}

			// Updates access code table - sets status to "remove" for old GACs
			err = bank.UpdateAccessCodes(utils.ConvertToStringSlice(updateRemoveGacs), claims.CurrentSite)
			if err != nil {
				http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				return
			}

			// Updates access code table - inserts new GACs, or updates their status to "setup"
			err = bank.UpdateAccessCodes(utils.ConvertToStringSlice(gacs), claims.CurrentSite)
			if err != nil {
				http.Error(w, "Internal server error updating access codes", http.StatusInternalServerError)
				return
			}

			revokeUnits, _ = utils.UniqueIntSlice(revokeUnits)
			// Uses command center to remove the old access code for the unit
			cc := bank.NewCommandCenterClient(claims.CurrentSite, r.Context())
			err = cc.RevokeAccessCodes(revokeUnits, make(map[string]struct{}, 0))
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to revoke previous access codes: %v", err), http.StatusInternalServerError)
				return
			}

			// Uses command center to set the new access code for the unit
			err = cc.SetAccessCodes(units, make(map[string]struct{}, 0))
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to set access codes: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	// Record access code update activity using activityEvents package
	// Record activity using the activity package
	_, err = clients.UserAccessCodeEdit(claims.UserID, user, claims.CurrentSiteUUID, time.Now())
	if err != nil {
		http.Error(w, fmt.Sprintf("Recording activity failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Record access code update activity using activityEvents package
	w.WriteHeader(http.StatusOK)
	return
}
