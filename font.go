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

const (
	ERROR_FONT_FOUND                       = "found"
	TTF_HEADER_LENGTH                      = 12
	TTF_MAGIC_NUMBER_1                     = 0x00010000
	TTF_MAGIC_NUMBER_2                     = 0x4F54544F
	TTC_MAGIC_NUMBER                       = 0x74746366
	TTC_HEADER_LENGTH                      = 12
	TABLE_ENTRY_LENGTH                     = 16
	NAME_HEADER_LENGTH                     = 6
	PLATFORM_ID_MICROSOFT                  = 3
	NAME_ID_FAMILY                         = 1
	FIRST_VARIANT_NAME_ID                  = 258
	LAST_VARIANT_NAME_ID                   = 271
	TABLE_TAG_FVAR                         = "fvar"
	TABLE_TAG_NAME                         = "name"
	HEADER_MAGIC_OFFSET                    = 0
	HEADER_MAGIC_LENGTH                    = 4
	HEADER_TABLE_COUNT_OFFSET              = 4
	HEADER_TABLE_COUNT_LENGTH              = 2
	HEADER_SEARCH_RANGE_OFFSET             = 6
	HEADER_SEARCH_RANGE_LENGTH             = 2
	HEADER_ENTRY_SELECTOR_OFFSET           = 8
	HEADER_ENTRY_SELECTOR_LENGTH           = 2
	HEADER_RANGE_SHIFT_OFFSET              = 10
	HEADER_RANGE_SHIFT_LENGTH              = 2
	TTC_VERSION_OFFSET                     = 4
	TTC_VERSION_LENGTH                     = 4
	TTC_NUM_FONTS_OFFSET                   = 8
	TTC_NUM_FONTS_LENGTH                   = 4
	TTC_FONT_OFFSET_TABLE_OFFSET           = 12
	TABLE_ENTRY_TAG_OFFSET                 = 0
	TABLE_ENTRY_TAG_LENGTH                 = 4
	TABLE_ENTRY_CHECKSUM_OFFSET            = 4
	TABLE_ENTRY_CHECKSUM_LENGTH            = 4
	TABLE_ENTRY_OFFSET_OFFSET              = 8
	TABLE_ENTRY_OFFSET_LENGTH              = 4
	TABLE_ENTRY_LENGTH_OFFSET              = 12
	TABLE_ENTRY_LENGTH_LENGTH              = 4
	NAME_RECORD_PLATFORM_ID_OFFSET         = 0
	NAME_RECORD_PLATFORM_ID_LENGTH         = 2
	NAME_RECORD_ENCODING_ID_OFFSET         = 2
	NAME_RECORD_ENCODING_ID_LENGTH         = 2
	NAME_RECORD_LANGUAGE_ID_OFFSET         = 4
	NAME_RECORD_LANGUAGE_ID_LENGTH         = 2
	NAME_RECORD_NAME_ID_OFFSET             = 6
	NAME_RECORD_NAME_ID_LENGTH             = 2
	NAME_RECORD_LENGTH_OFFSET              = 8
	NAME_RECORD_LENGTH_LENGTH              = 2
	NAME_RECORD_OFFSET_OFFSET              = 10
	NAME_RECORD_OFFSET_LENGTH              = 2
	NAME_RECORD_LENGTH                     = 12
	FVAR_VERSION_OFFSET                    = 0
	FVAR_VERSION_LENGTH                    = 4
	FVAR_AXIS_COUNT_OFFSET                 = 4
	FVAR_AXIS_COUNT_LENGTH                 = 2
	FVAR_INSTANCE_COUNT_OFFSET             = 6
	FVAR_INSTANCE_COUNT_LENGTH             = 2
	FVAR_INSTANCE_SUBFAMILY_NAME_ID_OFFSET = 0
	FVAR_INSTANCE_SUBFAMILY_NAME_ID_LENGTH = 2
	FVAR_INSTANCE_FLAGS_OFFSET             = 2
	FVAR_INSTANCE_FLAGS_LENGTH             = 2
	FVAR_INSTANCE_COORDINATES_OFFSET       = 8
	UTF16_CHAR_SIZE                        = 2
	UTF16_NULL_CHAR                        = 0
)

type FONT_CACHE struct {
	Name     string
	FilePath string
	Font     *sfnt.Font
	Language uint16
	NameIDs  map[string]string
}

type FONT struct {
	systemFontsCache []FONT_CACHE
	caseInsensitive  bool
	cacheMutex       sync.RWMutex
}

