package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/ollama"
)

var ollamaCmd = &cobra.Command{
	Use:   "ollama",
	Short: "Start a logging reverse proxy in front of a local Ollama server",
	Long: `Proxies requests to a local Ollama server (localhost:11434) and logs them.

This is a passthrough proxy for visibility. It does not sandbox the model, block
requests, or alter what the model is permitted to do.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		
		fmt.Println("🚀 Booting StuntDouble Local AI Interceptor...")
		if err := ollama.StartProxy(port); err != nil {
			fmt.Println("❌ Error starting Ollama proxy:", err)
		}
	},
}

func init() {
	ollamaCmd.Flags().StringP("port", "p", "11435", "Port to run the protected proxy on")
	rootCmd.AddCommand(ollamaCmd)
}
