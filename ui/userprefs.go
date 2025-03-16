package ui

import (
	"log"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

func createUserPrefsPage() *gtk.Box {
	page := gtk.NewBox(gtk.OrientationVertical, 10)
	page.SetMarginTop(12)
	page.SetMarginBottom(12)
	page.SetMarginStart(12)
	page.SetMarginEnd(12)

	// Create notebook for tabs
	notebook := gtk.NewNotebook()
	notebook.SetVExpand(true)
	page.Append(notebook)

	// Environment Variables tab
	envPage := createEnvTab()
	envLabel := gtk.NewLabel(i18n.Tr("env_variables"))
	notebook.AppendPage(envPage, envLabel)

	// Exec-once tab
	execPage := createExecTab()
	execLabel := gtk.NewLabel(i18n.Tr("exec_once"))
	notebook.AppendPage(execPage, execLabel)

	return page
}

func createEnvTab() *gtk.Box {
	page := gtk.NewBox(gtk.OrientationVertical, 10)
	page.SetMarginTop(12)
	page.SetMarginBottom(12)
	page.SetMarginStart(12)
	page.SetMarginEnd(12)

	// List box for env variables
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetVExpand(true)
	page.Append(scrolled)

	listBox := gtk.NewListBox()
	scrolled.SetChild(listBox)

	var rows []*gtk.ListBoxRow

	// Load configuration
	cfg, err := config.GetUserPrefsConfig()
	if err != nil {
		log.Printf("Failed to load user preferences: %v", err)
		return page
	}

	// Create env variable rows
	for _, env := range cfg.Env {
		row := createEnvRow(env, &rows)
		rows = append(rows, row)
		listBox.Append(row)
	}

	// Add new env variable button
	addButton := gtk.NewButtonWithLabel(i18n.Tr("add_env"))
	addButton.SetMarginTop(12)
	addButton.ConnectClicked(func() {
		env := config.EnvVar{
			Name:  "",
			Value: "",
		}
		row := createEnvRow(env, &rows)
		rows = append(rows, row)
		listBox.Append(row)
	})
	page.Append(addButton)

	// Save button
	saveButton := gtk.NewButtonWithLabel(i18n.Tr("save_changes"))
	saveButton.SetMarginTop(12)
	saveButton.ConnectClicked(func() {
		// Update config
		cfg.Env = make([]config.EnvVar, 0)

		// Collect env variables from UI
		for _, row := range rows {
			box := row.Child().(*gtk.Box)
			nameEntry := box.FirstChild().(*gtk.Entry)
			valueEntry := nameEntry.NextSibling().(*gtk.Entry)

			env := config.EnvVar{
				Name:  nameEntry.Text(),
				Value: valueEntry.Text(),
			}
			cfg.Env = append(cfg.Env, env)
		}

		// Save config
		if err := config.SaveUserPrefsConfig(cfg); err != nil {
			log.Printf("Failed to save user preferences: %v", err)
			showErrorDialog(nil, "save_error")
		} else {
			showInfoDialog(nil, "save_success")
		}
	})
	page.Append(saveButton)

	return page
}

func createExecTab() *gtk.Box {
	page := gtk.NewBox(gtk.OrientationVertical, 10)
	page.SetMarginTop(12)
	page.SetMarginBottom(12)
	page.SetMarginStart(12)
	page.SetMarginEnd(12)

	// List box for exec-once entries
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetVExpand(true)
	page.Append(scrolled)

	listBox := gtk.NewListBox()
	scrolled.SetChild(listBox)

	var rows []*gtk.ListBoxRow

	// Load configuration
	cfg, err := config.GetUserPrefsConfig()
	if err != nil {
		log.Printf("Failed to load user preferences: %v", err)
		return page
	}

	// Create exec-once rows
	for _, exec := range cfg.ExecOnce {
		row := createExecRow(exec, &rows)
		rows = append(rows, row)
		listBox.Append(row)
	}

	// Add new exec-once button
	addButton := gtk.NewButtonWithLabel(i18n.Tr("add_exec"))
	addButton.SetMarginTop(12)
	addButton.ConnectClicked(func() {
		exec := config.ExecEntry{
			Command: "",
		}
		row := createExecRow(exec, &rows)
		rows = append(rows, row)
		listBox.Append(row)
	})
	page.Append(addButton)

	// Save button
	saveButton := gtk.NewButtonWithLabel(i18n.Tr("save_changes"))
	saveButton.SetMarginTop(12)
	saveButton.ConnectClicked(func() {
		// Update config
		cfg.ExecOnce = make([]config.ExecEntry, 0)

		// Collect exec-once entries from UI
		for _, row := range rows {
			box := row.Child().(*gtk.Box)
			commandEntry := box.FirstChild().(*gtk.Entry)

			exec := config.ExecEntry{
				Command: commandEntry.Text(),
			}
			cfg.ExecOnce = append(cfg.ExecOnce, exec)
		}

		// Save config
		if err := config.SaveUserPrefsConfig(cfg); err != nil {
			log.Printf("Failed to save user preferences: %v", err)
			showErrorDialog(nil, "save_error")
		} else {
			showInfoDialog(nil, "save_success")
		}
	})
	page.Append(saveButton)

	return page
}

func createEnvRow(env config.EnvVar, rows *[]*gtk.ListBoxRow) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()
	box := gtk.NewBox(gtk.OrientationHorizontal, 6)
	box.SetMarginStart(6)
	box.SetMarginEnd(6)
	box.SetMarginTop(6)
	box.SetMarginBottom(6)

	// Name entry
	nameEntry := gtk.NewEntry()
	nameEntry.SetText(env.Name)
	nameEntry.SetWidthChars(20)
	nameEntry.SetPlaceholderText(i18n.Tr("env_name_placeholder"))
	box.Append(nameEntry)

	// Value entry
	valueEntry := gtk.NewEntry()
	valueEntry.SetText(env.Value)
	valueEntry.SetHExpand(true)
	valueEntry.SetPlaceholderText(i18n.Tr("env_value_placeholder"))
	box.Append(valueEntry)

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

func createExecRow(exec config.ExecEntry, rows *[]*gtk.ListBoxRow) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()
	box := gtk.NewBox(gtk.OrientationHorizontal, 6)
	box.SetMarginStart(6)
	box.SetMarginEnd(6)
	box.SetMarginTop(6)
	box.SetMarginBottom(6)

	// Command entry
	commandEntry := gtk.NewEntry()
	commandEntry.SetText(exec.Command)
	commandEntry.SetHExpand(true)
	commandEntry.SetPlaceholderText(i18n.Tr("exec_command_placeholder"))
	box.Append(commandEntry)

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
