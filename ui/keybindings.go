package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

type KeybindingsPage struct {
	*gtk.Box
	listBox      *gtk.ListBox
	addButton    *gtk.Button
	saveButton   *gtk.Button
	rules        *config.KeyBindingsConfig
	mainWindow   *gtk.Window
	rows         []*bindingRow
	searchBox    *gtk.Box
	searchEntry  *gtk.Entry
	searchButton *gtk.Button
}

type bindingRow struct {
	*gtk.ListBoxRow
	original  *config.Binding
	flagEntry *gtk.Entry
	modEntry  *gtk.Entry
	keyEntry  *gtk.Entry
	descEntry *gtk.Entry
	cmdEntry  *gtk.Entry
	argsEntry *gtk.Entry
}

func NewKeybindingsPage() *KeybindingsPage {
	page := &KeybindingsPage{
		Box: gtk.NewBox(gtk.OrientationVertical, 10),
	}
	page.SetMarginBottom(12)
	page.SetMarginTop(12)
	page.SetMarginEnd(12)
	page.SetMarginStart(12)
	page.SetCSSClasses([]string{"preferences-page"})

	// Header
	// header := gtk.NewLabel(i18n.Tr("keybindings_title"))
	// header.SetXAlign(0)
	// header.SetCSSClasses([]string{"header"})
	// page.Append(header)

	// Search Box
	page.searchBox = gtk.NewBox(gtk.OrientationHorizontal, 6)
	page.searchBox.SetMarginBottom(12)
	page.searchEntry = gtk.NewEntry()
	page.searchEntry.SetPlaceholderText(i18n.Tr("search_placeholder"))
	page.searchEntry.SetHExpand(true)
	page.searchEntry.ConnectChanged(page.onSearchChanged)

	page.searchButton = gtk.NewButtonFromIconName("edit-find-symbolic")
	page.searchButton.ConnectClicked(func() {
		page.searchBox.SetVisible(!page.searchBox.Visible())
		if page.searchBox.Visible() {
			page.searchEntry.GrabFocus()
		}
	})

	page.searchBox.Append(page.searchEntry)
	page.searchBox.SetVisible(false)
	page.Append(page.searchBox)

	// Scrolled Window
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetVExpand(true)
	page.listBox = gtk.NewListBox()
	page.listBox.SetSelectionMode(gtk.SelectionNone)
	scrolled.SetChild(page.listBox)
	page.Append(scrolled)

	// Control Buttons
	controls := gtk.NewBox(gtk.OrientationHorizontal, 6)
	controls.SetHAlign(gtk.AlignEnd)

	controls.Append(page.searchButton)

	page.addButton = gtk.NewButtonWithLabel(i18n.Tr("add_binding"))
	page.addButton.SetIconName("list-add-symbolic")
	page.addButton.ConnectClicked(page.showAddDialog)

	page.saveButton = gtk.NewButtonWithLabel(i18n.Tr("save_changes"))
	page.saveButton.SetName(i18n.Tr("save_changes"))
	//page.saveButton.SetIconName("document-save-symbolic")
	page.saveButton.ConnectClicked(page.onSaveClicked)

	controls.Append(page.addButton)
	controls.Append(page.saveButton)
	page.Append(controls)

	// Setup keyboard shortcuts
	keyController := gtk.NewEventControllerKey()
	keyController.ConnectKeyPressed(page.onKeyPress)
	page.AddController(keyController)

	// Load data
	go page.loadBindings()

	return page
}

func (p *KeybindingsPage) loadBindings() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".config/hypr/keybindings.conf")
	rules, err := config.ReadConfig(cfgPath)
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	// Update UI in main thread
	glib.IdleAdd(func() {
		p.rules = rules
		p.listBox.RemoveAll()
		p.rows = make([]*bindingRow, 0)
		for i := range rules.Bindings {
			row := p.createRuleRow(&rules.Bindings[i])
			p.rows = append(p.rows, row)
			p.listBox.Append(row)
		}
	})
}

