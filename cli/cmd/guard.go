package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuntdouble/cli/pkg/guard"
)

var guardCmd = &cobra.Command{
	Use:   "guard",
	Short: "Indirect Prompt Injection & adversarial vector detector",
}

var guardPromptCmd = &cobra.Command{
	Use:   "prompt <target>",
	Short: "Inspect an input prompt string or markdown context file for hidden prompt injection attacks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		var content string

		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			data, err := os.ReadFile(target)
			if err != nil {
				return fmt.Errorf("reading prompt file %s: %w", target, err)
			}
			content = string(data)
			fmt.Printf("🔍 [StuntDouble Guard] Analyzing prompt file context: %s\n", target)
		} else {
			content = target
			fmt.Println("🔍 [StuntDouble Guard] Analyzing prompt string...")
		}

		pg := guard.NewPromptGuard()
		res := pg.Analyze(content)

		if !res.IsSuspicious {
			fmt.Println("✅ [StuntDouble Guard] PASSED: Zero indirect prompt injection vectors detected.")
			return nil
		}

		fmt.Printf("\n🚨 [StuntDouble Guard] PROMPT INJECTION VECTORS DETECTED (%d findings):\n", len(res.Findings))
		for _, f := range res.Findings {
			fmt.Printf("  • [%s] %s: %s\n", f.RiskLevel, f.VectorName, f.MatchedText)
		}

		sanitizeFlag, _ := cmd.Flags().GetBool("sanitize")
		if sanitizeFlag {
			fmt.Println("\n🛡️ [StuntDouble Guard] Sanitized Safe Prompt Output:")
			fmt.Println(pg.Sanitize(content))
		}

		return fmt.Errorf("prompt guard failed: %d adversarial vector(s) detected", len(res.Findings))
	},
}

func init() {
	guardPromptCmd.Flags().Bool("sanitize", false, "Print prompt with detected injection vectors automatically stripped")
	guardCmd.AddCommand(guardPromptCmd)
	rootCmd.AddCommand(guardCmd)
}
