package ui

import (
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.con/falser101/hypr-gtk/config"
	"log"
	"strconv"
)

func createMonitorsPage() *gtk.Box {
	mainBox := gtk.NewBox(gtk.OrientationVertical, 10)
	mainBox.SetMarginTop(10)
	mainBox.SetMarginStart(10)
	mainBox.SetMarginEnd(10)

	// 标题
	header := gtk.NewLabel("显示器配置")
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
			log.Printf("读取显示器配置失败: %v", err)
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

	resEntry := gtk.NewEntry()
	resEntry.SetText(cfg.Resolution)
	resEntry.SetWidthChars(15)

	refreshEntry := gtk.NewEntry()
	refreshEntry.SetText(fmt.Sprintf("%.2f", cfg.RefreshRate))
	refreshEntry.SetWidthChars(8)

	posEntry := gtk.NewEntry()
	posEntry.SetText(cfg.Position)
	posEntry.SetWidthChars(8)

	scaleEntry := gtk.NewEntry()
	scaleEntry.SetText(fmt.Sprintf("%.1f", cfg.Scale))
	scaleEntry.SetWidthChars(5)

	enableSwitch := gtk.NewSwitch()
	enableSwitch.SetActive(cfg.Enabled)

	saveButton := gtk.NewButtonWithLabel("保存")
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
			log.Printf("更新配置失败: %v", err)
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
