package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record [command]",
	Short: "Records database and API traffic with Keploy to generate mocks",
	Long: `Runs the given command under Keploy's recording engine in a container so its
outbound calls are captured as replayable mocks.

The Keploy container runs privileged with --pid=host and --net=host because its
traffic capture requires it. That is a broad grant on your host: only record
commands you trust.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appCommand := args[0]
		fmt.Printf("Recording traffic for: %s\n", appCommand)

		cwd, _ := os.Getwd()

		// Run Keploy in record mode. Keploy itself does the capture; StuntDouble
		// only assembles the container invocation.
		recordArgs := []string{
			"run", "--rm", "-it",
			"--name", "stunt-keploy-record",
			"--privileged",
			"--pid=host",
			"--net=host",
			"-v", fmt.Sprintf("%s:/workspace", cwd),
			"-v", "/sys/fs/cgroup:/sys/fs/cgroup",
			"-v", "/sys/kernel/debug:/sys/kernel/debug",
			"-w", "/workspace",
			"ghcr.io/keploy/keploy:v2",
			"record", "-c", appCommand,
		}

		if len(args) > 1 {
			recordArgs = append(recordArgs, args[1:]...)
		}

		execCmd := exec.Command("docker", recordArgs...)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			fmt.Println("\n⚠️ Recording session ended or was terminated:", err)
		} else {
			fmt.Println("\n✅ Mocks recorded to ./keploy/tests")
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
