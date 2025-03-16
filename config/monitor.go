package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type MonitorConfig struct {
	Name        string  `json:"name"`
	Resolution  string  `json:"resolution"`
	RefreshRate float64 `json:"refresh_rate"`
	Position    string  `json:"position"`
	Scale       float64 `json:"scale"`
	Enabled     bool    `json:"enabled"`
}

type HyprctlMonitor struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Make           string   `json:"make"`
	Model          string   `json:"model"`
	Serial         string   `json:"serial"`
	Scale          float64  `json:"scale"`
	Transform      int      `json:"transform"`
	Focused        bool     `json:"focused"`
	DpmsStatus     bool     `json:"dpmsStatus"`
	Disabled       bool     `json:"disabled"`
	AvailableModes []string `json:"availableModes"`
}

func GetAvailableModes(monitorName string) ([]string, error) {
	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run hyprctl: %v", err)
	}

	var monitors []struct {
		Name           string   `json:"name"`
		AvailableModes []string `json:"availableModes"`
	}
	if err := json.Unmarshal(output, &monitors); err != nil {
		return nil, fmt.Errorf("failed to parse hyprctl output: %v", err)
	}

	for _, m := range monitors {
		if m.Name == monitorName {
			return m.AvailableModes, nil
		}
	}
	return nil, nil
}

func GetMonitors() ([]MonitorConfig, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "hypr", "monitors.conf")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []MonitorConfig{}, nil
		}
		return nil, err
	}

	var configs []MonitorConfig
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "monitor=") {
			parts := strings.Split(strings.TrimPrefix(line, "monitor="), ",")
			if len(parts) >= 4 {
				name := strings.TrimSpace(parts[0])

				// Parse resolution and refresh rate
				resParts := strings.Split(parts[1], "@")
				resolution := resParts[0]
				refreshRate := 60.0
				if len(resParts) > 1 {
					rateStr := strings.TrimSuffix(resParts[1], "Hz")
					if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
						refreshRate = rate
					}
				}

				// Parse scale
				scale := 1.0
				if len(parts) >= 4 {
					if s, err := strconv.ParseFloat(parts[3], 64); err == nil {
						scale = s
					}
				}

				configs = append(configs, MonitorConfig{
					Name:        name,
					Resolution:  resolution,
					RefreshRate: refreshRate,
					Position:    parts[2],
					Scale:       scale,
					Enabled:     !strings.HasPrefix(line, "#"),
				})
			}
		}
	}

	return configs, nil
}

func UpdateMonitorConfig(config MonitorConfig) error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "hypr", "monitors.conf")
	content, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Build monitor line
	monitorLine := fmt.Sprintf("monitor=%s,%s@%.2f,%s,%.2f",
		config.Name,
		config.Resolution,
		config.RefreshRate,
		config.Position,
		config.Scale,
	)

	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "monitor="+config.Name) {
			if config.Enabled {
				lines[i] = monitorLine
			} else {
				lines[i] = "# " + monitorLine
			}
			found = true
			break
		}
	}

	if !found {
		if config.Enabled {
			lines = append(lines, monitorLine)
		} else {
			lines = append(lines, "# "+monitorLine)
		}
	}

	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
}
