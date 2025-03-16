package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

func createMonitorsPage() *gtk.Box {
	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)

	// 滚动容器
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetVExpand(true)

	// 配置列表
	listBox := gtk.NewListBox()
	scrolled.SetChild(listBox)
	mainBox.Append(scrolled)

	// 加载数据
	go func() {
		configs, err := config.GetMonitors()
		if err != nil {
			log.Printf(i18n.Tr("monitor_load_error")+": %v", err)
			return
		}

		for _, cfg := range configs {
			row := createMonitorRow(cfg)
			listBox.Append(row)
		}
	}()

	return mainBox
}

func createMonitorRow(cfg config.MonitorConfig) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()
	row.SetSelectable(false)
	box := gtk.NewBox(gtk.OrientationHorizontal, 10)
	box.SetMarginStart(10)
	box.SetMarginEnd(10)
	box.SetMarginTop(5)
	box.SetMarginBottom(5)

	// Monitor name
	nameLabel := gtk.NewLabel(i18n.Tr("monitor_name") + ":")
	nameEntry := gtk.NewEntry()
	nameEntry.SetText(cfg.Name)
	nameEntry.SetWidthChars(10)

	// Get available modes for this monitor
	modes, err := config.GetAvailableModes(cfg.Name)
	if err != nil {
		log.Printf("Failed to get available modes: %v", err)
		modes = []string{}
	}

	// Resolution dropdown
	resLabel := gtk.NewLabel(i18n.Tr("resolution") + ":")
	strList := gtk.NewStringList(modes)
	resCombo := gtk.NewDropDown(strList, nil)
	resCombo.SetHExpand(true)

	// Set current mode
	currentMode := fmt.Sprintf("%s@%.2fHz", cfg.Resolution, cfg.RefreshRate)
	for i, mode := range modes {
		if mode == currentMode {
			resCombo.SetSelected(uint(i))
			break
		}
	}

	// Position
	posLabel := gtk.NewLabel(i18n.Tr("position") + ":")
	posEntry := gtk.NewEntry()
	posEntry.SetText(cfg.Position)
	posEntry.SetWidthChars(8)

	// Scale
	scaleLabel := gtk.NewLabel(i18n.Tr("scale") + ":")
	scaleEntry := gtk.NewEntry()
	scaleEntry.SetText(fmt.Sprintf("%.2f", cfg.Scale))
	scaleEntry.SetWidthChars(5)

	// Enable switch
	enableLabel := gtk.NewLabel(i18n.Tr("enabled") + ":")
	enableSwitch := gtk.NewSwitch()
	enableSwitch.SetActive(cfg.Enabled)

	// Save button
	saveButton := gtk.NewButtonWithLabel(i18n.Tr("save"))
	saveButton.ConnectClicked(func() {
		selectedMode := modes[resCombo.Selected()]
		parts := strings.Split(selectedMode, "@")
		resolution := parts[0]
		refreshRate := 60.0

		if len(parts) > 1 {
			rateStr := strings.TrimSuffix(parts[1], "Hz")
			if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
				refreshRate = rate
			}
		}

		newCfg := config.MonitorConfig{
			Name:        nameEntry.Text(),
			Resolution:  resolution,
			RefreshRate: refreshRate,
			Position:    posEntry.Text(),
			Enabled:     enableSwitch.Active(),
		}

		if scale, err := strconv.ParseFloat(scaleEntry.Text(), 64); err == nil {
			newCfg.Scale = scale
		}

		if err := config.UpdateMonitorConfig(newCfg); err != nil {
			log.Printf(i18n.Tr("monitor_update_error")+": %v", err)
			showErrorDialog(nil, "save_error")
		} else {
			showInfoDialog(nil, "save_success")
		}
	})

	// Layout
	box.Append(nameLabel)
	box.Append(nameEntry)
	box.Append(resLabel)
	box.Append(resCombo)
	box.Append(posLabel)
	box.Append(posEntry)
	box.Append(scaleLabel)
	box.Append(scaleEntry)
	box.Append(enableLabel)
	box.Append(enableSwitch)
	box.Append(saveButton)
	row.SetChild(box)

	return row
}
