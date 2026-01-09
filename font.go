package boost

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/image/font/sfnt"
)

type FONT struct{}

func NewFont() *FONT {
	ret := &FONT{}
	return ret
}

func (f *FONT) CheckFontCharacters(inputString string, fontName string) ([]rune, string, error) {
	retUndisplayable := []rune{}
	retDisplayable := ""
	retErr := error(nil)

	exists, checkErr := f.IsExists(fontName)
	if checkErr != nil {
		retErr = fmt.Errorf("check font existence failed: %w", checkErr)
		return retUndisplayable, retDisplayable, retErr
	}
	if !exists {
		retErr = fmt.Errorf("font %s not found in system", fontName)
		return retUndisplayable, retDisplayable, retErr
	}

	undisplayableTemp := []rune{}
	displayableBuilder := strings.Builder{}
	for _, char := range inputString {
		isValid, checkErr := f.CheckCharacter(fontName, char)
		if checkErr != nil {
			retErr = fmt.Errorf("check character %q (code:%d) failed: %w", char, char, checkErr)
			break
		}
		if !isValid {
			undisplayableTemp = append(undisplayableTemp, char)
		} else {
			displayableBuilder.WriteRune(char)
		}
	}

	uniqueUndisplayable := make(map[rune]struct{})
	for _, c := range undisplayableTemp {
		uniqueUndisplayable[c] = struct{}{}
	}
	retUndisplayable = make([]rune, 0, len(uniqueUndisplayable))
	for c := range uniqueUndisplayable {
		retUndisplayable = append(retUndisplayable, c)
	}
	sort.Slice(retUndisplayable, func(i, j int) bool {
		return retUndisplayable[i] < retUndisplayable[j]
	})
	retDisplayable = displayableBuilder.String()

	return retUndisplayable, retDisplayable, retErr
}

func (f *FONT) GetSystemFonts() ([]string, error) {
	retFontNames := []string{}
	retErr := error(nil)

	switch runtime.GOOS {
	case "windows":
		retFontNames, retErr = f.getWindowsFonts()
	case "darwin":
		retFontNames, retErr = f.getMacFonts()
	case "linux":
		retFontNames, retErr = f.getLinuxFonts()
	default:
		retErr = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if retErr != nil {
		retErr = fmt.Errorf("get system fonts failed: %w", retErr)
		return retFontNames, retErr
	}

	uniqueFonts := make(map[string]bool)
	for _, name := range retFontNames {
		if name != "" {
			uniqueFonts[name] = true
		}
	}
	retFontNames = make([]string, 0, len(uniqueFonts))
	for name := range uniqueFonts {
		retFontNames = append(retFontNames, name)
	}
	sort.Strings(retFontNames)

	return retFontNames, retErr
}

func (f *FONT) IsExists(fontName string) (bool, error) {
	retExists := false
	retErr := error(nil)

	fontFilePath, err := f.findFont(fontName)
	if err != nil {
		retErr = fmt.Errorf("find font %s failed: %w", fontName, err)
		retExists = false
		return retExists, retErr
	}
	if fontFilePath == "" {
		retExists = false
		retErr = nil
		return retExists, retErr
	}

	_, statErr := os.Stat(fontFilePath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			retExists = false
			retErr = nil
		} else {
			retErr = fmt.Errorf("stat font file %s failed: %w", fontFilePath, statErr)
			retExists = false
		}
		return retExists, retErr
	}

	retExists = true
	retErr = nil
	return retExists, retErr
}

func (f *FONT) CheckCharacter(fontName string, char rune) (bool, error) {
	retSupported := false
	retErr := error(nil)

	fontFilePath, err := f.findFont(fontName)
	if err != nil {
		retErr = fmt.Errorf("find font %s failed: %w", fontName, err)
		retSupported = false
		return retSupported, retErr
	}
	if fontFilePath == "" {
		retErr = fmt.Errorf("font %s not found", fontName)
		retSupported = false
		return retSupported, retErr
	}

	// 使用 sfnt 包解析字体文件
	fontData, err := os.ReadFile(fontFilePath)
	if err != nil {
		retErr = fmt.Errorf("read font file failed: %w", err)
		retSupported = false
		return retSupported, retErr
	}

	font, err := sfnt.Parse(fontData)
	if err != nil {
		retErr = fmt.Errorf("parse font failed: %w", err)
		retSupported = false
		return retSupported, retErr
	}

	// 使用 GlyphIndex 方法检测字符是否在字体中
	index, err := font.GlyphIndex(nil, char)
	if err != nil {
		retErr = fmt.Errorf("check glyph index failed: %w", err)
		retSupported = false
		return retSupported, retErr
	}

	// 如果 glyph index 为 0，表示字符不被支持
	retSupported = index != 0
	retErr = nil

	return retSupported, retErr
}



