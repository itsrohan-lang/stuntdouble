package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/docker"
	"github.com/stuntdouble/cli/pkg/ebpf"
	"github.com/stuntdouble/cli/pkg/snapshot"
)

const unenforcedWarning = `⚠️  Network egress filtering is NOT active.

    Kernel-level enforcement is not implemented on any platform yet, so this
    sandbox provides container isolation only (dropped capabilities, memory and
    CPU limits, a scoped workspace mount). The agent CAN still open outbound
    network connections.

    Re-run with --allow-unenforced-network to proceed anyway.`

var runCmd = &cobra.Command{
	Use:   "run [agent]",
	Short: "Runs an AI agent inside an isolated StuntDouble container",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		// Attempt kernel-level egress filtering. This is unimplemented on every
		// platform today, so the common path is the ErrUnsupported branch below.
		// Never continue silently: without the interceptor the agent has
		// unrestricted network access, and the caller has to know that.
		allowUnenforced, _ := cmd.Flags().GetBool("allow-unenforced-network")
		ebpfHook, err := ebpf.AttachInterceptor("/sys/fs/cgroup/")
		switch {
		case err == nil:
			defer ebpfHook.Detach()
			fmt.Println("🔒 Kernel egress filter attached.")
		case errors.Is(err, ebpf.ErrUnsupported) && allowUnenforced:
			fmt.Println("⚠️  Running WITHOUT network egress filtering (--allow-unenforced-network).")
		case errors.Is(err, ebpf.ErrUnsupported):
			return errors.New(unenforcedWarning)
		default:
			return fmt.Errorf("attaching egress filter: %w", err)
		}

		fmt.Printf("🚀 Starting StuntDouble sandbox for agent: %s\n", agentName)

		agentCmdStr := resolveAgentCommand(agentName, args[1:])
		agentCmd := []string{"sh", "-c", agentCmdStr}

		// Capture a snapshot of the workspace before the agent touches it.
		if err := snapshot.Create(cwd); err != nil {
			fmt.Println("⚠️ Failed to create safety snapshot:", err)
		}

		sdClient, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("initializing Docker client: %w", err)
		}

		startTime := time.Now()
		envImage, _ := cmd.Flags().GetString("env")
		runErr := sdClient.SpawnIsolatedAgent(cmd.Context(), agentCmd, cwd, envImage)
		if runErr != nil {
			fmt.Println("\n⚠️ Agent session ended or was terminated:", runErr)
		} else {
			fmt.Println("\n✅ Agent session completed.")
		}
		fmt.Printf("⏱️  Completed in %v\n", time.Since(startTime).Round(time.Millisecond))

		recordRun()
		return nil
	},
}

// resolveAgentCommand maps an agent name to the command executed inside the
// container. Anything unrecognised is treated as an npm package.
func resolveAgentCommand(agentName string, extraArgs []string) string {
	var cmdStr string
	switch agentName {
	case "claude":
		cmdStr = "npx -y @anthropic-ai/claude-code"
	case "sh", "bash":
		cmdStr = agentName
	default:
		cmdStr = "npx -y " + agentName
	}
	for _, arg := range extraArgs {
		cmdStr += " " + arg
	}
	return cmdStr
}

// recordRun appends to the local run counter used by `stuntdouble stats` and the
// dashboard. It counts runs only: there is no blocked-connection count to report
// while egress filtering is unimplemented.
func recordRun() {
	const file = ".stuntdouble.telemetry.json"

	stats := struct {
		TotalRuns int       `json:"total_runs"`
		LastRun   time.Time `json:"last_run"`
	}{}

	if data, err := os.ReadFile(file); err == nil {
		json.Unmarshal(data, &stats)
	}

	stats.TotalRuns++
	stats.LastRun = time.Now()

	if data, err := json.MarshalIndent(stats, "", "  "); err == nil {
		os.WriteFile(file, data, 0644)
	}
}

func init() {
	runCmd.Flags().Bool("allow-unenforced-network", false,
		"Proceed even though kernel-level network egress filtering is unavailable")
	runCmd.Flags().StringP("env", "e", "node:20-alpine",
		"Docker runtime image for the agent (e.g. python:3.11-alpine, rust:alpine)")

	// Stop parsing flags at the first positional argument so that everything
	// after the agent name is forwarded to it verbatim. Without this, cobra
	// claims the agent's own flags and `sd run sh -c "echo hi"` fails with
	// "unknown shorthand flag: 'c'".
	runCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(runCmd)
}
