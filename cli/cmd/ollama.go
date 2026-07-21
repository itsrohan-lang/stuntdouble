package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/ollama"
)

var ollamaCmd = &cobra.Command{
	Use:   "ollama",
	Short: "Start the StuntDouble proxy for local Ollama models",
	Long:  "Intercepts local AI traffic to localhost:11434, injects Universal Stunt Protocol attestations, and ensures local models execute inside the sandbox.",
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
