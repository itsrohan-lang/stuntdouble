package config

import (
	"bytes"
	"testing"
)

// FuzzParsePolicy ensures the StuntDouble policy parser doesn't crash on maliciously crafted inputs
func FuzzParsePolicy(f *testing.F) {
	// Add legitimate seeds
	f.Add([]byte(`
version: 1
policy:
  enforcement_mode: audit
  network:
    allow: ["github.com"]
`))
	f.Add([]byte(`
version: 1
policy:
  enforcement_mode: block
  network:
    allow: ["*"]
`))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Attempt to parse the random fuzzy bytes as a StuntDouble YAML policy
		// We expect this to fail gracefully (return an error), but NEVER panic or crash.
		_, err := ParsePolicy(bytes.NewReader(data))
		if err != nil {
			// This is expected for malformed input
			return
		}
	})
}
