// Package boost
// File:        filepath.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: FilePath provides utility methods for file path operations,
//
//	while FilePathBuilder is used for operations on specific file paths.
package boost

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type (
	FILEPATH struct {
		_filepath string
	}
)

func NewFilePath(filePath string) *FILEPATH {
	return &FILEPATH{_filepath: filePath}
}

func (f FILEPATH) Abs() (FILEPATH, error) {
	var result FILEPATH
	var err error
	if absPath, err := filepath.Abs(f._filepath); err == nil {
		result = FILEPATH{_filepath: absPath}
	}
	return result, err
}

func (f FILEPATH) Base() string {
	return filepath.Base(f._filepath)
}

func (f FILEPATH) Clean() FILEPATH {
	return FILEPATH{_filepath: filepath.Clean(f._filepath)}
}

func (f FILEPATH) Dir() string {
	return filepath.Dir(f._filepath)
}

func (f FILEPATH) EvalSymlinks() (FILEPATH, error) {
	var result FILEPATH
	var err error
	if evalPath, err := filepath.EvalSymlinks(f._filepath); err == nil {
		result = FILEPATH{_filepath: evalPath}
	}
	return result, err
}

func (f FILEPATH) Ext() string {
	return filepath.Ext(f._filepath)
}

func (f FILEPATH) FromSlash() FILEPATH {
	return FILEPATH{_filepath: filepath.FromSlash(f._filepath)}
}

func (f FILEPATH) GetFileExtension() string {
	return filepath.Ext(f._filepath)
}

func (f FILEPATH) GetFileName() string {
	return filepath.Base(f._filepath)
}

func (f FILEPATH) GetFileNameWithoutExtension() string {
	filename := filepath.Base(f._filepath)
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

func (f FILEPATH) GetParentDirectory() string {
	return filepath.Dir(f._filepath)
}

func (f FILEPATH) Glob() ([]string, error) {
	return filepath.Glob(f._filepath)
}

func (f FILEPATH) HasPrefix(prefix string) bool {
	return filepath.HasPrefix(f._filepath, prefix)
}

func (f FILEPATH) IsAbs() bool {
	return filepath.IsAbs(f._filepath)
}

func (f FILEPATH) IsLocal() bool {
	return filepath.IsLocal(f._filepath)
}

func (f FILEPATH) Join(elem ...string) FILEPATH {
	args := make([]string, len(elem)+1)
	args[0] = f._filepath
	copy(args[1:], elem)
	return FILEPATH{_filepath: filepath.Join(args...)}
}

func (f FILEPATH) Localize() (FILEPATH, error) {
	var result FILEPATH
	var err error
	if localPath, err := filepath.Localize(f._filepath); err == nil {
		result = FILEPATH{_filepath: localPath}
	}
	return result, err
}

func (f FILEPATH) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, f._filepath)
}

func (f *FILEPATH) Path(filePath ...string) string {
	if len(filePath) > 0 {
		f._filepath = filePath[0]
	}
	return f._filepath
}

func (f FILEPATH) Rel(basepath string) (FILEPATH, error) {
	var result FILEPATH
	var err error
	if relPath, err := filepath.Rel(basepath, f._filepath); err == nil {
		result = FILEPATH{_filepath: relPath}
	}
	return result, err
}

func (f FILEPATH) Split() (string, string) {
	return filepath.Split(f._filepath)
}

func (f FILEPATH) SplitList() []string {
	return filepath.SplitList(f._filepath)
}

func (f FILEPATH) ToSlash() FILEPATH {
	return FILEPATH{_filepath: filepath.ToSlash(f._filepath)}
}

func (f FILEPATH) VolumeName() string {
	return filepath.VolumeName(f._filepath)
}

func (f FILEPATH) Walk(fn filepath.WalkFunc) error {
	return filepath.Walk(f._filepath, fn)
}

func (f FILEPATH) WalkDir(fn fs.WalkDirFunc) error {
	return filepath.WalkDir(f._filepath, fn)
}
