package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type TelemetryStats struct {
	TotalRuns       int       `json:"total_runs"`
	BlockedCommands int       `json:"blocked_commands"`
	LastRun         time.Time `json:"last_run"`
	Status          string    `json:"status"`
}

// StartServer boots up the local StuntDouble Control Plane API
func StartServer(port string) error {
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/health", handleHealth)

	// Enable CORS for the local dashboard
	handler := corsMiddleware(http.DefaultServeMux)

	fmt.Printf("🚀 [StuntDouble] Control Plane API live on http://localhost:%s\n", port)
	return http.ListenAndServe(":"+port, handler)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "StuntDouble Engine Online", "version": "3.0.0"})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Read telemetry from the local .stuntdouble.telemetry.json
	file := ".stuntdouble.telemetry.json"
	stats := TelemetryStats{
		Status: "Secure",
	}

	if data, err := os.ReadFile(file); err == nil {
		json.Unmarshal(data, &stats)
	}

	json.NewEncoder(w).Encode(stats)
}

// corsMiddleware allows the Next.js dashboard to securely hit this local API
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
