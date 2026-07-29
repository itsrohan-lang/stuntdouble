package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// workflowTemplate only references commands that exist. Earlier versions
// generated a `sd sync-logs` step for a command the CLI never had, so the
// generated workflow always failed.
const workflowTemplate = `name: StuntDouble Agent Check

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  agent-audit:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Install StuntDouble
        run: npm install -g stuntdouble-sandbox-cli

      # Network egress filtering is not implemented, so the sandbox provides
      # container isolation only. --allow-unenforced-network acknowledges that.
      - name: Run agent in sandbox
        run: sd run sh --allow-unenforced-network -c "echo audited"
`

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Generates a GitHub Actions workflow that runs agents in the sandbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		workflowDir := filepath.Join(".github", "workflows")
		if err := os.MkdirAll(workflowDir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", workflowDir, err)
		}

		path := filepath.Join(workflowDir, "stuntdouble-ci.yml")
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; remove it first to regenerate", path)
		}

		if err := os.WriteFile(path, []byte(workflowTemplate), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}

		fmt.Printf("✅ Wrote %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ciCmd)
}
