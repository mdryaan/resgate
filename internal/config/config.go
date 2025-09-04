package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	StateFile string
	Verbose   bool
}

func DefaultStateFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".resgate", "state.json")
}

func Load() *Config {
	v := viper.New()
	v.SetDefault("state_file", DefaultStateFile())
	v.SetDefault("verbose", false)
	v.AutomaticEnv()

	return &Config{
		StateFile: v.GetString("state_file"),
		Verbose:   v.GetBool("verbose"),
	}
}
