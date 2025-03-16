package ui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"github.con/falser101/hypr-gtk/i18n"
)

func createMonitorsPage() *gtk.Box {
	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)

	// 标题
	header := gtk.NewLabel(i18n.Tr("monitors_title"))
	header.SetXAlign(0)
	header.SetCSSClasses([]string{"header"})
	mainBox.Append(header)

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

	// 配置项控件
	nameEntry := gtk.NewEntry()
	nameEntry.SetText(cfg.Name)
	nameEntry.SetWidthChars(10)
	nameEntry.SetPlaceholderText(i18n.Tr("monitor_name"))

	resEntry := gtk.NewEntry()
	resEntry.SetText(cfg.Resolution)
	resEntry.SetWidthChars(15)
	resEntry.SetPlaceholderText(i18n.Tr("resolution"))

	refreshEntry := gtk.NewEntry()
	refreshEntry.SetText(fmt.Sprintf("%.2f", cfg.RefreshRate))
	refreshEntry.SetWidthChars(8)
	refreshEntry.SetPlaceholderText(i18n.Tr("refresh_rate"))

	posEntry := gtk.NewEntry()
	posEntry.SetText(cfg.Position)
	posEntry.SetWidthChars(8)
	posEntry.SetPlaceholderText(i18n.Tr("position"))

	scaleEntry := gtk.NewEntry()
	scaleEntry.SetText(fmt.Sprintf("%.1f", cfg.Scale))
	scaleEntry.SetWidthChars(5)
	scaleEntry.SetPlaceholderText(i18n.Tr("scale"))

	enableSwitch := gtk.NewSwitch()
	enableSwitch.SetActive(cfg.Enabled)

	saveButton := gtk.NewButtonWithLabel(i18n.Tr("save"))
	saveButton.ConnectClicked(func() {
		newCfg := config.MonitorConfig{
			Name:       nameEntry.Text(),
			Resolution: resEntry.Text(),
			Position:   posEntry.Text(),
			Enabled:    enableSwitch.Active(),
		}

		// 解析数值类型
		if refresh, err := strconv.ParseFloat(refreshEntry.Text(), 64); err == nil {
			newCfg.RefreshRate = refresh
		}
		if scale, err := strconv.ParseFloat(scaleEntry.Text(), 64); err == nil {
			newCfg.Scale = scale
		}

		if err := config.UpdateMonitorConfig(newCfg); err != nil {
			log.Printf(i18n.Tr("monitor_update_error")+": %v", err)
		} else {
			// 更新显示名称
			row.SetChild(createMonitorRow(newCfg))
		}
	})

	// 布局控件
	box.Append(nameEntry)
	box.Append(resEntry)
	box.Append(refreshEntry)
	box.Append(posEntry)
	box.Append(scaleEntry)
	box.Append(enableSwitch)
	box.Append(saveButton)
	row.SetChild(box)

	return row
}
