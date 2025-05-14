package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AIClient represents a client for interacting with the AI service
type AIClient struct {
	baseURL    string
	httpClient *http.Client
}

// AnalysisResult represents the response from the AI service
type AnalysisResult struct {
	Category    string  `json:"category"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// NewAIClient creates a new AI client with the given base URL
func NewAIClient(baseURL string) *AIClient {
	return &AIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// AnalyzeTransaction sends a transaction description to the AI service for analysis
func (c *AIClient) AnalyzeTransaction(description string) (*AnalysisResult, error) {
	// Prepare request body
	reqBody := map[string]string{
		"description": description,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/analyze", c.baseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error   string `json:"error"`
			Details string `json:"details,omitempty"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %w", err)
		}
		return nil, fmt.Errorf("AI service error: %s", errorResp.Error)
	}

	// Decode response
	var result AnalysisResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// HealthCheck checks if the AI service is healthy
func (c *AIClient) HealthCheck() error {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/health", c.baseURL))
	if err != nil {
		return fmt.Errorf("error checking health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service is not healthy: status code %d", resp.StatusCode)
	}

	return nil
}
