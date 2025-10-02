package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Mock implementations for testing

type MockBank struct {
	GetBUserByIDFunc           func(BUserID int) (*BUser, error)
	V2UnitGetByIdFunc          func(unitID int, siteID int) (*Unit, error)
	GetCodesForUnitsFunc       func(units []int, siteID int) ([]GateAccessCode, error)
	UpdateAccessCodesFunc      func(codes []string, siteID int) error
	NewCommandCenterClientFunc func(siteID int, ctx context.Context) CommandCenterInterface
}

func (m *MockBank) GetBUserByID(BUserID int) (*BUser, error) {
	if m.GetBUserByIDFunc != nil {
		return m.GetBUserByIDFunc(BUserID)
	}
	return &BUser{}, nil
}

func (m *MockBank) V2UnitGetById(unitID int, siteID int) (*Unit, error) {
	if m.V2UnitGetByIdFunc != nil {
		return m.V2UnitGetByIdFunc(unitID, siteID)
	}
	return &Unit{}, nil
}

func (m *MockBank) GetCodesForUnits(units []int, siteID int) ([]GateAccessCode, error) {
	if m.GetCodesForUnitsFunc != nil {
		return m.GetCodesForUnitsFunc(units, siteID)
	}
	return []GateAccessCode{}, nil
}

func (m *MockBank) UpdateAccessCodes(codes []string, siteID int) error {
	if m.UpdateAccessCodesFunc != nil {
		return m.UpdateAccessCodesFunc(codes, siteID)
	}
	return nil
}

func (m *MockBank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface {
	if m.NewCommandCenterClientFunc != nil {
		return m.NewCommandCenterClientFunc(siteID, ctx)
	}
	return &MockCommandCenter{}
}

type MockCommandCenter struct {
	RevokeAccessCodesFunc func(revokeUnits []int, options map[string]struct{}) error
	SetAccessCodesFunc    func(units []int, options map[string]struct{}) error
}

func (m *MockCommandCenter) RevokeAccessCodes(revokeUnits []int, options map[string]struct{}) error {
	if m.RevokeAccessCodesFunc != nil {
		return m.RevokeAccessCodesFunc(revokeUnits, options)
	}
	return nil
}

func (m *MockCommandCenter) SetAccessCodes(units []int, options map[string]struct{}) error {
	if m.SetAccessCodesFunc != nil {
		return m.SetAccessCodesFunc(units, options)
	}
	return nil
}

type MockActivityRecorder struct {
	UserAccessCodeEditFunc func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error)
}

func (m *MockActivityRecorder) UserAccessCodeEdit(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
	if m.UserAccessCodeEditFunc != nil {
		return m.UserAccessCodeEditFunc(userID, user, siteUUID, timestamp)
	}
	return nil, nil
}

// Test helper functions

func createTestService(mockBank *MockBank, mockActivityRecorder *MockActivityRecorder) *AccessCodeEditService {
	if mockBank == nil {
		mockBank = &MockBank{}
	}
	if mockActivityRecorder == nil {
		mockActivityRecorder = &MockActivityRecorder{}
	}
	return NewAccessCodeEditService(mockBank, mockActivityRecorder)
}

func createTestClaims() *Claims {
	return &Claims{
		CompanyUUID:     "test-company-uuid",
		CurrentSite:     123,
		CurrentSiteUUID: "test-site-uuid",
		UserID:          456,
	}
}

func createTestUser() *BUser {
	return &BUser{
		CompanyUUID: "test-company-uuid",
		Id:          789,
		Sites:       []string{"123"},
	}
}

func createTestUnit() *Unit {
	return &Unit{
		SiteID:      123,
		RentalState: "active",
	}
}

// Test ConvertUUIDs function

func TestConvertUUIDs_InvalidUserUUID(t *testing.T) {
	service := createTestService(nil, nil)
	input := &AccessCodeEditRequest{
		UserUUID: "invalid-uuid",
	}

	err := service.ConvertUUIDs(input)
	if err == nil {
		t.Error("Expected error for invalid UserUUID, got nil")
	}
	if err.Error() != "invalid UUID: invalid-uuid" {
		t.Errorf("Expected 'invalid UUID: invalid-uuid', got '%s'", err.Error())
	}
}

func TestConvertUUIDs_InvalidUnitUUID(t *testing.T) {
	service := createTestService(nil, nil)
	input := &AccessCodeEditRequest{
		UserID:    123,
		UnitUUIDs: []string{"456", "invalid-uuid"},
	}

	err := service.ConvertUUIDs(input)
	if err == nil {
		t.Error("Expected error for invalid UnitUUID, got nil")
	}
	if err.Error() != "invalid UUID: invalid-uuid" {
		t.Errorf("Expected 'invalid UUID: invalid-uuid', got '%s'", err.Error())
	}
}

