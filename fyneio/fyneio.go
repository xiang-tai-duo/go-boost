// Package main
// File:        fyneio.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/fyne/fyne.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: FYNE provides functions to create and manage Fyne UI applications.
// --------------------------------------------------------------------------------

package fyneio

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

//goland:noinspection SpellCheckingInspection,GoUnusedExportedType
type (
	THEME struct {
		Theme                        fyne.Theme
		ThemeFont                    fyne.Resource
		ThemeVariant                 fyne.ThemeVariant
		ColorNameBackground          color.Color
		ColorNameButton              color.Color
		ColorNameDisabledButton      color.Color
		ColorNameDisabled            color.Color
		ColorNameError               color.Color
		ColorNameFocus               color.Color
		ColorNameForeground          color.Color
		ColorNameForegroundOnError   color.Color
		ColorNameForegroundOnPrimary color.Color
		ColorNameForegroundOnSuccess color.Color
		ColorNameForegroundOnWarning color.Color
		ColorNameHeaderBackground    color.Color
		ColorNameHover               color.Color
		ColorNameHyperlink           color.Color
		ColorNameInputBackground     color.Color
		ColorNameInputBorder         color.Color
		ColorNameMenuBackground      color.Color
		ColorNameOverlayBackground   color.Color
		ColorNamePlaceHolder         color.Color
		ColorNamePressed             color.Color
		ColorNamePrimary             color.Color
		ColorNameScrollBar           color.Color
		ColorNameScrollBarBackground color.Color
		ColorNameSelection           color.Color
		ColorNameSeparator           color.Color
		ColorNameShadow              color.Color
		ColorNameSuccess             color.Color
		ColorNameWarning             color.Color
		SizeNameCaptionText          float32
		SizeNameHeadingText          float32
		SizeNameInlineIcon           float32
		SizeNameInputBorder          float32
		SizeNamePadding              float32
		SizeNameScrollBar            float32
		SizeNameScrollBarSmall       float32
		SizeNameSeparatorThickness   float32
		SizeNameSubHeadingText       float32
		SizeNameText                 float32
	}

	FYNE struct {
		app    fyne.App
		window fyne.Window
	}

	LABEL struct {
		canvas.Text
	}
)

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func NewTheme(themeVariant fyne.ThemeVariant) *THEME {
	fyneTheme := theme.DefaultTheme()
	return &THEME{
		Theme:                        fyneTheme,
		ThemeVariant:                 themeVariant,
		ColorNameBackground:          fyneTheme.Color(theme.ColorNameBackground, themeVariant),
		ColorNameButton:              fyneTheme.Color(theme.ColorNameButton, themeVariant),
		ColorNameDisabledButton:      fyneTheme.Color(theme.ColorNameDisabledButton, themeVariant),
		ColorNameDisabled:            fyneTheme.Color(theme.ColorNameDisabled, themeVariant),
		ColorNameError:               fyneTheme.Color(theme.ColorNameError, themeVariant),
		ColorNameFocus:               fyneTheme.Color(theme.ColorNameFocus, themeVariant),
		ColorNameForeground:          fyneTheme.Color(theme.ColorNameForeground, themeVariant),
		ColorNameForegroundOnError:   fyneTheme.Color(theme.ColorNameForegroundOnError, themeVariant),
		ColorNameForegroundOnPrimary: fyneTheme.Color(theme.ColorNameForegroundOnPrimary, themeVariant),
		ColorNameForegroundOnSuccess: fyneTheme.Color(theme.ColorNameForegroundOnSuccess, themeVariant),
		ColorNameForegroundOnWarning: fyneTheme.Color(theme.ColorNameForegroundOnWarning, themeVariant),
		ColorNameHeaderBackground:    fyneTheme.Color(theme.ColorNameHeaderBackground, themeVariant),
		ColorNameHover:               fyneTheme.Color(theme.ColorNameHover, themeVariant),
		ColorNameHyperlink:           fyneTheme.Color(theme.ColorNameHyperlink, themeVariant),
		ColorNameInputBackground:     fyneTheme.Color(theme.ColorNameInputBackground, themeVariant),
		ColorNameInputBorder:         fyneTheme.Color(theme.ColorNameInputBorder, themeVariant),
		ColorNameMenuBackground:      fyneTheme.Color(theme.ColorNameMenuBackground, themeVariant),
		ColorNameOverlayBackground:   fyneTheme.Color(theme.ColorNameOverlayBackground, themeVariant),
		ColorNamePlaceHolder:         fyneTheme.Color(theme.ColorNamePlaceHolder, themeVariant),
		ColorNamePressed:             fyneTheme.Color(theme.ColorNamePressed, themeVariant),
		ColorNamePrimary:             fyneTheme.Color(theme.ColorNamePrimary, themeVariant),
		ColorNameScrollBar:           fyneTheme.Color(theme.ColorNameScrollBar, themeVariant),
		ColorNameScrollBarBackground: fyneTheme.Color(theme.ColorNameScrollBarBackground, themeVariant),
		ColorNameSelection:           fyneTheme.Color(theme.ColorNameSelection, themeVariant),
		ColorNameSeparator:           fyneTheme.Color(theme.ColorNameSeparator, themeVariant),
		ColorNameShadow:              fyneTheme.Color(theme.ColorNameShadow, themeVariant),
		ColorNameSuccess:             fyneTheme.Color(theme.ColorNameSuccess, themeVariant),
		ColorNameWarning:             fyneTheme.Color(theme.ColorNameWarning, themeVariant),
		SizeNameCaptionText:          12,
		SizeNameHeadingText:          24,
		SizeNameInlineIcon:           20,
		SizeNameInputBorder:          2,
		SizeNamePadding:              4,
		SizeNameScrollBar:            16,
		SizeNameScrollBarSmall:       6,
		SizeNameSeparatorThickness:   1,
		SizeNameSubHeadingText:       20,
		SizeNameText:                 18,
	}
}

