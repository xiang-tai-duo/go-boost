// Package embed
// File:        embed.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/embed/embed.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Embed provides functionality for embedded file operations
// --------------------------------------------------------------------------------
package embed

import (
	"crypto/md5"
	__embed "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//goland:noinspection GoSnakeCaseUsage
const (
	ROOT_DIRECTORY = "."
)

func GetAllEmbedFilesPath(embedFS fs.FS) ([]string, error) {
	result := make([]string, 0)
	err := error(nil)
	if err = fs.WalkDir(embedFS, ROOT_DIRECTORY, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			result = append(result, path)
		}
		return err
	}); err != nil {
		err = fmt.Errorf("failed to walk embedded filesystem: %v", err)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func IsEmpty(embedFS fs.FS) (bool, error) {
	result := true
	err := fs.WalkDir(embedFS, ROOT_DIRECTORY, func(path string, d fs.DirEntry, walkErr error) error {
		err := walkErr
		if walkErr == nil && path != ROOT_DIRECTORY {
			result = false
			err = fs.SkipAll
		}
		return err
	})
	if err != nil && !errors.Is(err, fs.SkipAll) {
		err = fmt.Errorf("failed to walk embedded filesystem: %v", err)
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func RestoreAll(embedFs __embed.FS, preserveEmbedRelativeDirectoryPath bool) error {
	err := error(nil)
	embedFilesPath := make([]string, 0)
	if embedFilesPath, err = getAllEmbedFilesPath(embedFs); err == nil {
		err = RestoreFiles(embedFs, embedFilesPath, preserveEmbedRelativeDirectoryPath)
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func RestoreFile(embedFs __embed.FS, embedFilesPath string, preserveEmbedRelativeDirectoryPath bool) error {
	err := error(nil)
	err = RestoreFiles(embedFs, []string{embedFilesPath}, preserveEmbedRelativeDirectoryPath)
	return err
}

func RestoreFiles(embedFs __embed.FS, embedFilesPath []string, preserveEmbedRelativeDirectoryPath bool) error {
	err := error(nil)
	workingDirectory := ""
	if workingDirectory, err = os.Getwd(); err == nil {
		filesPath := make([]string, 0)
		if filesPath, err = getAllEmbedFilesPath(embedFs); err == nil {
			for _, fileName := range embedFilesPath {
				isFileNotFound := true
				for _, relativeFilePath := range filesPath {
					if strings.Contains(relativeFilePath, fileName) {
						isFileNotFound = false
						var absoluteFilePath string
						if preserveEmbedRelativeDirectoryPath {
							absoluteFilePath = filepath.Join(workingDirectory, relativeFilePath)
						} else {
							absoluteFilePath = filepath.Join(workingDirectory, filepath.Base(relativeFilePath))
						}
						if preserveEmbedRelativeDirectoryPath {
							dir := filepath.Dir(absoluteFilePath)
							if err = os.MkdirAll(dir, 0755); err != nil {
								err = fmt.Errorf("failed to create directory %s: %v", dir, err)
								break
							}
						}
						isLatestVersion := false
						if _, statErr := os.Stat(absoluteFilePath); statErr == nil {
							var embedHash, fileHash string
							if embedHash, err = calculateEmbedFileHash(embedFs, relativeFilePath); err == nil {
								if fileHash, err = calculateFileHash(absoluteFilePath); err == nil {
									if embedHash == fileHash {
										isLatestVersion = true
									}
								}
							}
						}
						if !isLatestVersion {
							if err = copyFile(embedFs, relativeFilePath, absoluteFilePath); err != nil {
								err = fmt.Errorf("failed to write file %s: %v", absoluteFilePath, err)
								break
							}
						}
					}
				}
				if isFileNotFound {
					err = fmt.Errorf("embedded file not found: %s", fileName)
					break
				}
			}
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
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

//goland:noinspection GoUnhandledErrorResult
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

func calculateFileHashFromReader(reader io.Reader) (string, error) {
	result := ""
	err := error(nil)
	hash := md5.New()
	if _, err = io.Copy(hash, reader); err == nil {
		result = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func copyFile(embedFs fs.FS, from string, to string) error {
	err := fmt.Errorf("copy file failed")
	var src fs.File
	if src, err = embedFs.Open(from); err == nil {
		defer src.Close()
		var dst *os.File
		if dst, err = os.Create(to); err == nil {
			defer dst.Close()
			if _, err = io.Copy(dst, src); err == nil {
				dst.Chmod(0644)
			}
		}
	}
	return err
}

func getAllEmbedFilesPath(embedFS fs.FS) ([]string, error) {
	return GetAllEmbedFilesPath(embedFS)
}