func (f *FONT) findFont(fontName string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return f.findFontWindows(fontName)
	case "darwin":
		return f.findFontMac(fontName)
	case "linux":
		return f.findFontLinux(fontName)
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func (f *FONT) findFontWindows(fontName string) (string, error) {
	fontDir := "C:\\Windows\\Fonts"
	_, err := os.Stat(fontDir)
	if err != nil {
		return "", fmt.Errorf("font directory not found: %w", err)
	}

	var foundPath string
	err = filepath.Walk(fontDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".ttf" || ext == ".otf" || ext == ".ttc" {
				// 匹配字体名（忽略大小写）
				fileName := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ext)
				if strings.Contains(fileName, strings.ToLower(fontName)) {
					foundPath = path
					return fmt.Errorf("found") // 终止遍历
				}
			}
		}
		return nil
	})

	if err != nil && err.Error() == "found" {
		return foundPath, nil
	}
	return "", fmt.Errorf("font %s not found in Windows fonts directory", fontName)
}

func (f *FONT) findFontMac(fontName string) (string, error) {
	fontDirs := []string{
		"/Library/Fonts",
		"/System/Library/Fonts",
		filepath.Join(os.Getenv("HOME"), "Library/Fonts"),
	}

	for _, dir := range fontDirs {
		_, err := os.Stat(dir)
		if err != nil {
			continue
		}

		var foundPath string
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".ttf" || ext == ".otf" || ext == ".ttc" || ext == ".dfont" {
					fileName := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ext)
					if strings.Contains(fileName, strings.ToLower(fontName)) {
						foundPath = path
						return fmt.Errorf("found")
					}
				}
			}
			return nil
		})

		if err != nil && err.Error() == "found" {
			return foundPath, nil
		}
	}

	return "", fmt.Errorf("font %s not found in Mac fonts directories", fontName)
}

func (f *FONT) findFontLinux(fontName string) (string, error) {
	fontDirs := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		filepath.Join(os.Getenv("HOME"), ".local/share/fonts"),
		filepath.Join(os.Getenv("HOME"), ".fonts"),
	}

	for _, dir := range fontDirs {
		_, err := os.Stat(dir)
		if err != nil {
			continue
		}

		var foundPath string
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".ttf" || ext == ".otf" || ext == ".ttc" {
					fileName := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ext)
					if strings.Contains(fileName, strings.ToLower(fontName)) {
						foundPath = path
						return fmt.Errorf("found")
					}
				}
			}
			return nil
		})

		if err != nil && err.Error() == "found" {
			return foundPath, nil
		}
	}

	return "", fmt.Errorf("font %s not found in Linux fonts directories", fontName)
}

func (f *FONT) getWindowsFonts() ([]string, error) {
	cmd := exec.Command("powershell", "-Command", `
		Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts' | 
		Select-Object -ExpandProperty PSObject.Properties | 
		Where-Object {$_.Value -match '\.(ttf|otf|ttc)'} | 
		Select-Object -ExpandProperty Name
	`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("execute powershell command failed: %w", err)
	}

	fonts := []string{}
	lines := strings.Split(string(output), "\n")
	unique := make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !unique[line] {
			unique[line] = true
			fonts = append(fonts, line)
		}
	}

	return fonts, nil
}

func (f *FONT) getMacFonts() ([]string, error) {
	cmd := exec.Command("system_profiler", "SPFontsDataType", "-json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("execute system_profiler failed: %w", err)
	}

	fonts := []string{}
	unique := make(map[string]bool)
	outputStr := string(output)
	start := 0
	for {
		nameStart := strings.Index(outputStr[start:], `"name": "`)
		if nameStart == -1 {
			break
		}
		nameStart += start + 8
		nameEnd := strings.Index(outputStr[nameStart:], `"`)
		if nameEnd == -1 {
			break
		}
		nameEnd += nameStart
		fontName := outputStr[nameStart:nameEnd]
		if fontName != "" && !unique[fontName] {
			unique[fontName] = true
			fonts = append(fonts, fontName)
		}
		start = nameEnd + 1
	}

	return fonts, nil
}

func (f *FONT) getLinuxFonts() ([]string, error) {
	cmd := exec.Command("fc-list", ":format=%{family}\n")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("execute fc-list failed: %w", err)
	}

	fonts := []string{}
	unique := make(map[string]bool)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !unique[line] {
			unique[line] = true
			fonts = append(fonts, line)
		}
	}

	return fonts, nil
}
