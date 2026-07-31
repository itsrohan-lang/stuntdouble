package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/dlp"
)

var dlpCmd = &cobra.Command{
	Use:   "dlp",
	Short: "Enterprise Data Loss Prevention (DLP) & secret inspection suite",
}

var dlpScanCmd = &cobra.Command{
	Use:   "scan <target>",
	Short: "Scan a file or string payload for sensitive credentials and PII data leaks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		var content string

		// If target is a readable file, read file content; otherwise treat as string
		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			data, err := os.ReadFile(target)
			if err != nil {
				return fmt.Errorf("reading file %s: %w", target, err)
			}
			content = string(data)
			fmt.Printf("🔍 [StuntDouble DLP] Scanning file: %s\n", target)
		} else {
			content = target
			fmt.Println("🔍 [StuntDouble DLP] Scanning input string payload...")
		}

		scanner := dlp.NewScanner()
		findings := scanner.Scan(content)

		if len(findings) == 0 {
			fmt.Println("✅ [StuntDouble DLP] PASSED: Zero sensitive data leaks or PII findings detected.")
			return nil
		}

		fmt.Printf("\n🚨 [StuntDouble DLP] VIOLATIONS DETECTED (%d findings):\n", len(findings))
		for _, f := range findings {
			fmt.Printf("  • [%s] %s: %s\n", f.Severity, f.RuleName, f.MatchedText)
		}

		redactFlag, _ := cmd.Flags().GetBool("redact")
		if redactFlag {
			fmt.Println("\n🛡️ [StuntDouble DLP] Redacted Output:")
			fmt.Println(scanner.Redact(content))
		}

		return fmt.Errorf("DLP scan failed: %d sensitive finding(s) detected", len(findings))
	},
}

func init() {
	dlpScanCmd.Flags().Bool("redact", false, "Print payload with detected secrets and PII automatically redacted")
	dlpCmd.AddCommand(dlpScanCmd)
	rootCmd.AddCommand(dlpCmd)
}
