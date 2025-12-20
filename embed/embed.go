package embed

import (
	"crypto/sha3"
	__embed "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xiang-tai-duo/go-boost/file"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage
const (
	DIRECTORY_PERMISSION  = 0755
	EXECUTABLE_PERMISSION = 0755
	GOOS_LINUX            = "linux"
	MODULE_NAME_EMBED     = "embed"
	ROOT_DIRECTORY        = "."
	SH_EXTENSION          = ".sh"
)

func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_EMBED, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_EMBED, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_EMBED, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_EMBED, logger.SKIP_STACK_FRAMES_BASE)
}

func GetAllEmbedFilesPath(embedFS fs.FS) ([]string, error) {
	result := make([]string, 0)
	err := error(nil)
	if err = fs.WalkDir(embedFS, ROOT_DIRECTORY, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr == nil && !d.IsDir() {
			result = append(result, path)
		}
		return walkErr
	}); err != nil {
		err = fmt.Errorf("failed to walk embedded filesystem: %v", err)
	}
	return result, err
}

func IsEmpty(embedFS fs.FS) (bool, error) {
	result := true
	err := error(nil)
	if err = fs.WalkDir(embedFS, ROOT_DIRECTORY, func(path string, d fs.DirEntry, walkErr error) error {
		walkErrResult := walkErr
		if walkErr == nil && path != ROOT_DIRECTORY {
			result = false
			walkErrResult = fs.SkipAll
		}
		return walkErrResult
	}); err == nil {
	} else if !errors.Is(err, fs.SkipAll) {
		err = fmt.Errorf("failed to walk embedded filesystem: %v", err)
	} else {
		err = nil
	}
	return result, err
}

func RestoreAll(embedFs __embed.FS, preserveEmbedRelativeDirectoryPath bool) (int, error) {
	result := 0
	err := error(nil)
	embedFilesPath := make([]string, 0)
	if embedFilesPath, err = getAllEmbedFilesPath(embedFs); err == nil {
		result, err = RestoreFiles(embedFs, embedFilesPath, preserveEmbedRelativeDirectoryPath)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func RestoreFile(embedFs __embed.FS, embedFilesPath string, preserveEmbedRelativeDirectoryPath bool) (int, error) {
	return RestoreFiles(embedFs, []string{embedFilesPath}, preserveEmbedRelativeDirectoryPath)
}

func RestoreFiles(embedFs __embed.FS, embedFilesPath []string, preserveEmbedRelativeDirectoryPath bool) (int, error) {
	result := 0
	err := error(nil)
	workingDirectory := ""
	__debug(fmt.Sprintf("RestoreFiles started: %d target files, preserveEmbedRelativeDirectoryPath=%v", len(embedFilesPath), preserveEmbedRelativeDirectoryPath))
	if workingDirectory, err = os.Getwd(); err == nil {
		__debug(fmt.Sprintf("working directory: %s", workingDirectory))
		filesPath := make([]string, 0)
		if filesPath, err = getAllEmbedFilesPath(embedFs); err == nil {
			__debug(fmt.Sprintf("found %d embedded files", len(filesPath)))
			for _, fileName := range embedFilesPath {
				__debug(fmt.Sprintf("processing target file: %s", fileName))
				isFileNotFound := true
				for _, relativeFilePath := range filesPath {
					if strings.Contains(relativeFilePath, fileName) {
						isFileNotFound = false
						__debug(fmt.Sprintf("matched embedded file: %s", relativeFilePath))
						absoluteFilePath := ""
						if preserveEmbedRelativeDirectoryPath {
							absoluteFilePath = filepath.Join(workingDirectory, relativeFilePath)
							dir := filepath.Dir(absoluteFilePath)
							__debug(fmt.Sprintf("creating directory: %s", dir))
							if err = os.MkdirAll(dir, DIRECTORY_PERMISSION); err == nil {
							} else {
								err = fmt.Errorf("failed to create directory %s: %v", dir, err)
								__debug(err.Error())
								break
							}
						} else {
							absoluteFilePath = filepath.Join(workingDirectory, filepath.Base(relativeFilePath))
						}
						__debug(fmt.Sprintf("target absolute path: %s", absoluteFilePath))
						isLatestVersion := false
						if _, statErr := os.Stat(absoluteFilePath); statErr == nil {
							__debug(fmt.Sprintf("target file exists, comparing hash: %s", absoluteFilePath))
							embedHash := ""
							fileHash := ""
							if embedHash, err = calculateEmbedFileHash(embedFs, relativeFilePath); err == nil {
								if fileHash, err = calculateFileHash(absoluteFilePath); err == nil {
									if embedHash == fileHash {
										isLatestVersion = true
										__debug(fmt.Sprintf("hash matched, skip restore: %s", absoluteFilePath))
									} else {
										__debug(fmt.Sprintf("hash mismatch, embed=%s file=%s", embedHash, fileHash))
									}
								} else {
									__debug(fmt.Sprintf("failed to calculate file hash %s: %v", absoluteFilePath, err))
								}
							} else {
								__debug(fmt.Sprintf("failed to calculate embed file hash %s: %v", relativeFilePath, err))
							}
						} else {
							__debug(fmt.Sprintf("target file does not exist, will restore: %s", absoluteFilePath))
						}
						if !isLatestVersion {
							__debug(fmt.Sprintf("copying embedded file %s to %s", relativeFilePath, absoluteFilePath))
							if err = copyFile(embedFs, relativeFilePath, absoluteFilePath); err == nil {
								result++
								__debug(fmt.Sprintf("restored file: %s", absoluteFilePath))
							} else {
								err = fmt.Errorf("failed to write file %s: %v", absoluteFilePath, err)
								__debug(err.Error())
								break
							}
						}
					}
				}
				if isFileNotFound {
					err = fmt.Errorf("embedded file not found: %s", fileName)
					__debug(err.Error())
					break
				}
			}
		} else {
			__debug(fmt.Sprintf("failed to get embedded files path: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("failed to get working directory: %v", err))
	}
	__debug(fmt.Sprintf("RestoreFiles finished: %d files restored, err=%v", result, err))
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,GoImportUsedAsName
func calculateEmbedFileHash(embedFs fs.FS, filePath string) (string, error) {
	result := ""
	err := error(nil)
	var file fs.File
	if file, err = embedFs.Open(filePath); err == nil {
		defer file.Close()
		result, err = calculateFileHashFromReader(file)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult,GoImportUsedAsName
func calculateFileHash(filePath string) (string, error) {
	result := ""
	err := error(nil)
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		result, err = calculateFileHashFromReader(file)
	}
	return result, err
}

func calculateFileHashFromReader(reader io.Reader) (string, error) {
	result := ""
	err := error(nil)
	hash := sha3.New256()
	if _, err = io.Copy(hash, reader); err == nil {
		result = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func copyFile(embedFs fs.FS, from string, to string) error {
	err := error(nil)
	var src fs.File
	if src, err = embedFs.Open(from); err == nil {
		defer src.Close()
		var dst *os.File
		if dst, err = os.Create(to); err == nil {
			defer dst.Close()
			if _, err = io.Copy(dst, src); err == nil {
				if runtime.GOOS == GOOS_LINUX && (file.IsExecutable(to) || strings.HasSuffix(strings.ToLower(to), SH_EXTENSION)) {
					dst.Chmod(EXECUTABLE_PERMISSION)
					__debug(fmt.Sprintf("setting executable permission for file: %s", to))
				}
			}
		}
	}
	return err
}

func getAllEmbedFilesPath(embedFS fs.FS) ([]string, error) {
	return GetAllEmbedFilesPath(embedFS)
}