func (p *KeybindingsPage) createRuleRow(b *config.Binding) *bindingRow {
	row := &bindingRow{original: b}

	// Main container
	box := gtk.NewBox(gtk.OrientationHorizontal, 8)
	box.SetMarginBottom(8)
	box.SetMarginTop(8)
	box.SetMarginEnd(8)
	box.SetMarginStart(8)

	// Input Grid
	grid := gtk.NewGrid()
	grid.SetColumnSpacing(6)
	grid.SetRowSpacing(6)
	grid.SetHExpand(true)

	// Input fields with translated placeholders
	row.flagEntry = createEntry(i18n.Tr("flags"), "bind"+b.Flags)
	row.flagEntry.SetHExpand(true)
	row.modEntry = createEntry(i18n.Tr("modifiers"), b.Modifiers)
	row.modEntry.SetHExpand(true)
	row.keyEntry = createEntry(i18n.Tr("key"), b.Key)
	row.keyEntry.SetHExpand(true)
	row.descEntry = createEntry(i18n.Tr("description"), b.Description)
	row.descEntry.SetHExpand(true)
	row.cmdEntry = createEntry(i18n.Tr("command"), b.Command)
	row.cmdEntry.SetHExpand(true)
	row.argsEntry = createEntry(i18n.Tr("arguments"), b.Args)
	row.argsEntry.SetHExpand(true)

	// Add translated labels with alignment
	flagsLabel := gtk.NewLabel("Flags: ")
	flagsLabel.SetXAlign(1)
	flagsLabel.SetHExpand(true)
	modLabel := gtk.NewLabel(i18n.Tr("modifiers") + ":")
	modLabel.SetXAlign(1)
	modLabel.SetHExpand(true)
	keyLabel := gtk.NewLabel(i18n.Tr("key") + ":")
	keyLabel.SetXAlign(1)
	keyLabel.SetHExpand(true)
	descLabel := gtk.NewLabel(i18n.Tr("description") + ":")
	descLabel.SetXAlign(1)
	descLabel.SetHExpand(true)
	cmdLabel := gtk.NewLabel(i18n.Tr("command") + ":")
	cmdLabel.SetXAlign(1)
	cmdLabel.SetHExpand(true)
	argsLabel := gtk.NewLabel(i18n.Tr("arguments") + ":")
	argsLabel.SetXAlign(1)
	argsLabel.SetHExpand(true)

	// Layout in two rows for better space usage
	grid.Attach(flagsLabel, 0, 1, 1, 1)
	grid.Attach(row.flagEntry, 1, 1, 1, 1)
	grid.Attach(modLabel, 2, 1, 1, 1)
	grid.Attach(row.modEntry, 3, 1, 1, 1)
	grid.Attach(keyLabel, 4, 1, 1, 1)
	grid.Attach(row.keyEntry, 5, 1, 1, 1)

	grid.Attach(descLabel, 0, 2, 1, 1)
	grid.Attach(row.descEntry, 1, 2, 1, 1)
	grid.Attach(cmdLabel, 2, 2, 1, 1)
	grid.Attach(row.cmdEntry, 3, 2, 1, 1)
	grid.Attach(argsLabel, 4, 2, 1, 1)
	grid.Attach(row.argsEntry, 5, 2, 1, 1)

	// Button Box
	btnBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
	deleteBtn := gtk.NewButtonFromIconName("user-trash-symbolic")
	deleteBtn.SetTooltipText(i18n.Tr("delete_binding"))
	deleteBtn.ConnectClicked(func() { p.deleteBinding(row) })
	btnBox.Append(deleteBtn)

	box.Append(grid)
	box.Append(btnBox)

	row.ListBoxRow = gtk.NewListBoxRow()
	row.SetChild(box)
	return row
}

func (p *KeybindingsPage) deleteBinding(row *bindingRow) {
	// Find index in config
	for i, b := range p.rules.Bindings {
		if b.LineNumber == row.original.LineNumber {
			p.rules.DeleteBinding(i)
			break
		}
	}
	p.listBox.Remove(row.ListBoxRow)
}

