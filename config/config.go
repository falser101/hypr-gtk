package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type LangConfig struct {
	Language string `yaml:"language"`
}

var configPath string

func init() {
	home, _ := os.UserHomeDir()
	configPath = filepath.Join(home, ".config/hypr-gtk/config.yaml")
}

func LoadLangConfig() (*LangConfig, error) {
	// Create default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &LangConfig{
			Language: "en",
		}
		return defaultConfig, SaveConfig(defaultConfig)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config LangConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *LangConfig) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
