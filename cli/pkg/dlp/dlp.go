package dlp

import (
	"regexp"
	"strings"
)

type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
)

type Finding struct {
	RuleName    string   `json:"rule_name"`
	Severity    Severity `json:"severity"`
	MatchedText string   `json:"matched_text"`
}

type DLPRule struct {
	Name     string
	Severity Severity
	Pattern  *regexp.Regexp
}

var defaultRules = []DLPRule{
	{
		Name:     "AWS Access Key ID",
		Severity: SeverityCritical,
		Pattern:  regexp.MustCompile(`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}`),
	},
	{
		Name:     "Generic API Secret Key",
		Severity: SeverityCritical,
		Pattern:  regexp.MustCompile(`(?i)(api_key|apikey|secret_key|secretkey|auth_token)\s*[:=]\s*["']?([a-z0-9_\-]{20,})["']?`),
	},
	{
		Name:     "GitHub Personal Access Token",
		Severity: SeverityCritical,
		Pattern:  regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`),
	},
	{
		Name:     "Social Security Number (SSN)",
		Severity: SeverityHigh,
		Pattern:  regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
	},
	{
		Name:     "Credit Card Number",
		Severity: SeverityHigh,
		Pattern:  regexp.MustCompile(`\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13})\b`),
	},
	{
		Name:     "Internal IPv4 Address",
		Severity: SeverityMedium,
		Pattern:  regexp.MustCompile(`\b(?:10|172\.(?:1[6-9]|2[0-9]|3[01])|192\.168)\.\d{1,3}\.\d{1,3}\b`),
	},
}

// Scanner performs enterprise Data Loss Prevention checks on text payloads
type Scanner struct {
	rules []DLPRule
}

func NewScanner() *Scanner {
	return &Scanner{rules: defaultRules}
}

// Scan inspects text content and returns all detected sensitive data findings
func (s *Scanner) Scan(text string) []Finding {
	var findings []Finding
	for _, rule := range s.rules {
		matches := rule.Pattern.FindAllString(text, -1)
		for _, match := range matches {
			findings = append(findings, Finding{
				RuleName:    rule.Name,
				Severity:    rule.Severity,
				MatchedText: match,
			})
		}
	}
	return findings
}

// Redact replaces detected sensitive data in text with [REDACTED:<RULE>]
func (s *Scanner) Redact(text string) string {
	redacted := text
	for _, rule := range s.rules {
		redacted = rule.Pattern.ReplaceAllStringFunc(redacted, func(match string) string {
			return "[REDACTED:" + strings.ToUpper(strings.ReplaceAll(rule.Name, " ", "_")) + "]"
		})
	}
	return redacted
}
