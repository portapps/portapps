package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs64Arch(t *testing.T) {
	is64Arch := Is64Arch()
	assert.NotNil(t, is64Arch)
	t.Log(is64Arch)
}
