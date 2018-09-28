package portapps

import "os"

// OverrideEnv to override an env var
func OverrideEnv(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		Log.Fatalf("Cannot set %s env var: %v", key, err)
	}
}
