package ui

import (
	"log"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

func createSleepPage() *gtk.Box {
	page := gtk.NewBox(gtk.OrientationVertical, 10)
	page.SetMarginTop(12)
	page.SetMarginBottom(12)
	page.SetMarginStart(12)
	page.SetMarginEnd(12)

	// Header
	header := gtk.NewLabel(i18n.Tr("sleep_title"))
	header.SetXAlign(0)
	header.SetCSSClasses([]string{"title-1"})
	page.Append(header)

	// Listeners section
	listenersLabel := gtk.NewLabel(i18n.Tr("listeners"))
	listenersLabel.SetXAlign(0)
	listenersLabel.SetMarginTop(12)
	page.Append(listenersLabel)

	// Scrolled window for listeners
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetVExpand(true)
	page.Append(scrolled)

	// List box for listeners
	listBox := gtk.NewListBox()
	scrolled.SetChild(listBox)

	// Store rows in a slice
	var rows []*gtk.ListBoxRow

	// Load configuration
	cfg, err := config.GetHypridleConfig()
	if err != nil {
		log.Printf("Failed to load hypridle config: %v", err)
		return page
	}

	// Create listener rows
	for _, listener := range cfg.Listeners {
		row := createListenerRow(listener, &rows)
		rows = append(rows, row)
		listBox.Append(row)
	}

	// Add new listener button
	addButton := gtk.NewButtonWithLabel(i18n.Tr("add_listener"))
	addButton.SetMarginTop(12)
	addButton.ConnectClicked(func() {
		newListener := config.Listener{
			Timeout:   300,
			OnTimeout: "",
			OnResume:  "",
		}
		row := createListenerRow(newListener, &rows)
		rows = append(rows, row)
		listBox.Append(row)
	})
	page.Append(addButton)

	// Save button
	saveButton := gtk.NewButtonWithLabel(i18n.Tr("save_changes"))
	saveButton.SetMarginTop(12)
	saveButton.ConnectClicked(func() {
		// Update config
		cfg.Listeners = make([]config.Listener, 0)

		// Collect listeners from UI
		for _, row := range rows {
			box := row.Child().(*gtk.Box)
			timeoutEntry := box.FirstChild().(*gtk.Entry)
			onTimeoutEntry := timeoutEntry.NextSibling().(*gtk.Entry)
			onResumeEntry := onTimeoutEntry.NextSibling().(*gtk.Entry)

			timeout, _ := strconv.Atoi(timeoutEntry.Text())
			listener := config.Listener{
				Timeout:   timeout,
				OnTimeout: onTimeoutEntry.Text(),
				OnResume:  onResumeEntry.Text(),
			}
			cfg.Listeners = append(cfg.Listeners, listener)
		}

		// Save config
		if err := config.SaveHypridleConfig(cfg); err != nil {
			log.Printf("Failed to save hypridle config: %v", err)
			showErrorDialog(nil, "save_error")
		} else {
			showInfoDialog(nil, "save_success")
		}
	})
	page.Append(saveButton)

	return page
}

func createListenerRow(listener config.Listener, rows *[]*gtk.ListBoxRow) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()
	box := gtk.NewBox(gtk.OrientationHorizontal, 6)
	box.SetMarginStart(6)
	box.SetMarginEnd(6)
	box.SetMarginTop(6)
	box.SetMarginBottom(6)

	// Timeout entry
	timeoutEntry := gtk.NewEntry()
	timeoutEntry.SetText(strconv.Itoa(listener.Timeout))
	timeoutEntry.SetWidthChars(6)
	timeoutEntry.SetPlaceholderText(i18n.Tr("timeout_placeholder"))
	box.Append(timeoutEntry)

	// On-timeout entry
	onTimeoutEntry := gtk.NewEntry()
	onTimeoutEntry.SetText(listener.OnTimeout)
	onTimeoutEntry.SetHExpand(true)
	onTimeoutEntry.SetPlaceholderText(i18n.Tr("on_timeout_placeholder"))
	box.Append(onTimeoutEntry)

	// On-resume entry
	onResumeEntry := gtk.NewEntry()
	onResumeEntry.SetText(listener.OnResume)
	onResumeEntry.SetHExpand(true)
	onResumeEntry.SetPlaceholderText(i18n.Tr("on_resume_placeholder"))
	box.Append(onResumeEntry)

	// Delete button
	deleteButton := gtk.NewButtonFromIconName("user-trash-symbolic")
	deleteButton.ConnectClicked(func() {
		// Find and remove the row from the slice
		for i, r := range *rows {
			if r == row {
				*rows = append((*rows)[:i], (*rows)[i+1:]...)
				break
			}
		}
		row.Parent().(*gtk.ListBox).Remove(row)
	})
	box.Append(deleteButton)

	row.SetChild(box)
	return row
}
