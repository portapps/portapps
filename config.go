package portapps

import (
	"fmt"
	"os"

	"github.com/portapps/portapps/v3/pkg/utl"
	"gopkg.in/yaml.v3"
)

// Config holds portapp configuration details
type Config struct {
	Common Common      `yaml:"common" mapstructure:"common"`
	App    interface{} `yaml:"app,omitempty" mapstructure:"app"`
}

// Common holds common configuration data
type Common struct {
	DisableLog bool              `yaml:"disable_log" mapstructure:"disable_log"`
	Args       []string          `yaml:"args" mapstructure:"args"`
	Env        map[string]string `yaml:"env" mapstructure:"env"`
	AppPath    string            `yaml:"app_path" mapstructure:"app_path"`
}

// loadConfig load common and app configuration
func (app *App) loadConfig(appcfg interface{}) (err error) {
	cfgPath := utl.PathJoin(app.RootPath, fmt.Sprintf("%s.yml", app.ID))
	app.config = &Config{
		Common: Common{
			DisableLog: false,
			Args:       []string{},
			Env:        map[string]string{},
			AppPath:    "",
		},
		App: appcfg,
	}

	// Write sample config
	raw, err := yaml.Marshal(app.config)
	if err != nil {
		return err
	}
	err = os.WriteFile(utl.PathJoin(app.RootPath, fmt.Sprintf("%s.sample.yml", app.ID)), raw, 0o644)
	if err != nil {
		return err
	}

	// Skip if config file not found
	if _, err := os.Stat(cfgPath); err != nil {
		return nil
	}

	// Read config
	raw, err = os.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(raw, &app.config)
}
