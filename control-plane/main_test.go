package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testToken = "test-token"

func TestHandleTelemetry(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/telemetry",
		bytes.NewBufferString(`{"total_runs": 1}`))

	rr := httptest.NewRecorder()
	http.HandlerFunc(handleTelemetry).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestHandlePolicy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/policy", nil)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handlePolicy).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var got Policy
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if got.OrgID == "" {
		t.Error("policy response has empty org_id")
	}
}

// The control plane exposes the audit log and accepts policy writes, so every
// endpoint behind requireAuth must reject unauthenticated callers.
func TestRequireAuthRejectsMissingToken(t *testing.T) {
	authToken = testToken

	tests := []struct {
		name   string
		header string
		want   int
	}{
		{"no header", "", http.StatusUnauthorized},
		{"wrong scheme", "Basic " + testToken, http.StatusUnauthorized},
		{"wrong token", "Bearer nope", http.StatusUnauthorized},
		{"bare prefix", "Bearer ", http.StatusUnauthorized},
		{"valid token", "Bearer " + testToken, http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			rr := httptest.NewRecorder()
			requireAuth(handleStats).ServeHTTP(rr, req)

			if rr.Code != tc.want {
				t.Errorf("status = %d, want %d", rr.Code, tc.want)
			}
		})
	}
}

// A wildcard CORS origin would let any page in the user's browser read the
// audit log, so only the configured origin may be reflected.
func TestWithCORSOnlyAllowsConfiguredOrigin(t *testing.T) {
	corsOrigin = "http://localhost:3000"
	noop := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }

	tests := []struct {
		origin string
		want   string
	}{
		{"http://localhost:3000", "http://localhost:3000"},
		{"http://evil.example", ""},
		{"", ""},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}

		rr := httptest.NewRecorder()
		withCORS(noop).ServeHTTP(rr, req)

		if got := rr.Header().Get("Access-Control-Allow-Origin"); got != tc.want {
			t.Errorf("origin %q: Allow-Origin = %q, want %q", tc.origin, got, tc.want)
		}
	}
}

func TestHandleHealthReportsEnforcementHonestly(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)

	rr := httptest.NewRecorder()
	http.HandlerFunc(handleHealth).ServeHTTP(rr, req)

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if body["egress_enforcement"] != "unimplemented" {
		t.Errorf("egress_enforcement = %v, want \"unimplemented\"", body["egress_enforcement"])
	}
}
