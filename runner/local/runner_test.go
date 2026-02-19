package local

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viant/gosh/runner"
	"testing"
	"time"
)

func TestService_Run(t *testing.T) {
	runner := New()
	output, code, err := runner.Run(context.Background(), "ls /")
	assert.Nil(t, err)
	assert.Equal(t, 0, code)
	assert.Truef(t, len(output) > 0, "output was empty")
	assert.True(t, runner.PID() > 0)

}

func TestService_Run_PipelineWithRqEmptyResult(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX shell and wc")
	}
	if _, err := exec.LookPath("rq"); err != nil {
		t.Skip("rq is not installed")
	}

	r := New()
	defer func() { _ = r.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, code, err := r.Run(ctx, `printf '{"items":[]}' | rq '.items[]' | wc -l`, runner.WithTimeout(4000))
	require.NoError(t, err)
	assert.Equal(t, 0, code)
	assert.Equal(t, "0", strings.TrimSpace(output))
}
