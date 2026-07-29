package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the local StuntDouble telemetry API",
	Long: `Serves the local run counter over HTTP on loopback so the dashboard can read
it. This API reports telemetry; it does not enforce policy.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		
		fmt.Println("🌐 Booting StuntDouble Control Plane...")
		if err := api.StartServer(port); err != nil {
			fmt.Println("❌ Error starting API server:", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the API on")
	rootCmd.AddCommand(serveCmd)
}
