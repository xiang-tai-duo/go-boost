// Package zip
// File:        zip.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/zip/zip.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: ZIP provides utility methods for zip file operations, including listing and extracting files with optional password support
// --------------------------------------------------------------------------------
package zip

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/system"
	"github.com/yeka/zip"
)

type (
	ZIP_ENTRY_READ_CLOSER struct {
		Entry  io.ReadCloser
		Holder *zip.ReadCloser
	}
)

const (
	ERROR_CREATE_OUTPUT_DIRECTORY = "failed to create output directory: %s"
	ERROR_CREATE_TEMP_DIRECTORY   = "failed to create temporary directory"
	ERROR_INVALID_ZIP_ENTRY_CRC32 = "invalid zip entry crc32: %s"
	ERROR_INVALID_ZIP_ENTRY_SIZE  = "invalid zip entry size: %s"
	ERROR_UNSAFE_ZIP_ENTRY        = "unsafe zip entry path: %s"
	ERROR_WRITE_FILE              = "failed to write file: %s"
	MAX_UNIQUE_DIRECTORY_RETRIES  = 32
	MODULE_NAME_ZIP               = "zip"
	STREAM_BUFFER_SIZE            = 16 * 1024
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_ZIP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_ZIP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_ZIP, logger.SKIP_STACK_FRAMES_BASE)
}

func (instance *ZIP_ENTRY_READ_CLOSER) Close() error {
	result := error(nil)
	entryError := instance.Entry.Close()
	holderError := instance.Holder.Close()
	if entryError != nil {
		result = entryError
	} else {
		result = holderError
	}
	return result
}

func ExtractAll(path string, password *string) string {
	result := ""
	if path != "" {
		temporaryDirectory := system.CreateTemporaryDirectory()
		if temporaryDirectory != "" {
			var reader *zip.ReadCloser
			var err error
			if reader, err = zip.OpenReader(path); err == nil {
				defer reader.Close()
				if extractAllFiles(reader.File, password, temporaryDirectory) {
					result = temporaryDirectory
				} else {
					os.RemoveAll(temporaryDirectory)
				}
			} else {
				os.RemoveAll(temporaryDirectory)
			}
		}
	}
	return result
}

func ExtractFile(path string, innerFilePath string, password *string) string {
	return extractSingleFile(path, password, innerFilePath)
}

func ExtractFileFromZipByName(zipPath string, fileName string, outputDir *string) (string, error) {
	result := ""
	err := error(nil)
	if zipPath == "" || fileName == "" {
		err = os.ErrInvalid
	} else {
		var reader *zip.ReadCloser
		if reader, err = zip.OpenReader(zipPath); err == nil {
			defer reader.Close()
			found := false
			for _, file := range reader.File {
				if file.Name == fileName {
					if !file.FileInfo().IsDir() {
						targetDirectory := ""
						if outputDir != nil && *outputDir != "" {
							targetDirectory = *outputDir
							if err = os.MkdirAll(targetDirectory, os.ModePerm); err == nil {
							} else {
								err = fmt.Errorf(ERROR_CREATE_OUTPUT_DIRECTORY, err.Error())
								break
							}
						} else {
							targetDirectory = createUniqueTemporaryDirectory()
							if targetDirectory == "" {
								err = fmt.Errorf(ERROR_CREATE_TEMP_DIRECTORY)
								break
							}
						}
						outputPath := filepath.Join(targetDirectory, filepath.Base(file.Name))
						if _, joinErr := safeJoinPath(targetDirectory, filepath.Base(file.Name)); joinErr != nil {
							err = fmt.Errorf(ERROR_UNSAFE_ZIP_ENTRY, file.Name)
							break
						}
						if writeZipEntry(file, outputPath) {
							result = outputPath
						} else {
							if outputDir == nil {
								os.RemoveAll(targetDirectory)
							}
							err = fmt.Errorf(ERROR_WRITE_FILE, fileName)
						}
					}
					found = true
					break
				}
			}
			if !found && err == nil {
				err = os.ErrNotExist
			}
		}
	}
	return result, err
}

func ListFiles(path string, password *string) []*zip.File {
	return listFiles(path, password)
}

