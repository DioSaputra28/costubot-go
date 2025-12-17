package test

import (
	"bytes"
	"contact-management/src/apps"
	"contact-management/src/config"
	"contact-management/src/utils"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

// TestResponse represents a generic API response
type TestResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   interface{}     `json:"error,omitempty"` // Can be string or object
}

// makeRequest is a helper function to make HTTP requests for testing
func makeRequest(t *testing.T, router *httprouter.Router, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

// parseResponse parses the response body into TestResponse
func parseResponse(t *testing.T, rr *httptest.ResponseRecorder) TestResponse {
	t.Helper()

	var response TestResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v, body: %s", err, rr.Body.String())
	}

	return response
}

// getValidToken creates a test user and returns a valid JWT token
func getValidToken(t *testing.T, username string) string {
	t.Helper()

	token, err := utils.GenerateToken(username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	return token
}

// cleanupTestUser removes test user from database and Redis
func cleanupTestUser(t *testing.T, username string) {
	t.Helper()

	cfg := config.LoadConfig()
	db, err := apps.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Delete user from database
	_, err = db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		t.Logf("Warning: Failed to cleanup user %s: %v", username, err)
	}

	// Clear token from Redis (no cleanup needed - handled by expiration)
}

// cleanupTestCategory removes test category from database
func cleanupTestCategory(t *testing.T, categoryID int) {
	t.Helper()

	cfg := config.LoadConfig()
	db, err := apps.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM category WHERE category_id = ?", categoryID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup category %d: %v", categoryID, err)
	}
}

// cleanupTestBrandProduct removes test brand product from database
func cleanupTestBrandProduct(t *testing.T, brandProductID int) {
	t.Helper()

	cfg := config.LoadConfig()
	db, err := apps.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM brand_products WHERE brand_product_id = ?", brandProductID)
	if err != nil {
		t.Logf("Warning: Failed to cleanup brand product %d: %v", brandProductID, err)
	}
}

// assertStatusCode checks if the response status code matches expected
func assertStatusCode(t *testing.T, expected, got int) {
	t.Helper()
	if expected != got {
		t.Errorf("Expected status code %d, got %d", expected, got)
	}
}

// assertResponseStatus checks if the response status matches expected
func assertResponseStatus(t *testing.T, expected string, response TestResponse) {
	t.Helper()
	if response.Status != expected {
		t.Errorf("Expected status %s, got %s", expected, response.Status)
	}
}
