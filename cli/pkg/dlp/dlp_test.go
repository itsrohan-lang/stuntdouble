package dlp

import (
	"strings"
	"testing"
)

func TestDLPScannerAndRedaction(t *testing.T) {
	scanner := NewScanner()

	samplePayload := `Here is my AWS key: AKIAIOSFODNN7EXAMPLE and SSN 000-12-3456.`
	findings := scanner.Scan(samplePayload)

	if len(findings) < 2 {
		t.Fatalf("Expected at least 2 findings, got %d", len(findings))
	}

	redacted := scanner.Redact(samplePayload)
	if strings.Contains(redacted, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("AWS Key was not redacted: %s", redacted)
	}
	if strings.Contains(redacted, "000-12-3456") {
		t.Errorf("SSN was not redacted: %s", redacted)
	}

	if !strings.Contains(redacted, "[REDACTED:AWS_ACCESS_KEY_ID]") {
		t.Errorf("Expected AWS redaction tag, got: %s", redacted)
	}
}
