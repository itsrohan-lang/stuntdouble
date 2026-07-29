// Package api serves the local StuntDouble telemetry API consumed by the
// dashboard. It reads the run counter written by `stuntdouble run` and reports
// it; it does not enforce anything.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	envCORSOrigin     = "STUNTDOUBLE_CORS_ORIGIN"
	defaultCORSOrigin = "http://localhost:3000"
)

// TelemetryStats is the local run counter.
//
// There is no blocked-connection count: kernel-level egress filtering is
// unimplemented (see pkg/ebpf), so nothing blocks connections to count.
// EgressEnforcement says so explicitly rather than leaving callers to assume.
type TelemetryStats struct {
	TotalRuns         int       `json:"total_runs"`
	LastRun           time.Time `json:"last_run"`
	EgressEnforcement string    `json:"egress_enforcement"`
}

// StartServer boots the local telemetry API. It binds to loopback only: the
// server is unauthenticated, so it must not be reachable off-host.
func StartServer(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/stats", handleStats)
	mux.HandleFunc("/api/health", handleHealth)

	addr := "127.0.0.1:" + port
	srv := &http.Server{
		Addr:              addr,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Printf("[StuntDouble] Local telemetry API on http://%s\n", addr)
	return srv.ListenAndServe()
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":             "ok",
		"egress_enforcement": "unimplemented",
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := TelemetryStats{EgressEnforcement: "unimplemented"}
	if data, err := os.ReadFile(".stuntdouble.telemetry.json"); err == nil {
		json.Unmarshal(data, &stats)
		// The on-disk file never contains this field; keep the server's value.
		stats.EgressEnforcement = "unimplemented"
	}

	json.NewEncoder(w).Encode(stats)
}

// corsMiddleware permits one configured browser origin. A wildcard would let
// any page the user visits read this API.
func corsMiddleware(next http.Handler) http.Handler {
	allowed := os.Getenv(envCORSOrigin)
	if allowed == "" {
		allowed = defaultCORSOrigin
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		if origin := r.Header.Get("Origin"); origin != "" && origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", allowed)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
