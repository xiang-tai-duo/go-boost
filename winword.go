// Package boost
// File:        winword.go
// Url:         `https://github.com/xiang-tai-duo/go-boost/blob/master/winword.go`
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Word document operations, providing functions to read text from DOCX files by extracting them.
// --------------------------------------------------------------------------------

package boost

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

//goland:noinspection GoSnakeCaseUsage
const (
	WORD_DIRECTORY_NAME = "word"
)

//goland:noinspection GoSnakeCaseUsage
var (
	WORD_DOCUMENT_XML_PATH = fmt.Sprintf("%s/%s", WORD_DIRECTORY_NAME, "document.xml")
)

//goland:noinspection GoSnakeCaseUsage
type (
	WORD_DOCUMENT_TEXT struct {
		XMLName xml.Name `xml:"t"`
		Text    string   `xml:",chardata"`
	}

	WORD_DOCUMENT_RUN struct {
		XMLName xml.Name             `xml:"r"`
		Texts   []WORD_DOCUMENT_TEXT `xml:"t"`
	}

	WORD_DOCUMENT_PARAGRAPH struct {
		XMLName xml.Name            `xml:"p"`
		Runs    []WORD_DOCUMENT_RUN `xml:"r"`
	}

	WORD_DOCUMENT_BODY struct {
		XMLName    xml.Name                  `xml:"body"`
		Paragraphs []WORD_DOCUMENT_PARAGRAPH `xml:"p"`
	}

	WORD_DOCUMENT_XML struct {
		XMLName xml.Name           `xml:"document"`
		Body    WORD_DOCUMENT_BODY `xml:"body"`
	}

	WORD_DOCUMENT struct {
		Word struct {
			Document WORD_DOCUMENT_XML
		}
	}
)

func NewWordDocument() *WORD_DOCUMENT {
	return &WORD_DOCUMENT{}
}

//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
func (wd *WORD_DOCUMENT) Load(filePath string) error {
	var err error
	absoluteFilePath := ""
	if filepath.IsAbs(filePath) {
		absoluteFilePath = filePath
	} else {
		currentWorkingDirectory := ""
		if currentWorkingDirectory, err = os.Getwd(); err == nil {
			absoluteFilePath = filepath.Join(currentWorkingDirectory, filePath)
		}
	}
	if err == nil {
		err = wd.extractDocumentXml(absoluteFilePath)
	}
	return err
}

func (xml *WORD_DOCUMENT_XML) Text() []string {
	var texts []string
	for _, paragraph := range xml.Body.Paragraphs {
		paragraphText := strings.Builder{}
		for _, run := range paragraph.Runs {
			for _, text := range run.Texts {
				paragraphText.WriteString(text.Text)
			}
		}
		texts = append(texts, paragraphText.String())
	}
	return texts
}

//goland:noinspection SpellCheckingInspection
func (wd *WORD_DOCUMENT) encoding(data []byte, encoding string) ([]byte, error) {
	result := make([]byte, len(data))
	var err error
	encoding = strings.TrimSpace(strings.ToUpper(encoding))
	if encoding == "UTF-8" || encoding == "UTF8" {
		if utf8.Valid(data) {
			result = data
		} else {
			err = fmt.Errorf("invalid UTF-8 data")
		}
	} else {
		switch encoding {
		case "GB2312":
			result, _, err = transform.Bytes(simplifiedchinese.HZGB2312.NewDecoder(), data)
		case "GBK":
			result, _, err = transform.Bytes(simplifiedchinese.GBK.NewDecoder(), data)
		case "GB18030":
			result, _, err = transform.Bytes(simplifiedchinese.GB18030.NewDecoder(), data)
		case "BIG5", "BIG5-HKSCS":
			result, _, err = transform.Bytes(traditionalchinese.Big5.NewDecoder(), data)
		case "ISO-8859-1", "LATIN1":
			result, _, err = transform.Bytes(charmap.ISO8859_1.NewDecoder(), data)
		case "ISO-8859-2", "LATIN2":
			result, _, err = transform.Bytes(charmap.ISO8859_2.NewDecoder(), data)
		case "WINDOWS-1250":
			result, _, err = transform.Bytes(charmap.Windows1250.NewDecoder(), data)
		case "WINDOWS-1251":
			result, _, err = transform.Bytes(charmap.Windows1251.NewDecoder(), data)
		case "WINDOWS-1252":
			result, _, err = transform.Bytes(charmap.Windows1252.NewDecoder(), data)
		case "WINDOWS-1253":
			result, _, err = transform.Bytes(charmap.Windows1253.NewDecoder(), data)
		case "WINDOWS-1254":
			result, _, err = transform.Bytes(charmap.Windows1254.NewDecoder(), data)
		case "WINDOWS-1255":
			result, _, err = transform.Bytes(charmap.Windows1255.NewDecoder(), data)
		case "WINDOWS-1256":
			result, _, err = transform.Bytes(charmap.Windows1256.NewDecoder(), data)
		case "WINDOWS-1257":
			result, _, err = transform.Bytes(charmap.Windows1257.NewDecoder(), data)
		case "WINDOWS-1258":
			result, _, err = transform.Bytes(charmap.Windows1258.NewDecoder(), data)
		default:
			err = fmt.Errorf("unsupported encoding: %s", encoding)
		}
		if err != nil {
			err = fmt.Errorf("failed to decode %s data: %w", encoding, err)
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
func (wd *WORD_DOCUMENT) extractDocumentXml(absoluteFilePath string) error {
	var err error
	content := make([]byte, 0)
	if content, err = wd.extractXml(absoluteFilePath, WORD_DOCUMENT_XML_PATH); err == nil {
		err = xml.Unmarshal(content, &wd.Word.Document)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
func (wd *WORD_DOCUMENT) extractXml(absoluteFilePath string, xmlPath string) ([]byte, error) {
	content := make([]byte, 0)
	var err error
	var file *os.File
	if file, err = os.Open(absoluteFilePath); err == nil {
		defer file.Close()
		var zipReader *zip.Reader
		var fileInfo os.FileInfo
		if fileInfo, err = file.Stat(); err == nil {
			if zipReader, err = zip.NewReader(file, fileInfo.Size()); err == nil {
				for _, zipFile := range zipReader.File {
					if strings.Contains(zipFile.Name, xmlPath) {
						var fileReader io.ReadCloser
						if fileReader, err = zipFile.Open(); err == nil {
							defer fileReader.Close()
							rawData := make([]byte, 0)
							if rawData, err = io.ReadAll(fileReader); err == nil {
								encoding := ""
								var parseErr error
								if encoding, parseErr = wd.parseXMLEncoding(rawData); parseErr != nil {
									err = fmt.Errorf("invalid XML file: %w", parseErr)
								} else {
									content, err = wd.encoding(rawData, encoding)
								}
							}
						}
						break
					}
				}
			}
		}
	}
	return content, err
}

func (wd *WORD_DOCUMENT) parseXMLEncoding(data []byte) (string, error) {
	encoding := ""
	var err error
	encodingRegex := regexp.MustCompile(`<\?xml[^>]*encoding\s*=\s*['"]([^'"]+)['"][^>]*\?>`)
	matches := encodingRegex.FindSubmatch(data)
	if len(matches) >= 2 {
		encoding = string(matches[1])
	} else {
		err = fmt.Errorf("invalid XML declaration: encoding not found")
	}
	return encoding, err
}