func TestConvertUUIDs_ValidConversion(t *testing.T) {
	service := createTestService(nil, nil)
	input := &AccessCodeEditRequest{
		UserUUID:  "123",
		UnitUUIDs: []string{"456", "789"},
	}

	err := service.ConvertUUIDs(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if input.UserID != 123 {
		t.Errorf("Expected UserID 123, got %d", input.UserID)
	}
	if len(input.UnitIDs) != 2 || input.UnitIDs[0] != 456 || input.UnitIDs[1] != 789 {
		t.Errorf("Expected UnitIDs [456, 789], got %v", input.UnitIDs)
	}
}

// Test ValidateUserAccess function

func TestValidateUserAccess_UserNotFound(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return nil, errors.New("no_ob_found")
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	user, err := service.ValidateUserAccess(123, claims)
	if err == nil {
		t.Error("Expected error for user not found, got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got '%s'", err.Error())
	}
	if user != nil {
		t.Error("Expected nil user, got non-nil")
	}
}

func TestValidateUserAccess_DatabaseError(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return nil, errors.New("database connection failed")
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	user, err := service.ValidateUserAccess(123, claims)
	if err == nil {
		t.Error("Expected error for database error, got nil")
	}
	if err.Error() != "internal server error" {
		t.Errorf("Expected 'internal server error', got '%s'", err.Error())
	}
	if user != nil {
		t.Error("Expected nil user, got non-nil")
	}
}

func TestValidateUserAccess_UserNotInCompany(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{
				CompanyUUID: "different-company-uuid",
				Id:          123,
				Sites:       []string{"123"},
			}, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	user, err := service.ValidateUserAccess(123, claims)
	if err == nil {
		t.Error("Expected error for user not in company, got nil")
	}
	if err.Error() != "user not found in company" {
		t.Errorf("Expected 'user not found in company', got '%s'", err.Error())
	}
	if user != nil {
		t.Error("Expected nil user, got non-nil")
	}
}

func TestValidateUserAccess_UserNotAssociatedWithSite(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return &BUser{
				CompanyUUID: "test-company-uuid",
				Id:          123,
				Sites:       []string{"456", "789"}, // Not including site 123
			}, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	user, err := service.ValidateUserAccess(123, claims)
	if err == nil {
		t.Error("Expected error for user not associated with site, got nil")
	}
	if err.Error() != "invalid user, missing association with target site" {
		t.Errorf("Expected 'invalid user, missing association with target site', got '%s'", err.Error())
	}
	if user != nil {
		t.Error("Expected nil user, got non-nil")
	}
}

func TestValidateUserAccess_Success(t *testing.T) {
	testUser := createTestUser()
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return testUser, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	user, err := service.ValidateUserAccess(123, claims)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if user != testUser {
		t.Error("Expected returned user to match test user")
	}
}

// Test ValidateUnitAccess function

func TestValidateUnitAccess_UnitNotFound(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return nil, errors.New("unit not found")
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	unit, err := service.ValidateUnitAccess(123, claims)
	if err == nil {
		t.Error("Expected error for unit not found, got nil")
	}
	if err.Error() != "unit not found" {
		t.Errorf("Expected 'unit not found', got '%s'", err.Error())
	}
	if unit != nil {
		t.Error("Expected nil unit, got non-nil")
	}
}

func TestValidateUnitAccess_UnitNotAssociatedWithSite(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{
				SiteID:      456, // Different site
				RentalState: "active",
			}, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	unit, err := service.ValidateUnitAccess(123, claims)
	if err == nil {
		t.Error("Expected error for unit not associated with site, got nil")
	}
	if err.Error() != "invalid unit, missing association with target site" {
		t.Errorf("Expected 'invalid unit, missing association with target site', got '%s'", err.Error())
	}
	if unit != nil {
		t.Error("Expected nil unit, got non-nil")
	}
}

func TestValidateUnitAccess_UnitInOverlockState(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{
				SiteID:      123,
				RentalState: LOCK_STATE_OVERLOCK,
			}, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	unit, err := service.ValidateUnitAccess(123, claims)
	if err == nil {
		t.Error("Expected error for unit in overlock state, got nil")
	}
	expectedMsg := fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_OVERLOCK)
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}
	if unit != nil {
		t.Error("Expected nil unit, got non-nil")
	}
}

func TestValidateUnitAccess_UnitInGatelockState(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{
				SiteID:      123,
				RentalState: LOCK_STATE_GATELOCK,
			}, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	unit, err := service.ValidateUnitAccess(123, claims)
	if err == nil {
		t.Error("Expected error for unit in gatelock state, got nil")
	}
	expectedMsg := fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_GATELOCK)
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}
	if unit != nil {
		t.Error("Expected nil unit, got non-nil")
	}
}

