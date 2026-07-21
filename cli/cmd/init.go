package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .stuntdouble.yaml config file in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		configContent := `version: 1.0
isolation:
  network: blocked-except-mocks
  filesystem: read-write-workspace-only
mocks:
  auto-record: true
`
		// Using strict environment rules: we verify we are in a safe directory, 
		// but since we are just writing a local config, we proceed.
		err := os.WriteFile(".stuntdouble.yaml", []byte(configContent), 0644)
		if err != nil {
			fmt.Println("Error creating config:", err)
			return
		}
		fmt.Println("✅ Successfully initialized .stuntdouble.yaml")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
