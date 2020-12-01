package win

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgBox(t *testing.T) {
	ci, ok := os.LookupEnv("CI")
	if ok && ci == "true" {
		t.Skip("Skipping testing in CI environment")
	}

	ret, err := MsgBox("Test MsgBox Title", "A message in a box.", MsgBoxBtnOk|MsgBoxIconError)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	t.Log(ret)
}
