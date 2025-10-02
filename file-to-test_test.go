package main

// This test suite achieves full branch coverage of the refactored AccessCodeEdit logic.
// Strategy:
//  - The long HTTP handler was decomposed so we can call AccessCodeEdit directly with mocks.
//  - Each numbered section corresponds to a unique conditional / early return branch in order of appearance.
//  - We override side-effectful functions (validation + activity recording) via package level vars for isolation.
//  - Mocks capture error injection points for DB + command center interactions.
//  - Final case (19) validates the success path after all failure branches are exercised.

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

// mock implementations

type mockBank struct {
	user       *BUser
	userErr    error
	unit       *Unit
	unitErr    error
	codes      []GateAccessCode
	codesErr   error
	updateErrs []error // first call remove, second call add
	cc         *mockCC
}

func (m *mockBank) GetBUserByID(BUserID int) (*BUser, error)            { return m.user, m.userErr }
func (m *mockBank) V2UnitGetById(unitID int, siteID int) (*Unit, error) { return m.unit, m.unitErr }
func (m *mockBank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return m.codes, m.codesErr
}
func (m *mockBank) UpdateAccessCodes(codes []string, siteID int) error {
	if len(m.updateErrs) == 0 {
		return nil
	}
	e := m.updateErrs[0]
	m.updateErrs = m.updateErrs[1:]
	return e
}
func (m *mockBank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClientAPI {
	return m.cc
}

type mockCC struct {
	revokeErr error
	setErr    error
	revoked   []int
	setUnits  []int
}

func (c *mockCC) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	c.revoked = append(c.revoked, revokeUnits...)
	return c.revokeErr
}
func (c *mockCC) SetAccessCodes(units []int, options map[string]struct{}) error {
	c.setUnits = append(c.setUnits, units...)
	return c.setErr
}

// override helpers
func withUserAccessCodeEdit(f func()) {
	original := userAccessCodeEdit
	userAccessCodeEdit = func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) { return nil, nil }
	f()
	userAccessCodeEdit = original
}

func withValidateAccessCodes(fn func(g GateAccessCodes, b *Bank) error, f func()) {
	original := validateAccessCodes
	validateAccessCodes = fn
	f()
	validateAccessCodes = original
}

func baseClaims() *Claims {
	return &Claims{CompanyUUID: "co", CurrentSite: 10, CurrentSiteUUID: "site-uuid", UserID: 999}
}

