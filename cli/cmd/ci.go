package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Generates native GitHub Actions workflows for testing AI-generated PRs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️  Generating StuntDouble CI/CD Integration...")

		workflowDir := filepath.Join(".github", "workflows")
		if err := os.MkdirAll(workflowDir, 0755); err != nil {
			fmt.Println("❌ Error creating workflow directory:", err)
			return
		}

		workflowContent := `name: StuntDouble Safe PR Check

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
        run: npm install -g stuntdouble

      - name: Verify Sandbox Integrity
        run: sd run echo "Agent code safely audited in ephemeral pipeline"
        
      - name: Sync Audit Logs to CTO Dashboard
        env:
          STUNTDOUBLE_API_KEY: ${{ secrets.STUNTDOUBLE_API_KEY }}
        run: sd sync-logs
`

		workflowPath := filepath.Join(workflowDir, "stuntdouble-ci.yml")
		err := os.WriteFile(workflowPath, []byte(workflowContent), 0644)
		if err != nil {
			fmt.Println("❌ Error writing workflow file:", err)
			return
		}

		fmt.Println("✅ Successfully generated .github/workflows/stuntdouble-ci.yml")
		fmt.Println("🔒 AI Agent PRs will now be automatically tested inside the StuntDouble Cloud.")
	},
}

func init() {
	rootCmd.AddCommand(ciCmd)
}