type (
	FONT_TABLE_INFO struct {
		tag      string
		checkSum uint32
		offset   uint32
		length   uint32
	}

	FONT_VARIABLE_INSTANCE struct {
		subfamilyNameID uint16
		coordinates     []float32
	}

	FONT_NAME_ENTRY struct {
		Key   string
		Value string
	}
)

func NewFont() *FONT {
	ret := &FONT{
		systemFontsCache: []FONT_CACHE{},
		caseInsensitive:  false,
	}
	return ret
}

func (f *FONT) ValidateCharacters(inputString string, fontName string) ([]rune, []rune, error) {
	retUndisplayable := []rune{}
	retDisplayable := []rune{}
	retErr := error(nil)
	err := error(nil)
	exists, err := f.IsExists(fontName)
	if err == nil && exists {
		undisplayableTemp := []rune{}
		displayableTemp := []rune{}
		for _, char := range inputString {
			isValid := false
			if isValid, err = f.validateCharacter(fontName, char); err == nil {
				if isValid {
					displayableTemp = append(displayableTemp, char)
				} else {
					undisplayableTemp = append(undisplayableTemp, char)
				}
			} else {
				retErr = fmt.Errorf("check character %q (code:%d) failed: %w", char, char, err)
				break
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
		retDisplayable = displayableTemp
	} else if err != nil {
		retErr = fmt.Errorf("check font existence failed: %w", err)
	} else {
		retErr = fmt.Errorf("font %s not found in system", fontName)
	}
	return retUndisplayable, retDisplayable, retErr
}

func (f *FONT) getSystemFonts() error {
	retErr := error(nil)
	err := error(nil)
	f.cacheMutex.RLock()
	if len(f.systemFontsCache) > 0 {
		f.cacheMutex.RUnlock()
	} else {
		f.cacheMutex.RUnlock()
		var fontInfos []FONT_CACHE
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
		if err != nil {
			retErr = fmt.Errorf("get system fonts failed: %w", err)
		} else {
			f.cacheMutex.Lock()
			f.systemFontsCache = fontInfos
			f.cacheMutex.Unlock()
		}
	}
	return retErr
}

func (f *FONT) GetCaseInsensitive() bool {
	return f.caseInsensitive
}

func (f *FONT) IsExists(fontName string) (bool, error) {
	b := false
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		f.cacheMutex.RLock()
		for _, fontCache := range f.systemFontsCache {
			for _, name := range fontCache.NameIDs {
				if f.caseInsensitive {
					if strings.EqualFold(name, fontName) {
						b = true
						break
					}
				} else {
					if name == fontName {
						b = true
						break
					}
				}
			}
			if b {
				break
			}
		}
		f.cacheMutex.RUnlock()
	} else {
		err = fmt.Errorf("get system fonts failed: %w", err)
	}
	return b, err
}

func (f *FONT) SetCaseInsensitive(enabled bool) {
	f.caseInsensitive = enabled
}

func (f *FONT) getSystemFontsWindows() []FONT_CACHE {
	fontDirectories := []string{
		"C:\\Windows\\Fonts",
		filepath.Join(os.Getenv("USERPROFILE"), "AppData\\Local\\Microsoft\\Windows\\Fonts"),
	}
	return f.scanDirectories(fontDirectories)
}

func (f *FONT) getSystemFontsMac() []FONT_CACHE {
	fontDirectories := []string{
		"/Library/Fonts",
		"/System/Library/Fonts",
		filepath.Join(os.Getenv("HOME"), "Library/Fonts"),
	}
	return f.scanDirectories(fontDirectories)
}

func (f *FONT) getSystemFontsLinux() []FONT_CACHE {
	fontDirectories := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		filepath.Join(os.Getenv("HOME"), ".local/share/fonts"),
		filepath.Join(os.Getenv("HOME"), ".fonts"),
	}
	return f.scanDirectories(fontDirectories)
}

func (f *FONT) newFontCache(path string, font *sfnt.Font) *FONT_CACHE {
	return &FONT_CACHE{
		FilePath: path,
		Font:     font,
		NameIDs:  make(map[string]string),
	}
}

func (f *FONT) setFontName(cache *FONT_CACHE) {
	var fontName string
	if fontName == "" {
		for key, value := range cache.NameIDs {
			if strings.Contains(key, fmt.Sprintf("%d", int(sfnt.NameIDFamily))) {
				fontName = value
				break
			}
		}
	}
	if fontName == "" {
		for key, value := range cache.NameIDs {
			if strings.Contains(key, fmt.Sprintf("%d", int(sfnt.NameIDFull))) {
				fontName = value
				break
			}
		}
	}
	for key, value := range cache.NameIDs {
		if strings.Contains(key, fmt.Sprintf("%d", int(sfnt.NameIDCompatibleFull))) {
			fontName = value
			break
		}
	}
	if fontName != "" {
		cache.Name = fontName
	}
}

