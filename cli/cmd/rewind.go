package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/snapshot"
)

var rewindCmd = &cobra.Command{
	Use:   "rewind",
	Short: "Restores the workspace to the snapshot taken at the last run",
	Long: `Restores tracked files to the snapshot captured when 'stuntdouble run' last
started, and removes files created since.

This discards uncommitted work in the workspace, including your own changes made
after the snapshot was taken.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		fmt.Println("Restoring workspace from the last StuntDouble snapshot...")
		if err := snapshot.Restore(workspace); err != nil {
			return fmt.Errorf("rewinding workspace: %w", err)
		}

		fmt.Println("✅ Workspace restored to the last snapshot.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rewindCmd)
}
