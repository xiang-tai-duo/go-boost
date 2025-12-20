// Package qrcode
// File:        qrcode.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/qrcode/qrcode.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: QR code generation and recognition helpers.
// --------------------------------------------------------------------------------
package qrcode

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	gozxingqrcode "github.com/makiuchi-d/gozxing/qrcode"
	qrcodelib "github.com/skip2/go-qrcode"
	"github.com/xiang-tai-duo/go-boost/logger"
	"golang.design/x/clipboard"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	DEFAULT_QUALITY    = qrcodelib.Medium
	DEFAULT_SIZE       = 256
	MODULE_NAME_QRCODE = "qrcode"
	TRIM_CHARS         = "\"'"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_QRCODE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_QRCODE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_QRCODE, logger.SKIP_STACK_FRAMES_BASE)
}

func GenerateImage(content string) image.Image {
	result := image.Image(nil)
	if qr, err := qrcodelib.New(content, DEFAULT_QUALITY); err == nil {
		result = qr.Image(DEFAULT_SIZE)
	}
	return result
}

func ReadFromClipboard() string {
	result := ""
	err := error(nil)
	if err = clipboard.Init(); err == nil {
		imgData := clipboard.Read(clipboard.FmtImage)
		if imgData != nil {
			var img image.Image
			if img, _, err = image.Decode(bytes.NewReader(imgData)); err == nil {
				var bmp *gozxing.BinaryBitmap
				if bmp, err = gozxing.NewBinaryBitmapFromImage(img); err == nil {
					reader := gozxingqrcode.NewQRCodeReader()
					var qrResult *gozxing.Result
					if qrResult, err = reader.Decode(bmp, nil); err == nil {
						result = qrResult.GetText()
					}
				}
			}
		}
		if result == "" {
			textData := clipboard.Read(clipboard.FmtText)
			if textData != nil {
				filePath := string(textData)
				filePath = strings.TrimSpace(filePath)
				filePath = strings.Trim(filePath, TRIM_CHARS)
				if _, err = os.Stat(filePath); err == nil {
					result = ReadFromFile(filePath)
				}
			}
		}
	}
	return result
}

func ReadFromFile(filePath string) string {
	result := ""
	err := error(nil)
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		var img image.Image
		if img, _, err = image.Decode(file); err == nil {
			var bmp *gozxing.BinaryBitmap
			if bmp, err = gozxing.NewBinaryBitmapFromImage(img); err == nil {
				reader := gozxingqrcode.NewQRCodeReader()
				var qrResult *gozxing.Result
				if qrResult, err = reader.Decode(bmp, nil); err == nil {
					result = qrResult.GetText()
				}
			}
		}
	}
	return result
}

func SaveToFile(content string, filePath string) error {
	err := error(nil)
	img := GenerateImage(content)
	if img == nil {
		err = os.ErrInvalid
	} else {
		var file *os.File
		if file, err = os.Create(filePath); err == nil {
			defer file.Close()
			err = png.Encode(file, img)
		}
	}
	return err
}

func ScanScreen() string {
	result := ""
	err := error(nil)
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n && result == ""; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		var img *image.RGBA
		if img, err = screenshot.CaptureRect(bounds); err == nil {
			var bmp *gozxing.BinaryBitmap
			if bmp, err = gozxing.NewBinaryBitmapFromImage(img); err == nil {
				reader := gozxingqrcode.NewQRCodeReader()
				var qrResult *gozxing.Result
				if qrResult, err = reader.Decode(bmp, nil); err == nil {
					result = qrResult.GetText()
				}
			}
		}
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_QRCODE, logger.SKIP_STACK_FRAMES_BASE)
}
