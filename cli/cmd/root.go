package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sd",
	Short: "StuntDouble - run AI coding agents in an isolated container",
	Long: `StuntDouble runs autonomous coding agents inside a restricted Docker container
with dropped capabilities, memory and CPU limits, and a scoped workspace mount,
and can snapshot the workspace so agent changes can be reverted.

Network egress filtering is NOT implemented. Sandboxed agents can still reach
the network. See docs/ENFORCEMENT.md for what is and is not enforced.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
