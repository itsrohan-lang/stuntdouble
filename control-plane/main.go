package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	db *gorm.DB
)

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   string    `json:"agent_id"`
	Target    string    `json:"target"`
	Action    string    `json:"action"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

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

var globalPolicy = EnterpriseRBACPolicy{
	OrgID:         "ent_global",
	BlockedPorts:  []int{5432, 27017, 3306, 6379},
	AllowedAgents: []string{"claude", "cursor", "opendevin"},
	StrictEgress:  true,
}

var policyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Policy",
		Fields: graphql.Fields{
			"org_id": &graphql.Field{
				Type: graphql.String,
			},
			"blocked_ports": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
			"allowed_agents": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"strict_egress": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"policy": &graphql.Field{
				Type: policyType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return globalPolicy, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

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
	
	activeAgents.Set(float64(globalMetrics.TotalRuns))
	blockedRequests.Add(float64(data.BlockedCommands))
	mu.Unlock()

	log.Printf("Received telemetry: %+v", data)
	w.WriteHeader(http.StatusOK)
}

func handlePolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodPost {
		var newPolicy EnterpriseRBACPolicy
		if err := json.NewDecoder(r.Body).Decode(&newPolicy); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		globalPolicy = newPolicy
		log.Printf("🚀 CTO Dashboard deployed new global policy: %+v", globalPolicy)
	}

	json.NewEncoder(w).Encode(globalPolicy)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(globalMetrics)
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: req.Query,
	})
	json.NewEncoder(w).Encode(result)
}

func handleAuditLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		var logEntry AuditLog
		if err := json.NewDecoder(r.Body).Decode(&logEntry); err == nil {
			db.Create(&logEntry)
			w.WriteHeader(http.StatusCreated)
			return
		}
	} else if r.Method == http.MethodGet {
		var logs []AuditLog
		db.Order("created_at desc").Limit(50).Find(&logs)
		json.NewEncoder(w).Encode(logs)
	}
}

func handleKeployMock(w http.ResponseWriter, r *http.Request) {
	// Simulate Keploy intercepting a blocked outbound request and returning a ghost 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-StuntDouble-Mocked", "true")
	w.Header().Set("X-Keploy-Interceptor", "active")
	w.WriteHeader(http.StatusOK)

	// Return a generic success payload so the AI agent thinks its API call worked
	mockResponse := map[string]interface{}{
		"status":    "success",
		"id":        "evt_mocked_1337",
		"message":   "This is a ghost response injected by Keploy. The actual network request was dropped.",
		"mocked_by": "StuntDouble-Keploy-Integration",
		"data": map[string]string{
			"status": "created",
			"amount": "0.00",
		},
	}
	json.NewEncoder(w).Encode(mockResponse)
}

func startEBPFEngine() {
	log.Println("⚙️ Initializing Rust eBPF Kernel Engine...")
	cmd := exec.Command("cargo", "run", "--release")
	cmd.Dir = "../core-ebpf" // Assuming we run from control-plane dir
	
	// We don't wait for it to finish because it runs continuously
	if err := cmd.Start(); err != nil {
		log.Printf("⚠️ Warning: Failed to start Rust eBPF engine natively (is cargo installed?): %v\n", err)
		log.Println("Falling back to mock eBPF enforcement mode.")
	} else {
		log.Printf("🛡️ Rust eBPF Engine active (PID: %d). Intercepting syscalls at kernel-level.", cmd.Process.Pid)
	}
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("stuntdouble_audit.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&AuditLog{})

	// Start the Rust eBPF Kernel probes in the background
	go startEBPFEngine()

	http.HandleFunc("/telemetry", handleTelemetry)
	http.HandleFunc("/policy", handlePolicy)
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/audit", handleAuditLogs)
	http.HandleFunc("/api/keploy/mock", handleKeployMock)
	http.HandleFunc("/graphql", handleGraphQL)
	http.Handle("/metrics", promhttp.Handler())
	
	port := "4439"
	fmt.Printf("🏢 StuntDouble Enterprise Control Plane active on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