func (p *KeybindingsPage) onSaveClicked() {
	// Update all bindings from UI
	for _, row := range p.rows {
		row.original.Modifiers = row.modEntry.Text()
		row.original.Key = row.keyEntry.Text()
		row.original.Description = row.descEntry.Text()
		row.original.Command = row.cmdEntry.Text()
		row.original.Args = row.argsEntry.Text()
	}

	// Save to file
	cfgPath := filepath.Join(os.Getenv("HOME"), ".config/hypr/keybindings.conf")
	if err := p.rules.Save(cfgPath); err != nil {
		showErrorDialog(p.mainWindow, err.Error())
	} else {
		showInfoDialog(p.mainWindow, "配置已更新")
	}
}

func (p *KeybindingsPage) showAddDialog() {
	dialog := gtk.NewWindow()
	dialog.SetTitle(i18n.Tr("new_binding"))
	dialog.SetTransientFor(p.mainWindow)
	dialog.SetModal(true)

	content := gtk.NewBox(gtk.OrientationVertical, 6)
	content.SetMarginTop(12)
	content.SetMarginBottom(12)
	content.SetMarginStart(12)
	content.SetMarginEnd(12)
	dialog.SetChild(content)

	// Flag selection with translated strings
	model := gtk.NewStringList([]string{
		i18n.Tr("binding_type_desc"),
		i18n.Tr("binding_type_repeat"),
		i18n.Tr("binding_type_lock"),
	})
	flagCombo := gtk.NewDropDown(model, nil)
	flagCombo.SetSelected(0)

	// Input fields with translated placeholders
	modEntry := createEntry(i18n.Tr("modifiers"), "$mainMod")
	keyEntry := createEntry(i18n.Tr("key"), "F12")
	descEntry := createEntry(i18n.Tr("description"), "")
	cmdEntry := createEntry(i18n.Tr("command"), "exec")
	argsEntry := createEntry(i18n.Tr("arguments"), "./script.sh")

	grid := gtk.NewGrid()
	grid.SetColumnSpacing(6)
	grid.SetRowSpacing(6)

	// Create labels
	bindingLabel := gtk.NewLabel(i18n.Tr("binding_type") + ":")
	modLabel := gtk.NewLabel(i18n.Tr("modifiers") + ":")
	keyLabel := gtk.NewLabel(i18n.Tr("key") + ":")
	descLabel := gtk.NewLabel(i18n.Tr("description") + ":")
	cmdLabel := gtk.NewLabel(i18n.Tr("command") + ":")
	argsLabel := gtk.NewLabel(i18n.Tr("arguments") + ":")

	// Attach labels and entries to grid
	grid.Attach(bindingLabel, 0, 0, 1, 1)
	grid.Attach(flagCombo, 1, 0, 2, 1)

	grid.Attach(modLabel, 0, 1, 1, 1)
	grid.Attach(modEntry, 1, 1, 2, 1)

	grid.Attach(keyLabel, 0, 2, 1, 1)
	grid.Attach(keyEntry, 1, 2, 2, 1)

	grid.Attach(descLabel, 0, 3, 1, 1)
	grid.Attach(descEntry, 1, 3, 2, 1)

	grid.Attach(cmdLabel, 0, 4, 1, 1)
	grid.Attach(cmdEntry, 1, 4, 2, 1)

	grid.Attach(argsLabel, 0, 5, 1, 1)
	grid.Attach(argsEntry, 1, 5, 2, 1)

	content.Append(grid)

	// Buttons
	buttonBox := gtk.NewBox(gtk.OrientationHorizontal, 6)
	buttonBox.SetHAlign(gtk.AlignEnd)

	cancelBtn := gtk.NewButtonWithLabel(i18n.Tr("cancel"))
	cancelBtn.ConnectClicked(func() { dialog.Destroy() })

	confirmBtn := gtk.NewButtonWithLabel(i18n.Tr("confirm"))
	confirmBtn.ConnectClicked(func() {
		selected := flagCombo.Selected()
		flags := []string{"d", "e", "l"}[selected]
		binding := &config.Binding{
			Flags:       flags,
			Modifiers:   modEntry.Text(),
			Key:         keyEntry.Text(),
			Description: descEntry.Text(),
			Command:     cmdEntry.Text(),
			Args:        argsEntry.Text(),
		}
		p.rules.AddBinding(
			flags,
			binding.Modifiers,
			binding.Key,
			binding.Description,
			binding.Command,
			binding.Args,
		)

		// Create and add new row directly instead of reloading
		row := p.createRuleRow(binding)
		p.rows = append(p.rows, row)
		p.listBox.Append(row)

		dialog.Destroy()
	})

	buttonBox.Append(cancelBtn)
	buttonBox.Append(confirmBtn)
	content.Append(buttonBox)

	dialog.SetVisible(true)
}

