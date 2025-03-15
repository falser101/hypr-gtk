package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKeyBindings(t *testing.T) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config/hypr/keybindings.conf")
	config, err := ReadConfig(path)
	if err != nil {
		panic(err)
	}

	// 示例：新增一个绑定
	config.AddBinding("d", "$mainMod", "F12", "Test description", "exec", "test.sh")

	// 示例：修改第一个绑定
	if len(config.Bindings) > 0 {
		config.UpdateBinding(0, "$mainMod", "Q", "Updated description", "exit", "")
	}

	// 示例：删除最后一个绑定
	if len(config.Bindings) > 0 {
		config.DeleteBinding(len(config.Bindings) - 1)
	}

	if err := config.Save("hyprland-modified.conf"); err != nil {
		panic(err)
	}

}
