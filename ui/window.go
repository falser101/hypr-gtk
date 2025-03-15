// ui/window.go
package ui

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MainWindow struct {
	Window *gtk.ApplicationWindow
}

func NewMainWindow(app *gtk.Application) *MainWindow {
	w := gtk.NewApplicationWindow(app)
	w.SetTitle("Hyprland 配置工具")
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
	notebook.AppendPage(keybindingsPage, gtk.NewLabel("快捷键"))

	// 添加显示器页
	monitorsPage := createMonitorsPage()
	notebook.AppendPage(monitorsPage, gtk.NewLabel("显示器"))

	return &MainWindow{Window: w}
}

func loadCSS() {
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromPath("ui/style.css")
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		cssProvider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
