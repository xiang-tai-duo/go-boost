//go:build windows

// Package wingdi
// File:        wingdi.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/wingdi/wingdi.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: wingdi is a wrapper for Windows GDI operations, providing structures and methods for device mode management.
// --------------------------------------------------------------------------------
package wingdi

/*
#include <windows.h>
#include <wingdi.h>
#include <winreg.h>
#include <stdlib.h>
#cgo LDFLAGS: -lgdi32 -ladvapi32

#define CCHDEVICENAME 32
#define CCHFORMNAME 32
#define LF_FACESIZE 32
#define LF_FULLFACESIZE 64

extern int __stdcall EnumFontFamExProcCallback(ENUMLOGFONTEXW *lpelfe, NEWTEXTMETRICEXW *lpntme, DWORD FontType, LPARAM lParam);

static int EnumFontFamiliesExWrapper(HDC hdc, LOGFONTW *lpLogfont, LPARAM lParam) {
    return EnumFontFamiliesExW(hdc, lpLogfont, (FONTENUMPROCW)EnumFontFamExProcCallback, lParam, 0);
}
*/
import "C"
import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/windows/advapi32"
	"github.com/xiang-tai-duo/go-boost/windows/kernel32"
)

//goland:noinspection SpellCheckingInspection
type (
	DEVMODEA struct {
		DmDeviceName         [C.CCHDEVICENAME]byte
		DmFormName           [C.CCHFORMNAME]byte
		DmPosition           POINTL
		DmBitsPerPel         uint32
		DmDisplayFixedOutput uint32
		DmDisplayFlags       uint32
		DmDisplayFrequency   uint32
		DmDisplayOrientation uint32
		DmDitherType         uint32
		DmFields             uint32
		DmICMIntent          uint32
		DmICMMethod          uint32
		DmMediaType          uint32
		DmNup                uint32
		DmPanningHeight      uint32
		DmPanningWidth       uint32
		DmPelsHeight         uint32
		DmPelsWidth          uint32
		DmReserved1          uint32
		DmReserved2          uint32
		DmCollate            int16
		DmColor              int16
		DmCopies             int16
		DmDefaultSource      int16
		DmDuplex             int16
		DmOrientation        int16
		DmPaperLength        int16
		DmPaperSize          int16
		DmPaperWidth         int16
		DmPrintQuality       int16
		DmScale              int16
		DmTTOption           int16
		DmYResolution        int16
		DmDriverExtra        uint16
		DmDriverVersion      uint16
		DmLogPixels          uint16
		DmSize               uint16
		DmSpecVersion        uint16
	}

	DEVMODEW struct {
		DmDeviceName         [C.CCHDEVICENAME]uint16
		DmFormName           [C.CCHFORMNAME]uint16
		DmPosition           POINTL
		DmBitsPerPel         uint32
		DmDisplayFixedOutput uint32
		DmDisplayFlags       uint32
		DmDisplayFrequency   uint32
		DmDisplayOrientation uint32
		DmDitherType         uint32
		DmFields             uint32
		DmICMIntent          uint32
		DmICMMethod          uint32
		DmMediaType          uint32
		DmNup                uint32
		DmPanningHeight      uint32
		DmPanningWidth       uint32
		DmPelsHeight         uint32
		DmPelsWidth          uint32
		DmReserved1          uint32
		DmReserved2          uint32
		DmCollate            int16
		DmColor              int16
		DmCopies             int16
		DmDefaultSource      int16
		DmDuplex             int16
		DmOrientation        int16
		DmPaperLength        int16
		DmPaperSize          int16
		DmPaperWidth         int16
		DmPrintQuality       int16
		DmScale              int16
		DmTTOption           int16
		DmYResolution        int16
		DmDriverExtra        uint16
		DmDriverVersion      uint16
		DmLogPixels          uint16
		DmSize               uint16
		DmSpecVersion        uint16
	}

	ENUMLOGFONTEXW struct {
		ElfLogFont  LOGFONTW
		ElfFullName [C.LF_FULLFACESIZE]uint16
		ElfScript   [C.LF_FACESIZE]uint16
		ElfStyle    [C.LF_FACESIZE]uint16
	}

	FONTINFO struct {
		FaceName string
		FilePath string
		FullName string
		Script   string
		Style    string
		Weight   int32
		FontType uint32
		CharSet  byte
	}

	FONTSIGNATURE struct {
		FsCsb [2]uint32
		FsUsb [4]uint32
	}

	LOGFONTW struct {
		LfFaceName       [C.LF_FACESIZE]uint16
		LfEscapement     int32
		LfHeight         int32
		LfOrientation    int32
		LfWeight         int32
		LfWidth          int32
		LfCharSet        byte
		LfClipPrecision  byte
		LfItalic         byte
		LfOutPrecision   byte
		LfPitchAndFamily byte
		LfQuality        byte
		LfStrikeOut      byte
		LfUnderline      byte
	}

	LPDEVMODEA *DEVMODEA
	LPDEVMODEW *DEVMODEW

	NEWTEXTMETRICEXW struct {
		NtmFontSig FONTSIGNATURE
		NtmTm      NEWTEXTMETRICW
	}

	NEWTEXTMETRICW struct {
		TmAscent           int32
		TmAveCharWidth     int32
		TmDescent          int32
		TmDigitizedAspectX int32
		TmDigitizedAspectY int32
		TmExternalLeading  int32
		TmHeight           int32
		TmInternalLeading  int32
		TmMaxCharWidth     int32
		TmOverhang         int32
		TmWeight           int32
		NtmAvgWidth        uint32
		NtmCellHeight      uint32
		NtmFlags           uint32
		NtmSizeEM          uint32
		TmBreakChar        uint16
		TmDefaultChar      uint16
		TmFirstChar        uint16
		TmLastChar         uint16
		TmCharSet          byte
		TmItalic           byte
		TmPitchAndFamily   byte
		TmStruckOut        byte
		TmUnderlined       byte
	}

	NPDEVMODEA *DEVMODEA
	NPDEVMODEW *DEVMODEW
	PDEVMODEA  *DEVMODEA
	PDEVMODEW  *DEVMODEW

	POINTL struct {
		X int32
		Y int32
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,SpellCheckingInspection
const (
	AMPERSAND_SUFFIX       = "&"
	DEFAULT_DPI            = 96
	FONTS_DIRECTORY_SUFFIX = "\\Fonts\\"
	LOGPIXELSX             = 88
	LOGPIXELSY             = 90
	MM_PER_INCH            = 25.4
	MODULE_NAME_WINGDI     = "windows.wingdi"
	REG_FONT_PATH          = "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Fonts"
	SEMICOLON_SUFFIX       = ";"
	SPACE_CHARACTER        = " "
	TRUE_TYPE_SUFFIX       = "(TrueType)"
	WCHAR_SIZE             = 2
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_WINGDI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_WINGDI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_WINGDI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func CString(s string) *C.wchar_t {
	result := (*C.wchar_t)(nil)
	if utf16, err := syscall.UTF16FromString(s); err == nil {
		if ptr := C.malloc(C.size_t((len(utf16) + 1) * int(unsafe.Sizeof(uint16(0))))); ptr != nil {
			src := utf16
			dst := unsafe.Slice((*uint16)(ptr), len(utf16)+1)
			copy(dst, src)
			dst[len(utf16)] = 0
			result = (*C.wchar_t)(ptr)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func FreeCString(ptr *C.wchar_t) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func GoString(lpwsz interface{}) string {
	result := ""
	if ptr, ok := lpwsz.(*C.wchar_t); ok {
		if pwsz := (*uint16)(unsafe.Pointer(ptr)); pwsz != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(pwsz))[:])
		}
	} else if ptr, ok := lpwsz.(*uint16); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if ptr, ok := lpwsz.(C.LPWSTR); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if ptr, ok := lpwsz.(C.LPSTR); ok {
		if ptr != nil {
			result = syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
		}
	} else if slice, ok := lpwsz.([]uint16); ok {
		result = syscall.UTF16ToString(slice)
	} else {
		rv := reflect.ValueOf(lpwsz)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
			if rv.Kind() == reflect.Array && rv.Type().Elem().Kind() == reflect.Uint16 {
				slice := unsafe.Slice((*uint16)(unsafe.Pointer(rv.UnsafeAddr())), rv.Len())
				result = syscall.UTF16ToString(slice)
			} else {
				__error(fmt.Sprintf("GoString: unsupported type %T", lpwsz))
			}
		} else {
			__error(fmt.Sprintf("GoString: unsupported type %T", lpwsz))
		}
	}
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnusedParameter,GoVetUnsafePointer
//export EnumFontFamExProcCallback
func EnumFontFamExProcCallback(lpelfe *C.ENUMLOGFONTEXW, lpntme *C.NEWTEXTMETRICEXW, FontType C.DWORD, lParam C.LPARAM) C.int {
	result := (*[]FONTINFO)(unsafe.Pointer(uintptr(lParam)))
	fontInfo := FONTINFO{
		FaceName: GoString((*[32]uint16)(unsafe.Pointer(&lpelfe.elfLogFont.lfFaceName))),
		FullName: GoString((*[64]uint16)(unsafe.Pointer(&lpelfe.elfFullName))),
		Style:    GoString((*[32]uint16)(unsafe.Pointer(&lpelfe.elfStyle))),
		Script:   GoString((*[32]uint16)(unsafe.Pointer(&lpelfe.elfScript))),
		CharSet:  byte(lpelfe.elfLogFont.lfCharSet),
		Weight:   int32(lpelfe.elfLogFont.lfWeight),
		FontType: uint32(FontType),
	}
	*result = append(*result, fontInfo)
	return 1
}

//goland:noinspection GoUnusedExportedFunction
func EnumFontFamilies() []FONTINFO {
	result := make([]FONTINFO, 0)
	if hDC := C.GetDC(C.HWND(nil)); hDC != C.HDC(nil) {
		var logFont LOGFONTW
		logFont.LfCharSet = C.DEFAULT_CHARSET
		logFont.LfFaceName[0] = 0
		lParam := C.LPARAM(uintptr(unsafe.Pointer(&result)))
		C.EnumFontFamiliesExWrapper(hDC, (*C.LOGFONTW)(unsafe.Pointer(&logFont)), lParam)
		C.ReleaseDC(C.HWND(nil), hDC)
	}
	mapping := EnumFontFamiliesEx()
	for i := range result {
		if fileName, ok := mapping[result[i].FaceName]; ok {
			result[i].FilePath = fileName
		} else if fileName, ok := mapping[result[i].FullName]; ok {
			result[i].FilePath = fileName
		}
	}
	return result
}

//goland:noinspection GoUnreachableCode
func EnumFontFamiliesEx() map[string]string {
	result := make(map[string]string)
	windowsDir, _ := kernel32.GetWindowsDirectoryW()
	fontsDirectoryPath := windowsDir + FONTS_DIRECTORY_SUFFIX
	if hKey := advapi32.RegOpenKeyEx(advapi32.HKEY_LOCAL_MACHINE, REG_FONT_PATH, 0, advapi32.KEY_READ); hKey != 0 {
		eof := false
		for i := uint32(0); !eof; i++ {
			valueName, valueData, valueType := advapi32.RegEnumValue(hKey, i)
			switch valueType {
			case advapi32.REG_SZ:
				fontFileName := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(&valueData[0]))[:len(valueData)/2])
				fontFacesName := strings.Split(valueName, AMPERSAND_SUFFIX)
				for _, fontFaceName := range fontFacesName {
					for {
						isTrimmed := false
						for _, suffix := range []string{TRUE_TYPE_SUFFIX, SPACE_CHARACTER} {
							if strings.HasSuffix(fontFaceName, suffix) {
								fontFaceName = strings.TrimSuffix(fontFaceName, suffix)
								isTrimmed = true
							}
						}
						if !isTrimmed {
							break
						}
					}
					fontFaceName = strings.TrimSpace(fontFaceName)
					result[fontFaceName] = fontsDirectoryPath + fontFileName
				}
			case advapi32.REG_UNKNOWN:
				eof = true
			}
		}
		advapi32.RegCloseKey(hKey)
	}
	return result
}

func GetFontFilePath(fontName string) string {
	result := ""
	for fontFaceName, fontFilePath := range EnumFontFamiliesEx() {
		if strings.EqualFold(fontName, fontFaceName) {
			if _, err := os.Stat(fontFilePath); err == nil {
				result = fontFilePath
			}
			break
		}
	}
	return result
}

func GetFontFilePathFallback(faceName string, defaultFontName string, fontNamePrefix string) (string, string) {
	fontFaceName := faceName
	if fontFaceName == "" {
		fontFaceName = defaultFontName
	}
	fontFilePath := GetFontFilePath(fontFaceName)
	if fontFilePath == "" {
		fontFaceName = defaultFontName
		fontFilePath = GetFontFilePath(fontFaceName)
	}
	if fontFilePath == "" {
		for _fontFaceName, _fontFilePath := range EnumFontFamiliesEx() {
			if strings.HasPrefix(_fontFaceName, fontNamePrefix) {
				fontFaceName = _fontFaceName
				fontFilePath = _fontFilePath
				break
			}
		}
	}
	return fontFaceName, fontFilePath
}

func GetScreenDpiX() int {
	result := DEFAULT_DPI
	if hDC := C.GetDC(C.HWND(nil)); hDC != C.HDC(nil) {
		result = int(C.GetDeviceCaps(hDC, C.LOGPIXELSX))
		C.ReleaseDC(C.HWND(nil), hDC)
	}
	return result
}

func GetScreenDpiY() int {
	result := DEFAULT_DPI
	if hDC := C.GetDC(C.HWND(nil)); hDC != C.HDC(nil) {
		result = int(C.GetDeviceCaps(hDC, C.LOGPIXELSY))
		C.ReleaseDC(C.HWND(nil), hDC)
	}
	return result
}

func MmToPixelX(mm float64) float64 {
	dpiX := GetScreenDpiX()
	return mm * float64(dpiX) / MM_PER_INCH
}

func MmToPixelY(mm float64) float64 {
	dpiY := GetScreenDpiY()
	return mm * float64(dpiY) / MM_PER_INCH
}

//goland:noinspection GoUnusedExportedFunction
func PrintFontsInfo() {
	for fontFaceName, fontFilePath := range EnumFontFamiliesEx() {
		fmt.Printf("FontFaceName: %s, FontFilePath: %s\n", fontFaceName, fontFilePath)
	}
	for _, fontInfo := range EnumFontFamilies() {
		fmt.Printf("FontFaceName: %s, FontFilePath: %s\n", fontInfo.FaceName, fontInfo.FilePath)
	}
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_WINGDI, logger.SKIP_STACK_FRAMES_BASE)
}
