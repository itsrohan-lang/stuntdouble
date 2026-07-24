package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .stuntdouble.yaml config file in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		banner := `
   _____ __                  __  ____              __    __     
  / ___// /___  ______  / /_/ __ \____  __  __/ /_  / /__   
  \__ \/ __/ / / / __ \/ __/ / / / __ \/ / / / __ \/ / _ \  
 ___/ / /_/ /_/ / / / / /_/ /_/ / /_/ / /_/ / /_/ / /  __/  
/____/\__/\__,_/_/ /_/\__/_____/\____/\__,_/_.___/_/\___/   
                                                           
`
		fmt.Println(banner)
		
		configContent := `version: 1.0
isolation:
  network: blocked-except-mocks
  filesystem: read-write-workspace-only
mocks:
  auto-record: true
`
		err := os.WriteFile(".stuntdouble.yaml", []byte(configContent), 0644)
		if err != nil {
			fmt.Println("Error creating config:", err)
			return
		}
		
		// Native IDE Support: Generate .cursorrules for Cursor Agent
		cursorRules := `You are operating inside the StuntDouble secure sandbox.
Never attempt to bypass the Docker isolation.
Assume all database ports (5432, 27017) are intercepted and mocked by Keploy.
`
		err = os.WriteFile(".cursorrules", []byte(cursorRules), 0644)
		if err != nil {
			fmt.Println("Error creating .cursorrules:", err)
			return
		}

		fmt.Println("✅ Successfully initialized .stuntdouble.yaml and native .cursorrules")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
