package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/snapshot"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Manage named zero-copy workspace checkpoints",
}

var checkpointSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current workspace state to a named checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return snapshot.SaveCheckpoint(cwd, args[0])
	},
}

var checkpointRestoreCmd = &cobra.Command{
	Use:   "restore <name>",
	Short: "Rewind workspace to a named checkpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return snapshot.RestoreCheckpoint(cwd, args[0])
	},
}

var checkpointListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved workspace checkpoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		checkpoints, err := snapshot.ListCheckpoints(cwd)
		if err != nil {
			return err
		}
		if len(checkpoints) == 0 {
			fmt.Println("No saved checkpoints found.")
			return nil
		}
		fmt.Println("🔖 Saved Checkpoints:")
		for _, name := range checkpoints {
			fmt.Printf("  • %s\n", name)
		}
		return nil
	},
}

func init() {
	checkpointCmd.AddCommand(checkpointSaveCmd)
	checkpointCmd.AddCommand(checkpointRestoreCmd)
	checkpointCmd.AddCommand(checkpointListCmd)
	rootCmd.AddCommand(checkpointCmd)
}
