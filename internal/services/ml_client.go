package services

import (
	"bytes"
	"car-status-backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

type MLClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

func NewMLClient(baseURL string, timeout time.Duration, apiKey string) *MLClient {
	return &MLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		apiKey: apiKey,
	}
}

func (c *MLClient) PredictCarStatus(imagePath string) (*models.MLPredictionResponse, error) {
	absolutePath, err := filepath.Abs(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	payload := models.MLPredictionRequest{
		ImagePath: absolutePath,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/predict", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var result models.MLPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		result.Success = false
		if result.Error == "" {
			result.Error = fmt.Sprintf("ML service returned status %d", resp.StatusCode)
		}
		return &result, nil
	}

	result.Success = true
	return &result, nil
}

func (c *MLClient) HealthCheck() error {
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform health check: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ML service health check failed with status %d", resp.StatusCode)
	}

	return nil
}