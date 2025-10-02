package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock implementations for testing
type MockBUserService struct {
	GetBUserByIDFunc func(userID int) (*BUser, error)
}

func (m *MockBUserService) GetBUserByID(userID int) (*BUser, error) {
	return m.GetBUserByIDFunc(userID)
}

type MockUnitService struct {
	GetUnitByIDFunc func(unitID, siteID int) (*Unit, error)
}

func (m *MockUnitService) GetUnitByID(unitID, siteID int) (*Unit, error) {
	return m.GetUnitByIDFunc(unitID, siteID)
}

type MockAccessCodeService struct {
	GetCodesForUnitsFunc    func(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodesFunc   func(codes []string, siteID int) error
	ValidateAccessCodesFunc func(codes GateAccessCodes) error
}

func (m *MockAccessCodeService) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	return m.GetCodesForUnitsFunc(units, siteID)
}

func (m *MockAccessCodeService) UpdateAccessCodes(codes []string, siteID int) error {
	return m.UpdateAccessCodesFunc(codes, siteID)
}

func (m *MockAccessCodeService) ValidateAccessCodes(codes GateAccessCodes) error {
	// Always set codes as valid regardless of whether custom function is provided
	for i := range codes {
		codes[i].IsValid = true
	}

	// Then call custom function if it exists
	if m.ValidateAccessCodesFunc != nil {
		return m.ValidateAccessCodesFunc(codes)
	}
	return nil
}

type MockCommandCenterClient struct {
	RevokeAccessCodesFunc func(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodesFunc    func(units []int, options map[string]struct{}) error
}

func (m *MockCommandCenterClient) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	return m.RevokeAccessCodesFunc(revokeUnits, options)
}

func (m *MockCommandCenterClient) SetAccessCodes(units []int, options map[string]struct{}) error {
	return m.SetAccessCodesFunc(units, options)
}

type MockCommandCenterService struct {
	NewCommandCenterClientFunc func(siteID int, ctx context.Context) CommandCenterClient
}

func (m *MockCommandCenterService) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClient {
	if m.NewCommandCenterClientFunc != nil {
		return m.NewCommandCenterClientFunc(siteID, ctx)
	}
	// Default implementation returns a mock with no-op functions
	return &MockCommandCenterClient{
		RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
			return nil
		},
		SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
			return nil
		},
	}
}

type MockActivityService struct {
	LogUserAccessCodeEditFunc func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error)
}

func (m *MockActivityService) LogUserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	if m.LogUserAccessCodeEditFunc != nil {
		return m.LogUserAccessCodeEditFunc(userID, user, siteUUID, timestamp)
	}
	// Default implementation returns nil, nil (successful log)
	return nil, nil
}

// Helper function to create a request with claims context
func createRequestWithClaims(method, url string, body interface{}, claims *Claims) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	ctx := context.WithValue(req.Context(), claimsKey, claims)
	return req.WithContext(ctx)
}

