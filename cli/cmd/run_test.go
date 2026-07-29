package cmd

import (
	"reflect"
	"testing"
)

func TestResolveAgentCommandUsesArgv(t *testing.T) {
	tests := []struct {
		name      string
		agent     string
		extraArgs []string
		want      []string
	}{
		{
			name:      "claude maps to claude-code package",
			agent:     "claude",
			extraArgs: []string{"--dangerously-skip-permissions"},
			want:      []string{"npx", "-y", "@anthropic-ai/claude-code", "--dangerously-skip-permissions"},
		},
		{
			name:      "shell command remains direct argv",
			agent:     "sh",
			extraArgs: []string{"-c", "echo hi"},
			want:      []string{"sh", "-c", "echo hi"},
		},
		{
			name:      "unknown agent treated as npm package",
			agent:     "aider",
			extraArgs: []string{"--model", "sonnet"},
			want:      []string{"npx", "-y", "aider", "--model", "sonnet"},
		},
		{
			name:      "metacharacters stay literal",
			agent:     "some-agent",
			extraArgs: []string{"$(touch /tmp/pwned)", "a;rm -rf /"},
			want:      []string{"npx", "-y", "some-agent", "$(touch /tmp/pwned)", "a;rm -rf /"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveAgentCommand(tc.agent, tc.extraArgs)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("resolveAgentCommand() = %#v, want %#v", got, tc.want)
			}
		})
	}
}
