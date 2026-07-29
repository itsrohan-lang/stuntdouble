package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type CLIStats struct {
	TotalRuns int `json:"total_runs"`
}

// CLIPolicy is the policy document served by the control plane. The fields
// describe intent: nothing enforces them while egress filtering is
// unimplemented.
type CLIPolicy struct {
	OrgID         string   `json:"org_id"`
	BlockedPorts  []int    `json:"blocked_ports"`
	AllowedAgents []string `json:"allowed_agents"`
	StrictEgress  bool     `json:"strict_egress"`
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Live terminal view of the StuntDouble control plane",
	Long: `Polls the control plane and displays reported run telemetry and the active
policy document.

Requires STUNTDOUBLE_TOKEN to be set to the control plane's bearer token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := os.Getenv("STUNTDOUBLE_TOKEN")
		if token == "" {
			return fmt.Errorf("STUNTDOUBLE_TOKEN is not set; " +
				"export the same token the control plane was started with")
		}

		baseURL := os.Getenv("STUNTDOUBLE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:4439"
		}

		client := &http.Client{Timeout: 5 * time.Second}
		get := func(path string, out any) error {
			req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet, baseURL+path, nil)
			if err != nil {
				return err
			}
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("%s returned %s", path, resp.Status)
			}
			return json.NewDecoder(resp.Body).Decode(out)
		}

		fmt.Print("\033[2J")

		for {
			fmt.Print("\033[H\033[J")
			fmt.Println("STUNTDOUBLE MONITOR")
			fmt.Println("==================================================")
			fmt.Printf("Control plane: %s\n", baseURL)
			fmt.Printf("Updated:       %s\n\n", time.Now().Format(time.RFC1123))

			var stats CLIStats
			if err := get("/api/stats", &stats); err != nil {
				fmt.Printf("Control plane unreachable: %v\n", err)
				fmt.Println("\nRetrying in 2s. Press Ctrl+C to exit.")
				time.Sleep(2 * time.Second)
				continue
			}

			fmt.Println("REPORTED TELEMETRY")
			fmt.Printf("   Total agent runs:  %d\n\n", stats.TotalRuns)

			var policy CLIPolicy
			if err := get("/policy", &policy); err != nil {
				fmt.Printf("Could not read policy: %v\n\n", err)
			} else {
				fmt.Println("ACTIVE POLICY DOCUMENT (advisory — not enforced)")
				fmt.Printf("   Organization ID:    %s\n", policy.OrgID)
				fmt.Printf("   Strict egress:      %t\n", policy.StrictEgress)
				fmt.Printf("   Allowed agents:     %v\n", policy.AllowedAgents)
				fmt.Printf("   Blocked ports:      %v\n\n", policy.BlockedPorts)
			}

			fmt.Println("==================================================")
			fmt.Println("Egress filtering is not implemented; the policy above")
			fmt.Println("is distributed but not applied. Press Ctrl+C to exit.")

			select {
			case <-cmd.Context().Done():
				return nil
			case <-time.After(2 * time.Second):
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