// Test branches
func TestAccessCodeEdit_Branches(t *testing.T) {
	// 1. Unauthorized (nil claims)
	status, msg := AccessCodeEdit(context.Background(), &Bank{}, nil, accessCodeEditInput{})
	if status != http.StatusUnauthorized || msg != "Unauthorized" {
		t.Fatalf("expected unauthorized branch, got %d %s", status, msg)
	}

	// 2. Missing bank
	status, msg = AccessCodeEdit(context.Background(), nil, baseClaims(), accessCodeEditInput{})
	if status != http.StatusInternalServerError || msg != "Failed to retrieve bank from context" {
		t.Fatalf("expected bank missing branch")
	}

	// 3. Invalid user UUID parse
	status, msg = AccessCodeEdit(context.Background(), &Bank{}, baseClaims(), accessCodeEditInput{UserUUID: "abc"})
	if status != http.StatusBadRequest {
		t.Fatalf("expected bad request invalid uuid, got %d", status)
	}

	// 4. User not found (specific error)
	bankNF := &mockBank{userErr: errors.New("no_ob_found")}
	status, msg = AccessCodeEdit(context.Background(), bankNF, baseClaims(), accessCodeEditInput{UserUUID: "1"})
	if status != http.StatusNotFound || msg != "User not found" {
		t.Fatalf("expected user not found, got %d %s", status, msg)
	}

	// 5. Internal error getting user
	bankErr := &mockBank{userErr: errors.New("db down")}
	status, msg = AccessCodeEdit(context.Background(), bankErr, baseClaims(), accessCodeEditInput{UserUUID: "1"})
	if status != http.StatusInternalServerError {
		t.Fatalf("expected internal error, got %d", status)
	}

	// 6. Company mismatch
	bankCmp := &mockBank{user: &BUser{CompanyUUID: "other", Sites: []string{"10"}}}
	status, msg = AccessCodeEdit(context.Background(), bankCmp, baseClaims(), accessCodeEditInput{UserUUID: "1"})
	if status != http.StatusForbidden || msg != "User not found in company" {
		t.Fatalf("expected company mismatch forbidden")
	}

	// 7. Missing site association
	bankNoSite := &mockBank{user: &BUser{CompanyUUID: "co", Sites: []string{"11"}}}
	status, msg = AccessCodeEdit(context.Background(), bankNoSite, baseClaims(), accessCodeEditInput{UserUUID: "1"})
	if status != http.StatusForbidden || msg != "invalid user, missing association with target site" {
		t.Fatalf("expected missing site assoc forbidden")
	}

	// Prepare happy path user and unit base
	user := &BUser{CompanyUUID: "co", Sites: []string{"10"}}
	unit := &Unit{SiteID: 10, RentalState: "rented"}

	// 8. Unit not found
	bankUnitNF := &mockBank{user: user, unitErr: errors.New("no unit")}
	status, msg = AccessCodeEdit(context.Background(), bankUnitNF, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
	if status != http.StatusNotFound || msg != "Unit not found" {
		t.Fatalf("expected unit not found")
	}

	// 9. Unit site mismatch
	unit.SiteID = 11
	bankSiteMismatch := &mockBank{user: user, unit: unit}
	status, msg = AccessCodeEdit(context.Background(), bankSiteMismatch, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
	if status != http.StatusForbidden || msg != "invalid unit, missing association with target site" {
		t.Fatalf("expected unit site mismatch")
	}
	unit.SiteID = 10

	// 10. Forbidden rental state branches
	for _, st := range []string{LOCK_STATE_OVERLOCK, LOCK_STATE_GATELOCK, LOCK_STATE_PRELET} {
		unit.RentalState = st
		bankState := &mockBank{user: user, unit: unit}
		status, msg = AccessCodeEdit(context.Background(), bankState, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusForbidden {
			t.Fatalf("expected forbidden for state %s", st)
		}
	}
	unit.RentalState = "rented"

	// 11. Validation error (function error)
	bankValErr := &mockBank{user: user, unit: unit}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { return errors.New("val err") }, func() {
		status, _ = AccessCodeEdit(context.Background(), bankValErr, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusInternalServerError {
			t.Fatalf("expected validation internal error")
		}
	})

	// 12. Invalid (IsValid false duplicate)
	bankDup := &mockBank{user: user, unit: unit}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error {
		g[0].IsValid = false
		g[0].ValidationMessages = []string{AccessCodeMsgDuplicate}
		return nil
	}, func() {
		status, msg = AccessCodeEdit(context.Background(), bankDup, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusConflict || msg != "Duplicate access code" {
			t.Fatalf("expected duplicate conflict")
		}
	})

	// 13. Invalid (IsValid false generic)
	bankInv := &mockBank{user: user, unit: unit}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = false; return nil }, func() {
		status, msg = AccessCodeEdit(context.Background(), bankInv, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusBadRequest || msg != "Invalid access code" {
			t.Fatalf("expected invalid access code bad request")
		}
	})

	// 14. Early success (existing code same & active)
	bankEarly := &mockBank{user: user, unit: unit, codes: []GateAccessCode{{AccessCode: "1234", State: AccessCodeStateActive}}}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = true; return nil }, func() {
		status, _ = AccessCodeEdit(context.Background(), bankEarly, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusOK {
			t.Fatalf("expected early OK")
		}
	})

	// 15. UpdateAccessCodes remove error
	bankRemoveErr := &mockBank{user: user, unit: unit, codes: []GateAccessCode{{AccessCode: "old", State: AccessCodeStateActive}}, updateErrs: []error{errors.New("bad remove")}}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = true; return nil }, func() {
		status, _ = AccessCodeEdit(context.Background(), bankRemoveErr, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusInternalServerError {
			t.Fatalf("expected remove update error")
		}
	})

	// 16. UpdateAccessCodes add error
	bankAddErr := &mockBank{user: user, unit: unit, codes: []GateAccessCode{{AccessCode: "old", State: AccessCodeStateActive}}, updateErrs: []error{nil, errors.New("bad add")}}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = true; return nil }, func() {
		status, _ = AccessCodeEdit(context.Background(), bankAddErr, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusInternalServerError {
			t.Fatalf("expected add update error")
		}
	})

	// 17. Revoke error
	bankRevokeErr := &mockBank{user: user, unit: unit, codes: []GateAccessCode{{AccessCode: "old", State: AccessCodeStateActive}}, cc: &mockCC{revokeErr: errors.New("revoke fail")}}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = true; return nil }, func() {
		status, _ = AccessCodeEdit(context.Background(), bankRevokeErr, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusInternalServerError {
			t.Fatalf("expected revoke error")
		}
	})

	// 18. Set error
	bankSetErr := &mockBank{user: user, unit: unit, codes: []GateAccessCode{{AccessCode: "old", State: AccessCodeStateActive}}, cc: &mockCC{setErr: errors.New("set fail")}}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = true; return nil }, func() {
		status, _ = AccessCodeEdit(context.Background(), bankSetErr, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusInternalServerError {
			t.Fatalf("expected set error")
		}
	})

	// 19. Success path
	bankSuccess := &mockBank{user: user, unit: unit, codes: []GateAccessCode{{AccessCode: "old", State: AccessCodeStateActive}}, cc: &mockCC{}}
	withValidateAccessCodes(func(g GateAccessCodes, b *Bank) error { g[0].IsValid = true; return nil }, func() {
		status, msg = AccessCodeEdit(context.Background(), bankSuccess, baseClaims(), accessCodeEditInput{UserUUID: "1", UnitUUIDs: []string{"5"}, AccessCode: "1234"})
		if status != http.StatusOK || msg != "" {
			t.Fatalf("expected success got %d %s", status, msg)
		}
	})
}
