package main

import (
	"log"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/i18n"
	"github.con/falser101/hypr-gtk/ui"
)

func main() {
	// Initialize i18n
	if err := i18n.Initialize("i18n/locales"); err != nil {
		log.Printf("Failed to initialize i18n: %v", err)
	}

	// Set initial language based on environment
	lang := os.Getenv("LANG")
	if lang != "" {
		if len(lang) >= 2 {
			i18n.SetLanguage(lang[:2]) // Use first two characters (e.g., "en" from "en_US.UTF-8")
		}
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
