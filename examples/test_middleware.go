package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	appContext "github.com/hayden-erickson/ai-evaluation/context"
	"github.com/hayden-erickson/ai-evaluation/middleware"
	"github.com/hayden-erickson/ai-evaluation/repository"
)

// testHandler demonstrates accessing the bank from context
func testHandler(w http.ResponseWriter, r *http.Request) {
	// Try to get bank from context
	bank, ok := appContext.BankFromContext(r.Context())
	if !ok {
		http.Error(w, "Bank not found in context", http.StatusInternalServerError)
		return
	}

	// Use the bank (this would normally do actual work)
	fmt.Fprintf(w, "Successfully accessed bank from context: %T", bank)
}

// TestBankMiddleware demonstrates that the middleware correctly adds bank to context
func TestBankMiddleware() {
	// Create a bank instance with test DSN
	// In a real test, you'd use a test database or mock
	testDSN := "test:test@tcp(localhost:3306)/test_db"
	bank, err := repository.NewBank(testDSN)
	if err != nil {
		fmt.Printf("❌ Test failed: could not create bank instance: %v\n", err)
		fmt.Println("Note: This test requires a MySQL database connection.")
		fmt.Println("For testing without a database, consider using a mock or interface.")
		return
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Wrap our test handler with the bank middleware
	handler := middleware.BankHandlerFunc(bank, testHandler)

	// Execute the request
	handler(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		fmt.Printf("❌ Test failed: expected status 200, got %d\n", w.Code)
		fmt.Printf("Response: %s\n", w.Body.String())
		return
	}

	fmt.Printf("✅ Test passed: Bank middleware successfully added bank to context\n")
	fmt.Printf("Response: %s\n", w.Body.String())
}