func createEntry(placeholder, text string) *gtk.Entry {
	entry := gtk.NewEntry()
	entry.SetPlaceholderText(placeholder)
	entry.SetText(text)
	entry.SetCSSClasses([]string{"entry-field"})
	return entry
}

// Helper functions for dialogs
func showErrorDialog(parent *gtk.Window, msg string) {
	dialog := gtk.NewWindow()
	dialog.SetTitle(i18n.Tr("error"))
	dialog.SetTransientFor(parent)
	dialog.SetModal(true)

	content := gtk.NewBox(gtk.OrientationVertical, 6)
	content.SetMarginTop(12)
	content.SetMarginBottom(12)
	content.SetMarginStart(12)
	content.SetMarginEnd(12)
	dialog.SetChild(content)

	label := gtk.NewLabel(i18n.Tr("error") + ": " + msg)
	label.SetUseMarkup(true)
	content.Append(label)

	button := gtk.NewButtonWithLabel(i18n.Tr("ok"))
	button.ConnectClicked(func() { dialog.Destroy() })
	content.Append(button)

	dialog.SetVisible(true)
}

func showInfoDialog(parent *gtk.Window, msg string) {
	dialog := gtk.NewWindow()
	dialog.SetTitle(i18n.Tr("success"))
	dialog.SetTransientFor(parent)
	dialog.SetModal(true)

	content := gtk.NewBox(gtk.OrientationVertical, 6)
	content.SetMarginTop(12)
	content.SetMarginBottom(12)
	content.SetMarginStart(12)
	content.SetMarginEnd(12)
	dialog.SetChild(content)

	label := gtk.NewLabel(i18n.Tr(msg))
	content.Append(label)

	button := gtk.NewButtonWithLabel(i18n.Tr("ok"))
	button.ConnectClicked(func() { dialog.Destroy() })
	content.Append(button)
	dialog.SetVisible(true)
}

func (p *KeybindingsPage) onKeyPress(keyval uint, _ uint, state gdk.ModifierType) bool {
	// Check for Ctrl+F
	if state&gdk.ControlMask != 0 && (keyval == gdk.KEY_f || keyval == gdk.KEY_F) {
		p.searchBox.SetVisible(!p.searchBox.Visible())
		if p.searchBox.Visible() {
			p.searchEntry.GrabFocus()
		}
		return true
	}
	// Check for Escape
	if keyval == gdk.KEY_Escape && p.searchBox.Visible() {
		p.searchBox.SetVisible(false)
		p.searchEntry.SetText("")
		return true
	}
	return false
}

func (p *KeybindingsPage) onSearchChanged() {
	searchText := strings.ToLower(p.searchEntry.Text())

	for _, row := range p.rows {
		visible := searchText == "" || // Show all when search is empty
			strings.Contains(strings.ToLower(row.original.Modifiers), searchText) ||
			strings.Contains(strings.ToLower(row.original.Key), searchText) ||
			strings.Contains(strings.ToLower(row.original.Description), searchText) ||
			strings.Contains(strings.ToLower(row.original.Command), searchText) ||
			strings.Contains(strings.ToLower(row.original.Args), searchText)

		row.SetVisible(visible)
	}
}
