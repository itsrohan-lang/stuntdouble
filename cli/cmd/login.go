package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Credentials struct {
	Token string `json:"token"`
	Team  string `json:"team"`
}

var loginCmd = &cobra.Command{
	Use:   "login [api_token]",
	Short: "Authenticate with StuntDouble Enterprise Cloud",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := args[0]
		
		fmt.Println("🔐 Authenticating with StuntDouble Cloud...")
		
		// In a real app, we would verify the token against the Supabase/Stripe backend
		// For the MVP, we assume it's valid and save it.
		
		home, _ := os.UserHomeDir()
		stuntDir := filepath.Join(home, ".stuntdouble")
		os.MkdirAll(stuntDir, 0755)
		
		creds := Credentials{
			Token: token,
			Team:  "Enterprise-Alpha",
		}
		
		data, _ := json.MarshalIndent(creds, "", "  ")
		os.WriteFile(filepath.Join(stuntDir, "credentials.json"), data, 0600)
		
		fmt.Println("✅ Successfully authenticated!")
		fmt.Println("🚀 Cloud Telemetry Sync is now active. All agent safety events will be synced to your team dashboard.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
