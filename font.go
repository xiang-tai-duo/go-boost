// Package boost
// File:        font.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/font.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description:
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

//goland:noinspection GoSnakeCaseUsage
type (
	FONT_INFO struct {
		Name     string
		FilePath string
		Font     *sfnt.Font
		Language uint16
		Names    map[string]string
	}

	FONT struct {
		systemFonts     []FONT_INFO
		caseInsensitive bool
		mutex           sync.RWMutex
	}

	FONT_TABLE_INFO struct {
		tag      string
		checkSum uint32
		offset   uint32
		length   uint32
	}

	FONT_VARIABLE_INSTANCE struct {
		subFamilyNameId uint16
		coordinates     []float32
	}

	FONT_KEY_PAIR struct {
		key   string
		value string
	}
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage,SpellCheckingInspection
const (
	FIRST_VARIANT_NAME_ID                  = 258
	FVAR_AXIS_COUNT_LENGTH                 = 2
	FVAR_AXIS_COUNT_OFFSET                 = 8
	FVAR_AXIS_SIZE_LENGTH                  = 2
	FVAR_AXIS_SIZE_OFFSET                  = 10
	FVAR_INSTANCE_COUNT_LENGTH             = 2
	FVAR_INSTANCE_COUNT_OFFSET             = 12
	FVAR_INSTANCE_COORDINATES_OFFSET       = 4
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
	ERROR_FONT_FOUND                       = "found"
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

func NewFont() *FONT {
	ret := &FONT{
		systemFonts:     []FONT_INFO{},
		caseInsensitive: false,
	}
	return ret
}

func (f *FONT) ValidateCharacters(inputString string, fontName string) ([]rune, []rune, error) {
	invisibleCharacters := make([]rune, 0)
	visibleCharacters := make([]rune, 0)
	err := error(nil)
	var p *FONT_INFO = nil
	if p, err = f.getSystemFont(fontName); err == nil {
		__invisibleCharacters := make([]rune, 0)
		__visibleCharacters := make([]rune, 0)
		for _, char := range inputString {
			index := sfnt.GlyphIndex(0)
			if index, err = p.Font.GlyphIndex(nil, char); err == nil {
				if index == 0 {
					__invisibleCharacters = append(__invisibleCharacters, char)
				} else {
					__visibleCharacters = append(__visibleCharacters, char)
				}
			} else {
				err = fmt.Errorf("check character %q (code:%d) failed: %w", char, char, err)
				break
			}
		}
		if err == nil {
			uniqueUndisplayable := make(map[rune]struct{})
			for _, c := range __invisibleCharacters {
				uniqueUndisplayable[c] = struct{}{}
			}
			invisibleCharacters = make([]rune, 0, len(uniqueUndisplayable))
			for c := range uniqueUndisplayable {
				invisibleCharacters = append(invisibleCharacters, c)
			}
			sort.Slice(invisibleCharacters, func(i, j int) bool {
				return invisibleCharacters[i] < invisibleCharacters[j]
			})
			visibleCharacters = __visibleCharacters
		}
	}
	return invisibleCharacters, visibleCharacters, err
}

func (f *FONT) getSystemFonts() error {
	err := error(nil)
	f.mutex.RLock()
	if len(f.systemFonts) > 0 {
		f.mutex.RUnlock()
	} else {
		f.mutex.RUnlock()
		var fontInfos []FONT_INFO
		switch runtime.GOOS {
		case "windows":
			fontInfos = f.getSystemFontsWindows()
		case "darwin":
			fontInfos = f.getSystemFontsMac()
		case "linux":
			fontInfos = f.getSystemFontsLinux()
		default:
			err = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}
		if err == nil {
			f.mutex.Lock()
			f.systemFonts = fontInfos
			f.mutex.Unlock()
		}
	}
	return err
}

func (f *FONT) GetCaseInsensitive() bool {
	return f.caseInsensitive
}

func (f *FONT) IsExists(fontName string) (bool, error) {
	font, err := f.getSystemFont(fontName)
	return font != nil, err
}

func (f *FONT) getSystemFont(fontName string) (*FONT_INFO, error) {
	var font *FONT_INFO = nil
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
						font = &systemFont
						break
					}
				} else {
					if systemFontName == fontName {
						font = &systemFont
						break
					}
				}
			}
			if font != nil {
				break
			}
		}
		f.mutex.RUnlock()
	} else {
		err = fmt.Errorf("get system fonts failed: %w", err)
	}
	return font, err
}

