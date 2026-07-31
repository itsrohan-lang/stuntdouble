package guard

import (
	"regexp"
	"strings"
)

type InjectionRiskLevel string

const (
	RiskCritical InjectionRiskLevel = "CRITICAL"
	RiskHigh     InjectionRiskLevel = "HIGH"
	RiskMedium   InjectionRiskLevel = "MEDIUM"
)

type InjectionFinding struct {
	VectorName  string             `json:"vector_name"`
	RiskLevel   InjectionRiskLevel `json:"risk_level"`
	MatchedText string             `json:"matched_text"`
}

type InjectionResult struct {
	IsSuspicious bool               `json:"is_suspicious"`
	Findings     []InjectionFinding `json:"findings"`
}

type PatternRule struct {
	Name      string
	RiskLevel InjectionRiskLevel
	Pattern   *regexp.Regexp
}

var defaultInjectionRules = []PatternRule{
	{
		Name:      "System Override Instruction",
		RiskLevel: RiskCritical,
		Pattern:   regexp.MustCompile(`(?i)(ignore\s+(all\s+)?previous\s+instructions|disregard\s+prior\s+prompts|override\s+system\s+rules|you\s+are\s+now\s+in\s+(developer|god)\s+mode)`),
	},
	{
		Name:      "Hidden HTML Comment Injection",
		RiskLevel: RiskHigh,
		Pattern:   regexp.MustCompile(`<!--[\s\S]*?(ignore|override|system|exec|curl|token)[\s\S]*?-->`),
	},
	{
		Name:      "Arbitrary File Exfiltration Probe",
		RiskLevel: RiskCritical,
		Pattern:   regexp.MustCompile(`(?i)(cat\s+/etc/passwd|read\s+\.env|curl\s+-X\s+POST\s+http|exfiltrate\s+token)`),
	},
	{
		Name:      "Role-Play Hijack Pattern",
		RiskLevel: RiskHigh,
		Pattern:   regexp.MustCompile(`(?i)(act\s+as\s+an?\s+unrestricted|bypass\s+all\s+safety\s+filters)`),
	},
}

// PromptGuard analyzes agent input prompts for indirect prompt injection vectors
type PromptGuard struct {
	rules []PatternRule
}

func NewPromptGuard() *PromptGuard {
	return &PromptGuard{rules: defaultInjectionRules}
}

// Analyze inspects a prompt string and returns any detected prompt injection findings
func (g *PromptGuard) Analyze(prompt string) *InjectionResult {
	var findings []InjectionFinding

	for _, rule := range g.rules {
		matches := rule.Pattern.FindAllString(prompt, -1)
		for _, match := range matches {
			findings = append(findings, InjectionFinding{
				VectorName:  rule.Name,
				RiskLevel:   rule.RiskLevel,
				MatchedText: match,
			})
		}
	}

	return &InjectionResult{
		IsSuspicious: len(findings) > 0,
		Findings:     findings,
	}
}

// Sanitize strips malicious prompt override instructions and hidden HTML comment payloads
func (g *PromptGuard) Sanitize(prompt string) string {
	sanitized := prompt
	for _, rule := range g.rules {
		sanitized = rule.Pattern.ReplaceAllString(sanitized, "[STRIPPED_PROMPT_INJECTION_VECTOR]")
	}
	// Strip zero-width characters
	zeroWidthRe := regexp.MustCompile(`[\x{200B}-\x{200D}\x{FEFF}]`)
	sanitized = zeroWidthRe.ReplaceAllString(sanitized, "")
	return strings.TrimSpace(sanitized)
}
