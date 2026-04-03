package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmdCapturesExitCodeWithoutReturningError(t *testing.T) {
	result, err := Cmd(CmdOptions{
		Command:    "cmd",
		Args:       []string{"/c", "echo stdout&& echo stderr 1>&2&& exit /b 7"},
		HideWindow: true,
	})

	require.NoError(t, err)
	assert.Equal(t, uint32(7), result.ExitCode)
	assert.Equal(t, "stdout", result.Stdout)
	assert.Equal(t, "stderr", result.Stderr)
}

func TestQuickCmdReturnsErrorOnNonZeroExitCode(t *testing.T) {
	err := QuickCmd("cmd", []string{"/c", "echo stderr 1>&2&& exit /b 9"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exit code 9")
	assert.Contains(t, err.Error(), "stderr")
}