func (f *FONT) SetCaseInsensitive(enabled bool) {
	f.caseInsensitive = enabled
}

func (f *FONT) getSystemFontsWindows() []FONT_INFO {
	fontDirectories := []string{
		"C:\\Windows\\Fonts",
		filepath.Join(os.Getenv("USERPROFILE"), "AppData\\Local\\Microsoft\\Windows\\Fonts"),
	}
	return f.init(fontDirectories)
}

func (f *FONT) getSystemFontsMac() []FONT_INFO {
	fontDirectories := []string{
		"/Library/Fonts",
		"/System/Library/Fonts",
		filepath.Join(os.Getenv("HOME"), "Library/Fonts"),
	}
	return f.init(fontDirectories)
}

func (f *FONT) getSystemFontsLinux() []FONT_INFO {
	fontDirectories := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		filepath.Join(os.Getenv("HOME"), ".local/share/fonts"),
		filepath.Join(os.Getenv("HOME"), ".fonts"),
	}
	return f.init(fontDirectories)
}

//goland:noinspection SpellCheckingInspection
func (f *FONT) init(directories []string) []FONT_INFO {
	result := make([]FONT_INFO, 0)
	for _, directory := range directories {
		if _, err := os.Stat(directory); err == nil {
			err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err == nil {
					if !info.IsDir() {
						ext := strings.ToLower(filepath.Ext(path))
						if ext == ".ttf" || ext == ".otf" || ext == ".ttc" || ext == ".dfont" {
							if file, err := os.ReadFile(path); err == nil {
								if collection, err := sfnt.ParseCollection(file); err == nil {
									multilanguageNames := make([]FONT_KEY_PAIR, 0)
									multilanguages := make(map[uint16]string)
									if multilanguageNames, err = f.getMultilanguageNameTables(path); err == nil {
										for _, element := range multilanguageNames {
											parts := strings.Split(element.key, ":")
											if len(parts) == 2 {
												var langID uint16
												if _, err := fmt.Sscanf(parts[1], "0x%x", &langID); err == nil {
													multilanguages[langID] = element.value
												}
											}
										}
									}
									for i := 0; i < collection.NumFonts(); i++ {
										if sfntData, err := collection.Font(i); err == nil {
											familyName := ""
											font := FONT_INFO{
												FilePath: path,
												Font:     sfntData,
												Language: 0,
												Names:    make(map[string]string),
											}
											for _, id := range NAME_ID {
												if value, err := sfntData.Name(nil, id); err == nil && value != "" {
													if id == sfnt.NameIDFamily {
														font.Name = value
														familyName = value
													}
													font.Names[fmt.Sprintf("%d:0x000", id)] = value
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
												langId := fmt.Sprintf("0x%x", language)
												for _, element := range multilanguageNames {
													if parts := strings.Split(element.key, ":"); len(parts) == 2 {
														if id, err := strconv.Atoi(parts[0]); err == nil && id == int(sfnt.NameIDFamily) {
															if language == LANGUAGE_ENGLISH {
																font.Name = element.value
															}
														}
														if parts[1] == langId {
															font.Names[element.key] = element.value
														}
													}
												}
												result = append(result, font)
											}
										}
									}
								}
								if ext == ".ttf" {
									if sfnt, err := sfnt.Parse(file); err == nil {
										if variableNames, err := f.getTrueTypeFontVariableNames(path); err == nil {
											for _, variableName := range variableNames {
												result = append(result, FONT_INFO{
													Name:     variableName,
													Font:     sfnt,
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
				return err
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
						if name==orphan.Name{
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

func (f *FONT) getFontInfo(fontName string) (FONT_INFO, error) {
	var font FONT_INFO
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		found := false
		f.mutex.RLock()
		for _, fc := range f.systemFonts {
			if f.caseInsensitive {
				if strings.EqualFold(fc.Name, fontName) {
					font = fc
					found = true
					break
				}
			} else {
				if fc.Name == fontName {
					font = fc
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
	return font, err
}

func (f *FONT) decodeUTF16(data []byte) string {
	if len(data)%UTF16_CHAR_SIZE != 0 {
		data = append(data, 0)
	}
	var result string
	for i := 0; i < len(data); i += UTF16_CHAR_SIZE {
		char := binary.BigEndian.Uint16(data[i : i+UTF16_CHAR_SIZE])
		if char != UTF16_NULL_CHAR {
			result += string(rune(char))
		}
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func (f *FONT) getMultilanguageNameTables(path string) ([]FONT_KEY_PAIR, error) {
	var err error
	var file *os.File
	var fileInfo os.FileInfo
	var fileSize int64
	var header []byte
	var magic uint32
	result := make([]FONT_KEY_PAIR, 0)
	if file, err = os.Open(path); err == nil {
		defer file.Close()
		if fileInfo, err = file.Stat(); err == nil {
			fileSize = fileInfo.Size()
			if fileSize >= TTF_HEADER_LENGTH {
				header = make([]byte, TTF_HEADER_LENGTH)
				if _, err = file.Read(header); err == nil {
					magic = binary.BigEndian.Uint32(header[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
					if magic == TTC_MAGIC_NUMBER {
						numFonts := binary.BigEndian.Uint32(header[TTC_NUM_FONTS_OFFSET : TTC_NUM_FONTS_OFFSET+TTC_NUM_FONTS_LENGTH])
						if numFonts > 0 {
							fontOffsets := make([]uint32, numFonts)
							for i := 0; i < int(numFonts); i++ {
								offsetPos := TTC_FONT_OFFSET_TABLE_OFFSET + i*SIZE_DWROD
								if int64(offsetPos+SIZE_DWROD) <= fileSize {
									if _, err = file.Seek(int64(offsetPos), 0); err == nil {
										offsetData := make([]byte, SIZE_DWROD)
										if _, err = file.Read(offsetData); err == nil {
											fontOffsets[i] = binary.BigEndian.Uint32(offsetData)
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
												ttfMagic := binary.BigEndian.Uint32(ttfHeader[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
												if ttfMagic == TTF_MAGIC_NUMBER_1 || ttfMagic == TTF_MAGIC_NUMBER_2 {
													if fontNameTable, err := f.getMultilanguageNameTable(file, fileSize, int64(fontOffset), ttfHeader); err == nil {
														result = append(result, fontNameTable...)
													}
												}
											}
										}
									}
								}
							}
						}
					} else if magic == TTF_MAGIC_NUMBER_1 || magic == TTF_MAGIC_NUMBER_2 {
						if fontNameTable, err := f.getMultilanguageNameTable(file, fileSize, 0, header); err == nil {
							result = append(result, fontNameTable...)
						}
					}
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		key1 := result[i].key
		key2 := result[j].key
		for len(key1) < 8 {
			key1 = "0" + key1
		}
		for len(key2) < 8 {
			key2 = "0" + key2
		}
		return key1 < key2
	})
	return result, err
}

//goland:noinspection DuplicatedCode
func (f *FONT) getMultilanguageNameTable(file *os.File, fileSize int64, fontOffset int64, header []byte) ([]FONT_KEY_PAIR, error) {
	result := make([]FONT_KEY_PAIR, 0)
	var err error
	var numTables uint16
	var tableDirOffset int64
	var tables map[string]FONT_TABLE_INFO
	numTables = binary.BigEndian.Uint16(header[HEADER_TABLE_COUNT_OFFSET : HEADER_TABLE_COUNT_OFFSET+HEADER_TABLE_COUNT_LENGTH])
	tableDirOffset = fontOffset + int64(TTF_HEADER_LENGTH)
	if tableDirOffset+int64(numTables*TABLE_ENTRY_LENGTH) <= fileSize {
		if _, err = file.Seek(tableDirOffset, 0); err == nil {
			tables = make(map[string]FONT_TABLE_INFO)
			for i := 0; i < int(numTables); i++ {
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
				if nameInfo, ok := tables[TABLE_TAG_NAME]; ok {
					if _, err = file.Seek(int64(nameInfo.offset), 0); err == nil {
						nameHeader := make([]byte, NAME_HEADER_LENGTH)
						if _, err = file.Read(nameHeader); err == nil {
							nameCount := binary.BigEndian.Uint16(nameHeader[2:4])
							stringOffset := binary.BigEndian.Uint16(nameHeader[4:6])
							for i := 0; i < int(nameCount); i++ {
								nameRecord := make([]byte, NAME_RECORD_LENGTH)
								if _, err = file.Read(nameRecord); err == nil {
									platformID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_PLATFORM_ID_OFFSET : NAME_RECORD_PLATFORM_ID_OFFSET+NAME_RECORD_PLATFORM_ID_LENGTH])
									languageID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_LANGUAGE_ID_OFFSET : NAME_RECORD_LANGUAGE_ID_OFFSET+NAME_RECORD_LANGUAGE_ID_LENGTH])
									nameID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_NAME_ID_OFFSET : NAME_RECORD_NAME_ID_OFFSET+NAME_RECORD_NAME_ID_LENGTH])
									length := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_LENGTH_OFFSET : NAME_RECORD_LENGTH_OFFSET+NAME_RECORD_LENGTH_LENGTH])
									offset := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_OFFSET_OFFSET : NAME_RECORD_OFFSET_OFFSET+NAME_RECORD_OFFSET_LENGTH])
									if platformID == PLATFORM_ID_MICROSOFT {
										var currentPos int64
										if currentPos, err = file.Seek(0, 1); err == nil {
											if _, err = file.Seek(int64(nameInfo.offset)+int64(stringOffset)+int64(offset), 0); err == nil {
												stringData := make([]byte, length)
												if _, err = file.Read(stringData); err == nil {
													stringValue := f.decodeUTF16(stringData)
													key := fmt.Sprintf("%d:0x%x", int(nameID), languageID)
													result = append(result, FONT_KEY_PAIR{
														key:   key,
														value: stringValue,
													})
												}
												if _, err = file.Seek(currentPos, 0); err != nil {
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
			if i, err := strconv.Atoi(parts[0]); err == nil && i > NAME_ID_LAST {
				key := fmt.Sprintf("1:%s", parts[1])
				for _, __element := range result {
					if __element.key == key {
						result[index].value = fmt.Sprintf("%s %s", __element.value, element.value)
						break
					}
				}
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,SpellCheckingInspection,DuplicatedCode
func (f *FONT) getTrueTypeFontVariableNames(fontPath string) ([]string, error) {
	err := error(nil)
	var variantNames []string
	var file *os.File
	var fileInfo os.FileInfo
	var fileSize int64
	var header []byte
	var magic uint32
	var numTables uint16
	var tableDirOffset int64
	var tables map[string]FONT_TABLE_INFO
	var instances []FONT_VARIABLE_INSTANCE
	if file, err = os.Open(fontPath); err == nil {
		defer file.Close()
		if fileInfo, err = file.Stat(); err == nil {
			fileSize = fileInfo.Size()
			if fileSize >= TTF_HEADER_LENGTH {
				header = make([]byte, TTF_HEADER_LENGTH)
				if _, err = file.Read(header); err == nil {
					magic = binary.BigEndian.Uint32(header[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
					if magic == TTF_MAGIC_NUMBER_1 || magic == TTF_MAGIC_NUMBER_2 {
						numTables = binary.BigEndian.Uint16(header[HEADER_TABLE_COUNT_OFFSET : HEADER_TABLE_COUNT_OFFSET+HEADER_TABLE_COUNT_LENGTH])
						tableDirOffset = int64(TTF_HEADER_LENGTH)
						if tableDirOffset+int64(numTables*TABLE_ENTRY_LENGTH) <= fileSize {
							if _, err = file.Seek(tableDirOffset, 0); err == nil {
								tables = make(map[string]FONT_TABLE_INFO)
								for i := 0; i < int(numTables); i++ {
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
										if nameInfo, ok := tables[TABLE_TAG_NAME]; ok {
											if _, err = file.Seek(int64(nameInfo.offset), 0); err == nil {
												nameHeader := make([]byte, NAME_HEADER_LENGTH)
												if _, err = file.Read(nameHeader); err == nil {
													nameCount := binary.BigEndian.Uint16(nameHeader[2:4])
													stringOffset := binary.BigEndian.Uint16(nameHeader[4:6])
													for i := 0; i < int(nameCount); i++ {
														nameRecord := make([]byte, NAME_RECORD_LENGTH)
														if _, err = file.Read(nameRecord); err == nil {
															platformID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_PLATFORM_ID_OFFSET : NAME_RECORD_PLATFORM_ID_OFFSET+NAME_RECORD_PLATFORM_ID_LENGTH])
															nameID := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_NAME_ID_OFFSET : NAME_RECORD_NAME_ID_OFFSET+NAME_RECORD_NAME_ID_LENGTH])
															length := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_LENGTH_OFFSET : NAME_RECORD_LENGTH_OFFSET+NAME_RECORD_LENGTH_LENGTH])
															offset := binary.BigEndian.Uint16(nameRecord[NAME_RECORD_OFFSET_OFFSET : NAME_RECORD_OFFSET_OFFSET+NAME_RECORD_OFFSET_LENGTH])
															var currentPos int64
															if currentPos, err = file.Seek(0, 1); err == nil {
																if _, err = file.Seek(int64(nameInfo.offset)+int64(stringOffset)+int64(offset), 0); err == nil {
																	stringData := make([]byte, length)
																	if _, err = file.Read(stringData); err == nil {
																		var stringValue string
																		switch platformID {
																		case PLATFORM_ID_UNICODE:
																			stringValue = f.decodeUTF16(stringData)
																		case PLATFORM_ID_MACINTOSH:
																			var name strings.Builder
																			for j := 0; j < int(length); j++ {
																				if j >= len(stringData) {
																					break
																				}
																				char := stringData[j]
																				if char == 0 {
																					break
																				}
																				name.WriteRune(rune(char))
																			}
																			stringValue = name.String()
																		case PLATFORM_ID_ISO:
																			var name strings.Builder
																			for j := 0; j < int(length); j++ {
																				if j >= len(stringData) {
																					break
																				}
																				char := stringData[j]
																				if char == 0 {
																					break
																				}
																				name.WriteRune(rune(char))
																			}
																			stringValue = name.String()
																		case PLATFORM_ID_MICROSOFT:
																			stringValue = f.decodeUTF16(stringData)
																		case PLATFORM_ID_CUSTOM:
																			var name strings.Builder
																			for j := 0; j < int(length); j++ {
																				if j >= len(stringData) {
																					break
																				}
																				char := stringData[j]
																				if char == 0 {
																					break
																				}
																				name.WriteRune(rune(char))
																			}
																			stringValue = name.String()
																		default:
																			var name strings.Builder
																			for j := 0; j < int(length); j++ {
																				if j >= len(stringData) {
																					break
																				}
																				char := stringData[j]
																				if char == 0 {
																					break
																				}
																				name.WriteRune(rune(char))
																			}
																			stringValue = name.String()
																		}
																		nameTable[nameID] = stringValue
																		if _, err = file.Seek(currentPos, 0); err != nil {
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
											if fvar, ok := tables[TABLE_TAG_FVAR]; ok {
												if _, err = file.Seek(int64(fvar.offset), 0); err == nil {
													fvarBytes := make([]byte, fvar.length)
													if _, err = file.Read(fvarBytes); err == nil {
														offsetToData := binary.BigEndian.Uint16(fvarBytes[FVAR_OFFSET_TO_DATA_OFFSET : FVAR_OFFSET_TO_DATA_OFFSET+FVAR_VERSION_LENGTH])
														axisCount := binary.BigEndian.Uint16(fvarBytes[FVAR_AXIS_COUNT_OFFSET : FVAR_AXIS_COUNT_OFFSET+FVAR_AXIS_COUNT_LENGTH])
														axisSize := binary.BigEndian.Uint16(fvarBytes[FVAR_AXIS_SIZE_OFFSET : FVAR_AXIS_SIZE_OFFSET+FVAR_AXIS_SIZE_LENGTH])
														instanceCount := binary.BigEndian.Uint16(fvarBytes[FVAR_INSTANCE_COUNT_OFFSET : FVAR_INSTANCE_COUNT_OFFSET+FVAR_INSTANCE_COUNT_LENGTH])
														instanceSize := binary.BigEndian.Uint16(fvarBytes[FVAR_INSTANCE_SIZE_OFFSET : FVAR_INSTANCE_SIZE_OFFSET+FVAR_INSTANCE_SIZE_LENGTH])
														instances = make([]FONT_VARIABLE_INSTANCE, 0)
														pos := int(offsetToData)
														for i := 0; i < int(axisCount); i++ {
															pos += int(axisSize)
														}
														for i := 0; i < int(instanceCount); i++ {
															if pos+int(instanceSize) <= len(fvarBytes) {
																subFamilyNameId := binary.BigEndian.Uint16(fvarBytes[pos : pos+FVAR_INSTANCE_SUBFAMILY_NAME_ID_LENGTH])
																instances = append(instances, FONT_VARIABLE_INSTANCE{
																	subFamilyNameId: subFamilyNameId,
																	coordinates:     nil,
																})
																pos += int(instanceSize)
															}
														}
													}
												}
											}
										}
										if err == nil {
											familyNames := make([]string, 0)
											if name, ok := nameTable[uint16(sfnt.NameIDFamily)]; ok {
												familyNames = append(familyNames, name)
											}
											if name, ok := nameTable[uint16(sfnt.NameIDTypographicFamily)]; ok {
												familyNames = append(familyNames, name)
											}
											variantNames = make([]string, 0, len(instances))
											for _, familyName := range familyNames {
												for _, instance := range instances {
													if subFamilyName, ok := nameTable[instance.subFamilyNameId]; ok {
														if familyName == "" {
															variantNames = append(variantNames, subFamilyName)
														} else {
															variantNames = append(variantNames, fmt.Sprintf("%s %s", familyName, subFamilyName))
														}
													}
													//} else {
													//	if familyName == "" {
													//		variantNames = append(variantNames, fmt.Sprintf("%d", instance.subFamilyNameId))
													//	} else {
													//		variantNames = append(variantNames, fmt.Sprintf("%s %d", familyName, instance.subFamilyNameId))
													//	}
													//}
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
	return variantNames, err
}

func (f *FONT) ExportSystemFontsToFile(filename string) error {
	var err error
	var file *os.File
	if err = f.getSystemFonts(); err == nil {
		if file, err = os.Create(filename); err == nil {
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