func (f *FONT) createFontCache(path string, font *sfnt.Font, nameIDs []sfnt.NameID, index int) []*FONT_CACHE {
	result := make([]*FONT_CACHE, 0)
	if nameEntries, err := f.getNameTables(path); err == nil && len(nameEntries) > 0 {
		nameMap := make(map[string]string)
		nameLanguages := make(map[uint16]bool)
		for _, entry := range nameEntries {
			nameMap[entry.Key] = entry.Value
			parts := strings.Split(entry.Key, ":")
			if len(parts) == 2 {
				var langID uint16
				if _, err := fmt.Sscanf(parts[1], "0x%x", &langID); err == nil {
					nameLanguages[langID] = true
				}
			}
		}
		if len(nameLanguages) > 0 {
			for language := range nameLanguages {
				cache := f.newFontCache(path, font)
				cache.Language = language
				for _, nameID := range nameIDs {
					if value, ok := nameMap[fmt.Sprintf("%d:0x%x", int(nameID), language)]; ok && value != "" {
						cache.NameIDs[fmt.Sprintf("%d", nameID)] = value
					}
				}
				f.setFontName(cache)
				result = append(result, cache)
			}
		}
	}
	cache := f.newFontCache(path, font)
	cache.Language = 0
	for _, id := range nameIDs {
		value, err := font.Name(nil, id)
		if err == nil && value != "" {
			cache.NameIDs[fmt.Sprintf("%d", id)] = value
		}
	}
	f.setFontName(cache)
	result = append(result, cache)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Language < result[j].Language
	})
	return result
}

func (f *FONT) scanDirectories(directories []string) []FONT_CACHE {
	fontsMap := make(map[string]*FONT_CACHE)
	nameIDs := []sfnt.NameID{
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
	for _, directory := range directories {
		if _, err := os.Stat(directory); err == nil {
			err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err == nil {
					if !info.IsDir() {
						ext := strings.ToLower(filepath.Ext(path))
						if ext == ".ttf" || ext == ".otf" || ext == ".ttc" || ext == ".dfont" {
							if path == "C:\\Windows\\Fonts\\SitkaVF.ttf" {
								fmt.Printf("%s", path)
							}
							if fontData, readErr := os.ReadFile(path); readErr == nil {
								tempMap := make(map[string]*FONT_CACHE)
								var parsedFont *sfnt.Font
								collection, collectionErr := sfnt.ParseCollection(fontData)
								if collectionErr == nil {
									for i := 0; i < collection.NumFonts(); i++ {
										if font, err := collection.Font(i); err == nil {
											for j, fontCache := range f.createFontCache(path, font, nameIDs, i) {
												tempMap[fmt.Sprintf("%s#%d#%d", path, i, j)] = fontCache
												if i == 0 && j == 0 {
													parsedFont = font
												}
											}
										}
									}
								} else {
									font, parseErr := sfnt.Parse(fontData)
									if parseErr != nil {
										return nil
									}
									parsedFont = font
									for j, fontCache := range f.createFontCache(path, font, nameIDs, -1) {
										tempMap[fmt.Sprintf("%s#%d", path, j)] = fontCache
									}
								}
								if variantNames, variantErr := f.getVariableNames(path); variantErr == nil && len(variantNames) > 0 {
									for _, variantName := range variantNames {
										found := false
										for _, cache := range tempMap {
											for _, name := range cache.NameIDs {
												if name == variantName {
													found = true
													break
												}
											}
											if found {
												break
											}
										}
										if !found {
											if parsedFont != nil {
												variantCache := &FONT_CACHE{
													Name:     variantName,
													FilePath: path,
													Font:     parsedFont,
													NameIDs:  make(map[string]string),
												}
												for _, id := range nameIDs {
													value, err := parsedFont.Name(nil, id)
													if err == nil && value != "" {
														variantCache.NameIDs[fmt.Sprintf("%d:0x%x", int(id), 0x0409)] = value // 使用英文语言ID
													}
												}
												variantCache.NameIDs[fmt.Sprintf("%d:0x%x", int(sfnt.NameIDFull), 0x0409)] = variantName // 使用英文语言ID
												tempMap[fmt.Sprintf("%s#%s", path, variantName)] = variantCache
											}
										}
									}
								}
								for key, cache := range tempMap {
									fontsMap[key] = cache
								}
							}
						}
					}
				}
				return err
			})
		}
	}
	fontsCache := make([]FONT_CACHE, 0, len(fontsMap))
	for _, fontCache := range fontsMap {
		fontsCache = append(fontsCache, *fontCache)
	}
	sort.Slice(fontsCache, func(i, j int) bool {
		return fontsCache[i].Name < fontsCache[j].Name
	})
	for _, fontCache := range fontsCache {
		if fontCache.Name !="" {
			fmt.Printf("Name: %s, FilePath: %s\n", fontCache.Name, fontCache.FilePath)
		}
	}
	return fontsCache
}

