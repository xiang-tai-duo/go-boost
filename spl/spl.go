// Package spl
// File:        spl.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/spl/spl.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SPL provides functions to read and modify PJL job attributes in spool/PRN files.
// --------------------------------------------------------------------------------
package spl

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	FILE_PERMISSION           = 0644
	LINE_FEED                 = "\n"
	MODULE_NAME_SPL           = "spl"
	PJL_CARRIAGE_RETURN       = "\r"
	PJL_CARRIAGE_RETURN_BYTE  = '\r'
	PJL_DOUBLE_QUOTE          = "\""
	PJL_ENTER_LANGUAGE_PREFIX = "@PJL ENTER LANGUAGE="
	PJL_JOBATTR_PREFIX_ACNA   = "@PJL SET JOBATTR=\"@ACNA="
	PJL_JOBATTR_PREFIX_JCID   = "@PJL SET JOBATTR=\"@JCID="
	PJL_LINE_ENDING_BYTE      = '\n'
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SPL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SPL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SPL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SPL, logger.SKIP_STACK_FRAMES_BASE)
}

func bytesTrimSuffixString(value string, suffix string) string {
	result := value
	__debug(fmt.Sprintf("bytesTrimSuffixString start: valueLength=%d, suffix=%q", len(value), suffix))
	if strings.HasSuffix(result, suffix) {
		result = result[:len(result)-len(suffix)]
		__debug(fmt.Sprintf("bytesTrimSuffixString suffix trimmed: resultLength=%d", len(result)))
	} else {
		__debug("bytesTrimSuffixString suffix not found")
	}
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func GetJOBATTR(filePath string, attributePrefix string) string {
	result := ""
	content := make([]byte, 0)
	readErr := error(nil)
	__debug(fmt.Sprintf("GetJOBATTR start: filePath=%s, prefix=%q", filePath, attributePrefix))
	if content, readErr = os.ReadFile(filePath); readErr == nil {
		__debug(fmt.Sprintf("GetJOBATTR read file success: bytes=%d", len(content)))
		prefix := []byte(attributePrefix)
		attributeStart := bytes.Index(content, prefix)
		__debug(fmt.Sprintf("GetJOBATTR prefix search result: index=%d", attributeStart))
		if attributeStart >= 0 {
			attributeEnd := bytes.IndexByte(content[attributeStart:], PJL_LINE_ENDING_BYTE)
			if attributeEnd >= 0 {
				attributeEnd = attributeStart + attributeEnd
				__debug(fmt.Sprintf("GetJOBATTR line end found: start=%d, end=%d", attributeStart, attributeEnd))
			} else {
				attributeEnd = len(content)
				__debug(fmt.Sprintf("GetJOBATTR line end not found, using content end: start=%d, end=%d", attributeStart, attributeEnd))
			}
			attributeLine := content[attributeStart:attributeEnd]
			attributeValue := string(bytes.TrimPrefix(attributeLine, prefix))
			__debug(fmt.Sprintf("GetJOBATTR raw value: %s", attributeValue))
			attributeValue = bytesTrimSuffixString(attributeValue, PJL_CARRIAGE_RETURN)
			attributeValue = bytesTrimSuffixString(attributeValue, PJL_DOUBLE_QUOTE)
			if attributeValue != "" {
				result = attributeValue
				__debug(fmt.Sprintf("GetJOBATTR parsed value: %s", result))
			} else {
				__debug("GetJOBATTR value is empty after trimming")
			}
		}
	} else {
		__error(fmt.Sprintf("GetJOBATTR os.ReadFile failed: filePath=%s, error=%v", filePath, readErr))
	}
	if result == "" {
		__debug(fmt.Sprintf("GetJOBATTR value not found: filePath=%s, prefix=%q", filePath, attributePrefix))
	} else {
		__debug(fmt.Sprintf("GetJOBATTR value found: filePath=%s, prefix=%q, value=%s", filePath, attributePrefix, result))
	}
	return result
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func GetACNA(filePath string) string {
	return GetJOBATTR(filePath, PJL_JOBATTR_PREFIX_ACNA)
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func GetJCID(filePath string) string {
	return GetJOBATTR(filePath, PJL_JOBATTR_PREFIX_JCID)
}

//goland:noinspection SpellCheckingInspection,GoUnusedExportedFunction
func SetJCID(prnFilePath string, jcid string) error {
	err := error(nil)
	content := make([]byte, 0)
	__debug(fmt.Sprintf("prnFilePath=%s, jcid=%s", prnFilePath, jcid))
	if content, err = os.ReadFile(prnFilePath); err == nil {
		__debug(fmt.Sprintf("Read file success: bytes=%d", len(content)))
		command := PJL_JOBATTR_PREFIX_JCID + jcid + PJL_DOUBLE_QUOTE
		prefix := []byte(PJL_JOBATTR_PREFIX_JCID)
		jcidStart := bytes.Index(content, prefix)
		updated := make([]byte, 0, len(content)+len(command)+2)
		if jcidStart >= 0 {
			jcidEnd := bytes.IndexByte(content[jcidStart:], PJL_LINE_ENDING_BYTE)
			if jcidEnd >= 0 {
				jcidEnd = jcidStart + jcidEnd
				if jcidEnd > jcidStart && content[jcidEnd-1] == PJL_CARRIAGE_RETURN_BYTE {
					jcidEnd--
				}
			} else {
				jcidEnd = len(content)
			}
			__debug(fmt.Sprintf("Existing JCID found, replace in place: start=%d, end=%d", jcidStart, jcidEnd))
			updated = append(updated, content[:jcidStart]...)
			updated = append(updated, []byte(command)...)
			updated = append(updated, content[jcidEnd:]...)
			err = os.WriteFile(prnFilePath, updated, FILE_PERMISSION)
		} else {
			enterLanguageStart := bytes.Index(content, []byte(PJL_ENTER_LANGUAGE_PREFIX))
			if enterLanguageStart >= 0 {
				__debug(fmt.Sprintf("JCID not found, insert before ENTER LANGUAGE: index=%d", enterLanguageStart))
				insertion := command + LINE_FEED
				updated = append(updated, content[:enterLanguageStart]...)
				updated = append(updated, []byte(insertion)...)
				updated = append(updated, content[enterLanguageStart:]...)
				err = os.WriteFile(prnFilePath, updated, FILE_PERMISSION)
			} else {
				err = fmt.Errorf("JCID attribute nor ENTER LANGUAGE command found in PRN file: %s", prnFilePath)
			}
		}
		if err == nil {
			__debug(fmt.Sprintf("Write file success: prnFilePath=%s, bytes=%d", prnFilePath, len(updated)))
		}
	} else {
		__debug(fmt.Sprintf("os.ReadFile failed: prnFilePath=%s, error=%v", prnFilePath, err))
	}
	__debug(fmt.Sprintf("prnFilePath=%s, error=%v", prnFilePath, err))
	return err
}
