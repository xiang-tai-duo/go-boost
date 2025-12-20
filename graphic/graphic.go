// Package graphic
// File:        graphic.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/graphic/graphic.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: GRAPHIC is a wrapper for image and graphics operations, providing watermark and image manipulation functions.
// --------------------------------------------------------------------------------
package graphic

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/windows/wingdi"
)

//goland:noinspection GoSnakeCaseUsage
type (
	WATERMARK_INFO struct {
		Text         string
		FontFilePath string
		FontSize     float64
		Rotation     float64
		Spacing      float64
		ColorRatio   int
	}
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	BASE_HEIGHT           = 297.0
	BASE_WIDTH            = 210.0
	MAX_ALPHA_VALUE       = 255
	MODULE_NAME_GRAPHIC   = "graphic"
	WATERMARK_COLOR_ALPHA = 50
	WATERMARK_FONT_SIZE   = 30
	WATERMARK_ROTATION    = -45
	WATERMARK_SPACING     = 100
)

//goland:noinspection GoUnusedGlobalVariable,GoSnakeCaseUsage
var (
	darkGrayColor = color.RGBA{R: 64, G: 64, B: 64, A: 255}
	whiteColor    = color.RGBA{R: 255, G: 255, B: 255, A: 255}
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_GRAPHIC, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_GRAPHIC, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_GRAPHIC, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func DrawWatermark(img image.Image, text string, fontFilePath string) image.Image {
	result := DrawWatermarkEx(img, WATERMARK_INFO{
		Text:         text,
		FontFilePath: fontFilePath,
		FontSize:     WATERMARK_FONT_SIZE,
		Rotation:     WATERMARK_ROTATION,
		Spacing:      WATERMARK_SPACING,
		ColorRatio:   WATERMARK_COLOR_ALPHA,
	})
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func DrawWatermarkEx(img image.Image, config WATERMARK_INFO) image.Image {
	var result image.Image
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	baseWidthPixels := wingdi.MmToPixelX(BASE_WIDTH)
	baseHeightPixels := wingdi.MmToPixelY(BASE_HEIGHT)
	scaleX := float64(width) / baseWidthPixels
	scaleY := float64(height) / baseHeightPixels
	scaledFontSize := config.FontSize * scaleX
	watermarkDC := gg.NewContext(width, height)
	watermarkDC.LoadFontFace(config.FontFilePath, scaledFontSize)
	watermarkDC.Push()
	watermarkDC.Translate(float64(width)/2, float64(height)/2)
	watermarkDC.Rotate(gg.Radians(config.Rotation))
	watermarkDC.Translate(-float64(width), -float64(height))
	textWidth, textHeight := watermarkDC.MeasureString(config.Text)
	stepX := textWidth + config.Spacing*scaleX
	stepY := textHeight + config.Spacing*scaleY
	watermarkColorAlpha := config.ColorRatio
	if watermarkColorAlpha < 0 {
		watermarkColorAlpha = 0
	}
	if watermarkColorAlpha > MAX_ALPHA_VALUE {
		watermarkColorAlpha = MAX_ALPHA_VALUE
	}
	watermarkDC.SetColor(color.Black)
	for y := float64(-height); y < float64(height)*2; y += stepY {
		for x := float64(-width); x < float64(width)*2; x += stepX {
			watermarkDC.DrawString(config.Text, x, y)
		}
	}
	watermarkDC.Pop()
	watermarkImg := watermarkDC.Image()
	rgbaResult := image.NewRGBA(bounds)
	alphaRatio := float64(watermarkColorAlpha) / MAX_ALPHA_VALUE
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origR, origG, origB, _ := img.At(x, y).RGBA()
			wmR, wmG, wmB, wmA := watermarkImg.At(x, y).RGBA()
			if wmA == 0 {
				rgbaResult.SetRGBA(x, y, color.RGBA{
					R: uint8(origR >> 8),
					G: uint8(origG >> 8),
					B: uint8(origB >> 8),
					A: MAX_ALPHA_VALUE,
				})
			} else {
				wmAlpha := float64(wmA>>8) / MAX_ALPHA_VALUE * alphaRatio
				mixedR := uint8(float64(origR>>8)*(1-wmAlpha) + float64(wmR>>8)*wmAlpha)
				mixedG := uint8(float64(origG>>8)*(1-wmAlpha) + float64(wmG>>8)*wmAlpha)
				mixedB := uint8(float64(origB>>8)*(1-wmAlpha) + float64(wmB>>8)*wmAlpha)
				rgbaResult.SetRGBA(x, y, color.RGBA{R: mixedR, G: mixedG, B: mixedB, A: MAX_ALPHA_VALUE})
			}
		}
	}
	result = rgbaResult
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_GRAPHIC, logger.SKIP_STACK_FRAMES_BASE)
}
