package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UserPrefs struct {
	Env      []EnvVar    `yaml:"env"`
	ExecOnce []ExecEntry `yaml:"exec-once"`
}

type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ExecEntry struct {
	Command string `yaml:"command"`
}

func GetUserPrefsConfig() (*UserPrefs, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "hypr", "userprefs.conf")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &UserPrefs{
				Env:      make([]EnvVar, 0),
				ExecOnce: make([]ExecEntry, 0),
			}, nil
		}
		return nil, err
	}

	// Parse existing config
	config := &UserPrefs{
		Env:      make([]EnvVar, 0),
		ExecOnce: make([]ExecEntry, 0),
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "env =") {
			parts := strings.SplitN(strings.TrimPrefix(line, "env ="), ",", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				config.Env = append(config.Env, EnvVar{
					Name:  name,
					Value: value,
				})
			}
		} else if strings.HasPrefix(line, "exec-once =") {
			command := strings.TrimSpace(strings.TrimPrefix(line, "exec-once ="))
			config.ExecOnce = append(config.ExecOnce, ExecEntry{
				Command: command,
			})
		}
	}

	return config, nil
}

func SaveUserPrefsConfig(config *UserPrefs) error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "hypr", "userprefs.conf")

	// Build config content
	var content strings.Builder

	// Write env variables
	for _, env := range config.Env {
		fmt.Fprintf(&content, "env = %s,%s\n", env.Name, env.Value)
	}
	if len(config.Env) > 0 {
		content.WriteString("\n")
	}

	// Write exec-once entries
	for _, exec := range config.ExecOnce {
		fmt.Fprintf(&content, "exec-once = %s\n", exec.Command)
	}

	// Write to file
	return os.WriteFile(configPath, []byte(content.String()), 0644)
}
