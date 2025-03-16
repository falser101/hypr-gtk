package ui

import (
	"log"
	"os"
	"os/exec"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

func createSettingsPage() *gtk.Box {
	page := gtk.NewBox(gtk.OrientationVertical, 10)
	page.SetMarginTop(12)
	page.SetMarginBottom(12)
	page.SetMarginStart(12)
	page.SetMarginEnd(12)

	// Language section
	langBox := gtk.NewBox(gtk.OrientationHorizontal, 6)
	langBox.SetMarginTop(12)

	langLabel := gtk.NewLabel(i18n.Tr("language") + ":")
	langLabel.SetHAlign(gtk.AlignStart)
	langBox.Append(langLabel)

	// Create language dropdown
	languages := map[string]string{
		"en": "English",
		"zh": "中文",
	}

	// Create string list for dropdown
	var langNames []string
	langCodes := make([]string, 0, len(languages))
	for code, name := range languages {
		langNames = append(langNames, name)
		langCodes = append(langCodes, code)
	}

	model := gtk.NewStringList(langNames)
	langDropdown := gtk.NewDropDown(model, nil)

	// Set current language from config
	cfg, err := config.LoadLangConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	} else {
		for i, code := range langCodes {
			if code == cfg.Language {
				langDropdown.SetSelected(uint(i))
				break
			}
		}
	}

	// Connect change handler
	langDropdown.Connect("notify::selected", func() {
		selected := langDropdown.Selected()
		if selected < uint(len(langCodes)) {
			newLang := langCodes[selected]

			// Update config file
			cfg, err := config.LoadLangConfig()
			if err != nil {
				log.Printf("Failed to load config: %v", err)
				return
			}

			cfg.Language = newLang
			if err := config.SaveConfig(cfg); err != nil {
				log.Printf("Failed to save config: %v", err)
				return
			}

			// Restart application
			if err := exec.Command(os.Args[0], os.Args[1:]...).Start(); err != nil {
				log.Printf("Failed to restart application: %v", err)
			}
			os.Exit(0)
		}
	})

	langBox.Append(langDropdown)
	page.Append(langBox)

	return page
}
