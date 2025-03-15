package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MonitorConfig struct {
	Name        string
	Resolution  string
	RefreshRate float64
	Position    string
	Scale       float64
	Enabled     bool
}

func GetMonitors() ([]MonitorConfig, error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config/hypr/monitors.conf")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configs []MonitorConfig
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "monitor=") {
			parts := strings.Split(line[8:], ",")
			if len(parts) >= 4 {
				// 解析分辨率@刷新率
				resParts := strings.Split(parts[1], "@")
				refreshRate := 0.0
				if len(resParts) > 1 {
					refreshRate, _ = strconv.ParseFloat(resParts[1], 64)
				}

				// 解析缩放比例
				scale := 1.0
				if len(parts) >= 4 {
					scale, _ = strconv.ParseFloat(parts[3], 64)
				}

				configs = append(configs, MonitorConfig{
					Name:        strings.TrimSpace(parts[0]),
					Resolution:  resParts[0],
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
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config/hypr/monitors.conf")

	// 读取并更新配置
	lines, err := readLines(path)
	if err != nil {
		return err
	}

	// 查找并更新对应配置行
	for i, line := range lines {
		if strings.Contains(line, "monitor="+config.Name) {
			newLine := buildMonitorLine(config)
			if config.Enabled {
				lines[i] = newLine
			} else {
				lines[i] = "# " + newLine
			}
			break
		}
	}

	return writeLines(path, lines)
}

func buildMonitorLine(config MonitorConfig) string {
	refresh := ""
	if config.RefreshRate > 0 {
		refresh = "@" + strconv.FormatFloat(config.RefreshRate, 'f', -1, 64)
	}
	return strings.Join([]string{
		"monitor=" + config.Name,
		config.Resolution + refresh,
		config.Position,
		strconv.FormatFloat(config.Scale, 'f', -1, 64),
	}, ",")
}

// 辅助函数
func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func writeLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}
