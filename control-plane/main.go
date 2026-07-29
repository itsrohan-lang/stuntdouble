// Command control-plane runs the StuntDouble control plane: it collects run
// telemetry from CLI instances, stores an audit log, and serves the policy
// document consumed by the dashboard.
//
// Scope note: the control plane distributes policy and records what agents
// report. It does not itself enforce anything — kernel-level egress filtering
// is unimplemented (see cli/pkg/ebpf and docs/ENFORCEMENT.md). Treat the audit
// log as self-reported by the CLI, not as a tamper-proof ledger.
package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	envToken          = "STUNTDOUBLE_TOKEN"
	envListenAddr     = "STUNTDOUBLE_LISTEN_ADDR"
	envCORSOrigin     = "STUNTDOUBLE_CORS_ORIGIN"
	defaultListen     = "127.0.0.1:4439"
	defaultCORSOrigin = "http://localhost:3000"
)

var (
	agentRuns = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "stuntdouble_agent_runs_total",
		Help: "Total agent runs reported by StuntDouble CLI instances",
	})
	db *gorm.DB

	// authToken gates every endpoint except /api/health. Set via STUNTDOUBLE_TOKEN.
	authToken string
	// corsOrigin is the single browser origin permitted to call this API.
	corsOrigin string
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
	prometheus.MustRegister(agentRuns)
}

// Policy is the sandbox policy distributed to CLI instances.
//
// BlockedPorts and StrictEgress describe intent only. Nothing enforces them
// today; they are advisory until the egress filter in cli/pkg/ebpf is
// implemented.
type Policy struct {
	OrgID         string   `json:"org_id"`
	BlockedPorts  []int    `json:"blocked_ports"`
	AllowedAgents []string `json:"allowed_agents"`
	StrictEgress  bool     `json:"strict_egress"`
}

var globalPolicy = Policy{
	OrgID:         "default",
	BlockedPorts:  []int{5432, 27017, 3306, 6379},
	AllowedAgents: []string{"claude", "cursor", "opendevin"},
	StrictEgress:  true,
}

var policyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Policy",
		Fields: graphql.Fields{
			"org_id":         &graphql.Field{Type: graphql.String},
			"blocked_ports":  &graphql.Field{Type: graphql.NewList(graphql.Int)},
			"allowed_agents": &graphql.Field{Type: graphql.NewList(graphql.String)},
			"strict_egress":  &graphql.Field{Type: graphql.Boolean},
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

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{Query: queryType})

// TelemetryData is the run counter reported by each CLI instance. There is no
// blocked-connection count: nothing blocks connections yet.
type TelemetryData struct {
	TotalRuns int       `json:"total_runs"`
	LastRun   time.Time `json:"last_run"`
}

var (
	mu            sync.Mutex
	globalMetrics TelemetryData
)

// requireAuth rejects requests without a matching bearer token. Comparison is
// constant-time so the handler does not leak the token through timing.
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}
		const prefix = "Bearer "
		header := r.Header.Get("Authorization")
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix ||
			subtle.ConstantTimeCompare([]byte(header[len(prefix):]), []byte(authToken)) != 1 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="stuntdouble"`)
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// withCORS permits exactly one configured browser origin. It deliberately does
// not echo arbitrary origins back: with credentials in play, a wildcard would
// let any page read this API's responses.
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		if origin := r.Header.Get("Origin"); origin != "" && origin == corsOrigin {
			w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

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
	globalMetrics.LastRun = data.LastRun
	agentRuns.Set(float64(globalMetrics.TotalRuns))
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func handlePolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodPost {
		var newPolicy Policy
		if err := json.NewDecoder(r.Body).Decode(&newPolicy); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		globalPolicy = newPolicy
		log.Printf("policy updated: %+v", globalPolicy)
	}

	json.NewEncoder(w).Encode(globalPolicy)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(globalMetrics)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":             "ok",
		"egress_enforcement": "unimplemented",
	})
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
	result := graphql.Do(graphql.Params{Schema: schema, RequestString: req.Query})
	json.NewEncoder(w).Encode(result)
}

func handleAuditLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var logEntry AuditLog
		if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := db.Create(&logEntry).Error; err != nil {
			http.Error(w, "Failed to record audit log", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		var logs []AuditLog
		db.Order("created_at desc").Limit(50).Find(&logs)
		json.NewEncoder(w).Encode(logs)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleKeployMock returns a synthetic success payload for agents whose
// outbound call was intercepted by a Keploy sidecar.
//
// This endpoint only produces the canned response. It does not itself
// intercept or drop anything — the caller has to have been routed here.
func handleKeployMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-StuntDouble-Mocked", "true")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"id":        "evt_mocked_1337",
		"message":   "Synthetic response from the StuntDouble mock endpoint. No upstream request was made.",
		"mocked_by": "StuntDouble-Keploy-Integration",
		"data": map[string]string{
			"status": "created",
			"amount": "0.00",
		},
	})
}

func main() {
	authToken = os.Getenv(envToken)
	if authToken == "" {
		log.Fatalf("%s is not set.\n\n"+
			"The control plane serves the audit log and accepts policy writes, so it\n"+
			"refuses to start without a token. Generate one and export it:\n\n"+
			"  export %s=$(openssl rand -hex 32)\n",
			envToken, envToken)
	}

	corsOrigin = os.Getenv(envCORSOrigin)
	if corsOrigin == "" {
		corsOrigin = defaultCORSOrigin
	}
	listenAddr := os.Getenv(envListenAddr)
	if listenAddr == "" {
		listenAddr = defaultListen
	}

	var err error
	db, err = gorm.Open(sqlite.Open("stuntdouble_audit.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.AutoMigrate(&AuditLog{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", withCORS(handleHealth))
	mux.HandleFunc("/telemetry", withCORS(requireAuth(handleTelemetry)))
	mux.HandleFunc("/policy", withCORS(requireAuth(handlePolicy)))
	mux.HandleFunc("/api/stats", withCORS(requireAuth(handleStats)))
	mux.HandleFunc("/api/audit", withCORS(requireAuth(handleAuditLogs)))
	mux.HandleFunc("/api/keploy/mock", withCORS(requireAuth(handleKeployMock)))
	mux.HandleFunc("/graphql", withCORS(requireAuth(handleGraphQL)))
	mux.Handle("/metrics", requireAuth(promhttp.Handler().ServeHTTP))

	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Printf("StuntDouble control plane listening on http://%s\n", listenAddr)
	fmt.Printf("CORS origin: %s\n", corsOrigin)
	fmt.Println("Note: kernel-level egress enforcement is not implemented; policy is advisory.")
	log.Fatal(srv.ListenAndServe())
}
