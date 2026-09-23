package win

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestPermEnv(t *testing.T) {
	path := `Software\portapps-test-` + rand.Text()
	root, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.ALL_ACCESS)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, registry.DeleteKey(registry.CURRENT_USER, path))
	})
	t.Cleanup(func() {
		assert.NoError(t, root.Close())
	})
	env, _, err := registry.CreateKey(root, "Environment", registry.ALL_ACCESS)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, registry.DeleteKey(root, "Environment"))
	})
	require.NoError(t, env.Close())

	const keyName = "TEST_PERM_ENV"
	require.NoError(t, SetPermEnv(root, keyName, "portapps"))

	keyValue, err := GetPermEnv(root, keyName)
	require.NoError(t, err)
	assert.Equal(t, "portapps", keyValue)

	require.NoError(t, DeletePermEnv(root, keyName))
}

func TestGetPermEnvReturnsOpenKeyError(t *testing.T) {
	value, err := GetPermEnv(registry.NONE, "TEST_PERM_ENV")

	assert.Error(t, err)
	assert.Empty(t, value)
}
