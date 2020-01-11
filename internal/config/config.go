package config

import (
	"github.com/amitizle/muffin/pkg/checks"
	"github.com/amitizle/muffin/pkg/notifiers"
	"github.com/spf13/viper"
)

// Config is the struct that holds the entire configuration
// for the app
type Config struct {
	Checks    []*checkInstance   `yaml:"checks"`
	Notifiers []*notifierInstace `yaml:"notifiers"`
	Log       *LogConfig         `yaml:"log"`
}

type checkInstance struct {
	Check checks.Check

	Type   string `yaml:"type"`
	Cron   string `yaml:"cron"`
	Name   string `yaml:"name"`
	Config map[string]interface{}
}

type notifierInstace struct {
	Notifier notifiers.Notifier

	Type   string                 `yaml:"type"`
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

// LogConfig is the struct that holds the configuration for the logger
type LogConfig struct {
	Level string `yaml:"level"`
}

func init() {
	viper.SetDefault("log.level", "debug")
}

// New return a new `*Config` with `Checks` slice initialized
// with 0 checks
func New() *Config {
	return &Config{
		Checks: []*checkInstance{},
	}
}