func TestValidateUnitAccess_UnitInPreletState(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return &Unit{
				SiteID:      123,
				RentalState: LOCK_STATE_PRELET,
			}, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	unit, err := service.ValidateUnitAccess(123, claims)
	if err == nil {
		t.Error("Expected error for unit in prelet state, got nil")
	}
	expectedMsg := fmt.Sprintf("access code changes not allowed - unit in %s", LOCK_STATE_PRELET)
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}
	if unit != nil {
		t.Error("Expected nil unit, got non-nil")
	}
}

func TestValidateUnitAccess_Success(t *testing.T) {
	testUnit := createTestUnit()
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return testUnit, nil
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	unit, err := service.ValidateUnitAccess(123, claims)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if unit != testUnit {
		t.Error("Expected returned unit to match test unit")
	}
}

// Test ProcessAccessCodeForUnit function - focuses on early termination scenarios

func TestProcessAccessCodeForUnit_SkipsZeroUnitID(t *testing.T) {
	service := createTestService(nil, nil)
	claims := createTestClaims()

	err := service.ProcessAccessCodeForUnit(context.Background(), 0, 123, "testcode", claims)
	if err != nil {
		t.Errorf("Expected no error for unitID 0, got: %v", err)
	}
}

func TestProcessAccessCodeForUnit_ValidationError(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
	}

	// This test demonstrates the validation path
	// Note: With current mock setup, validation always succeeds (IsValid = true)
	// In a real implementation, you'd want to make Validate injectable for better testing
	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	err := service.ProcessAccessCodeForUnit(context.Background(), 123, 456, "testcode", claims)
	if err != nil {
		t.Errorf("Expected no error with current mock setup, got: %v", err)
	}
}

func TestProcessAccessCodeForUnit_DuplicateAccessCode(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{
				{
					AccessCode: "testcode",
					UnitID:     123,
					UserID:     456,
					SiteID:     123,
					State:      AccessCodeStateActive,
				},
			}, nil
		},
	}

	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	// When existing code matches new code and is active, should return nil (no error, early termination)
	err := service.ProcessAccessCodeForUnit(context.Background(), 123, 456, "testcode", claims)
	if err != nil {
		t.Errorf("Expected no error for matching active code, got: %v", err)
	}
}

func TestProcessAccessCodeForUnit_GetCodesForUnitsError(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return nil, errors.New("database error")
		},
	}

	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	err := service.ProcessAccessCodeForUnit(context.Background(), 123, 456, "testcode", claims)
	if err == nil {
		t.Error("Expected error for GetCodesForUnits failure, got nil")
	}
	if err.Error() != "internal server error updating access codes" {
		t.Errorf("Expected 'internal server error updating access codes', got '%s'", err.Error())
	}
}

func TestProcessAccessCodeForUnit_UpdateAccessCodesError(t *testing.T) {
	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{}, nil
		},
		UpdateAccessCodesFunc: func(codes []string, siteID int) error {
			return errors.New("update failed")
		},
	}

	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	err := service.ProcessAccessCodeForUnit(context.Background(), 123, 456, "testcode", claims)
	if err == nil {
		t.Error("Expected error for UpdateAccessCodes failure, got nil")
	}
	if err.Error() != "internal server error updating access codes" {
		t.Errorf("Expected 'internal server error updating access codes', got '%s'", err.Error())
	}
}

func TestProcessAccessCodeForUnit_RevokeAccessCodesError(t *testing.T) {
	mockCommandCenter := &MockCommandCenter{
		RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
			return errors.New("revoke failed")
		},
	}

	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{
				{
					AccessCode: "oldcode",
					UnitID:     123,
					UserID:     456,
					SiteID:     123,
					State:      AccessCodeStateActive,
				},
			}, nil
		},
		NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterInterface {
			return mockCommandCenter
		},
	}

	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	err := service.ProcessAccessCodeForUnit(context.Background(), 123, 456, "newcode", claims)
	if err == nil {
		t.Error("Expected error for RevokeAccessCodes failure, got nil")
	}
	if !strings.Contains(err.Error(), "failed to revoke previous access codes") {
		t.Errorf("Expected error containing 'failed to revoke previous access codes', got '%s'", err.Error())
	}
}

func TestProcessAccessCodeForUnit_SetAccessCodesError(t *testing.T) {
	mockCommandCenter := &MockCommandCenter{
		RevokeAccessCodesFunc: func(revokeUnits []int, options map[string]struct{}) error {
			return nil // Revoke succeeds
		},
		SetAccessCodesFunc: func(units []int, options map[string]struct{}) error {
			return errors.New("set failed")
		},
	}

	mockBank := &MockBank{
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{
				{
					AccessCode: "oldcode",
					UnitID:     123,
					UserID:     456,
					SiteID:     123,
					State:      AccessCodeStateActive,
				},
			}, nil
		},
		NewCommandCenterClientFunc: func(siteID int, ctx context.Context) CommandCenterInterface {
			return mockCommandCenter
		},
	}

	service := createTestService(mockBank, nil)
	claims := createTestClaims()

	err := service.ProcessAccessCodeForUnit(context.Background(), 123, 456, "newcode", claims)
	if err == nil {
		t.Error("Expected error for SetAccessCodes failure, got nil")
	}
	if !strings.Contains(err.Error(), "failed to set access codes") {
		t.Errorf("Expected error containing 'failed to set access codes', got '%s'", err.Error())
	}
}

