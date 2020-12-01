package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver, err := GetVersion()
	assert.NoError(t, err)
	assert.NotEmpty(t, ver)
	t.Log(ver)
}
