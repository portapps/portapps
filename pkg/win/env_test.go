package win

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows/registry"
)

func TestRefreshEnv(t *testing.T) {
	err := RefreshEnv()
	assert.NoError(t, err)
}

func TestPermEnv(t *testing.T) {
	keyName := fmt.Sprintf("TEST_PERM_ENV_%d", time.Now().Unix())
	err := SetPermEnv(registry.CURRENT_USER, keyName, "portapps")
	assert.NoError(t, err)

	keyValue, err := GetPermEnv(registry.CURRENT_USER, keyName)
	assert.NoError(t, err)
	assert.Equal(t, keyValue, "portapps")

	err = DeletePermEnv(registry.CURRENT_USER, keyName)
	assert.NoError(t, err)
}
