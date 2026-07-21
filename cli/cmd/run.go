package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [agent]",
	Short: "Runs an AI agent inside the StuntDouble secure container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentName := args[0]
		fmt.Printf("🚀 Starting StuntDouble Sandbox for agent: %s\n", agentName)
		fmt.Println("🔒 Panic Mode: Network outbound connections blocked.")
		fmt.Println("🎭 Keploy Stunt Layer: Active (intercepting DB ports).")
		
		// TODO: Implement Docker SDK container provisioning here
		fmt.Printf(">> Spawning ephemeral container for %s...\n", agentName)
		fmt.Println("[MOCK OUTPUT] Container spawned successfully. Hooking IO streams...")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
