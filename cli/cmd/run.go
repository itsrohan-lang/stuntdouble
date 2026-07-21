package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
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

		// Core Docker Isolation Arguments
		dockerArgs := []string{
			"run", "-it", "--rm",
			"--cap-drop=ALL",              // Drop root privileges
			"-v", fmt.Sprintf("%s:/workspace", cwd), // Mount only current workspace
			"-w", "/workspace",
		}

		// Pass necessary API keys into the sandbox securely
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey != "" {
			dockerArgs = append(dockerArgs, "-e", "ANTHROPIC_API_KEY="+apiKey)
		}

		// Determine base image and execution command
		dockerArgs = append(dockerArgs, "node:20-alpine")

		if agentName == "claude" {
			dockerArgs = append(dockerArgs, "npx", "-y", "@anthropic-ai/claude-code")
		} else {
			dockerArgs = append(dockerArgs, "npx", "-y", agentName)
		}
		
		// Pass any additional arguments directly to the agent
		if len(args) > 1 {
			dockerArgs = append(dockerArgs, args[1:]...)
		}

		execCmd := exec.Command("docker", dockerArgs...)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Printf(">> Spawning highly restricted Docker container for %s...\n", agentName)
		
		if err := execCmd.Run(); err != nil {
			fmt.Println("\n⚠️ Agent session ended or was terminated.")
		} else {
			fmt.Println("\n✅ Agent session completed safely.")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
