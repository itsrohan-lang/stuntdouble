package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sd",
	Short: "StuntDouble - 1-click safe mode for AI agents",
	Long: `StuntDouble intercepts and isolates autonomous agents like Claude Code inside secure Docker containers with mocked eBPF network layers.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
