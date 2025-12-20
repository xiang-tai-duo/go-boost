// Package main
// File:        main.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ca/issueIPAddress/main.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Issue IP address certificate GUI helper.
// --------------------------------------------------------------------------------
package main

import (
	"fmt"
	"os/exec"
	"runtime"

	. "common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/mesa3d"
)

//goland:noinspection GoSnakeCaseUsage
const (
	BUTTON_TEXT_SAVE     = "Generate Certificate"
	EXPLORER             = "explorer"
	LABEL_TEXT_IPV4      = "IPv4 Address"
	MODULE_NAME_MAIN     = "ca.issueIPAddress.main"
	PLACEHOLDER_TEXT_SAN = "Enter IPv4 addresses, one per line if multiple"
	WINDOW_HEIGHT        = 240
	WINDOW_TITLE         = "issueIPAddress"
	WINDOW_WIDTH         = 420
)

func main() {
	mesa3d.ExtractEnvironment()
	clientApp := app.New()
	if runtime.GOOS == "linux" {
		clientApp.Settings().SetTheme(theme.LightTheme())
	}
	window := clientApp.NewWindow(WINDOW_TITLE)
	window.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT))
	window.CenterOnScreen()
	ipEntry := widget.NewMultiLineEntry()
	ipEntry.SetPlaceHolder(PLACEHOLDER_TEXT_SAN)
	ipEntry.SetText(SelectDefaultIPv4())
	saveButton := widget.NewButton(BUTTON_TEXT_SAVE, func() {
		ipTexts := CollectIPv4s(ipEntry.Text)
		if len(ipTexts) > 0 {
			if outputDirectory, err := IssueClientCertificate(ipTexts); err != nil {
				dialog.ShowError(err, window)
			} else if err = exec.Command(EXPLORER, outputDirectory).Start(); err != nil {
				dialog.ShowError(err, window)
			}
		} else {
			dialog.ShowError(fmt.Errorf("please enter a valid IPv4 address"), window)
		}
	})
	window.SetContent(container.NewVBox(
		widget.NewLabel(LABEL_TEXT_IPV4),
		ipEntry,
		saveButton,
	))
	window.ShowAndRun()
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}
