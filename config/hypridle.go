package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Listener struct {
	Timeout   int
	OnTimeout string
	OnResume  string
}

type HypridleConfig struct {
	LockScreen string
	Listeners  []Listener
}

func GetHypridleConfig() (*HypridleConfig, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config/hypr/hypridle.conf")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &HypridleConfig{}
	scanner := bufio.NewScanner(file)
	currentListener := &Listener{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Parse LOCKSCREEN variable
		if strings.HasPrefix(line, "$LOCKSCREEN = ") {
			config.LockScreen = strings.TrimPrefix(line, "$LOCKSCREEN = ")
			continue
		}

		// Parse listener blocks
		if strings.HasPrefix(line, "listener {") {
			if currentListener.Timeout > 0 {
				config.Listeners = append(config.Listeners, *currentListener)
			}
			currentListener = &Listener{}
			continue
		}

		// Parse listener properties
		if strings.HasPrefix(line, "timeout = ") {
			timeout, _ := strconv.Atoi(strings.TrimPrefix(line, "timeout = "))
			currentListener.Timeout = timeout
		} else if strings.HasPrefix(line, "on-timeout = ") {
			currentListener.OnTimeout = strings.TrimPrefix(line, "on-timeout = ")
		} else if strings.HasPrefix(line, "on-resume = ") {
			currentListener.OnResume = strings.TrimPrefix(line, "on-resume = ")
		}
	}

	// Add the last listener if it exists
	if currentListener.Timeout > 0 {
		config.Listeners = append(config.Listeners, *currentListener)
	}

	return config, scanner.Err()
}

func SaveHypridleConfig(config *HypridleConfig) error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config/hypr/hypridle.conf")

	var content strings.Builder
	content.WriteString("#!     ░▒▒▒░░░▓▓             ___________\n")
	content.WriteString("#!    ░░▒▒▒░░░░░▓▓        //___________/\n")
	content.WriteString("#!   ░░▒▒▒░░░░░▓▓     _   _ _    _ _____\n")
	content.WriteString("#!   ░░▒▒░░░░░▓▓▓▓▓▓ | | | | |  | |  __/\n")
	content.WriteString("#!    ░▒▒░░░░░▓▓   ▓▓ | |_| | |_/ /| |___\n")
	content.WriteString("#!     ░▒▒░░▓▓   ▓▓   \\__  |____/ |____/\n")
	content.WriteString("#!       ░▒▓▓   ▓▓  //____/\n\n")

	content.WriteString("$LOCKSCREEN = " + config.LockScreen + "\n\n")

	content.WriteString("general {\n")
	content.WriteString("    lock_cmd = $LOCKSCREEN\n")
	content.WriteString("    unlock_cmd = #notify-send \"unlock!\"      # same as above, but unlock\n")
	content.WriteString("    before_sleep_cmd = $LOCKSCREEN    # command ran before sleep\n")
	content.WriteString("    after_sleep_cmd = # notify-send \"Awake!\"  # command ran after sleep\n")
	content.WriteString("    ignore_dbus_inhibit = 0\n")
	content.WriteString("}\n\n")

	for _, listener := range config.Listeners {
		content.WriteString("listener {\n")
		content.WriteString("    timeout = " + strconv.Itoa(listener.Timeout) + "\n")
		if listener.OnTimeout != "" {
			content.WriteString("    on-timeout = " + listener.OnTimeout + "\n")
		}
		if listener.OnResume != "" {
			content.WriteString("    on-resume = " + listener.OnResume + "\n")
		}
		content.WriteString("}\n\n")
	}

	content.WriteString("# hyprlang noerror true\n")
	content.WriteString("# Source anything  from this path if you want to add your own listener\n")
	content.WriteString("# source command actually do not exist yet\n")
	content.WriteString("source = ~/.config/hypridle/*\n")
	content.WriteString("# hyprlang noerror false\n")

	return os.WriteFile(configPath, []byte(content.String()), 0644)
}
