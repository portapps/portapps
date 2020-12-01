package mutex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	mu, err := Create("TestMutex")
	defer func() {
		err = Release(mu)
		assert.NoError(t, err)
	}()

	assert.NoError(t, err)
	assert.NotEmpty(t, mu)
	t.Log(mu)
}

func TestAlreadyRunning(t *testing.T) {
	mu, err := Create("TestMutex")
	defer func() {
		err = Release(mu)
		assert.NoError(t, err)
	}()

	assert.NoError(t, err)
	assert.NotEmpty(t, mu)
	t.Log(mu)

	_, err = Create("TestMutex")
	assert.Error(t, err)
	t.Log(err)
}
