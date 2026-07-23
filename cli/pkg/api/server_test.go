package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleHealth)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if response["status"] != "StuntDouble Engine Online" {
		t.Errorf("handler returned unexpected body: got %v", response["status"])
	}
}

func TestHandleStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var stats TelemetryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	// Because we might not have a .stuntdouble.telemetry.json in the test directory,
	// it should default to Status: "Secure"
	if stats.Status != "Secure" {
		t.Errorf("handler returned unexpected status: got %v want Secure", stats.Status)
	}
}