func TestAccessCodeEditHandler(t *testing.T) {
	// Test cases
	t.Run("missing claims in context", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request without claims
		req := httptest.NewRequest(http.MethodPost, "/accesscode/edit", nil)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		}
		if w.Body.String() != "Unauthorized\n" {
			t.Errorf("Expected body %q, got %q", "Unauthorized\n", w.Body.String())
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims but invalid body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", "invalid-json", claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid UUID conversion", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body with invalid UUID
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserUUID:   "not-a-number",
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		expectedError := fmt.Sprintf("invalid UUID: %s", input.UserUUID)
		if w.Body.String() != expectedError+"\n" {
			t.Errorf("Expected body %q, got %q", expectedError+"\n", w.Body.String())
		}
	})

	t.Run("invalid unit UUID conversion", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body with invalid Unit UUID
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitUUIDs:  []string{"not-a-number"},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		expectedError := fmt.Sprintf("invalid UUID: %s", input.UnitUUIDs[0])
		if w.Body.String() != expectedError+"\n" {
			t.Errorf("Expected body %q, got %q", expectedError+"\n", w.Body.String())
		}
	})

	t.Run("user not found", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return nil, errors.New("no_ob_found")
			},
		}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
		if w.Body.String() != "user not found\n" {
			t.Errorf("Expected body %q, got %q", "user not found\n", w.Body.String())
		}
	})

	t.Run("user from different company", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "different-company",
					Id:          789,
					Sites:       []string{"123"},
				}, nil
			},
		}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
		}
		if w.Body.String() != "user not found in company\n" {
			t.Errorf("Expected body %q, got %q", "user not found in company\n", w.Body.String())
		}
	})

	t.Run("user not associated with site", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"999"}, // Different site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
		}
		if w.Body.String() != "invalid user, missing association with target site\n" {
			t.Errorf("Expected body %q, got %q", "invalid user, missing association with target site\n", w.Body.String())
		}
	})

	t.Run("unit not found", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return nil, errors.New("unit not found")
			},
		}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
		if w.Body.String() != "unit not found\n" {
			t.Errorf("Expected body %q, got %q", "unit not found\n", w.Body.String())
		}
	})

	t.Run("unit not associated with site", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      999, // Different site
					RentalState: "available",
				}, nil
			},
		}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
		}
		if w.Body.String() != "invalid unit, missing association with target site\n" {
			t.Errorf("Expected body %q, got %q", "invalid unit, missing association with target site\n", w.Body.String())
		}
	})

	t.Run("unit in overlock state", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123,                 // Matching site
					RentalState: LOCK_STATE_OVERLOCK, // Locked state
				}, nil
			},
		}
		mockAccessCodeService := &MockAccessCodeService{}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
		}
		expectedError := fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_OVERLOCK)
		if w.Body.String() != expectedError+"\n" {
			t.Errorf("Expected body %q, got %q", expectedError+"\n", w.Body.String())
		}
	})

	t.Run("access code validation failure", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}
		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				for i := range codes {
					codes[i].IsValid = false
					codes[i].ValidationMessages = []string{"invalid_code"}
				}
				return nil
			},
		}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		if w.Body.String() != "invalid access code\n" {
			t.Errorf("Expected body %q, got %q", "invalid access code\n", w.Body.String())
		}
	})

	t.Run("duplicate access code", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}
		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				for i := range codes {
					codes[i].IsValid = false
					codes[i].ValidationMessages = []string{AccessCodeMsgDuplicate}
				}
				return nil
			},
		}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status code %d, got %d", http.StatusConflict, w.Code)
		}
		if w.Body.String() != "duplicate access code\n" {
			t.Errorf("Expected body %q, got %q", "duplicate access code\n", w.Body.String())
		}
	})

	t.Run("error getting existing access codes", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}
		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				// Set codes as valid
				for i := range codes {
					codes[i].IsValid = true
				}
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				return nil, errors.New("database error")
			},
		}
		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
		if w.Body.String() != "error retrieving existing access codes: database error\n" {
			t.Errorf("Expected body %q, got %q", "error retrieving existing access codes: database error\n", w.Body.String())
		}
	})

	t.Run("successful update with no existing codes", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}
		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				// Set codes as valid
				for i := range codes {
					codes[i].IsValid = true
				}
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				return []GateAccessCode{}, nil // No existing codes
			},
			UpdateAccessCodesFunc: func(codes []string, siteID int) error {
				return nil // Successful update
			},
		}

		mockCommandClient := &MockCommandCenterClient{
			SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
				return nil // Successful set
			},
		}

		mockCommandCenterService := &MockCommandCenterService{
			NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterClient {
				return mockCommandClient
			},
		}

		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil // Successful log
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: "1234",
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("skip update for existing active code", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}

		accessCode := "1234"

		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				// Set codes as valid
				for i := range codes {
					codes[i].IsValid = true
				}
				return nil
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				// Code already exists and is active
				return []GateAccessCode{
					{
						AccessCode: accessCode,
						UnitID:     101,
						UserID:     789,
						SiteID:     123,
						State:      AccessCodeStateActive,
					},
				}, nil
			},
		}

		mockCommandCenterService := &MockCommandCenterService{}
		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil // Successful log
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: accessCode,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("error revoking existing access codes", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}

		accessCode := "1234"

		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				// Set codes as valid
				for i := range codes {
					codes[i].IsValid = true
				}
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				// Existing code with different value
				return []GateAccessCode{
					{
						AccessCode: "5678", // Different from the new code
						UnitID:     101,
						UserID:     789,
						SiteID:     123,
						State:      AccessCodeStateActive,
					},
				}, nil
			},
			UpdateAccessCodesFunc: func(codes []string, siteID int) error {
				return nil // Successful update
			},
		}

		mockCommandClient := &MockCommandCenterClient{
			RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
				return errors.New("failed to revoke")
			},
		}

		mockCommandCenterService := &MockCommandCenterService{
			NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterClient {
				return mockCommandClient
			},
		}

		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: accessCode,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
		if w.Body.String() != "failed to revoke previous access codes: failed to revoke\n" {
			t.Errorf("Expected body %q, got %q", "failed to revoke previous access codes: failed to revoke\n", w.Body.String())
		}
	})

	t.Run("error setting new access codes", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}

		accessCode := "1234"

		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				// Existing code with different value
				return []GateAccessCode{
					{
						AccessCode: "5678", // Different from the new code
						UnitID:     101,
						UserID:     789,
						SiteID:     123,
						State:      AccessCodeStateActive,
					},
				}, nil
			},
			UpdateAccessCodesFunc: func(codes []string, siteID int) error {
				return nil // Successful update
			},
		}

		mockCommandClient := &MockCommandCenterClient{
			RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
				return nil
			},
			SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
				return errors.New("failed to set")
			},
		}

		mockCommandCenterService := &MockCommandCenterService{
			NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterClient {
				return mockCommandClient
			},
		}

		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: accessCode,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
		if w.Body.String() != "failed to set access codes: failed to set\n" {
			t.Errorf("Expected body %q, got %q", "failed to set access codes: failed to set\n", w.Body.String())
		}
	})

	t.Run("error logging activity", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}

		accessCode := "1234"

		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				// No existing codes
				return []GateAccessCode{}, nil
			},
			UpdateAccessCodesFunc: func(codes []string, siteID int) error {
				return nil // Successful update
			},
		}

		mockCommandClient := &MockCommandCenterClient{
			SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
				return nil
			},
		}

		mockCommandCenterService := &MockCommandCenterService{
			NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterClient {
				return mockCommandClient
			},
		}

		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, errors.New("activity log failure")
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: accessCode,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}
		if w.Body.String() != "Recording activity failed: activity log failure\n" {
			t.Errorf("Expected body %q, got %q", "Recording activity failed: activity log failure\n", w.Body.String())
		}
	})

	t.Run("complete successful flow", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}
		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: "available",
				}, nil
			},
		}

		accessCode := "1234"

		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				// Existing code with different value
				return []GateAccessCode{
					{
						AccessCode: "5678", // Different from the new code
						UnitID:     101,
						UserID:     789,
						SiteID:     123,
						State:      AccessCodeStateActive,
					},
				}, nil
			},
			UpdateAccessCodesFunc: func(codes []string, siteID int) error {
				return nil // Successful update
			},
		}

		mockCommandClient := &MockCommandCenterClient{
			RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
				return nil // Successful revoke
			},
			SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
				return nil // Successful set
			},
		}

		mockCommandCenterService := &MockCommandCenterService{
			NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterClient {
				return mockCommandClient
			},
		}

		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil // Successful log
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101},
			AccessCode: accessCode,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("multiple units", func(t *testing.T) {
		// Setup
		mockUserService := &MockBUserService{
			GetBUserByIDFunc: func(userID int) (*BUser, error) {
				return &BUser{
					CompanyUUID: "company-1",
					Id:          789,
					Sites:       []string{"123"}, // Matching site
				}, nil
			},
		}

		unitStates := map[int]string{
			101: "available",
			102: "available",
		}

		mockUnitService := &MockUnitService{
			GetUnitByIDFunc: func(unitID, siteID int) (*Unit, error) {
				state, ok := unitStates[unitID]
				if !ok {
					return nil, errors.New("unit not found")
				}
				return &Unit{
					SiteID:      123, // Matching site
					RentalState: state,
				}, nil
			},
		}

		accessCode := "1234"

		mockAccessCodeService := &MockAccessCodeService{
			ValidateAccessCodesFunc: func(codes GateAccessCodes) error {
				return nil // Valid access codes
			},
			GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
				// No existing codes
				return []GateAccessCode{}, nil
			},
			UpdateAccessCodesFunc: func(codes []string, siteID int) error {
				return nil // Successful update
			},
		}

		mockCommandClient := &MockCommandCenterClient{
			SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
				return nil // Successful set
			},
		}

		mockCommandCenterService := &MockCommandCenterService{
			NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterClient {
				return mockCommandClient
			},
		}

		mockActivityService := &MockActivityService{
			LogUserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
				return nil, nil // Successful log
			},
		}

		handler := NewAccessCodeHandler(
			mockUserService,
			mockUnitService,
			mockAccessCodeService,
			mockCommandCenterService,
			mockActivityService,
		)

		// Create request with claims and body
		claims := &Claims{
			CompanyUUID:     "company-1",
			CurrentSite:     123,
			CurrentSiteUUID: "site-123",
			UserID:          456,
		}

		input := AccessCodeEditInput{
			UserID:     789,
			UnitIDs:    []int{101, 102}, // Multiple units
			AccessCode: accessCode,
		}

		req := createRequestWithClaims(http.MethodPost, "/accesscode/edit", input, claims)
		w := httptest.NewRecorder()

		// Execute
		handler.AccessCodeEditHandler(w, req)

		// Verify
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})
}
