//go:build interactive

package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgBox(t *testing.T) {
	ret, err := MsgBox("Test MsgBox Title", "A message in a box.", MsgBoxBtnOk|MsgBoxIconError)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	t.Log(ret)
}
