package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// TelemetryData mirrors .stuntdouble.telemetry.json, written by `stuntdouble run`.
type TelemetryData struct {
	TotalRuns int       `json:"total_runs"`
	LastRun   time.Time `json:"last_run"`
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Shows local agent run counts",
	Long: `Shows the run counter recorded by 'stuntdouble run' in the current directory.

This reports how many sandboxed sessions ran. It does not report blocked network
connections: kernel-level egress filtering is not implemented, so no connections
are blocked or counted.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		const telemetryFile = ".stuntdouble.telemetry.json"

		data, err := os.ReadFile(telemetryFile)
		if os.IsNotExist(err) {
			fmt.Println("No runs recorded in this directory yet.")
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading %s: %w", telemetryFile, err)
		}

		var stats TelemetryData
		if err := json.Unmarshal(data, &stats); err != nil {
			return fmt.Errorf("parsing %s: %w", telemetryFile, err)
		}

		fmt.Println("StuntDouble local run history")
		fmt.Println("=============================")
		fmt.Printf("Total agent sessions: %d\n", stats.TotalRuns)
		if !stats.LastRun.IsZero() {
			fmt.Printf("Last session:         %s\n", stats.LastRun.Format(time.RFC822))
		}
		fmt.Println()
		fmt.Println("Note: network egress filtering is not implemented; sessions are")
		fmt.Println("isolated by container limits only.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
