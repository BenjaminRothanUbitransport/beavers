package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// loadConfig loads the configuration from a specific file path or default locations.
// 1. Provided config flag
// 2. Local beavers.yaml
// 3. ~/.beavers/config.yaml
func loadConfig(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		// If a config file is explicitly provided, use it.
		v.SetConfigFile(configPath)
	} else {
		// Check for local beavers.yaml or beavers.yml explicitly to avoid picking up the 'beavers' binary
		localConfig := ""
		if _, err := os.Stat("beavers.yaml"); err == nil {
			localConfig = "beavers.yaml"
		} else if _, err := os.Stat("beavers.yml"); err == nil {
			localConfig = "beavers.yml"
		}

		if localConfig != "" {
			v.SetConfigFile(localConfig)
		} else {
			// Search for ~/.beavers/config.yaml or ~/.beavers/config.yml
			home, err := os.UserHomeDir()
			if err == nil {
				globalBase := filepath.Join(home, ".beavers", "config")
				if _, err := os.Stat(globalBase + ".yaml"); err == nil {
					v.SetConfigFile(globalBase + ".yaml")
				} else if _, err := os.Stat(globalBase + ".yml"); err == nil {
					v.SetConfigFile(globalBase + ".yml")
				}
			}
		}
	}

	// Read the config file (if we found one)
	if v.ConfigFileUsed() != "" || configPath != "" {
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
