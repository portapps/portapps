package portapps

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/portapps/portapps/pkg/utl"
	"gopkg.in/yaml.v2"
)

// Config holds portapp configuration details
type Config struct {
	Common Common      `yaml:"common" mapstructure:"common"`
	App    interface{} `yaml:"app,omitempty" mapstructure:"app"`
}

// Common holds common configuration data
type Common struct {
	Args []string `yaml:"args" mapstructure:"cmd_switches"`
}

// loadConfig load common and app configuration
func (app *App) loadConfig(appcfg interface{}) (err error) {
	cfgPath := utl.PathJoin(app.RootPath, fmt.Sprintf("%s.yml", app.ID))
	app.config = &Config{
		Common: Common{
			Args: []string{},
		},
		App: appcfg,
	}

	// Write sample config
	raw, err := yaml.Marshal(app.config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(utl.PathJoin(app.RootPath, fmt.Sprintf("%s.sample.yml", app.ID)), raw, 0644)
	if err != nil {
		return err
	}

	// Skip if config file not found
	if _, err := os.Stat(cfgPath); err != nil {
		return nil
	}

	// Read config
	raw, err = ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(raw, &app.config)
}
