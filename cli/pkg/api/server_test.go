package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handleHealth).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if response["status"] != "ok" {
		t.Errorf("status = %q, want \"ok\"", response["status"])
	}
}

// The API must not imply that egress filtering is active, whether or not a
// telemetry file exists in the working directory.
func TestHandleStatsReportsEnforcementHonestly(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handleStats).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var stats TelemetryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if stats.EgressEnforcement != "unimplemented" {
		t.Errorf("egress_enforcement = %q, want \"unimplemented\"", stats.EgressEnforcement)
	}
}

// A wildcard origin would let any page in the user's browser read local
// telemetry, so only the configured origin may be reflected.
func TestCORSOnlyAllowsConfiguredOrigin(t *testing.T) {
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		origin string
		want   string
	}{
		{defaultCORSOrigin, defaultCORSOrigin},
		{"http://evil.example", ""},
		{"", ""},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if got := rr.Header().Get("Access-Control-Allow-Origin"); got != tc.want {
			t.Errorf("origin %q: Allow-Origin = %q, want %q", tc.origin, got, tc.want)
		}
	}
}
