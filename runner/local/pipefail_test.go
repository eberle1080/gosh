package local

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/eberle1080/gosh/runner"
)

func TestPipefailAcrossShells(t *testing.T) {
	for _, tc := range []struct {
		shell            string
		pipelineExitCode int
	}{
		{shell: "dash", pipelineExitCode: 0},
		{shell: "bash", pipelineExitCode: 1},
	} {
		t.Run(tc.shell, func(t *testing.T) {
			shell, err := exec.LookPath(tc.shell)
			if err != nil {
				t.Skipf("%s is unavailable: %v", tc.shell, err)
			}
			session := New(runner.WithShell(shell))
			t.Cleanup(func() { _ = session.Close() })
			ctx := context.Background()
			for _, command := range []struct {
				text     string
				exitCode int
				output   string
			}{
				{text: "false", exitCode: 1},
				{text: "false | true", exitCode: tc.pipelineExitCode},
				{text: "printf '%s\\n' session-alive", output: "session-alive"},
			} {
				output, code, err := session.Run(ctx, command.text, runner.WithTimeout(2000))
				if err != nil || code != command.exitCode || strings.TrimSpace(output) != command.output {
					t.Fatalf("%q returned (%q, %d, %v), want (%q, %d, nil)", command.text, output, code, err, command.output, command.exitCode)
				}
			}
		})
	}
}
