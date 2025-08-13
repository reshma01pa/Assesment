package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	_ "log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	called := m.Called(query, args)
	return called.Get(0).(sql.Result), called.Error(1)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	called := m.Called(query, args)
	return called.Get(0).(*sql.Rows), called.Error(1)
}

func TestWriteAlertHandler(t *testing.T) {
	// Initialize the SQLite database
	initializeDatabase()

	// Create a test request body
	requestBody := `{
        "alert_id": "test_alert_id",
        "service_id": "test_service_id",
        "service_name": "test_service_name",
        "model": "test_model",
        "alert_type": "test_type",
        "alert_ts": "123456",
        "severity": "test_severity",
        "team_slack": "test_slack"
    }`

	// Create a test request
	req, err := http.NewRequest("POST", "/alerts", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the WriteAlert function with the test request
	WriteAlert(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}
}

// Similarly, you can create tests for ReadAlertsHandler using the MockDB.
func TestReadAlertsHandler(t *testing.T) {
	// Initialize the SQLite database
	initializeDatabase()

	// Insert some test data into the database
	testAlert := Alert{
		AlertID:     "test_alert_id",
		ServiceID:   "test_service_id",
		ServiceName: "test_service_name",
		Model:       "test_model",
		AlertType:   "test_type",
		AlertTS:     "123456",
		Severity:    "test_severity",
		TeamSlack:   "test_slack",
	}
	storeAlertInDatabase(testAlert)

	// Create a test request
	req, err := http.NewRequest("GET", "/alerts?service_id=test_service_id&start_ts=0&end_ts=999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the ReadAlerts function with the test request
	ReadAlerts(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	// Define the updated expected JSON response to match the actual structure
	expectedResponse := map[string]interface{}{
		"service_id":   "test_service_id",
		"service_name": "test_service_name",
		"alerts": []map[string]string{
			{
				"alert_id":   "test_alert_id",
				"model":      "test_model",
				"alert_type": "test_type",
				"alert_ts":   "123456",
				"severity":   "test_severity",
				"team_slack": "test_slack",
			},
		},
	}

	// Parse the response body and compare it with the expected response
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}
	actualJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to encode actual response to JSON: %v", err)
	}

	expectedJSON, err := json.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Failed to encode expected response to JSON: %v", err)
	}
	if string(actualJSON) != string(expectedJSON) {
		t.Errorf("Response does not match the expected response.")
	}
}

func TestArianaPageHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/ariana", nil)
	rr := httptest.NewRecorder()

	ArianaPageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
	if ctype := rr.Header().Get("Content-Type"); ctype != "text/html; charset=utf-8" {
		t.Fatalf("unexpected content type: %s", ctype)
	}
	if body := rr.Body.String(); len(body) == 0 || !contains(body, "Ariana Grande") {
		t.Fatalf("unexpected body: %q", body)
	}
}

// contains is a tiny helper to avoid importing strings just for one assertion
func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
