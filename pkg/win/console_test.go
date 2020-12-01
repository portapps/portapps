package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConsoleTitle(t *testing.T) {
	title, err := GetConsoleTitle()
	assert.NoError(t, err)
	assert.NotEmpty(t, title)
	t.Log(title)
}

func TestSetConsoleTitle(t *testing.T) {
	err := SetConsoleTitle("Test Console Title")
	assert.NoError(t, err)

	title, err := GetConsoleTitle()
	assert.NoError(t, err)
	assert.NotEmpty(t, title)
	assert.Equal(t, title, "Test Console Title")
	t.Log(title)
}
