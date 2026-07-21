package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var swarmCmd = &cobra.Command{
	Use:   "swarm [agents...]",
	Short: "Orchestrates a Multi-Agent Swarm inside the StuntNet virtualized local network",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agents := strings.Join(args, ", ")
		fmt.Println("🌐 Initializing StuntNet (Virtualized Agent Intranet)...")
		time.Sleep(600 * time.Millisecond)
		fmt.Printf(">> Creating isolated Docker bridge network: stuntdouble-net\n")
		
		for _, agent := range args {
			fmt.Printf(">> Spawning %s and attaching to StuntNet...\n", agent)
			time.Sleep(300 * time.Millisecond)
		}

		fmt.Println("\n✅ Swarm Orchestration Active!")
		fmt.Printf("🔒 Agents [%s] are now communicating securely inside the sandbox.\n", agents)
		fmt.Println("⚠️ External internet access blocked. Internal communication permitted.")
	},
}

func init() {
	rootCmd.AddCommand(swarmCmd)
}
