// Package ghostscript
// File:        ghostscript.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ghostscript/ghostscript.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: GhostScript wrapper for PCL/PRN to PDF/PNG conversion
// --------------------------------------------------------------------------------

package ghostscript

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xiang-tai-duo/go-boost/embed"
	"github.com/xiang-tai-duo/go-boost/file"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/osutil"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	COMMAND_FAILED_FORMAT             = "GhostPCL invocation failed: %v%s, executable: %s, subprocess output: %q"
	DEVICE_PDF_IMAGE24                = "pdfimage24"
	DEVICE_PNG_16M                    = "png16m"
	ERROR_FAILED_TO_RESTORE_FORMAT    = "failed to restore embedded GhostPCL distribution files: %v"
	ERROR_INPUT_EXTENSION_NOT_ALLOWED = "input file extension not allowed (only .prn is supported): %s"
	ERROR_INPUT_FILE_NOT_PRN          = "input file is not a valid PRN file (PCL magic header missing): %s"
	ERROR_EXECUTABLE_PERMISSION       = "failed to set executable permission: %v"
	GPCL_NOT_INITIALIZED_ERROR        = "GhostPCL is not initialized (failed to restore embedded distribution files; please check directory permissions and disk space)"
	LOG_FAILED_TO_RESTORE_FILES       = "Failed to restore embedded GhostPCL distribution files: %v"
	LOG_FAILED_TO_SET_PERMISSION      = "Failed to set executable permission: %v"
	MAGIC_HEADER_PRN                  = "\x1B%-12345X"
	MODULE_NAME_GHOSTSCRIPT           = "ghostscript"
	OPT_D_NOCACHE                     = "-dNOCACHE"
	OPT_D_NOPAUSE                     = "-dNOPAUSE"
	OPT_R_600                         = "-r600"
	OPT_S_DEVICE                      = "-sDEVICE="
	OPT_S_OUTPUT                      = "-sOutputFile="
	RESULT_DEFAULT                    = -1
	RESULT_SUCCESS                    = 0
)

// Platform-specific constants
var (
	GHOSTPCL_DISTRIBUTION_DIRECTORY string
	GPCL_EXE_NAME                   string
)

//goland:noinspection GoUnusedGlobalVariable
var (
	initialized            bool
	allowedInputExtensions = []string{file.Extensions.PRN}
)

