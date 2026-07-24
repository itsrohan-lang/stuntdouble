package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type CLIStats struct {
	TotalRuns       int `json:"total_runs"`
	BlockedCommands int `json:"blocked_commands"`
}

type CLIPolicy struct {
	OrgID         string   `json:"org_id"`
	BlockedPorts  []int    `json:"blocked_ports"`
	AllowedAgents []string `json:"allowed_agents"`
	StrictEgress  bool     `json:"strict_egress"`
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Live terminal view of the StuntDouble Control Plane",
	Long:  "Polls the central enterprise Control Plane and displays a terminal-based dashboard of live security telemetry and active Zero-Trust policies.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[2J") // Clear screen
		
		for {
			fmt.Print("\033[H") // Move cursor to top-left
			fmt.Println("🛡️  STUNTDOUBLE TERMINAL SOC (Live Monitoring)")
			fmt.Println("==================================================")
			fmt.Printf("Last Updated: %s\n\n", time.Now().Format(time.RFC1123))

			// Fetch Live Stats
			statsResp, err := http.Get("http://localhost:4439/api/stats")
			if err != nil {
				fmt.Println("❌  Control Plane Status: OFFLINE (Is it running on port 4439?)")
				fmt.Println("Waiting for connection...")
				time.Sleep(2 * time.Second)
				continue
			}
			defer statsResp.Body.Close()

			var stats CLIStats
			json.NewDecoder(statsResp.Body).Decode(&stats)

			// Fetch Live Policy
			policyResp, err := http.Get("http://localhost:4439/policy")
			var policy CLIPolicy
			if err == nil {
				json.NewDecoder(policyResp.Body).Decode(&policy)
				policyResp.Body.Close()
			}

			fmt.Println("📊 GLOBAL TELEMETRY")
			fmt.Printf("   Total Agent Runs:    %d\n", stats.TotalRuns)
			fmt.Printf("   Blocked Outbound:    \033[31m%d\033[0m (Critical Overrides)\n\n", stats.BlockedCommands)

			fmt.Println("📜 ACTIVE ZERO-TRUST POLICY")
			fmt.Printf("   Organization ID:     %s\n", policy.OrgID)
			fmt.Printf("   Enforcement Mode:    %t (Strict Egress)\n", policy.StrictEgress)
			fmt.Printf("   Whitelisted Agents:  %v\n", policy.AllowedAgents)
			fmt.Printf("   Blocked Ports:       %v\n\n", policy.BlockedPorts)

			fmt.Println("==================================================")
			fmt.Println("Press Ctrl+C to exit monitor mode.")
			
			time.Sleep(2 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
