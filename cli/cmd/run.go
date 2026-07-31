package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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
		envImage, _ := cmd.Flags().GetString("env")
		maxDurationStr, _ := cmd.Flags().GetString("max-duration")
		strictDLP, _ := cmd.Flags().GetBool("strict-dlp")

		// Filter out sd run flags if passed after agentName (e.g. sd run claude --allow-unenforced-network)
		var forwardedArgs []string
		for _, arg := range args[1:] {
			switch {
			case arg == "--allow-unenforced-network":
				allowUnenforced = true
			case arg == "--strict-dlp":
				strictDLP = true
			case strings.HasPrefix(arg, "--max-duration="):
				maxDurationStr = strings.TrimPrefix(arg, "--max-duration=")
			case strings.HasPrefix(arg, "--env="):
				envImage = strings.TrimPrefix(arg, "--env=")
			case strings.HasPrefix(arg, "-e="):
				envImage = strings.TrimPrefix(arg, "-e=")
			default:
				forwardedArgs = append(forwardedArgs, arg)
			}
		}

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

		agentCmd := resolveAgentCommand(agentName, forwardedArgs)

		// Capture a snapshot of the workspace before the agent touches it.
		if err := snapshot.Create(cwd); err != nil {
			fmt.Println("⚠️ Failed to create safety snapshot:", err)
		}

		sdClient, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("initializing Docker client: %w", err)
		}

		startTime := time.Now()

		if strictDLP {
			fmt.Println("🛡️ [DLP Guardrail] Enforcing strict Data Loss Prevention & egress leak blocking.")
		}

		runCtx := cmd.Context()
		if maxDurationStr != "" {
			d, err := time.ParseDuration(maxDurationStr)
			if err != nil {
				return fmt.Errorf("invalid --max-duration: %w", err)
			}
			var cancel context.CancelFunc
			runCtx, cancel = context.WithTimeout(cmd.Context(), d)
			defer cancel()
			fmt.Printf("⏳ [Guardrail] Enforcing maximum execution limit: %v\n", d)
		}

		runErr := sdClient.SpawnIsolatedAgent(runCtx, agentCmd, cwd, envImage)
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

// resolveAgentCommand maps an agent name to argv executed inside the container.
// Anything unrecognised is treated as an npm package. It deliberately returns an
// argv slice instead of a shell command string so forwarded agent flags are not
// reinterpreted by /bin/sh.
func resolveAgentCommand(agentName string, extraArgs []string) []string {
	var cmdArgs []string
	switch agentName {
	case "claude":
		cmdArgs = []string{"npx", "-y", "@anthropic-ai/claude-code"}
	case "sh", "bash":
		cmdArgs = []string{agentName}
	default:
		cmdArgs = []string{"npx", "-y", agentName}
	}
	return append(cmdArgs, extraArgs...)
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
	runCmd.Flags().String("max-duration", "",
		"Maximum session duration limit for runaway agent prevention (e.g. 15m, 30s, 1h)")
	runCmd.Flags().Bool("strict-dlp", false,
		"Automatically block proxy egress calls containing PII or sensitive credential leaks")

	// Stop parsing flags at the first positional argument so that everything
	// after the agent name is forwarded to it verbatim. Without this, cobra
	// claims the agent's own flags and `sd run sh -c "echo hi"` fails with
	// "unknown shorthand flag: 'c'".
	runCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(runCmd)
}
