package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultConfig = `version: 1
config:
  # Advisory only: nothing enforces this yet. See docs/ENFORCEMENT.md.
  enforcement_mode: warn
  network:
    allow: []
    deny: []
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Writes a .stuntdouble.yaml config file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		const path = ".stuntdouble.yaml"

		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; remove it first to regenerate", path)
		}

		if err := os.WriteFile(path, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}

		fmt.Printf("✅ Wrote %s\n", path)
		fmt.Println("Note: the network allow/deny lists are not enforced yet.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
