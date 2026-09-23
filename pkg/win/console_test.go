package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConsoleTitle(t *testing.T) {
	title, err := GetConsoleTitle()
	assert.NoError(t, err)
	assert.NotEmpty(t, title)
	t.Log(title)
}

func TestSetConsoleTitle(t *testing.T) {
	original, err := GetConsoleTitle()
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, SetConsoleTitle(original))
	})

	err = SetConsoleTitle("Test Console Title")
	assert.NoError(t, err)

	title, err := GetConsoleTitle()
	assert.NoError(t, err)
	assert.NotEmpty(t, title)
	assert.Equal(t, title, "Test Console Title")
	t.Log(title)
}
