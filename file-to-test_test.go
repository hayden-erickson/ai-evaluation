package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockBank struct {
	GetBUserByIDFunc           func(BUserID int) (*BUser, error)
	V2UnitGetByIdFunc          func(unitID int, siteID int) (*Unit, error)
	GetCodesForUnitsFunc       func(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodesFunc      func(codes []string, siteID int) error
	NewCommandCenterClientFunc func(siteID int, ctx context.Context) CommandCenter
}

func (m *MockBank) GetBUserByID(BUserID int) (*BUser, error) {
	return m.GetBUserByIDFunc(BUserID)
}

func (m *MockBank) V2UnitGetById(unitID int, siteID int) (*Unit, error) {
	return m.V2UnitGetByIdFunc(unitID, siteID)
}

func (m *MockBank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return m.GetCodesForUnitsFunc(units, siteID)
}

func (m *MockBank) UpdateAccessCodes(codes []string, siteID int) error {
	return m.UpdateAccessCodesFunc(codes, siteID)
}

func (m *MockBank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenter {
	return m.NewCommandCenterClientFunc(siteID, ctx)
}

type MockCommandCenter struct {
	RevokeAccessCodesFunc func(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodesFunc    func(units []int, options map[string]struct{}) error
}

func (m *MockCommandCenter) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	return m.RevokeAccessCodesFunc(revokeUnits, options)
}

func (m *MockCommandCenter) SetAccessCodes(units []int, options map[string]struct{}) error {
	return m.SetAccessCodesFunc(units, options)
}

func TestAccessCodeEditHandler_Unauthorized(t *testing.T) {
	req, err := http.NewRequest("POST", "/access-code-edit", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AccessCodeEditHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestAccessCodeEditHandler_NoBankInContext(t *testing.T) {
	req, err := http.NewRequest("POST", "/access-code-edit", nil)
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID:     "test-company",
		CurrentSite:     1,
		CurrentSiteUUID: "site-uuid-1",
		UserID:          100,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AccessCodeEditHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "Failed to retrieve bank from context"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAccessCodeEditHandler_UserNotFound(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return nil, errors.New("no_ob_found")
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"]}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestAccessCodeEditHandler_UserNotInCompany(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "other-company"}, nil
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"]}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestAccessCodeEditHandler_UserNotAssociatedWithSite(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "test-company", Sites: []string{"2"}}, nil
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"]}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestAccessCodeEditHandler_UnitNotFound(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "test-company", Sites: []string{"1"}}, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return nil, errors.New("unit not found")
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"]}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestAccessCodeEditHandler_UnitNotAssociatedWithSite(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "test-company", Sites: []string{"1"}}, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: 2}, nil
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"]}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestAccessCodeEditHandler_UnitOverlocked(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "test-company", Sites: []string{"1"}}, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: 1, RentalState: LOCK_STATE_OVERLOCK}, nil
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"]}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestAccessCodeEditHandler_DuplicateAccessCode(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "test-company", Sites: []string{"1"}}, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: 1, RentalState: "available"}, nil
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"], "accessCode": "1234"}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	originalValidate := ValidateGateAccessCodes
	ValidateGateAccessCodes = func(gacs GateAccessCodes, bank Banker) error {
		gacs[0].IsValid = false
		gacs[0].ValidationMessages = []string{AccessCodeMsgDuplicate}
		return nil
	}
	defer func() { ValidateGateAccessCodes = originalValidate }()

	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestAccessCodeEditHandler_InvalidAccessCode(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{CompanyUUID: "test-company", Sites: []string{"1"}}, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: 1, RentalState: "available"}, nil
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"], "accessCode": "1234"}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID: "test-company",
		CurrentSite: 1,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	originalValidate := ValidateGateAccessCodes
	ValidateGateAccessCodes = func(gacs GateAccessCodes, bank Banker) error {
		gacs[0].IsValid = false
		gacs[0].ValidationMessages = []string{"some other error"}
		return nil
	}
	defer func() { ValidateGateAccessCodes = originalValidate }()

	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestAccessCodeEditHandler_HappyPath(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{Id: 1, CompanyUUID: "test-company", Sites: []string{"1"}}, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{SiteID: 1, RentalState: "available"}, nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{
				{UnitID: 101, AccessCode: "5678", State: AccessCodeStateActive},
			}, nil
		},
		UpdateAccessCodesFunc: func(codes []string, siteID int) error {
			return nil
		},
		NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenter {
			return &MockCommandCenter{
				RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
					return nil
				},
				SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
					return nil
				},
			}
		},
	}

	req, err := http.NewRequest("POST", "/access-code-edit", strings.NewReader(`{"userId": 1, "unitUUIDs": ["101"], "accessCode": "1234"}`))
	if err != nil {
		t.Fatal(err)
	}

	claims := &Claims{
		CompanyUUID:     "test-company",
		CurrentSite:     1,
		CurrentSiteUUID: "site-uuid-1",
		UserID:          100,
	}
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	ctx = NewBankContext(ctx, mockBank)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	originalValidate := ValidateGateAccessCodes
	ValidateGateAccessCodes = func(gacs GateAccessCodes, bank Banker) error {
		gacs[0].IsValid = true
		return nil
	}
	defer func() { ValidateGateAccessCodes = originalValidate }()

	handler := &AccessCodeHandler{Bank: mockBank}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
