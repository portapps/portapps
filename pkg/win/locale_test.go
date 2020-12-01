package win

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocale(t *testing.T) {
	locale := Locale()
	assert.NotEmpty(t, locale)
	t.Log(locale)
}
