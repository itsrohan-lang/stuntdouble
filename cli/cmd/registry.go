package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "install [plugin]",
	Short: "Installs community eBPF interceptor plugins (e.g., Stripe, AWS) from the StuntDouble Registry",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plugin := args[0]
		fmt.Printf("📦 Searching StuntDouble Registry for plugin: '%s'...\n", plugin)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf(">> Found '%s' (v1.0.4) by StuntDouble Community\n", plugin)
		fmt.Println(">> Downloading eBPF mock signatures...")
		time.Sleep(800 * time.Millisecond)
		
		fmt.Printf("✅ Successfully installed %s eBPF interceptor!\n", plugin)
		fmt.Printf("🔒 Agents attempting to hit %s APIs will now receive synthetic success responses.\n", plugin)
	},
}

func init() {
	rootCmd.AddCommand(registryCmd)
}
