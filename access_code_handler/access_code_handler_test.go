package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper functions for testing
func createTestContext(claims *Claims) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, claimsKey, claims)
}

func createTestRequest(ctx context.Context, input interface{}) *http.Request {
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/access-code-edit", bytes.NewReader(body))
	return req.WithContext(ctx)
}

// Mock implementation for testing
type mockCommandCenterClient struct {
	revokeErr error
	setErr    error
}

func (m *mockCommandCenterClient) RevokeAccessCodes(units []int, options map[string]struct{}) error {
	return m.revokeErr
}

func (m *mockCommandCenterClient) SetAccessCodes(units []int, options map[string]struct{}) error {
	return m.setErr
}

func newCommandCenterClient(m *mockCommandCenterClient) *CommandCenterClient {
	return &CommandCenterClient{m}
}

type mockBank struct {
	userErr   error
	user      *BUser
	unitErr   error
	unit      *Unit
	codes     []GateAccessCode
	codesErr  error
	updateErr error
	cc        *mockCommandCenterClient
}

func (m *mockBank) GetBUserByID(id int) (*BUser, error) {
	return m.user, m.userErr
}

func (m *mockBank) V2UnitGetById(unitID, siteID int) (*Unit, error) {
	return m.unit, m.unitErr
}

func (m *mockBank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return m.codes, m.codesErr
}

func (m *mockBank) UpdateAccessCodes(codes []string, siteID int) error {
	return m.updateErr
}

func (m *mockBank) NewCommandCenterClient(siteID int, ctx context.Context) *CommandCenterClient {
	if m.cc == nil {
		m.cc = &mockCommandCenterClient{}
	}
	return newCommandCenterClient(m.cc)
}

// Test cases for the handler
func TestAccessCodeEditHandler_ClaimsNotFound(t *testing.T) {
	ctx := context.Background()
	req := httptest.NewRequest(http.MethodPost, "/access-code-edit", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	AccessCodeEditHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAccessCodeEditHandler_BankNotFound(t *testing.T) {
	claims := &Claims{
		CompanyUUID:     "test-company",
		CurrentSite:     1,
		CurrentSiteUUID: "site-1",
		UserID:          1,
	}
	ctx := createTestContext(claims)
	req := httptest.NewRequest(http.MethodPost, "/access-code-edit", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	AccessCodeEditHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestAccessCodeEditHandler_InvalidJSON(t *testing.T) {
	claims := &Claims{
		CompanyUUID:     "test-company",
		CurrentSite:     1,
		CurrentSiteUUID: "site-1",
		UserID:          1,
	}
	ctx := createTestContext(claims)
	bank := &mockBank{}
	ctx = context.WithValue(ctx, bankKey, bank)

	req := httptest.NewRequest(http.MethodPost, "/access-code-edit", bytes.NewReader([]byte("invalid json")))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	AccessCodeEditHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAccessCodeService_ValidateUser(t *testing.T) {
	tests := []struct {
		name        string
		bank        *mockBank
		userID      int
		companyUUID string
		siteID      int
		wantErr     bool
		errMsg      string
	}{
		{
			name: "User not found",
			bank: &mockBank{
				userErr: fmt.Errorf("no_ob_found"),
			},
			userID:      1,
			companyUUID: "test-company",
			siteID:      1,
			wantErr:     true,
			errMsg:      "user not found",
		},
		{
			name: "Wrong company",
			bank: &mockBank{
				user: &BUser{CompanyUUID: "different-company"},
			},
			userID:      1,
			companyUUID: "test-company",
			siteID:      1,
			wantErr:     true,
			errMsg:      "user not found in company",
		},
		{
			name: "Site not associated",
			bank: &mockBank{
				user: &BUser{
					CompanyUUID: "test-company",
					Sites:       []string{"2", "3"},
				},
			},
			userID:      1,
			companyUUID: "test-company",
			siteID:      1,
			wantErr:     true,
			errMsg:      "invalid user, missing association with target site",
		},
		{
			name: "Valid user",
			bank: &mockBank{
				user: &BUser{
					CompanyUUID: "test-company",
					Sites:       []string{"1", "2"},
				},
			},
			userID:      1,
			companyUUID: "test-company",
			siteID:      1,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAccessCodeService(tt.bank)
			_, err := service.validateUser(tt.userID, tt.companyUUID, tt.siteID)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateUser() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestAccessCodeService_ValidateUnit(t *testing.T) {
	tests := []struct {
		name    string
		bank    *mockBank
		unitID  int
		siteID  int
		wantErr bool
		errMsg  string
	}{
		{
			name: "Unit not found",
			bank: &mockBank{
				unitErr: fmt.Errorf("not found"),
			},
			unitID:  1,
			siteID:  1,
			wantErr: true,
			errMsg:  "unit not found: not found",
		},
		{
			name: "Wrong site",
			bank: &mockBank{
				unit: &Unit{SiteID: 2},
			},
			unitID:  1,
			siteID:  1,
			wantErr: true,
			errMsg:  "invalid unit, missing association with target site",
		},
		{
			name: "Unit overlocked",
			bank: &mockBank{
				unit: &Unit{
					SiteID:      1,
					RentalState: LOCK_STATE_OVERLOCK,
				},
			},
			unitID:  1,
			siteID:  1,
			wantErr: true,
			errMsg:  "access code changes not allowed - unit in overlock",
		},
		{
			name: "Valid unit",
			bank: &mockBank{
				unit: &Unit{
					SiteID:      1,
					RentalState: "active",
				},
			},
			unitID:  1,
			siteID:  1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAccessCodeService(tt.bank)
			err := service.validateUnit(tt.unitID, tt.siteID)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateUnit() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestAccessCodeService_UpdateUnitAccessCode(t *testing.T) {
	tests := []struct {
		name       string
		bank       *mockBank
		unitID     int
		userID     int
		siteID     int
		accessCode string
		wantErr    bool
		errMsg     string
	}{
		{
			name: "Invalid access code",
			bank: &mockBank{
				codes: nil,
				cc:    &mockCommandCenterClient{},
			},
			unitID:     1,
			userID:     1,
			siteID:     1,
			accessCode: "123",
			wantErr:    true,
			errMsg:     "invalid access code",
		},
		{
			name: "Successful update",
			bank: &mockBank{
				codes: []GateAccessCode{
					{
						AccessCode: "654321",
						State:      AccessCodeStateActive,
					},
				},
				cc: &mockCommandCenterClient{},
			},
			unitID:     1,
			userID:     1,
			siteID:     1,
			accessCode: "123456",
			wantErr:    false,
		},
		{
			name: "Command center failure",
			bank: &mockBank{
				codes: []GateAccessCode{
					{
						AccessCode: "654321",
						State:      AccessCodeStateActive,
					},
				},
				cc: &mockCommandCenterClient{
					setErr: fmt.Errorf("command center error"),
				},
			},
			unitID:     1,
			userID:     1,
			siteID:     1,
			accessCode: "123456",
			wantErr:    true,
			errMsg:     "failed to set access codes: command center error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAccessCodeService(tt.bank)
			err := service.updateUnitAccessCode(context.Background(), tt.unitID, tt.userID, tt.siteID, tt.accessCode)

			if tt.wantErr && err == nil {
				t.Errorf("updateUnitAccessCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("updateUnitAccessCode() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