func ListFilesName(path string, password *string) []string {
	files := listFiles(path, password)
	result := make([]string, 0, len(files))
	for _, f := range files {
		result = append(result, f.Name)
	}
	return result
}

func OpenEntry(path string, innerFilePath string, password *string) (io.ReadCloser, int64, error) {
	var resultReader io.ReadCloser
	resultSize := int64(0)
	err := error(nil)
	if path == "" {
		err = os.ErrInvalid
	} else if innerFilePath == "" {
		err = os.ErrInvalid
	} else {
		var reader *zip.ReadCloser
		if reader, err = zip.OpenReader(path); err == nil {
			found := false
			for _, file := range reader.File {
				if file.Name == innerFilePath {
					found = true
					if file.FileInfo().IsDir() {
						err = os.ErrInvalid
					} else {
						if password != nil && *password != "" && file.IsEncrypted() {
							file.SetPassword(*password)
						}
						var entryReader io.ReadCloser
						if entryReader, err = file.Open(); err == nil {
							resultReader = &ZIP_ENTRY_READ_CLOSER{Entry: entryReader, Holder: reader}
							resultSize = int64(file.UncompressedSize64)
						}
					}
					break
				}
			}
			if !found && err == nil {
				err = os.ErrNotExist
			}
			if resultReader == nil {
				reader.Close()
			}
		}
	}
	return resultReader, resultSize, err
}

func OpenFirstEntry(path string, password *string) (io.ReadCloser, int64, string, uint32, error) {
	var resultReader io.ReadCloser
	resultSize := int64(0)
	resultFileName := ""
	resultCrc32 := uint32(0)
	err := error(nil)
	files := ListFiles(path, password)
	if len(files) == 0 {
		err = os.ErrNotExist
	} else {
		firstFile := files[0]
		resultFileName = firstFile.Name
		resultCrc32 = firstFile.CRC32
		resultReader, resultSize, err = OpenEntry(path, resultFileName, password)
	}
	return resultReader, resultSize, resultFileName, resultCrc32, err
}

func (instance *ZIP_ENTRY_READ_CLOSER) Read(p []byte) (int, error) {
	return instance.Entry.Read(p)
}

func StreamFile(path string, innerFilePath string, callback func(data []byte) bool, password *string) {
	if path != "" && innerFilePath != "" && callback != nil {
		var reader *zip.ReadCloser
		var err error
		if reader, err = zip.OpenReader(path); err == nil {
			defer reader.Close()
			for _, file := range reader.File {
				if file.Name == innerFilePath {
					if !file.FileInfo().IsDir() {
						streamZipEntry(file, callback, password)
					}
					break
				}
			}
		}
	}
}

func Validate(path string, password *string) error {
	result := error(nil)
	if path == "" {
		result = os.ErrInvalid
	} else {
		var reader *zip.ReadCloser
		var err error
		if reader, err = zip.OpenReader(path); err == nil {
			defer reader.Close()
			for _, file := range reader.File {
				if file.FileInfo().IsDir() {
					continue
				}
				if _, joinErr := safeJoinPath("zip", file.Name); joinErr != nil {
					result = joinErr
					break
				}
				if result = validateZipEntry(file, password); result != nil {
					break
				}
			}
		} else {
			result = err
		}
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_ZIP, logger.SKIP_STACK_FRAMES_BASE)
}

func createUniqueTemporaryDirectory() string {
	result := ""
	for index := 0; index < MAX_UNIQUE_DIRECTORY_RETRIES; index++ {
		uniqueDirectory := filepath.Join(os.TempDir(), uuid.New().String())
		if _, err := os.Stat(uniqueDirectory); err != nil && os.IsNotExist(err) {
			if err := os.MkdirAll(uniqueDirectory, os.ModePerm); err == nil {
				result = uniqueDirectory
				break
			}
		}
	}
	return result
}

func extractAllFiles(files []*zip.File, password *string, temporaryDirectory string) bool {
	result := true
	for _, file := range files {
		if password != nil && *password != "" && file.IsEncrypted() {
			file.SetPassword(*password)
		}
		outputPath, joinErr := safeJoinPath(temporaryDirectory, file.Name)
		if joinErr != nil {
			result = false
			break
		}
		if file.FileInfo().IsDir() {
			var err error
			if err = os.MkdirAll(outputPath, os.ModePerm); err != nil {
				result = false
				break
			}
		} else {
			parentDirectory := filepath.Dir(outputPath)
			var err error
			if err = os.MkdirAll(parentDirectory, os.ModePerm); err == nil {
				if !writeZipEntry(file, outputPath) {
					result = false
					break
				}
			} else {
				result = false
				break
			}
		}
	}
	return result
}

