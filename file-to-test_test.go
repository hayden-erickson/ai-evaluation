package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockBank struct {
	user      *BUser
	userErr   error
	unit      *Unit
	unitErr   error
	codes     []GateAccessCode
	codesErr  error
	updateErr error
	cc        CommandCenterClientInterface
}

func (m *mockBank) GetBUserByID(BUserID int) (*BUser, error) {
	return m.user, m.userErr
}
func (m *mockBank) V2UnitGetById(unitID int, siteID int) (*Unit, error) {
	return m.unit, m.unitErr
}
func (m *mockBank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return m.codes, m.codesErr
}
func (m *mockBank) UpdateAccessCodes(codes []string, siteID int) error {
	return m.updateErr
}
func (m *mockBank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientInterface {
	return m.cc
}

type mockCC struct {
	revokeErr error
	setErr    error
}

func (m *mockCC) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	return m.revokeErr
}
func (m *mockCC) SetAccessCodes(units []int, options map[string]struct{}) error {
	return m.setErr
}

func TestAccessCodeEditHandlerWithDeps(t *testing.T) {
	claims := &Claims{
		CompanyUUID:     "company-1",
		CurrentSite:     123,
		CurrentSiteUUID: "site-uuid-1",
		UserID:          42,
	}
	ctx := context.WithValue(context.Background(), claimsKey, claims)

	// Helper to run handler and get status
	runHandler := func(bank BankInterface, input inputData) int {
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte{}))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		decodeAndValidateInput = func(r *http.Request) (inputData, string) { return input, "" }
		AccessCodeEditHandlerWithDeps(w, req, bank)
		return w.Result().StatusCode
	}

	// 1. Unauthorized (claims missing)
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	AccessCodeEditHandlerWithDeps(w, req, &mockBank{})
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("expected unauthorized, got %d", w.Result().StatusCode)
	}

	// 2. User not found
	bank := &mockBank{userErr: errors.New("no_ob_found")}
	input := inputData{UserID: 1}
	status := runHandler(bank, input)
	if status != http.StatusNotFound {
		t.Errorf("expected not found, got %d", status)
	}

	// 3. User not in company
	bank = &mockBank{user: &BUser{CompanyUUID: "other"}}
	status = runHandler(bank, input)
	if status != http.StatusForbidden {
		t.Errorf("expected forbidden, got %d", status)
	}

	// 4. User not associated with site
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{"999"}}}
	status = runHandler(bank, input)
	if status != http.StatusForbidden {
		t.Errorf("expected forbidden, got %d", status)
	}

	// 5. Unit not found
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unitErr: errors.New("not found")}
	input.UnitIDs = []int{1}
	status = runHandler(bank, input)
	if status != http.StatusNotFound {
		t.Errorf("expected unit not found, got %d", status)
	}

	// 6. Unit not associated with site
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unit: &Unit{SiteID: 999}}
	status = runHandler(bank, input)
	if status != http.StatusForbidden {
		t.Errorf("expected forbidden, got %d", status)
	}

	// 7. Unit in forbidden state
	forbiddenStates := []string{LOCK_STATE_OVERLOCK, LOCK_STATE_GATELOCK, LOCK_STATE_PRELET}
	for _, state := range forbiddenStates {
		bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unit: &Unit{SiteID: claims.CurrentSite, RentalState: state}}
		status = runHandler(bank, input)
		if status != http.StatusForbidden {
			t.Errorf("expected forbidden for state %s, got %d", state, status)
		}
	}

	// 8. Validation error
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unit: &Unit{SiteID: claims.CurrentSite}}
	// Simulate validation error by using a custom mockBank and a custom GateAccessCodes type
	// Instead, test with empty GateAccessCodes to trigger validation error
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unit: &Unit{SiteID: claims.CurrentSite}}
	// Patch the handler to use empty GateAccessCodes
	// This is not possible with the current design, so skip this branch or refactor the code for better testability if needed

	// 9. Duplicate access code
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unit: &Unit{SiteID: claims.CurrentSite}, codes: []GateAccessCode{{State: AccessCodeStateActive, AccessCode: "code"}}}
	input.AccessCode = "code"
	status = runHandler(bank, input)
	if status != http.StatusOK {
		t.Errorf("expected OK for duplicate code, got %d", status)
	}

	// 10. Success path
	bank = &mockBank{user: &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, unit: &Unit{SiteID: claims.CurrentSite}, codes: []GateAccessCode{{State: AccessCodeStateInactive, AccessCode: "other"}}, cc: &mockCC{}}
	input.AccessCode = "newcode"
	status = runHandler(bank, input)
	if status != http.StatusOK {
		t.Errorf("expected OK, got %d", status)
	}
}
