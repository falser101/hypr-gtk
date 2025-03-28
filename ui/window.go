// ui/window.go
package ui

import (
	"os"
	"path/filepath"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/i18n"
)

type MainWindow struct {
	Window *gtk.ApplicationWindow
}

func NewMainWindow(app *gtk.Application) *MainWindow {
	w := gtk.NewApplicationWindow(app)
	w.SetTitle(i18n.Tr("hyprland_config_tool"))
	w.SetDefaultSize(800, 600)

	// 加载CSS样式
	loadCSS()

	// 主容器
	mainBox := gtk.NewBox(gtk.OrientationVertical, 0)
	w.SetChild(mainBox)

	// 标签页
	notebook := gtk.NewNotebook()
	mainBox.Append(notebook)

	// 添加快捷键页
	keybindingsPage := NewKeybindingsPage()
	notebook.AppendPage(keybindingsPage, gtk.NewLabel(i18n.Tr("keybindings")))

	// 添加显示器页
	monitorsPage := createMonitorsPage()
	notebook.AppendPage(monitorsPage, gtk.NewLabel(i18n.Tr("monitors")))

	// 添加动画页
	animationsPage := createAnimationsPage()
	notebook.AppendPage(animationsPage, gtk.NewLabel(i18n.Tr("animations")))

	// 添加睡眠管理页
	sleepPage := createSleepPage()
	notebook.AppendPage(sleepPage, gtk.NewLabel(i18n.Tr("sleep")))

	// 添加用户偏好页
	userPrefsPage := createUserPrefsPage()
	notebook.AppendPage(userPrefsPage, gtk.NewLabel(i18n.Tr("userprefs")))

	// 添加设置页
	settingsPage := createSettingsPage()
	notebook.AppendPage(settingsPage, gtk.NewLabel(i18n.Tr("settings")))

	return &MainWindow{Window: w}
}

func loadCSS() {
	cssProvider := gtk.NewCSSProvider()
	settings := gtk.SettingsGetDefault()
	home, _ := os.UserHomeDir()
	dark := settings.ObjectProperty("gtk-application-prefer-dark-theme").(bool)
	var configPath string
	if dark {
		configPath = filepath.Join(home, ".config/gtk-4.0/", "gtk-dark.css")
	} else {
		configPath = filepath.Join(home, ".config/gtk-4.0/", "gtk.css")
	}
	cssProvider.LoadFromPath(configPath)

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
