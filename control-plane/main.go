package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// EnterpriseRBACPolicy represents global sandbox policies set by a CTO.
type EnterpriseRBACPolicy struct {
	OrgID         string   `json:"org_id"`
	BlockedPorts  []int    `json:"blocked_ports"`
	AllowedAgents []string `json:"allowed_agents"`
	StrictEgress  bool     `json:"strict_egress"`
}

func handleTelemetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	fmt.Println("☁️  [Control Plane] Received Telemetry payload from local StuntDouble CLI")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"synced"}`))
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
	
	port := "4439"
	fmt.Printf("🏢 StuntDouble Enterprise Control Plane active on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
