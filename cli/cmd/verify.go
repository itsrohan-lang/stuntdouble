package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Automated CI/CD security gate verifying workspace compliance against policy",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		strict, _ := cmd.Flags().GetBool("strict")
		fmt.Println("🔍 [StuntDouble Security Gate] Scanning workspace for policy violations...")

		violations := []string{}

		// 1. Check for uncommitted secret/credential files (.env, id_rsa, .pem, etc.)
		sensitivePatterns := []string{".env", ".env.local", ".env.production", "id_rsa", "id_ed25519", ".pem", ".pfx"}
		for _, pattern := range sensitivePatterns {
			path := filepath.Join(cwd, pattern)
			if _, err := os.Stat(path); err == nil {
				violations = append(violations, fmt.Sprintf("Sensitive credential file detected in workspace: %s", pattern))
			}
		}

		// 2. Check for untracked executable files created by agents
		untracked, err := getUntrackedFiles(cwd)
		if err == nil {
			for _, rel := range untracked {
				if strings.HasSuffix(rel, ".sh") || strings.HasSuffix(rel, ".exe") || strings.HasSuffix(rel, ".bat") {
					if strict {
						violations = append(violations, fmt.Sprintf("Untracked executable script created in workspace: %s", rel))
					}
				}
			}
		}

		if len(violations) > 0 {
			fmt.Println("\n❌ [StuntDouble Security Gate] POLICY VIOLATIONS DETECTED:")
			for _, v := range violations {
				fmt.Printf("  • %s\n", v)
			}
			return fmt.Errorf("security verification failed: %d violation(s) found", len(violations))
		}

		fmt.Println("\n✅ [StuntDouble Security Gate] PASSED: Workspace is compliant with zero-trust enterprise policy.")
		return nil
	},
}

func getUntrackedFiles(workspace string) ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	cmd.Dir = workspace
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, strings.TrimSpace(line))
		}
	}
	return result, nil
}

func init() {
	verifyCmd.Flags().Bool("strict", false, "Enable strict mode to flag all untracked executable script artifacts")
	rootCmd.AddCommand(verifyCmd)
}
