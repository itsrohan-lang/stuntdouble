package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type TelemetryData struct {
	TotalRuns       int       `json:"total_runs"`
	BlockedCommands int       `json:"blocked_commands"`
	LastRun         time.Time `json:"last_run"`
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays a telemetry dashboard of agent actions and blocked destructive commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 StuntDouble Telemetry Dashboard")
		fmt.Println("=====================================")

		telemetryFile := ".stuntdouble.telemetry.json"
		data, err := os.ReadFile(telemetryFile)

		if err != nil {
			// If file doesn't exist, just show zeros
			fmt.Println("Total Agent Sessions:  0")
			fmt.Println("Destructive Commands Blocked: 0")
			fmt.Println("Status: All clear. No rogue actions detected.")
			return
		}

		var stats TelemetryData
		if err := json.Unmarshal(data, &stats); err != nil {
			fmt.Println("Error parsing telemetry data")
			return
		}

		fmt.Printf("Total Agent Sessions:  %d\n", stats.TotalRuns)
		fmt.Printf("Destructive Commands Blocked: %d\n", stats.BlockedCommands)
		fmt.Printf("Last Session Time: %s\n", stats.LastRun.Format(time.RFC822))
		fmt.Println("=====================================")

		if stats.BlockedCommands > 0 {
			fmt.Println("⚠️  Your agents attempted destructive host actions that were safely intercepted.")
		} else {
			fmt.Println("✅  Your agents have behaved safely.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