func init() {
	// Initialize platform-specific constants
	if runtime.GOOS == "windows" {
		GHOSTPCL_DISTRIBUTION_DIRECTORY = "ghostpcl-10.07.1-win32"
		GPCL_EXE_NAME = "gpcl6win32.exe"
	} else {
		GHOSTPCL_DISTRIBUTION_DIRECTORY = "ghostpdl-10.07.1/bin"
		GPCL_EXE_NAME = "gpcl6"
	}

	err := error(nil)
	isEmpty := false
	if isEmpty, err = embed.IsEmpty(GHOSTSCRIPT_DISTRIBUTION_FILES); err == nil && !isEmpty {
		if err = restoreFiles(); err != nil {
			__error(fmt.Sprintf(LOG_FAILED_TO_RESTORE_FILES, err))
		} else if runtime.GOOS != "windows" {
			// Set executable permissions on Linux
			if err = setExecutablePermissions(); err != nil {
				__error(fmt.Sprintf(LOG_FAILED_TO_SET_PERMISSION, err))
			}
		}
	}
	if err == nil {
		initialized = true
	}
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_GHOSTSCRIPT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_GHOSTSCRIPT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_GHOSTSCRIPT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_GHOSTSCRIPT, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection,DuplicatedCode
func ConvertPRNFileToPDFFile(prnFilePath string, pdfFilePath string) (int, error) {
	result := RESULT_DEFAULT
	err := error(nil)
	isPRNFile := false
	if initialized {
		if !isAllowedInputExtension(prnFilePath) {
			err = fmt.Errorf(ERROR_INPUT_EXTENSION_NOT_ALLOWED, prnFilePath)
		} else if isPRNFile, err = IsPRNFile(prnFilePath); err == nil {
			if isPRNFile {
				exeAbsPath := ""
				if exeAbsPath, err = filepath.Abs(filepath.Join(GHOSTPCL_DISTRIBUTION_DIRECTORY, GPCL_EXE_NAME)); err == nil {
					if _, err = os.Stat(exeAbsPath); err == nil {
						cmd := exec.Command(
							exeAbsPath,
							OPT_D_NOPAUSE,
							OPT_S_DEVICE+DEVICE_PDF_IMAGE24,
							OPT_D_NOCACHE,
							OPT_R_600,
							OPT_S_OUTPUT+pdfFilePath,
							prnFilePath,
						)
						osutil.SetHideWindow(cmd)
						output := make([]byte, 0)
						if output, err = cmd.CombinedOutput(); err == nil {
							result = RESULT_SUCCESS
						} else {
							err = fmt.Errorf(COMMAND_FAILED_FORMAT, err, describeExitCodeHint(err), exeAbsPath, string(output))
						}
					}
				}
			} else {
				err = fmt.Errorf(ERROR_INPUT_FILE_NOT_PRN, prnFilePath)
			}
		}
	} else {
		err = fmt.Errorf(GPCL_NOT_INITIALIZED_ERROR)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection,DuplicatedCode
func ConvertPRNToPNG(prnFilePath string, pngFilePattern string) (int, error) {
	result := RESULT_DEFAULT
	err := error(nil)
	if initialized {
		if !isAllowedInputExtension(prnFilePath) {
			err = fmt.Errorf(ERROR_INPUT_EXTENSION_NOT_ALLOWED, prnFilePath)
		} else {
			absoluteGhostscriptFilePath := ""
			if absoluteGhostscriptFilePath, err = filepath.Abs(filepath.Join(GHOSTPCL_DISTRIBUTION_DIRECTORY, GPCL_EXE_NAME)); err == nil {
				if _, err = os.Stat(absoluteGhostscriptFilePath); err == nil {
					cmd := exec.Command(
						absoluteGhostscriptFilePath,
						OPT_D_NOPAUSE,
						OPT_S_DEVICE+DEVICE_PNG_16M,
						OPT_D_NOCACHE,
						OPT_R_600,
						OPT_S_OUTPUT+pngFilePattern,
						prnFilePath,
					)
					osutil.SetHideWindow(cmd)
					output := make([]byte, 0)
					if output, err = cmd.CombinedOutput(); err == nil {
						result = RESULT_SUCCESS
					} else {
						err = fmt.Errorf(COMMAND_FAILED_FORMAT, err, describeExitCodeHint(err), absoluteGhostscriptFilePath, string(output))
					}
				}
			}
		}
	} else {
		err = fmt.Errorf(GPCL_NOT_INITIALIZED_ERROR)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection,GoUnhandledErrorResult
func IsPRNFile(filePath string) (bool, error) {
	result := false
	err := error(nil)
	prnFile := (*os.File)(nil)
	if prnFile, err = os.Open(filePath); err == nil {
		defer prnFile.Close()
		magicHeader := []byte(MAGIC_HEADER_PRN)
		buffer := make([]byte, len(magicHeader))
		if _, err = io.ReadFull(prnFile, buffer); err == nil {
			result = string(buffer) == MAGIC_HEADER_PRN
		} else if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			err = nil
		}
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func describeExitCodeHint(cmdErr error) string {
	if runtime.GOOS == "windows" {
		return describeWindowsExitCodeHint(cmdErr)
	}
	return describeLinuxExitCodeHint(cmdErr)
}

//goland:noinspection SpellCheckingInspection
func describeWindowsExitCodeHint(cmdErr error) string {
	hint := ""
	exitErr := (*exec.ExitError)(nil)
	if errors.As(cmdErr, &exitErr) && exitErr != nil && exitErr.ProcessState != nil {
		code := uint32(exitErr.ProcessState.ExitCode())
		switch code {
		case 0xC0000135:
			hint = " (meaning: STATUS_DLL_NOT_FOUND, the subprocess failed to load a required DLL; common causes: the target machine is missing the Microsoft Visual C++ 2015-2022 Redistributable (x86), or the ghostpcl-10.07.1-win32 directory is missing dependencies such as vcruntime140.dll / msvcp140.dll / gpcl6dll32.dll)"
		case 0xC0000142:
			hint = " (meaning: STATUS_DLL_INIT_FAILED, a DLL loaded by the subprocess failed to initialize)"
		case 0xC000007B:
			hint = " (meaning: STATUS_INVALID_IMAGE_FORMAT, the executable or one of its DLLs has a mismatched bitness; please verify 32/64-bit runtime alignment)"
		case 0xC0000005:
			hint = " (meaning: STATUS_ACCESS_VIOLATION, the subprocess triggered a memory access violation)"
		}
	}
	return hint
}

func describeLinuxExitCodeHint(cmdErr error) string {
	hint := ""
	exitErr := (*exec.ExitError)(nil)
	if errors.As(cmdErr, &exitErr) && exitErr != nil && exitErr.ProcessState != nil {
		code := exitErr.ProcessState.ExitCode()
		switch code {
		case 126:
			hint = " (meaning: Cannot execute command, possible cause: missing executable permission or not a binary)"
		case 127:
			hint = " (meaning: Command not found, possible cause: the executable path is incorrect)"
		case 139:
			hint = " (meaning: Segmentation fault, the subprocess triggered a memory access violation)"
		}
	}
	return hint
}

func isAllowedInputExtension(filePath string) bool {
	result := false
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, allowed := range allowedInputExtensions {
		if ext == allowed {
			result = true
			break
		}
	}
	return result
}

func restoreFiles() error {
	err := error(nil)
	if _, err = embed.RestoreAll(GHOSTSCRIPT_DISTRIBUTION_FILES, true); err != nil {
		err = fmt.Errorf(ERROR_FAILED_TO_RESTORE_FORMAT, err)
	}
	return err
}

func setExecutablePermissions() error {
	// Set executable permissions for all binaries in the distribution directory
	binaries := []string{"gpcl6", "gpdf", "gpdl", "gs", "gxps"}
	for _, bin := range binaries {
		binPath := filepath.Join(GHOSTPCL_DISTRIBUTION_DIRECTORY, bin)
		if _, err := os.Stat(binPath); err == nil {
			if err := os.Chmod(binPath, 0755); err != nil {
				return fmt.Errorf(ERROR_EXECUTABLE_PERMISSION, err)
			}
		}
	}
	return nil
}
