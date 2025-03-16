package main

import (
	"log"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
	"github.con/falser101/hypr-gtk/ui"
)

func main() {
	// Initialize i18n
	if err := i18n.Initialize("i18n/locales"); err != nil {
		log.Printf("Failed to initialize i18n: %v", err)
	}

	// Load language from config
	cfg, err := config.LoadLangConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	} else {
		i18n.SetLanguage(cfg.Language)
	}

	app := gtk.NewApplication("com.github.hypr-config", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		win := ui.NewMainWindow(app)
		win.Window.SetVisible(true)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
