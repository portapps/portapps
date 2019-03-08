package utl

import (
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
)

// FindElectronAppFolder retrieved the app electron folder
func FindElectronAppFolder(prefix string, source string) string {
	log.Info().Msgf("Lookup electron app folder in: %s", source)
	rootFiles, _ := ioutil.ReadDir(source)
	for _, f := range rootFiles {
		if strings.HasPrefix(f.Name(), prefix) && f.IsDir() {
			log.Info().Msgf("Electron app folder found: %s", f.Name())
			return f.Name()
		}
	}
	log.Fatal().Msgf("Electron main path does not exist with prefix '%s' in %s", prefix, source)
	return ""
}