func extractSingleFile(path string, password *string, innerFilePath string) string {
	result := ""
	if path != "" && innerFilePath != "" {
		var reader *zip.ReadCloser
		var err error
		if reader, err = zip.OpenReader(path); err == nil {
			defer reader.Close()
			for _, file := range reader.File {
				if file.Name == innerFilePath {
					if !file.FileInfo().IsDir() {
						if password != nil && *password != "" && file.IsEncrypted() {
							file.SetPassword(*password)
						}
						uniqueDirectory := createUniqueTemporaryDirectory()
						if uniqueDirectory != "" {
							outputPath, joinErr := safeJoinPath(uniqueDirectory, filepath.Base(file.Name))
							if joinErr != nil {
								os.RemoveAll(uniqueDirectory)
								break
							}
							if writeZipEntry(file, outputPath) {
								result = outputPath
							} else {
								os.RemoveAll(uniqueDirectory)
							}
						}
					}
					break
				}
			}
		}
	}
	return result
}

func listFiles(path string, password *string) []*zip.File {
	result := make([]*zip.File, 0)
	if path != "" {
		var reader *zip.ReadCloser
		var err error
		if reader, err = zip.OpenReader(path); err == nil {
			defer reader.Close()
			result = make([]*zip.File, 0, len(reader.File))
			for _, file := range reader.File {
				if password != nil && *password != "" && file.IsEncrypted() {
					file.SetPassword(*password)
				}
				if !file.FileInfo().IsDir() {
					result = append(result, file)
				}
			}
		}
	}
	return result
}

func safeJoinPath(baseDir string, entryName string) (string, error) {
	cleanBase := filepath.Clean(baseDir)
	joined := filepath.Join(cleanBase, entryName)
	rel, err := filepath.Rel(cleanBase, joined)
	if err != nil {
		return "", fmt.Errorf(ERROR_UNSAFE_ZIP_ENTRY, entryName)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf(ERROR_UNSAFE_ZIP_ENTRY, entryName)
	}
	return joined, nil
}

func streamZipEntry(file *zip.File, callback func(data []byte) bool, password *string) {
	if password != nil && *password != "" && file.IsEncrypted() {
		file.SetPassword(*password)
	}
	var reader io.ReadCloser
	var err error
	if reader, err = file.Open(); err == nil {
		defer reader.Close()
		buffer := make([]byte, STREAM_BUFFER_SIZE)
		for {
			count, err := reader.Read(buffer)
			if count > 0 {
				chunk := make([]byte, count)
				copy(chunk, buffer[:count])
				if !callback(chunk) {
					break
				}
			}
			if err != nil {
				break
			}
		}
	}
}

func validateZipEntry(file *zip.File, password *string) error {
	result := error(nil)
	if password != nil && *password != "" && file.IsEncrypted() {
		file.SetPassword(*password)
	}
	var reader io.ReadCloser
	var err error
	if reader, err = file.Open(); err == nil {
		defer reader.Close()
		hash := crc32.NewIEEE()
		count := int64(0)
		if count, err = io.Copy(hash, reader); err == nil {
			if count != int64(file.UncompressedSize64) {
				result = fmt.Errorf(ERROR_INVALID_ZIP_ENTRY_SIZE, file.Name)
			} else if hash.Sum32() != file.CRC32 {
				result = fmt.Errorf(ERROR_INVALID_ZIP_ENTRY_CRC32, file.Name)
			}
		} else {
			result = err
		}
	} else {
		result = err
	}
	return result
}

func writeZipEntry(file *zip.File, outputPath string) bool {
	result := false
	var reader io.ReadCloser
	var err error
	if reader, err = file.Open(); err == nil {
		defer reader.Close()
		var writer *os.File
		if writer, err = os.Create(outputPath); err == nil {
			defer writer.Close()
			if _, err = io.Copy(writer, reader); err == nil {
				result = true
			}
		}
	}
	return result
}
