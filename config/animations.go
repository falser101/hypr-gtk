package config

import (
	"os"
	"path/filepath"
	"strings"
)

type AnimationConfig struct {
	Theme     string
	ThemePath string
}

func GetAnimationConfig() (*AnimationConfig, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config/hypr/animations.conf")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	config := &AnimationConfig{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "$ANIMATION=") {
			config.Theme = strings.TrimPrefix(line, "$ANIMATION=")
		} else if strings.HasPrefix(line, "$ANIMATION_PATH=") {
			config.ThemePath = strings.TrimPrefix(line, "$ANIMATION_PATH=")
		}
	}

	return config, nil
}

func GetAvailableThemes() ([]string, error) {
	home, _ := os.UserHomeDir()
	themesDir := filepath.Join(home, ".config/hypr/animations")

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return nil, err
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".conf") {
			themes = append(themes, strings.TrimSuffix(entry.Name(), ".conf"))
		}
	}

	return themes, nil
}

func UpdateAnimationTheme(theme string) error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config/hypr/animations.conf")
	themePath := filepath.Join(home, ".config/hypr/animations", theme+".conf")

	content := []byte("$ANIMATION=" + theme + "\n$ANIMATION_PATH=" + themePath + "\n")
	return os.WriteFile(configPath, content, 0644)
}
