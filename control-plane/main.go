package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	activeAgents = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "stuntdouble_active_agents_total",
		Help: "The total number of active StuntDouble agents globally",
	})
	blockedRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "stuntdouble_blocked_network_requests_total",
		Help: "The total number of network requests blocked by StuntDouble eBPF",
	})
)

func init() {
	prometheus.MustRegister(activeAgents)
	prometheus.MustRegister(blockedRequests)
}

// EnterpriseRBACPolicy represents global sandbox policies set by a CTO.
type EnterpriseRBACPolicy struct {
	OrgID         string   `json:"org_id"`
	BlockedPorts  []int    `json:"blocked_ports"`
	AllowedAgents []string `json:"allowed_agents"`
	StrictEgress  bool     `json:"strict_egress"`
}

type TelemetryData struct {
	TotalRuns       int       `json:"total_runs"`
	BlockedCommands int       `json:"blocked_commands"`
	LastRun         time.Time `json:"last_run"`
}

var (
	mu            sync.Mutex
	globalMetrics TelemetryData
)

func handleTelemetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data TelemetryData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	globalMetrics.TotalRuns += data.TotalRuns
	globalMetrics.BlockedCommands += data.BlockedCommands
	globalMetrics.LastRun = data.LastRun
	
	// Update Prometheus metrics
	activeAgents.Set(float64(globalMetrics.TotalRuns)) // Simplified mapping
	blockedRequests.Add(float64(data.BlockedCommands))
	mu.Unlock()

	log.Printf("Received telemetry: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handlePolicy(w http.ResponseWriter, r *http.Request) {
	policy := EnterpriseRBACPolicy{
		OrgID:         "ent_global",
		BlockedPorts:  []int{5432, 27017, 3306, 6379},
		AllowedAgents: []string{"claude", "cursor", "opendevin"},
		StrictEgress:  true,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

func main() {
	http.HandleFunc("/telemetry", handleTelemetry)
	http.HandleFunc("/policy", handlePolicy)
	http.Handle("/metrics", promhttp.Handler())
	
	port := "4439"
	fmt.Printf("🏢 StuntDouble Enterprise Control Plane active on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
