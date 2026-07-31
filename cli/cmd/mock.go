package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/mock"
)

var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Synthetic offline API and database response generator",
}

var mockGenCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate synthetic mock response template (.stuntdouble/mocks.json)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		gen := mock.NewGenerator(cwd)
		path, err := gen.GenerateConfigFile()
		if err != nil {
			return fmt.Errorf("failed to generate mocks: %w", err)
		}

		fmt.Printf("🎭 [StuntDouble Mock] Synthetic mock rules generated at: %s\n", path)
		return nil
	},
}

var mockServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start standalone synthetic offline mock server on port 9090",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetString("port")
		addr := ":" + port
		handler := mock.NewMockHandler(nil)

		fmt.Printf("🎭 [StuntDouble Mock] Synthetic mock server listening on http://localhost%s\n", addr)
		return http.ListenAndServe(addr, handler)
	},
}

func init() {
	mockServeCmd.Flags().StringP("port", "p", "9090", "Port to bind synthetic mock server")
	mockCmd.AddCommand(mockGenCmd)
	mockCmd.AddCommand(mockServeCmd)
	rootCmd.AddCommand(mockCmd)
}
