package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type mockBank struct {
	getBUserByID           func(int) (*BUser, error)
	v2UnitGetById          func(int, int) (*Unit, error)
	getCodesForUnits       func([]int, int) ([]GateAccessCode, error)
	updateAccessCodes      func([]string, int) error
	newCommandCenterClient func(int, context.Context) CommandCenter
}

func (m *mockBank) GetBUserByID(id int) (*BUser, error) {
	if m.getBUserByID != nil {
		return m.getBUserByID(id)
	}
	return &BUser{CompanyUUID: "company-1", Id: id, Sites: []string{"100"}}, nil
}

func (m *mockBank) V2UnitGetById(unitID int, siteID int) (*Unit, error) {
	if m.v2UnitGetById != nil {
		return m.v2UnitGetById(unitID, siteID)
	}
	return &Unit{SiteID: siteID, RentalState: ""}, nil
}

func (m *mockBank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	if m.getCodesForUnits != nil {
		return m.getCodesForUnits(units, siteID)
	}
	return []GateAccessCode{}, nil
}

func (m *mockBank) UpdateAccessCodes(codes []string, siteID int) error {
	if m.updateAccessCodes != nil {
		return m.updateAccessCodes(codes, siteID)
	}
	return nil
}

func (m *mockBank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenter {
	if m.newCommandCenterClient != nil {
		return m.newCommandCenterClient(siteID, ctx)
	}
	return &mockCommandCenter{}
}

type mockCommandCenter struct {
	revoke func([]int, map[string]struct{}) error
	set    func([]int, map[string]struct{}) error
}

func (m *mockCommandCenter) RevokeAccessCodes(units []int, options map[string]struct{}) error {
	if m.revoke != nil {
		return m.revoke(units, options)
	}
	return nil
}

func (m *mockCommandCenter) SetAccessCodes(units []int, options map[string]struct{}) error {
	if m.set != nil {
		return m.set(units, options)
	}
	return nil
}

func newRequest(t *testing.T, body string, claims *Claims, bank BankService) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	var buf *bytes.Reader
	if body != "" {
		buf = bytes.NewReader([]byte(body))
	} else {
		buf = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(http.MethodPost, "/access-code", buf)
	ctx := req.Context()
	if claims != nil {
		ctx = context.WithValue(ctx, claimsKey, claims)
	}
	if bank != nil {
		ctx = NewBankContext(ctx, bank)
	}
	req = req.WithContext(ctx)
	return req, httptest.NewRecorder()
}

func withDefaultSuccessDeps(t *testing.T) (*Claims, *mockBank, *mockCommandCenter) {
	t.Helper()
	bank := &mockBank{}
	cc := &mockCommandCenter{}
	bank.newCommandCenterClient = func(int, context.Context) CommandCenter { return cc }
	claims := &Claims{CompanyUUID: "company-1", CurrentSite: 100, CurrentSiteUUID: "site-uuid", UserID: 77}
	return claims, bank, cc
}

func TestAccessCodeEditHandler(t *testing.T) {
	originalValidator := accessCodeValidator
	originalEditFn := userAccessCodeEditFn
	originalNow := nowFunc
	t.Cleanup(func() {
		accessCodeValidator = originalValidator
		userAccessCodeEditFn = originalEditFn
		nowFunc = originalNow
	})

	baseBody := `{"userUuid":"1","unitUUIDs":["201"],"accessCode":"9999"}`

	t.Run("missingClaimsReturnsUnauthorized", func(t *testing.T) {
		bank := &mockBank{}
		req, rr := newRequest(t, baseBody, nil, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("missingBankReturnsError", func(t *testing.T) {
		claims := &Claims{CompanyUUID: "company-1"}
		req, rr := newRequest(t, baseBody, claims, nil)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("malformedJSONReturnsBadRequest", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		req, rr := newRequest(t, "{bad json", claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalidUserUUIDReturnsBadRequest", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		req, rr := newRequest(t, `{"userUuid":"not-int","unitUUIDs":["201"],"accessCode":"9999"}`, claims, bank)
		bank.getBUserByID = func(int) (*BUser, error) {
			t.Fatal("GetBUserByID should not be called")
			return nil, nil
		}
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalidUnitUUIDReturnsBadRequest", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		req, rr := newRequest(t, `{"userUuid":"1","unitUUIDs":["abc"],"accessCode":"9999"}`, claims, bank)
		bank.getBUserByID = func(int) (*BUser, error) {
			t.Fatal("GetBUserByID should not be called before unit conversion")
			return nil, nil
		}
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("userNotFoundReturnsNotFound", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return nil, errors.New("no_ob_found")
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("userLookupFailureReturnsServerError", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return nil, errors.New("db failure")
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("userCompanyMismatchReturnsForbidden", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: "other", Sites: []string{"100"}}, nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("userMissingSiteAssociationReturnsForbidden", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{"999"}}, nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("unitLookupFailureReturnsNotFound", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.v2UnitGetById = func(int, int) (*Unit, error) {
			return nil, errors.New("not found")
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("unitSiteMismatchReturnsForbidden", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.v2UnitGetById = func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: siteID + 1}, nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("unitLockedReturnsForbidden", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.v2UnitGetById = func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: siteID, RentalState: LOCK_STATE_OVERLOCK}, nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("validatorErrorReturnsServerError", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		accessCodeValidator = stubValidator{err: errors.New("validation failure")}
		t.Cleanup(func() { accessCodeValidator = originalValidator })
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("validatorDuplicateReturnsConflict", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		accessCodeValidator = stubValidator{valid: false, messages: []string{AccessCodeMsgDuplicate}}
		t.Cleanup(func() { accessCodeValidator = originalValidator })
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusConflict {
			t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
		}
	})

	t.Run("validatorInvalidReturnsBadRequest", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		accessCodeValidator = stubValidator{valid: false, messages: []string{"other"}}
		t.Cleanup(func() { accessCodeValidator = originalValidator })
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("existingActiveMatchingCodeShortCircuits", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{AccessCode: "9999", State: AccessCodeStateActive}}, nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("getCodesFailureReturnsServerError", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return nil, errors.New("codes failure")
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("updateAccessCodesRemoveFailureReturnsServerError", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{UnitID: 201, State: AccessCodeStateActive, AccessCode: "old"}}, nil
		}
		bank.updateAccessCodes = func(codes []string, siteID int) error {
			if len(codes) > 0 && codes[0] == "old" {
				return errors.New("remove failure")
			}
			return nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("updateAccessCodesAddFailureReturnsServerError", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{UnitID: 201, State: AccessCodeStateActive, AccessCode: "old"}}, nil
		}
		bank.updateAccessCodes = func(codes []string, siteID int) error {
			if len(codes) > 0 && codes[0] == "9999" {
				return errors.New("insert failure")
			}
			return nil
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("revokeAccessCodesFailureReturnsServerError", func(t *testing.T) {
		claims, bank, cc := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{UnitID: 201, State: AccessCodeStateActive, AccessCode: "old"}}, nil
		}
		cc.revoke = func([]int, map[string]struct{}) error {
			return errors.New("revoke failure")
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("setAccessCodesFailureReturnsServerError", func(t *testing.T) {
		claims, bank, cc := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{UnitID: 201, State: AccessCodeStateActive, AccessCode: "old"}}, nil
		}
		cc.set = func([]int, map[string]struct{}) error {
			return errors.New("set failure")
		}
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("activityRecordFailureReturnsServerError", func(t *testing.T) {
		claims, bank, _ := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{UnitID: 201, State: AccessCodeStateActive, AccessCode: "old"}}, nil
		}
		userAccessCodeEditFn = func(int, *BUser, string, time.Time) (*string, error) {
			return nil, errors.New("activity failure")
		}
		t.Cleanup(func() { userAccessCodeEditFn = originalEditFn })
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("successfulRequestReturnsOK", func(t *testing.T) {
		claims, bank, cc := withDefaultSuccessDeps(t)
		bank.getBUserByID = func(int) (*BUser, error) {
			return &BUser{CompanyUUID: claims.CompanyUUID, Sites: []string{strconv.Itoa(claims.CurrentSite)}}, nil
		}
		bank.getCodesForUnits = func([]int, int) ([]GateAccessCode, error) {
			return []GateAccessCode{{UnitID: 201, State: AccessCodeStateActive, AccessCode: "old"}}, nil
		}
		cc.revoke = func([]int, map[string]struct{}) error { return nil }
		cc.set = func([]int, map[string]struct{}) error { return nil }
		userAccessCodeEditFn = func(int, *BUser, string, time.Time) (*string, error) {
			result := "ok"
			return &result, nil
		}
		nowFunc = func() time.Time { return time.Unix(0, 0) }
		t.Cleanup(func() {
			userAccessCodeEditFn = originalEditFn
			nowFunc = originalNow
		})
		req, rr := newRequest(t, baseBody, claims, bank)
		AccessCodeEditHandler(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type stubValidator struct {
	err      error
	valid    bool
	messages []string
}

func (s stubValidator) Validate(bank BankService, codes GateAccessCodes) error {
	if s.err != nil {
		return s.err
	}
	for i := range codes {
		codes[i].IsValid = s.valid
		if s.messages != nil {
			codes[i].ValidationMessages = append([]string{}, s.messages...)
		}
	}
	return nil
}
