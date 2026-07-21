package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/snapshot"
)

var rewindCmd = &cobra.Command{
	Use:   "rewind",
	Short: "Instantly rewinds the workspace state to undo destructive agent actions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⏪ Initiating StuntDouble State Rewind...")
		
		fmt.Println(">> Reverting filesystem from isolated Git tree snapshot...")
		
		workspace, _ := os.Getwd()
		if err := snapshot.Restore(workspace); err != nil {
			fmt.Println("❌ Error rewinding workspace:", err)
			return
		}

		fmt.Println("\n✅ Workspace successfully restored to previous safe state.")
		fmt.Println("🛡️  Disaster averted. The AI agents' mistakes have been erased.")
	},
}

func init() {
	rootCmd.AddCommand(rewindCmd)
}
