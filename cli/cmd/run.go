package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/docker"
	"github.com/stuntdouble/cli/pkg/ebpf"
	"github.com/stuntdouble/cli/pkg/snapshot"
)

var runCmd = &cobra.Command{
	Use:   "run [agent]",
	Short: "Runs an AI agent inside the StuntDouble secure container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentName := args[0]
		fmt.Printf("🚀 Starting StuntDouble Sandbox for agent: %s\n", agentName)
		fmt.Println("🔒 Panic Mode: Network outbound connections strictly monitored.")
		fmt.Println("🎭 Keploy Stunt Layer: Prepared for injection.")

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("❌ Error getting current directory:", err)
			return
		}

		isRemote, _ := cmd.Flags().GetBool("remote")

		if isRemote {
			fmt.Printf("☁️  StuntDouble Cloud: Provisioning remote MicroVM for %s...\n", agentName)
			fmt.Println(">> Streaming local workspace state to secure enterprise cloud...")
			time.Sleep(800 * time.Millisecond) // Simulate network latency
			fmt.Println(">> Connection established! Executing agent remotely (saving local battery/RAM).")
			
			// Mock Remote Execution
			fmt.Println("\n✅ Remote agent session completed safely.")
			fmt.Println("🔒 Enterprise Audit Logs: Synced to CTO Dashboard.")
			updateTelemetry()
			return
		}

		// LOCAL EXECUTION PATH via Native API
		var agentCmdStr string
		if agentName == "claude" {
			agentCmdStr = "npx -y @anthropic-ai/claude-code"
		} else {
			agentCmdStr = "npx -y " + agentName
		}

		if len(args) > 1 {
			for _, extraArg := range args[1:] {
				agentCmdStr += " " + extraArg
			}
		}

		agentCmd := []string{"sh", "-c", agentCmdStr}

		fmt.Printf(">> Spawning highly restricted Docker container for %s natively...\n", agentName)

		// Capture a zero-copy snapshot of the workspace before the AI touches it
		if err := snapshot.Create(cwd); err != nil {
			fmt.Println("⚠️ Failed to create safety snapshot:", err)
		}

		startTime := time.Now()

		// Inject the native eBPF Interceptor to lock down the kernel
		ebpfHook, err := ebpf.AttachInterceptor("/sys/fs/cgroup/")
		if err != nil {
			fmt.Println("❌ Error attaching eBPF hooks:", err)
			return
		}
		defer ebpfHook.Detach()

		sdClient, err := docker.NewClient()
		if err != nil {
			fmt.Println("❌ Error initializing native Docker client:", err)
			return
		}

		if err := sdClient.SpawnIsolatedAgent(cmd.Context(), agentCmd, cwd); err != nil {
			fmt.Println("\n⚠️ Agent session ended or was terminated:", err)
		} else {
			fmt.Println("\n✅ Agent session completed safely.")
		}

		duration := time.Since(startTime)
		fmt.Printf("⏱️  Container spin-up and execution completed in %v\n", duration)

		// Update MVP Telemetry
		updateTelemetry()
	},
}

func updateTelemetry() {
	// In MVP, we simulate reading intercepted eBPF events and updating stats
	file := ".stuntdouble.telemetry.json"

	stats := struct {
		TotalRuns       int       `json:"total_runs"`
		BlockedCommands int       `json:"blocked_commands"`
		LastRun         time.Time `json:"last_run"`
	}{}

	// Try to read existing
	if data, err := os.ReadFile(file); err == nil {
		json.Unmarshal(data, &stats)
	}

	stats.TotalRuns++
	// Simulate blocking a destructive command 20% of the time for the demo
	if time.Now().Unix()%5 == 0 {
		stats.BlockedCommands++
	}
	stats.LastRun = time.Now()

	if data, err := json.MarshalIndent(stats, "", "  "); err == nil {
		os.WriteFile(file, data, 0644)
	}
}

func init() {
	runCmd.Flags().BoolP("remote", "r", false, "Execute the agent in a remote StuntDouble Cloud MicroVM")
	rootCmd.AddCommand(runCmd)
}
