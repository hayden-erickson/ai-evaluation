package commandcenter

import (
	"context"
	"testing"
)

func TestFactory_NewClient(t *testing.T) {
	factory := NewFactory()
	client := factory.NewClient(123)
	
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
}

func TestFactory_NewClientWithContext(t *testing.T) {
	factory := NewFactory()
	ctx := context.Background()
	client := factory.(*Factory).NewClientWithContext(123, ctx)
	
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
}

func TestClient_RevokeAccessCodes(t *testing.T) {
	client := &Client{siteID: 123, ctx: context.Background()}
	
	err := client.RevokeAccessCodes([]int{1, 2, 3}, make(map[string]struct{}))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_SetAccessCodes(t *testing.T) {
	client := &Client{siteID: 123, ctx: context.Background()}
	
	err := client.SetAccessCodes([]int{1, 2, 3}, make(map[string]struct{}))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClient_RevokeAccessCodes_EmptyUnits(t *testing.T) {
	client := &Client{siteID: 123, ctx: context.Background()}
	
	err := client.RevokeAccessCodes([]int{}, make(map[string]struct{}))
	if err == nil {
		t.Error("Expected error for empty units, got nil")
	}
}

func TestClient_SetAccessCodes_EmptyUnits(t *testing.T) {
	client := &Client{siteID: 123, ctx: context.Background()}
	
	err := client.SetAccessCodes([]int{}, make(map[string]struct{}))
	if err == nil {
		t.Error("Expected error for empty units, got nil")
	}
}