func (f *FONT) validateCharacter(fontName string, char rune) (bool, error) {
	b := false
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		var fontCache *FONT_CACHE
		f.cacheMutex.RLock()
		for _, fc := range f.systemFontsCache {
			if f.caseInsensitive {
				if strings.EqualFold(fc.Name, fontName) {
					fontCache = &fc
					break
				}
			} else {
				if fc.Name == fontName {
					fontCache = &fc
					break
				}
			}
		}
		f.cacheMutex.RUnlock()
		if fontCache == nil || fontCache.Font == nil {
			err = fmt.Errorf("font not found")
		} else {
			index := sfnt.GlyphIndex(0)
			if index, err = fontCache.Font.GlyphIndex(nil, char); err == nil {
				if index != 0 {
					b = true
				}
			} else {
				err = fmt.Errorf("check glyph index failed: %w", err)
			}
		}
	}
	return b, err
}

func (f *FONT) getFontCache(fontName string) (FONT_CACHE, error) {
	var fontCache FONT_CACHE
	err := error(nil)
	if err = f.getSystemFonts(); err == nil {
		found := false
		f.cacheMutex.RLock()
		for _, fc := range f.systemFontsCache {
			if f.caseInsensitive {
				if strings.EqualFold(fc.Name, fontName) {
					fontCache = fc
					found = true
					break
				}
			} else {
				if fc.Name == fontName {
					fontCache = fc
					found = true
					break
				}
			}
		}
		f.cacheMutex.RUnlock()
		if !found {
			err = fmt.Errorf("font %s not found in system", fontName)
		}
	}
	return fontCache, err
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

func (f *FONT) getNameTables(path string) ([]FONT_NAME_ENTRY, error) {
	var err error
	var file *os.File
	var fileInfo os.FileInfo
	var fileSize int64
	var header []byte
	var magic uint32
	var nameEntries []FONT_NAME_ENTRY
	if file, err = os.Open(path); err == nil {
		defer file.Close()
		if fileInfo, err = file.Stat(); err == nil {
			fileSize = fileInfo.Size()
			if fileSize >= TTF_HEADER_LENGTH {
				header = make([]byte, TTF_HEADER_LENGTH)
				if _, err = file.Read(header); err == nil {
					magic = binary.BigEndian.Uint32(header[HEADER_MAGIC_OFFSET : HEADER_MAGIC_OFFSET+HEADER_MAGIC_LENGTH])
					if magic == TTC_MAGIC_NUMBER {
						if fileSize >= TTC_HEADER_LENGTH {
							numFonts := binary.BigEndian.Uint32(header[TTC_NUM_FONTS_OFFSET : TTC_NUM_FONTS_OFFSET+TTC_NUM_FONTS_LENGTH])
							if numFonts > 0 {
								fontOffsets := make([]uint32, numFonts)
								for i := 0; i < int(numFonts); i++ {
									offsetPos := TTC_FONT_OFFSET_TABLE_OFFSET + i*4
									if int64(offsetPos+4) <= fileSize {
										if _, err = file.Seek(int64(offsetPos), 0); err == nil {
											offsetData := make([]byte, 4)
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
														fontNameEntries, fontErr := f.getNameTable(file, fileSize, int64(fontOffset), ttfHeader)
														if fontErr == nil {
															entryMap := make(map[string]string)
															for _, entry := range nameEntries {
																entryMap[entry.Key] = entry.Value
															}
															for _, entry := range fontNameEntries {
																if val, exists := entryMap[entry.Key]; !exists || val != entry.Value {
																	nameEntries = append(nameEntries, entry)
																	entryMap[entry.Key] = entry.Value
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
					} else if magic == TTF_MAGIC_NUMBER_1 || magic == TTF_MAGIC_NUMBER_2 {
						nameEntries, err = f.getNameTable(file, fileSize, 0, header)
					}
				}
			}
		}
	}
	sort.Slice(nameEntries, func(i, j int) bool {
		iKey := nameEntries[i].Key
		jKey := nameEntries[j].Key
		iNameIDStr := ""
		jNameIDStr := ""
		for _, c := range iKey {
			if c == ':' {
				break
			}
			iNameIDStr += string(c)
		}
		for _, c := range jKey {
			if c == ':' {
				break
			}
			jNameIDStr += string(c)
		}
		iNameID, _ := strconv.Atoi(iNameIDStr)
		jNameID, _ := strconv.Atoi(jNameIDStr)
		iPadded := fmt.Sprintf("%08d", iNameID)
		jPadded := fmt.Sprintf("%08d", jNameID)
		return iPadded < jPadded
	})
	return nameEntries, err
}

func (f *FONT) getNameTable(file *os.File, fileSize int64, fontOffset int64, header []byte) ([]FONT_NAME_ENTRY, error) {
	nameEntries := []FONT_NAME_ENTRY{}
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
													if platformID == PLATFORM_ID_MICROSOFT {
														key := fmt.Sprintf("%d:0x%x", int(nameID), languageID)
														nameEntries = append(nameEntries, FONT_NAME_ENTRY{Key: key, Value: stringValue})
													}
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
	return nameEntries, err
}

func (f *FONT) getVariableNames(fontPath string) ([]string, error) {
	var variantNames []string
	var err error
	var file *os.File
	var fileInfo os.FileInfo
	var fileSize int64
	var header []byte
	var magic uint32
	var numTables uint16
	var tableDirOffset int64
	var tables map[string]FONT_TABLE_INFO
	var nameTable map[uint16]string
	var instances []FONT_VARIABLE_INSTANCE
	var fontFamily string
	var variantNameIDs []uint16
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
										nameTable = make(map[uint16]string)
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
															if platformID == PLATFORM_ID_MICROSOFT {
																var currentPos int64
																if currentPos, err = file.Seek(0, 1); err == nil {
																	if _, err = file.Seek(int64(nameInfo.offset)+int64(stringOffset)+int64(offset), 0); err == nil {
																		stringData := make([]byte, length)
																		if _, err = file.Read(stringData); err == nil {
																			stringValue := f.decodeUTF16(stringData)
																			nameTable[nameID] = stringValue

																			if _, err = file.Seek(currentPos, 0); err != nil {
																				break
																			}
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
											if fvarInfo, ok := tables[TABLE_TAG_FVAR]; ok {
												if _, err = file.Seek(int64(fvarInfo.offset), 0); err == nil {
													fvarData := make([]byte, fvarInfo.length)
													if _, err = file.Read(fvarData); err == nil {
														axisCount := binary.BigEndian.Uint16(fvarData[FVAR_AXIS_COUNT_OFFSET : FVAR_AXIS_COUNT_OFFSET+FVAR_AXIS_COUNT_LENGTH])
														instanceCount := binary.BigEndian.Uint16(fvarData[FVAR_INSTANCE_COUNT_OFFSET : FVAR_INSTANCE_COUNT_OFFSET+FVAR_INSTANCE_COUNT_LENGTH])
														instanceOffset := 8 + int(axisCount)*16
														instances = make([]FONT_VARIABLE_INSTANCE, 0)
														for i := 0; i < int(instanceCount); i++ {
															offset := instanceOffset + i*(8+int(axisCount)*4)
															if offset+8+int(axisCount)*4 <= len(fvarData) {
																subfamilyNameID := binary.BigEndian.Uint16(fvarData[offset : offset+2])
																instances = append(instances, FONT_VARIABLE_INSTANCE{
																	subfamilyNameID: subfamilyNameID,
																	coordinates:     nil,
																})
															}
														}
													}
												}
											}
										}
										if err == nil {
											if familyName, ok := nameTable[NAME_ID_FAMILY]; ok {
												fontFamily = familyName
												if fontFamily != "" {
													variantNameIDs = make([]uint16, 0, LAST_VARIANT_NAME_ID-FIRST_VARIANT_NAME_ID+1)
													for id := uint16(FIRST_VARIANT_NAME_ID); id <= LAST_VARIANT_NAME_ID; id++ {
														variantNameIDs = append(variantNameIDs, id)
													}
													variantNames = make([]string, 0)
													for _, id := range variantNameIDs {
														if subfamilyName, ok := nameTable[id]; ok {
															fullName := fontFamily + " " + subfamilyName
															variantNames = append(variantNames, fullName)
														}
													}
													if len(variantNames) == 0 {
														for _, instance := range instances {
															if subfamilyName, ok := nameTable[instance.subfamilyNameID]; ok {
																fullName := fontFamily + " " + subfamilyName
																variantNames = append(variantNames, fullName)
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
	return variantNames, err
}