//goland:noinspection GoUnusedExportedFunction
func NewLabel(text string, themeVariant fyne.ThemeVariant) *canvas.Text {
	result := canvas.NewText(text, color.Black)
	if themeVariant == theme.VariantDark {
		result = canvas.NewText(text, color.White)
	}
	return result
}

func (t *THEME) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	result := t.Theme.Color(name, variant)
	switch name {
	case theme.ColorNameBackground:
		result = t.ColorNameBackground
	case theme.ColorNameButton:
		result = t.ColorNameButton
	case theme.ColorNameDisabledButton:
		result = t.ColorNameDisabledButton
	case theme.ColorNameDisabled:
		result = t.ColorNameDisabled
	case theme.ColorNameError:
		result = t.ColorNameError
	case theme.ColorNameFocus:
		result = t.ColorNameFocus
	case theme.ColorNameForeground:
		result = t.ColorNameForeground
	case theme.ColorNameForegroundOnError:
		result = t.ColorNameForegroundOnError
	case theme.ColorNameForegroundOnPrimary:
		result = t.ColorNameForegroundOnPrimary
	case theme.ColorNameForegroundOnSuccess:
		result = t.ColorNameForegroundOnSuccess
	case theme.ColorNameForegroundOnWarning:
		result = t.ColorNameForegroundOnWarning
	case theme.ColorNameHeaderBackground:
		result = t.ColorNameHeaderBackground
	case theme.ColorNameHover:
		result = t.ColorNameHover
	case theme.ColorNameHyperlink:
		result = t.ColorNameHyperlink
	case theme.ColorNameInputBackground:
		result = t.ColorNameInputBackground
	case theme.ColorNameInputBorder:
		result = t.ColorNameInputBorder
	case theme.ColorNameMenuBackground:
		result = t.ColorNameMenuBackground
	case theme.ColorNameOverlayBackground:
		result = t.ColorNameOverlayBackground
	case theme.ColorNamePlaceHolder:
		result = t.ColorNamePlaceHolder
	case theme.ColorNamePressed:
		result = t.ColorNamePressed
	case theme.ColorNamePrimary:
		result = t.ColorNamePrimary
	case theme.ColorNameScrollBar:
		result = t.ColorNameScrollBar
	case theme.ColorNameScrollBarBackground:
		result = t.ColorNameScrollBarBackground
	case theme.ColorNameSelection:
		result = t.ColorNameSelection
	case theme.ColorNameSeparator:
		result = t.ColorNameSeparator
	case theme.ColorNameShadow:
		result = t.ColorNameShadow
	case theme.ColorNameSuccess:
		result = t.ColorNameSuccess
	case theme.ColorNameWarning:
		result = t.ColorNameWarning
	default:
		result = t.Theme.Color(name, variant)
	}
	return result
}

func (t *THEME) Font(style fyne.TextStyle) fyne.Resource {
	result := t.Theme.Font(style)
	if t.ThemeFont != nil {
		result = t.ThemeFont
	}
	return result
}

func (t *THEME) Icon(name fyne.ThemeIconName) fyne.Resource {
	result := t.Theme.Icon(name)
	return result
}

func (t *THEME) Size(name fyne.ThemeSizeName) float32 {
	result := t.Theme.Size(name)
	switch name {
	case theme.SizeNameCaptionText:
		result = t.SizeNameCaptionText
	case theme.SizeNameHeadingText:
		result = t.SizeNameHeadingText
	case theme.SizeNameInlineIcon:
		result = t.SizeNameInlineIcon
	case theme.SizeNameInputBorder:
		result = t.SizeNameInputBorder
	case theme.SizeNamePadding:
		result = t.SizeNamePadding
	case theme.SizeNameScrollBar:
		result = t.SizeNameScrollBar
	case theme.SizeNameScrollBarSmall:
		result = t.SizeNameScrollBarSmall
	case theme.SizeNameSeparatorThickness:
		result = t.SizeNameSeparatorThickness
	case theme.SizeNameSubHeadingText:
		result = t.SizeNameSubHeadingText
	case theme.SizeNameText:
		result = t.SizeNameText
	default:
		result = t.Theme.Size(name)
	}
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func LoadFont(name string, content []byte) fyne.Resource {
	return fyne.NewStaticResource(name, content)
}
