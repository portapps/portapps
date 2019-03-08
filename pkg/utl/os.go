package utl

import (
	"os"

	"github.com/rs/zerolog/log"
)

// OverrideEnv to override an env var
func OverrideEnv(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Error().Err(err).Msgf("Cannot set %s env var", key)
	}
}
