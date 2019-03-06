package portapps

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration holds portapp configuration details
type Config struct {
	Common Common      `yaml:"common" mapstructure:"common"`
	App    interface{} `yaml:"app,omitempty" mapstructure:"app"`
}

// Common holds common configuration data
type Common struct {
	Args []string `yaml:"args" mapstructure:"cmd_switches"`
}

// loadConfig load common and app configuration
func loadConfig(appcfg interface{}) error {
	cfgPath := PathJoin(Papp.Path, fmt.Sprintf("%s.yml", Papp.ID))
	Papp.config = &Config{
		Common: Common{
			Args: []string{},
		},
		App: appcfg,
	}

	// Write sample config
	raw, err := yaml.Marshal(Papp.config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(PathJoin(Papp.Path, fmt.Sprintf("%s.sample.yml", Papp.ID)), raw, 0644)
	if err != nil {
		return err
	}

	// Check config exists
	if _, err := os.Stat(cfgPath); err != nil {
		return nil
	}

	Log.Info("Loading configuration...")
	raw, err = ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(raw, &Papp.config)
}
