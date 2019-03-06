package portapps

import (
	"io/ioutil"
	"strings"
)

// FindElectronAppFolder retrieved the app electron folder
func FindElectronAppFolder(prefix string, source string) string {
	Log.Infof("Lookup app folder in: %s", source)
	rootFiles, _ := ioutil.ReadDir(source)
	for _, f := range rootFiles {
		if strings.HasPrefix(f.Name(), prefix) && f.IsDir() {
			Log.Infof("Electron app folder found: %s", f.Name())
			return f.Name()
		}
	}

	LogFatalf("Electron main path does not exist with prefix '%s' in %s", prefix, source)
	return ""
}
