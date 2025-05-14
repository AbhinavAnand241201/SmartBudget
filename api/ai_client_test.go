package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAIClient_AnalyzeTransaction(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/analyze" {
			t.Errorf("Expected /analyze path, got %s", r.URL.Path)
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}

		// Check request body
		if reqBody["description"] != "Test transaction" {
			t.Errorf("Expected description 'Test transaction', got %s", reqBody["description"])
		}

		// Send response
		response := AnalysisResult{
			Category:    "groceries",
			Confidence:  0.95,
			Description: "Test transaction",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client
	client := NewAIClient(server.URL)

	// Test successful analysis
	result, err := client.AnalyzeTransaction("Test transaction")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Category != "groceries" {
		t.Errorf("Expected category 'groceries', got %s", result.Category)
	}
	if result.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", result.Confidence)
	}
	if result.Description != "Test transaction" {
		t.Errorf("Expected description 'Test transaction', got %s", result.Description)
	}

	// Test error response
	server.Close()
	_, err = client.AnalyzeTransaction("Test transaction")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAIClient_HealthCheck(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/health" {
			t.Errorf("Expected /health path, got %s", r.URL.Path)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}))
	defer server.Close()

	// Create client
	client := NewAIClient(server.URL)

	// Test successful health check
	err := client.HealthCheck()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test error response
	server.Close()
	err = client.HealthCheck()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