// Test ProcessAccessCodeEdit function

func TestProcessAccessCodeEdit_ConvertUUIDsError(t *testing.T) {
	service := createTestService(nil, nil)
	claims := createTestClaims()
	input := &AccessCodeEditRequest{
		UserUUID: "invalid-uuid",
	}

	err := service.ProcessAccessCodeEdit(context.Background(), input, claims)
	if err == nil {
		t.Error("Expected error for invalid UUID, got nil")
	}
	if !strings.Contains(err.Error(), "invalid UUID") {
		t.Errorf("Expected error containing 'invalid UUID', got '%s'", err.Error())
	}
}

func TestProcessAccessCodeEdit_ValidateUserAccessError(t *testing.T) {
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return nil, errors.New("no_ob_found")
		},
	}
	service := createTestService(mockBank, nil)
	claims := createTestClaims()
	input := &AccessCodeEditRequest{
		UserID:     123,
		UnitIDs:    []int{456},
		AccessCode: "testcode",
	}

	err := service.ProcessAccessCodeEdit(context.Background(), input, claims)
	if err == nil {
		t.Error("Expected error for user validation failure, got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got '%s'", err.Error())
	}
}

func TestProcessAccessCodeEdit_ActivityRecorderError(t *testing.T) {
	mockActivityRecorder := &MockActivityRecorder{
		UserAccessCodeEditFunc: func(userID int, user *BUser, siteUUID string, timestamp time.Time) (*string, error) {
			return nil, errors.New("activity recording failed")
		},
	}

	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return createTestUser(), nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{}, nil
		},
	}

	service := createTestService(mockBank, mockActivityRecorder)
	claims := createTestClaims()
	input := &AccessCodeEditRequest{
		UserID:     123,
		UnitIDs:    []int{456},
		AccessCode: "testcode",
	}

	err := service.ProcessAccessCodeEdit(context.Background(), input, claims)
	if err == nil {
		t.Error("Expected error for activity recording failure, got nil")
	}
	if !strings.Contains(err.Error(), "recording activity failed") {
		t.Errorf("Expected error containing 'recording activity failed', got '%s'", err.Error())
	}
}

func TestProcessAccessCodeEdit_Success(t *testing.T) {
	testUser := createTestUser()
	mockBank := &MockBank{
		GetBUserByIDFunc: func(BUserID int) (*BUser, error) {
			return testUser, nil
		},
		V2UnitGetByIdFunc: func(unitID int, siteID int) (*Unit, error) {
			return createTestUnit(), nil
		},
		GetCodesForUnitsFunc: func(units []int, siteID int) ([]GateAccessCode, error) {
			return []GateAccessCode{}, nil
		},
	}

	service := createTestService(mockBank, nil)
	claims := createTestClaims()
	input := &AccessCodeEditRequest{
		UserID:     123,
		UnitIDs:    []int{456},
		AccessCode: "testcode",
	}

	err := service.ProcessAccessCodeEdit(context.Background(), input, claims)
	if err != nil {
		t.Errorf("Expected no error for successful processing, got: %v", err)
	}
}

// Test HTTP Handler scenarios

func TestAccessCodeEditHandler_ClaimsNotFound(t *testing.T) {
	// Create a context without claims
	ctx := context.Background()

	// This would test the first conditional branch in the handler
	// In practice, you'd use httptest.NewRecorder() and httptest.NewRequest()
	// to create proper HTTP test scenarios, but this shows the approach
	claims, err := getClaimsFromContext(ctx)
	if err == nil {
		t.Error("Expected error for missing claims, got nil")
	}
	if err.Error() != "claims not found in context" {
		t.Errorf("Expected 'claims not found in context', got '%s'", err.Error())
	}
	if claims != nil {
		t.Error("Expected nil claims, got non-nil")
	}
}

func TestAccessCodeEditHandler_BankNotFound(t *testing.T) {
	// Create a context with claims but without bank
	ctx := context.WithValue(context.Background(), claimsKey, &Claims{})

	bank, ok := BankFromContext(ctx)
	if ok {
		t.Error("Expected false for missing bank, got true")
	}
	if bank != nil {
		t.Error("Expected nil bank, got non-nil")
	}
}
