package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record [command]",
	Short: "Records database and API traffic to generate StuntDouble mocks",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appCommand := args[0]
		fmt.Printf("🎙️  Starting StuntDouble Record Mode for: %s\n", appCommand)
		fmt.Println(">> Injecting eBPF proxy (Keploy) to capture network traffic...")

		// For the MVP, we simulate wrapping the command in a Keploy record container
		// In a full implementation, this would invoke Keploy's recording engine.
		
		recordArgs := []string{
			"run", "--rm", "-it",
			"--privileged", // Required for eBPF recording
			"-v", fmt.Sprintf("%s:/workspace", os.Getenv("PWD")),
			"-w", "/workspace",
			"keploy/keploy:latest",
			"record", "-c", appCommand,
		}

		if len(args) > 1 {
			recordArgs = append(recordArgs, args[1:]...)
		}

		execCmd := exec.Command("docker", recordArgs...)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Println(">> Listening for outbound database connections (Postgres, Mongo, etc.)...")
		if err := execCmd.Run(); err != nil {
			fmt.Println("\n⚠️ Recording session ended or was terminated.")
		} else {
			fmt.Println("\n✅ Mocks recorded successfully! Saved to ./keploy/tests")
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
