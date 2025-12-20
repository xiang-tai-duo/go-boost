// Package boost
// File:        font.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/font/font.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: BOOST provides cross-platform font management functions including system font scanning, character validation, and glyph export.
// --------------------------------------------------------------------------------

package boost

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/font/sfnt"
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
type (
	CMAP_FONT_INFO struct {
		FontName       string
		GlyphToUnicode map[sfnt.GlyphIndex]rune
	}

	FONT struct {
		caseInsensitive bool
		mutex           sync.RWMutex
		systemFonts     []FONT_INFO
	}

	FONT_INFO struct {
		FilePath string
		Font     *sfnt.Font
		Language uint16
		Name     string
		Names    map[string]string
	}

	FONT_KEY_PAIR struct {
		key   string
		value string
	}

	FONT_TABLE_INFO struct {
		checkSum uint32
		length   uint32
		offset   uint32
		tag      string
	}

	FONT_VARIABLE_INSTANCE struct {
		coordinates     []float32
		subFamilyNameId uint16
	}
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage,SpellCheckingInspection
const (
	FIRST_VARIANT_NAME_ID                  = 258
	FVAR_AXIS_COUNT_LENGTH                 = 2
	FVAR_AXIS_COUNT_OFFSET                 = 8
	FVAR_AXIS_SIZE_LENGTH                  = 2
	FVAR_AXIS_SIZE_OFFSET                  = 10
	FVAR_INSTANCE_COORDINATES_OFFSET       = 4
	FVAR_INSTANCE_COUNT_LENGTH             = 2
	FVAR_INSTANCE_COUNT_OFFSET             = 12
	FVAR_INSTANCE_FLAGS_LENGTH             = 2
	FVAR_INSTANCE_FLAGS_OFFSET             = 2
	FVAR_INSTANCE_SIZE_LENGTH              = 2
	FVAR_INSTANCE_SIZE_OFFSET              = 14
	FVAR_INSTANCE_SUBFAMILY_NAME_ID_LENGTH = 2
	FVAR_INSTANCE_SUBFAMILY_NAME_ID_OFFSET = 0
	FVAR_OFFSET_TO_DATA_OFFSET             = 4
	FVAR_RESERVED_OFFSET                   = 6
	FVAR_VERSION_LENGTH                    = 4
	FVAR_VERSION_OFFSET                    = 0
	HEADER_ENTRY_SELECTOR_LENGTH           = 2
	HEADER_ENTRY_SELECTOR_OFFSET           = 8
	HEADER_MAGIC_LENGTH                    = 4
	HEADER_MAGIC_OFFSET                    = 0
	HEADER_RANGE_SHIFT_LENGTH              = 2
	HEADER_RANGE_SHIFT_OFFSET              = 10
	HEADER_SEARCH_RANGE_LENGTH             = 2
	HEADER_SEARCH_RANGE_OFFSET             = 6
	HEADER_TABLE_COUNT_LENGTH              = 2
	HEADER_TABLE_COUNT_OFFSET              = 4
	LAST_VARIANT_NAME_ID                   = 271
	LANGUAGE_ENGLISH                       = 0x409
	NAME_HEADER_LENGTH                     = 6
	NAME_ID_LAST                           = int(sfnt.NameIDVariationsPostScriptPrefix)
	NAME_RECORD_ENCODING_ID_LENGTH         = 2
	NAME_RECORD_ENCODING_ID_OFFSET         = 2
	NAME_RECORD_LANGUAGE_ID_LENGTH         = 2
	NAME_RECORD_LANGUAGE_ID_OFFSET         = 4
	NAME_RECORD_LENGTH                     = 12
	NAME_RECORD_LENGTH_LENGTH              = 2
	NAME_RECORD_LENGTH_OFFSET              = 8
	NAME_RECORD_NAME_ID_LENGTH             = 2
	NAME_RECORD_NAME_ID_OFFSET             = 6
	NAME_RECORD_OFFSET_LENGTH              = 2
	NAME_RECORD_OFFSET_OFFSET              = 10
	NAME_RECORD_PLATFORM_ID_LENGTH         = 2
	NAME_RECORD_PLATFORM_ID_OFFSET         = 0
	PLATFORM_ID_CUSTOM                     = 4
	PLATFORM_ID_ISO                        = 2
	PLATFORM_ID_MACINTOSH                  = 1
	PLATFORM_ID_MICROSOFT                  = 3
	PLATFORM_ID_UNICODE                    = 0
	SIZE_DWROD                             = 4
	TABLE_ENTRY_CHECKSUM_LENGTH            = 4
	TABLE_ENTRY_CHECKSUM_OFFSET            = 4
	TABLE_ENTRY_LENGTH                     = 16
	TABLE_ENTRY_LENGTH_LENGTH              = 4
	TABLE_ENTRY_LENGTH_OFFSET              = 12
	TABLE_ENTRY_OFFSET_LENGTH              = 4
	TABLE_ENTRY_OFFSET_OFFSET              = 8
	TABLE_ENTRY_TAG_LENGTH                 = 4
	TABLE_ENTRY_TAG_OFFSET                 = 0
	TTC_FONT_OFFSET_TABLE_OFFSET           = 12
	TTC_HEADER_LENGTH                      = 12
	TTC_MAGIC_NUMBER                       = 0x74746366
	TTC_NUM_FONTS_LENGTH                   = 4
	TTC_NUM_FONTS_OFFSET                   = 8
	TTC_VERSION_LENGTH                     = 4
	TTC_VERSION_OFFSET                     = 4
	TTF_HEADER_LENGTH                      = 12
	TTF_MAGIC_NUMBER_1                     = 0x00010000
	TTF_MAGIC_NUMBER_2                     = 0x4F54544F
	UTF16_CHAR_SIZE                        = 2
	UTF16_NULL_CHAR                        = 0
	TABLE_TAG_FVAR                         = "fvar"
	TABLE_TAG_NAME                         = "name"
)

//goland:noinspection GoSnakeCaseUsage
var (
	NAME_ID = []sfnt.NameID{
		sfnt.NameIDCopyright,
		sfnt.NameIDFamily,
		sfnt.NameIDSubfamily,
		sfnt.NameIDUniqueIdentifier,
		sfnt.NameIDFull,
		sfnt.NameIDVersion,
		sfnt.NameIDPostScript,
		sfnt.NameIDTrademark,
		sfnt.NameIDManufacturer,
		sfnt.NameIDDesigner,
		sfnt.NameIDDescription,
		sfnt.NameIDVendorURL,
		sfnt.NameIDDesignerURL,
		sfnt.NameIDLicense,
		sfnt.NameIDLicenseURL,
		sfnt.NameIDTypographicFamily,
		sfnt.NameIDTypographicSubfamily,
		sfnt.NameIDCompatibleFull,
		sfnt.NameIDSampleText,
		sfnt.NameIDPostScriptCID,
		sfnt.NameIDWWSFamily,
		sfnt.NameIDWWSSubfamily,
		sfnt.NameIDLightBackgroundPalette,
		sfnt.NameIDDarkBackgroundPalette,
		sfnt.NameIDVariationsPostScriptPrefix,
	}
)

//goland:noinspection GoUnusedExportedFunction
func NewFont() *FONT {
	result := &FONT{
		systemFonts:     []FONT_INFO{},
		caseInsensitive: false,
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func (f *FONT) ExportGlyphsToFile(fontName string) error {
	err := error(nil)
	var fontInfo *FONT_INFO
	if fontInfo, err = f.getSystemFont(fontName); err == nil {
		fileName := fontName + ".txt"
		var file *os.File
		if file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
			defer file.Close()
			fmt.Fprintf(file, "Font: %s\n", fontName)
			fmt.Fprintf(file, "File path: %s\n\n", fontInfo.FilePath)
			fmt.Fprintf(file, "Total Glyphs: %d\n\n", fontInfo.Font.NumGlyphs())
			fmt.Fprintln(file, "Glyphs:")
			numberOfGlyphs := fontInfo.Font.NumGlyphs()
			for index := 0; index < numberOfGlyphs; index++ {
				glyphIndex := sfnt.GlyphIndex(index)
				characterName := f.getGlyphCharacter(fontInfo.Font, glyphIndex)
				var name string
				if characterName != "" {
					runeValue := []rune(characterName)[0]
					name = fmt.Sprintf("0x%04X", runeValue)
				} else {
					name = fmt.Sprintf("0x%X", 0x1D18C)
				}
				var displayCharacter rune
				if code, parseError := strconv.ParseUint(name[2:], 16, 32); parseError == nil {
					displayCharacter = rune(code)
				} else {
					displayCharacter = rune(0x1D18C)
				}
				fmt.Fprintf(file, "Glyph Index: %d, Name: %c\n", index, displayCharacter)
			}
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func (f *FONT) ExportSystemFontsToFile(fileName string) error {
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		var file *os.File
		if file, err = os.Create(fileName); err == nil {
			defer file.Close()
			f.mutex.RLock()
			defer f.mutex.RUnlock()
			for _, font := range f.systemFonts {
				for key, value := range font.Names {
					fmt.Fprintf(file, "key：%s，value: %s, Language: %d, FilePath：%s\n", key, value, font.Language, font.FilePath)
				}
			}
		}
	}
	return err
}

func (f *FONT) GetCaseInsensitive() bool {
	result := f.caseInsensitive
	return result
}

func (f *FONT) IsExists(fontName string) (bool, error) {
	result := false
	err := error(nil)
	var fontInfo *FONT_INFO
	if fontInfo, err = f.getSystemFont(fontName); err == nil {
		result = fontInfo != nil
	}
	return result, err
}

func (f *FONT) SetCaseInsensitive(enabled bool) {
	f.caseInsensitive = enabled
}

func (f *FONT) ValidateCharacters(inputString string, fontName string) ([]rune, []rune, error) {
	invisibleCharacters := make([]rune, 0)
	visibleCharacters := make([]rune, 0)
	err := error(nil)
	var fontInfo *FONT_INFO
	if fontInfo, err = f.getSystemFont(fontName); err == nil {
		tempInvisibleCharacters := make([]rune, 0)
		tempVisibleCharacters := make([]rune, 0)
		for _, character := range inputString {
			glyphIndex := sfnt.GlyphIndex(0)
			if glyphIndex, err = fontInfo.Font.GlyphIndex(nil, character); err == nil {
				if glyphIndex == 0 {
					tempInvisibleCharacters = append(tempInvisibleCharacters, character)
				} else {
					tempVisibleCharacters = append(tempVisibleCharacters, character)
				}
			} else {
				err = fmt.Errorf("check character %q (code:%d) failed: %w", character, character, err)
				break
			}
		}
		if err == nil {
			uniqueUndisplayable := make(map[rune]struct{})
			for _, character := range tempInvisibleCharacters {
				uniqueUndisplayable[character] = struct{}{}
			}
			invisibleCharacters = make([]rune, 0, len(uniqueUndisplayable))
			for character := range uniqueUndisplayable {
				invisibleCharacters = append(invisibleCharacters, character)
			}
			sort.Slice(invisibleCharacters, func(i, j int) bool {
				return invisibleCharacters[i] < invisibleCharacters[j]
			})
			visibleCharacters = tempVisibleCharacters
		}
	}
	return invisibleCharacters, visibleCharacters, err
}

func (f *FONT) decodeUTF16(data []byte) string {
	result := ""
	if len(data)%UTF16_CHAR_SIZE != 0 {
		data = append(data, 0)
	}
	for index := 0; index < len(data); index += UTF16_CHAR_SIZE {
		character := binary.BigEndian.Uint16(data[index : index+UTF16_CHAR_SIZE])
		if character != UTF16_NULL_CHAR {
			result += string(rune(character))
		}
	}
	return result
}

func (f *FONT) getFontInfo(fontName string) (FONT_INFO, error) {
	result := FONT_INFO{}
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		found := false
		f.mutex.RLock()
		for _, fontCollection := range f.systemFonts {
			if f.caseInsensitive {
				if strings.EqualFold(fontCollection.Name, fontName) {
					result = fontCollection
					found = true
					break
				}
			} else {
				if fontCollection.Name == fontName {
					result = fontCollection
					found = true
					break
				}
			}
		}
		f.mutex.RUnlock()
		if !found {
			err = fmt.Errorf("font %s not found in system", fontName)
		}
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func (f *FONT) getGlyphCharacter(font *sfnt.Font, glyphIndex sfnt.GlyphIndex) string {
	result := ""
	var fontPath string
	f.mutex.RLock()
	for _, fontInformation := range f.systemFonts {
		if fontInformation.Font == font {
			fontPath = fontInformation.FilePath
			break
		}
	}
	f.mutex.RUnlock()
	if fontPath != "" {
		if cmapInformations, err := ParseCMapTable(fontPath); err == nil {
			for _, cmapInformation := range cmapInformations {
				if unicodeCharacter, ok := cmapInformation.GlyphToUnicode[glyphIndex]; ok {
					result = string(unicodeCharacter)
					break
				}
			}
		}
	}

	if result == "" {
		buffer := &sfnt.Buffer{}
		for runeValue := rune(0); runeValue <= 0xFFFF; runeValue++ {
			if index, err := font.GlyphIndex(buffer, runeValue); err == nil && index == glyphIndex {
				result = string(runeValue)
				break
			}
		}
	}
	return result
}

//goland:noinspection DuplicatedCode
func (f *FONT) getMultilanguageNameTable(file *os.File, fileSize int64, fontOffset int64, header []byte) ([]FONT_KEY_PAIR, error) {
	result := make([]FONT_KEY_PAIR, 0)
	err := error(nil)
	numberOfTables := binary.BigEndian.Uint16(header[HEADER_TABLE_COUNT_OFFSET : HEADER_TABLE_COUNT_OFFSET+HEADER_TABLE_COUNT_LENGTH])
	tableDirectoryOffset := fontOffset + int64(TTF_HEADER_LENGTH)
	if tableDirectoryOffset+int64(numberOfTables*TABLE_ENTRY_LENGTH) <= fileSize {
		if _, err = file.Seek(tableDirectoryOffset, 0); err == nil {
			tables := make(map[string]FONT_TABLE_INFO)
			for index := 0; index < int(numberOfTables); index++ {
				tableEntry := make([]byte, TABLE_ENTRY_LENGTH)
				if _, err = file.Read(tableEntry); err == nil {
					tag := string(tableEntry[TABLE_ENTRY_TAG_OFFSET : TABLE_ENTRY_TAG_OFFSET+TABLE_ENTRY_TAG_LENGTH])
					checkSum := binary.BigEndian.Uint32(tableEntry[TABLE_ENTRY_CHECKSUM_OFFSET : TABLE_ENTRY_CHECKSUM_OFFSET+TABLE_ENTRY_CHECKSUM_LENGTH])
					offset := binary.BigEndian.Uint32(tableEntry[TABLE_ENTRY_OFFSET_OFFSET : TABLE_ENTRY_OFFSET_OFFSET+TABLE_ENTRY_OFFSET_LENGTH])
					length := binary.BigEndian.Uint32(tableEntry[TABLE_ENTRY_LENGTH_OFFSET : TABLE_ENTRY_LENGTH_OFFSET+TABLE_ENTRY_LENGTH_LENGTH])
					tables[tag] = FONT_TABLE_INFO{
						tag:      tag,
						checkSum: checkSum,
						offset:   offset,
						length:   length,
					}
				} else {
					break
				}
			}
			if err == nil {
				if nameInformation, ok := tables[TABLE_TAG_NAME]; ok {
					if _, err = file.Seek(int64(nameInformation.offset), 0); err == nil {
						nameHeader := make([]byte, NAME_HEADER_LENGTH)
						if _, err = file.Read(nameHeader); err == nil {
							nameCount := binary.BigEndian.Uint16(nameHeader[2:4])
							stringOffset := binary.BigEndian.Uint16(nameHeader[4:6])
							for index := 0; index < int(nameCount); index++ {
								nameRecord := make([]byte, NAME_RECORD_LENGTH)
								if _, err = file.Read(nameRecord); err == nil {
									platformID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_PLATFORM_ID_OFFSET : NAME_RECORD_PLATFORM_ID_OFFSET+NAME_RECORD_PLATFORM_ID_LENGTH])
									languageID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_LANGUAGE_ID_OFFSET : NAME_RECORD_LANGUAGE_ID_OFFSET+NAME_RECORD_LANGUAGE_ID_LENGTH])
									nameID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_NAME_ID_OFFSET : NAME_RECORD_NAME_ID_OFFSET+NAME_RECORD_NAME_ID_LENGTH])
									length := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_LENGTH_OFFSET : NAME_RECORD_LENGTH_OFFSET+NAME_RECORD_LENGTH_LENGTH])
									offset := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_OFFSET_OFFSET : NAME_RECORD_OFFSET_OFFSET+NAME_RECORD_OFFSET_LENGTH])
									if platformID == PLATFORM_ID_MICROSOFT {
										var currentPosition int64
										if currentPosition, err = file.Seek(0, 1); err == nil {
											if _, err = file.Seek(int64(nameInformation.offset)+int64(stringOffset)+int64(offset), 0); err == nil {
												stringData := make([]byte, length)
												if _, err = file.Read(stringData); err == nil {
													stringValue := f.decodeUTF16(stringData)
													key := fmt.Sprintf("%d:0x%x", int(nameID), languageID)
													result = append(result, FONT_KEY_PAIR{
														key:   key,
														value: stringValue,
													})
												}
												if _, err = file.Seek(currentPosition, 0); err != nil {
													break
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	for index, element := range result {
		parts := strings.Split(element.key, ":")
		if len(parts) == 2 {
			if identifier, conversionError := strconv.Atoi(parts[0]); conversionError == nil && identifier > NAME_ID_LAST {
				key := fmt.Sprintf("1:%s", parts[1])
				for _, innerElement := range result {
					if innerElement.key == key {
						result[index].value = fmt.Sprintf("%s %s", innerElement.value, element.value)
						break
					}
				}
			}
		}
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection,GoUnhandledErrorResult
func (f *FONT) getMultilanguageNameTables(path string) ([]FONT_KEY_PAIR, error) {
	result := make([]FONT_KEY_PAIR, 0)
	err := error(nil)
	var file *os.File
	if file, err = os.Open(path); err == nil {
		defer file.Close()
		var fileInformation os.FileInfo
		if fileInformation, err = file.Stat(); err == nil {
			fileSize := fileInformation.Size()
			if fileSize >= TTF_HEADER_LENGTH {
				header := make([]byte, TTF_HEADER_LENGTH)
				if _, err = file.Read(header); err == nil {
					magicNumber := binary.BigEndian.Uint32(header[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
					if magicNumber == TTC_MAGIC_NUMBER {
						numberOfFonts := binary.BigEndian.Uint32(header[TTC_NUM_FONTS_OFFSET : TTC_NUM_FONTS_OFFSET+TTC_NUM_FONTS_LENGTH])
						if numberOfFonts > 0 {
							fontOffsets := make([]uint32, numberOfFonts)
							for index := 0; index < int(numberOfFonts); index++ {
								offsetPosition := TTC_FONT_OFFSET_TABLE_OFFSET + index*SIZE_DWROD
								if int64(offsetPosition+SIZE_DWROD) <= fileSize {
									if _, err = file.Seek(int64(offsetPosition), 0); err == nil {
										offsetData := make([]byte, SIZE_DWROD)
										if _, err = file.Read(offsetData); err == nil {
											fontOffsets[index] = binary.BigEndian.Uint32(offsetData)
										}
									}
								}
							}
							if len(fontOffsets) > 0 {
								for _, fontOffset := range fontOffsets {
									if int64(fontOffset)+TTF_HEADER_LENGTH <= fileSize {
										if _, err = file.Seek(int64(fontOffset), 0); err == nil {
											ttfHeader := make([]byte, TTF_HEADER_LENGTH)
											if _, err = file.Read(ttfHeader); err == nil {
												ttfMagicNumber := binary.BigEndian.Uint32(ttfHeader[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
												if ttfMagicNumber == TTF_MAGIC_NUMBER_1 || ttfMagicNumber == TTF_MAGIC_NUMBER_2 {
													var fontNameTable []FONT_KEY_PAIR
													if fontNameTable, err = f.getMultilanguageNameTable(file, fileSize, int64(fontOffset), ttfHeader); err == nil {
														result = append(result, fontNameTable...)
													}
												}
											}
										}
									}
								}
							}
						}
					} else if magicNumber == TTF_MAGIC_NUMBER_1 || magicNumber == TTF_MAGIC_NUMBER_2 {
						var fontNameTable []FONT_KEY_PAIR
						if fontNameTable, err = f.getMultilanguageNameTable(file, fileSize, 0, header); err == nil {
							result = append(result, fontNameTable...)
						}
					}
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		keyOne := result[i].key
		keyTwo := result[j].key
		for len(keyOne) < 8 {
			keyOne = "0" + keyOne
		}
		for len(keyTwo) < 8 {
			keyTwo = "0" + keyTwo
		}
		return keyOne < keyTwo
	})
	return result, err
}

func (f *FONT) getSystemFont(fontName string) (*FONT_INFO, error) {
	result := (*FONT_INFO)(nil)
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		f.mutex.RLock()
		for _, systemFont := range f.systemFonts {
			systemFontNames := make([]string, 0)
			systemFontNames = append(systemFontNames, systemFont.Name)
			for _, systemFontName := range systemFont.Names {
				systemFontNames = append(systemFontNames, systemFontName)
			}
			for _, systemFontName := range systemFontNames {
				if f.caseInsensitive {
					if strings.EqualFold(systemFontName, fontName) {
						result = &systemFont
						break
					}
				} else {
					if systemFontName == fontName {
						result = &systemFont
						break
					}
				}
			}
			if result != nil {
				break
			}
		}
		f.mutex.RUnlock()
	} else {
		err = fmt.Errorf("get system fonts failed: %w", err)
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func (f *FONT) getSystemFonts() error {
	err := error(nil)
	f.mutex.RLock()
	if len(f.systemFonts) > 0 {
		f.mutex.RUnlock()
	} else {
		f.mutex.RUnlock()
		var font []FONT_INFO
		switch runtime.GOOS {
		case "windows":
			font = f.getSystemFontsWindows()
		case "darwin":
			font = f.getSystemFontsMac()
		case "linux":
			font = f.getSystemFontsLinux()
		default:
			err = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}
		if err == nil {
			f.mutex.Lock()
			f.systemFonts = font
			f.mutex.Unlock()
		}
	}
	return err
}

func (f *FONT) getSystemFontsLinux() []FONT_INFO {
	fontDirectories := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		filepath.Join(os.Getenv("HOME"), ".local/share/fonts"),
		filepath.Join(os.Getenv("HOME"), ".fonts"),
	}
	result := f.init(fontDirectories)
	return result
}

func (f *FONT) getSystemFontsMac() []FONT_INFO {
	fontDirectories := []string{
		"/Library/Fonts",
		"/System/Library/Fonts",
		filepath.Join(os.Getenv("HOME"), "Library/Fonts"),
	}
	result := f.init(fontDirectories)
	return result
}

func (f *FONT) getSystemFontsWindows() []FONT_INFO {
	fontDirectories := []string{
		"C:\\Windows\\Fonts",
		filepath.Join(os.Getenv("USERPROFILE"), "AppData\\Local\\Microsoft\\Windows\\Fonts"),
	}
	result := f.init(fontDirectories)
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnhandledErrorResult,DuplicatedCode
func (f *FONT) getTrueTypeFontVariableNames(fontPath string) ([]string, error) {
	result := make([]string, 0)
	err := error(nil)
	var file *os.File
	if file, err = os.Open(fontPath); err == nil {
		defer file.Close()
		var fileInformation os.FileInfo
		if fileInformation, err = file.Stat(); err == nil {
			fileSize := fileInformation.Size()
			if fileSize >= TTF_HEADER_LENGTH {
				header := make([]byte, TTF_HEADER_LENGTH)
				if _, err = file.Read(header); err == nil {
					magicNumber := binary.BigEndian.Uint32(header[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
					if magicNumber == TTF_MAGIC_NUMBER_1 || magicNumber == TTF_MAGIC_NUMBER_2 {
						numberOfTables := binary.BigEndian.Uint16(header[HEADER_TABLE_COUNT_OFFSET : HEADER_TABLE_COUNT_OFFSET+HEADER_TABLE_COUNT_LENGTH])
						tableDirectoryOffset := int64(TTF_HEADER_LENGTH)
						if tableDirectoryOffset+int64(numberOfTables*TABLE_ENTRY_LENGTH) <= fileSize {
							if _, err = file.Seek(tableDirectoryOffset, 0); err == nil {
								tables := make(map[string]FONT_TABLE_INFO)
								for index := 0; index < int(numberOfTables); index++ {
									tableEntry := make([]byte, TABLE_ENTRY_LENGTH)
									if _, err = file.Read(tableEntry); err == nil {
										tag := string(tableEntry[TABLE_ENTRY_TAG_OFFSET : TABLE_ENTRY_TAG_OFFSET+TABLE_ENTRY_TAG_LENGTH])
										checkSum := binary.BigEndian.Uint32(tableEntry[TABLE_ENTRY_CHECKSUM_OFFSET : TABLE_ENTRY_CHECKSUM_OFFSET+TABLE_ENTRY_CHECKSUM_LENGTH])
										offset := binary.BigEndian.Uint32(tableEntry[TABLE_ENTRY_OFFSET_OFFSET : TABLE_ENTRY_OFFSET_OFFSET+TABLE_ENTRY_OFFSET_LENGTH])
										length := binary.BigEndian.Uint32(tableEntry[TABLE_ENTRY_LENGTH_OFFSET : TABLE_ENTRY_LENGTH_OFFSET+TABLE_ENTRY_LENGTH_LENGTH])
										tables[tag] = FONT_TABLE_INFO{
											tag:      tag,
											checkSum: checkSum,
											offset:   offset,
											length:   length,
										}
									} else {
										break
									}
								}
								if err == nil {
									if _, ok := tables[TABLE_TAG_FVAR]; ok {
										nameTable := make(map[uint16]string)
										if nameInformation, ok := tables[TABLE_TAG_NAME]; ok {
											if _, err = file.Seek(int64(nameInformation.offset), 0); err == nil {
												nameHeader := make([]byte, NAME_HEADER_LENGTH)
												if _, err = file.Read(nameHeader); err == nil {
													nameCount := binary.BigEndian.Uint16(nameHeader[2:4])
													stringOffset := binary.BigEndian.Uint16(nameHeader[4:6])
													for index := 0; index < int(nameCount); index++ {
														nameRecord := make([]byte, NAME_RECORD_LENGTH)
														if _, err = file.Read(nameRecord); err == nil {
															platformID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_PLATFORM_ID_OFFSET : NAME_RECORD_PLATFORM_ID_OFFSET+NAME_RECORD_PLATFORM_ID_LENGTH])
															nameID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_NAME_ID_OFFSET : NAME_RECORD_NAME_ID_OFFSET+NAME_RECORD_NAME_ID_LENGTH])
															length := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_LENGTH_OFFSET : NAME_RECORD_LENGTH_OFFSET+NAME_RECORD_LENGTH_LENGTH])
															offset := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_OFFSET_OFFSET : NAME_RECORD_OFFSET_OFFSET+NAME_RECORD_OFFSET_LENGTH])
															var currentPosition int64
															if currentPosition, err = file.Seek(0, 1); err == nil {
																if _, err = file.Seek(int64(nameInformation.offset)+int64(stringOffset)+int64(offset), 0); err == nil {
																	stringData := make([]byte, length)
																	if _, err = file.Read(stringData); err == nil {
																		var stringValue string
																		switch platformID {
																		case PLATFORM_ID_UNICODE:
																			stringValue = f.decodeUTF16(stringData)
																		case PLATFORM_ID_MACINTOSH:
																			var nameBuilder strings.Builder
																			for innerIndex := 0; innerIndex < int(length); innerIndex++ {
																				if innerIndex >= len(stringData) {
																					break
																				}
																				character := stringData[innerIndex]
																				if character == 0 {
																					break
																				}
																				nameBuilder.WriteRune(rune(character))
																			}
																			stringValue = nameBuilder.String()
																		case PLATFORM_ID_ISO:
																			var nameBuilder strings.Builder
																			for innerIndex := 0; innerIndex < int(length); innerIndex++ {
																				if innerIndex >= len(stringData) {
																					break
																				}
																				character := stringData[innerIndex]
																				if character == 0 {
																					break
																				}
																				nameBuilder.WriteRune(rune(character))
																			}
																			stringValue = nameBuilder.String()
																		case PLATFORM_ID_MICROSOFT:
																			stringValue = f.decodeUTF16(stringData)
																		case PLATFORM_ID_CUSTOM:
																			var nameBuilder strings.Builder
																			for innerIndex := 0; innerIndex < int(length); innerIndex++ {
																				if innerIndex >= len(stringData) {
																					break
																				}
																				character := stringData[innerIndex]
																				if character == 0 {
																					break
																				}
																				nameBuilder.WriteRune(rune(character))
																			}
																			stringValue = nameBuilder.String()
																		default:
																			var nameBuilder strings.Builder
																			for innerIndex := 0; innerIndex < int(length); innerIndex++ {
																				if innerIndex >= len(stringData) {
																					break
																				}
																				character := stringData[innerIndex]
																				if character == 0 {
																					break
																				}
																				nameBuilder.WriteRune(rune(character))
																			}
																			stringValue = nameBuilder.String()
																		}
																		nameTable[nameID] = stringValue
																		if _, err = file.Seek(currentPosition, 0); err != nil {
																			break
																		}
																	}
																}
															}
														} else {
															break
														}
													}
												}
											}
										}
										if err == nil {
											if fvarInformation, ok := tables[TABLE_TAG_FVAR]; ok {
												if _, err = file.Seek(int64(fvarInformation.offset), 0); err == nil {
													fvarBytes := make([]byte, fvarInformation.length)
													if _, err = file.Read(fvarBytes); err == nil {
														offsetToData := binary.BigEndian.Uint16(fvarBytes[FVAR_OFFSET_TO_DATA_OFFSET : FVAR_OFFSET_TO_DATA_OFFSET+FVAR_VERSION_LENGTH])
														axisCount := binary.BigEndian.Uint16(fvarBytes[FVAR_AXIS_COUNT_OFFSET : FVAR_AXIS_COUNT_OFFSET+FVAR_AXIS_COUNT_LENGTH])
														axisSize := binary.BigEndian.Uint16(fvarBytes[FVAR_AXIS_SIZE_OFFSET : FVAR_AXIS_SIZE_OFFSET+FVAR_AXIS_SIZE_LENGTH])
														instanceCount := binary.BigEndian.Uint16(fvarBytes[FVAR_INSTANCE_COUNT_OFFSET : FVAR_INSTANCE_COUNT_OFFSET+FVAR_INSTANCE_COUNT_LENGTH])
														instanceSize := binary.BigEndian.Uint16(fvarBytes[FVAR_INSTANCE_SIZE_OFFSET : FVAR_INSTANCE_SIZE_OFFSET+FVAR_INSTANCE_SIZE_LENGTH])
														instances := make([]FONT_VARIABLE_INSTANCE, 0)
														position := int(offsetToData)
														for index := 0; index < int(axisCount); index++ {
															position += int(axisSize)
														}
														for index := 0; index < int(instanceCount); index++ {
															if position+int(instanceSize) <= len(fvarBytes) {
																subFamilyNameId := binary.BigEndian.Uint16(fvarBytes[position : position+FVAR_INSTANCE_SUBFAMILY_NAME_ID_LENGTH])
																instances = append(instances, FONT_VARIABLE_INSTANCE{
																	subFamilyNameId: subFamilyNameId,
																	coordinates:     nil,
																})
																position += int(instanceSize)
															}
														}
														familyNames := make([]string, 0)
														if name, ok := nameTable[uint16(sfnt.NameIDFamily)]; ok {
															familyNames = append(familyNames, name)
														}
														if name, ok := nameTable[uint16(sfnt.NameIDTypographicFamily)]; ok {
															familyNames = append(familyNames, name)
														}
														result = make([]string, 0, len(instances))
														for _, familyName := range familyNames {
															for _, instance := range instances {
																if subFamilyName, ok := nameTable[instance.subFamilyNameId]; ok {
																	if familyName == "" {
																		result = append(result, subFamilyName)
																	} else {
																		result = append(result, fmt.Sprintf("%s %s", familyName, subFamilyName))
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func (f *FONT) init(directories []string) []FONT_INFO {
	result := make([]FONT_INFO, 0)
	for _, directory := range directories {
		if _, err := os.Stat(directory); err == nil {
			err = filepath.Walk(directory, func(path string, info os.FileInfo, walkError error) error {
				if walkError == nil {
					if !info.IsDir() {
						extension := strings.ToLower(filepath.Ext(path))
						if extension == ".ttf" || extension == ".otf" || extension == ".ttc" || extension == ".dfont" {
							var fileData []byte
							if fileData, walkError = os.ReadFile(path); walkError == nil {
								var collection *sfnt.Collection
								if collection, walkError = sfnt.ParseCollection(fileData); walkError == nil {
									multilanguageNames := make([]FONT_KEY_PAIR, 0)
									multilanguages := make(map[uint16]string)
									if multilanguageNames, walkError = f.getMultilanguageNameTables(path); walkError == nil {
										for _, element := range multilanguageNames {
											parts := strings.Split(element.key, ":")
											if len(parts) == 2 {
												var languageIdentifier uint16
												if _, parseError := fmt.Sscanf(parts[1], "0x%x", &languageIdentifier); parseError == nil {
													multilanguages[languageIdentifier] = element.value
												}
											}
										}
									}
									for index := 0; index < collection.NumFonts(); index++ {
										var sfntData *sfnt.Font
										if sfntData, walkError = collection.Font(index); walkError == nil {
											familyName := ""
											font := FONT_INFO{
												FilePath: path,
												Font:     sfntData,
												Language: 0,
												Names:    make(map[string]string),
											}
											for _, identifier := range NAME_ID {
												var value string
												if value, walkError = sfntData.Name(nil, identifier); walkError == nil && value != "" {
													if identifier == sfnt.NameIDFamily {
														font.Name = value
														familyName = value
													}
													font.Names[fmt.Sprintf("%d:0x000", identifier)] = value
												}
											}
											result = append(result, font)
											for language := range multilanguages {
												font := FONT_INFO{
													FilePath: path,
													Font:     sfntData,
													Language: language,
													Name:     familyName,
													Names:    make(map[string]string),
												}
												languageIdentifier := fmt.Sprintf("0x%x", language)
												for _, element := range multilanguageNames {
													if parts := strings.Split(element.key, ":"); len(parts) == 2 {
														if identifier, conversionError := strconv.Atoi(parts[0]); conversionError == nil && identifier == int(sfnt.NameIDFamily) {
															if language == LANGUAGE_ENGLISH {
																font.Name = element.value
															}
														}
														if parts[1] == languageIdentifier {
															font.Names[element.key] = element.value
														}
													}
												}
												result = append(result, font)
											}
										}
									}
								}
								if extension == ".ttf" {
									var sfntFont *sfnt.Font
									if sfntFont, walkError = sfnt.Parse(fileData); walkError == nil {
										var variableNames []string
										if variableNames, walkError = f.getTrueTypeFontVariableNames(path); walkError == nil {
											for _, variableName := range variableNames {
												result = append(result, FONT_INFO{
													Name:     variableName,
													Font:     sfntFont,
													FilePath: path,
												})
											}
										}
									}
								}
							}
						}
					}
				}
				return walkError
			})
		}
	}
	orphans := make([]FONT_INFO, 0)
	for _, font := range result {
		if font.Names == nil || len(font.Names) == 0 {
			orphans = append(orphans, font)
		}
	}
	for {
		duplicate := false
		for index, orphan := range orphans {
			for _, font := range result {
				if font.Name == orphan.Name && font.FilePath == orphan.FilePath {
					if font.Names != nil && len(font.Names) > 0 {
						duplicate = true
						break
					}
				} else if font.Names != nil && len(font.Names) > 0 && font.FilePath == orphan.FilePath {
					for _, name := range font.Names {
						if name == orphan.Name {
							duplicate = true
							break
						}
					}
				}
			}
			if duplicate {
				orphans = append(orphans[:index], orphans[index+1:]...)
				break
			}
		}
		if !duplicate {
			break
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func ParseCMapTable(fontPath string) ([]CMAP_FONT_INFO, error) {
	result := make([]CMAP_FONT_INFO, 0)
	err := error(nil)
	var fontData []byte
	if fontData, err = os.ReadFile(fontPath); err == nil {
		var collection *sfnt.Collection
		if collection, err = sfnt.ParseCollection(fontData); err == nil {
			result = make([]CMAP_FONT_INFO, 0, collection.NumFonts())
			for index := 0; index < collection.NumFonts(); index++ {
				var font *sfnt.Font
				if font, err = collection.Font(index); err == nil {
					fontName, _ := font.Name(nil, sfnt.NameIDFull)
					if fontName == "" {
						fontName, _ = font.Name(nil, sfnt.NameIDFamily)
					}
					glyphToUnicode := make(map[sfnt.GlyphIndex]rune)
					buffer := &sfnt.Buffer{}
					numberOfGlyphs := font.NumGlyphs()
					for runeValue := rune(0); runeValue <= 0x10FFFF; runeValue++ {
						var glyphIndex sfnt.GlyphIndex
						if glyphIndex, err = font.GlyphIndex(buffer, runeValue); err == nil && glyphIndex != 0 && int(glyphIndex) < numberOfGlyphs {
							glyphToUnicode[glyphIndex] = runeValue
						}
					}
					result = append(result, CMAP_FONT_INFO{
						FontName:       fontName,
						GlyphToUnicode: glyphToUnicode,
					})
				}
			}
		}
	}
	return result, err
}
