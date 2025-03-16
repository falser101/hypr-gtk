package ui

import (
	"log"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

func createAnimationsPage() *gtk.Box {
	page := gtk.NewBox(gtk.OrientationVertical, 10)
	page.SetMarginTop(12)
	page.SetMarginBottom(12)
	page.SetMarginStart(12)
	page.SetMarginEnd(12)

	// Header
	// header := gtk.NewLabel(i18n.Tr("animations_title"))
	// header.SetXAlign(0)
	// header.SetCSSClasses([]string{"header"})
	// page.Append(header)

	// Theme selection section
	themeBox := gtk.NewBox(gtk.OrientationHorizontal, 6)
	themeBox.SetMarginTop(12)

	themeLabel := gtk.NewLabel(i18n.Tr("animation_theme") + ":")
	themeLabel.SetHAlign(gtk.AlignStart)
	themeBox.Append(themeLabel)

	// Create theme dropdown
	themes, err := config.GetAvailableThemes()
	if err != nil {
		log.Printf("Failed to load themes: %v", err)
		return page
	}

	model := gtk.NewStringList(themes)
	themeDropdown := gtk.NewDropDown(model, nil)

	// Set current theme
	cfg, err := config.GetAnimationConfig()
	if err != nil {
		log.Printf("Failed to load animation config: %v", err)
		return page
	}

	for i, theme := range themes {
		if theme == cfg.Theme {
			themeDropdown.SetSelected(uint(i))
			break
		}
	}

	// Connect change handler
	themeDropdown.Connect("notify::selected", func() {
		selected := themeDropdown.Selected()
		if selected < uint(len(themes)) {
			newTheme := themes[selected]
			if err := config.UpdateAnimationTheme(newTheme); err != nil {
				log.Printf("Failed to update animation theme: %v", err)
			} else {
				showInfoDialog(nil, "animation_theme_updated")
			}
		}
	})

	themeBox.Append(themeDropdown)
	page.Append(themeBox)

	return page
}
