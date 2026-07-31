package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/docker"
)

var swarmCmd = &cobra.Command{
	Use:   "swarm [agents...]",
	Short: "Orchestrates a Multi-Agent Swarm inside the StuntNet virtualized local network",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌐 Initializing StuntNet (Virtualized Agent Intranet)...")
		
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		// Handle Ctrl+C gracefully
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			cancel()
		}()

		sdClient, err := docker.NewClient()
		if err != nil {
			fmt.Println("❌ Error initializing native Docker client:", err)
			return
		}

		workspace, _ := os.Getwd()
		redteam, _ := cmd.Flags().GetBool("redteam")

		if redteam {
			fmt.Println("⚔️ [Red-Team Swarm] Adding automated adversarial security verification node to StuntNet...")
			args = append(args, "redteam-security-bot")
		}

		if err := sdClient.SpawnSwarm(ctx, args, workspace); err != nil {
			fmt.Println("\n❌ Swarm orchestration failed:", err)
		}
	},
}

func init() {
	swarmCmd.Flags().Bool("redteam", false, "Spawn an automated adversarial Red-Team agent node inside StuntNet")
	rootCmd.AddCommand(swarmCmd)
}
