//go:build darwin

package wingdi

type FONTINFO struct {
	FaceName string
	FileName string
}

func CString(s string) *uint16 {
	return nil
}

func FreeCString(ptr *uint16) {
}

func GoString(lpwsz interface{}) string {
	return ""
}

func EnumFontFamilies() []FONTINFO {
	return nil
}

func EnumFontFamiliesEx() map[string]string {
	return nil
}

func GetFontFilePath(fontName string) string {
	return ""
}

func GetFontFilePathFallback(faceName string, defaultFontName string, fontNamePrefix string) (string, string) {
	return "", ""
}

func GetScreenDpiX() int {
	return 72
}

func GetScreenDpiY() int {
	return 72
}

func MmToPixelX(mm float64) float64 {
	return mm * 72 / 25.4
}

func MmToPixelY(mm float64) float64 {
	return mm * 72 / 25.4
}

func PrintFontsInfo() {
}
