// Package filepath
// File:        filepath.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/filepath/filepath.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: FilePath provides utility methods for file path operations, including absolute paths, cleaning, and joining
// --------------------------------------------------------------------------------
package filepath

import (
	"io/fs"
	__filepath "path/filepath"
	"strings"
)

//goland:noinspection GoSnakeCaseUsage
type (
	FILE_PATH struct {
		path string
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(filePath string) *FILE_PATH {
	return &FILE_PATH{path: filePath}
}

func (f *FILE_PATH) Abs() (FILE_PATH, error) {
	var result FILE_PATH
	var err error
	if absPath, err := __filepath.Abs(f.path); err == nil {
		result = FILE_PATH{path: absPath}
	}
	return result, err
}

func (f *FILE_PATH) Base() string {
	return __filepath.Base(f.path)
}

func (f *FILE_PATH) Clean() FILE_PATH {
	return FILE_PATH{path: __filepath.Clean(f.path)}
}

func (f *FILE_PATH) Dir() string {
	return __filepath.Dir(f.path)
}

func (f *FILE_PATH) EvalSymlinks() (FILE_PATH, error) {
	var result FILE_PATH
	var err error
	if evalPath, err := __filepath.EvalSymlinks(f.path); err == nil {
		result = FILE_PATH{path: evalPath}
	}
	return result, err
}

func (f *FILE_PATH) Ext() string {
	return __filepath.Ext(f.path)
}

func (f *FILE_PATH) FromSlash() FILE_PATH {
	return FILE_PATH{path: __filepath.FromSlash(f.path)}
}

func (f *FILE_PATH) GetFileExtension() string {
	return __filepath.Ext(f.path)
}

func (f *FILE_PATH) GetFileName() string {
	return __filepath.Base(f.path)
}

func (f *FILE_PATH) GetFileNameWithoutExtension() string {
	filename := __filepath.Base(f.path)
	extension := __filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

func (f *FILE_PATH) GetParentDirectory() string {
	return __filepath.Dir(f.path)
}

func (f *FILE_PATH) Glob() ([]string, error) {
	return __filepath.Glob(f.path)
}

func (f *FILE_PATH) HasPrefix(prefix string) bool {
	return strings.HasPrefix(f.path, prefix)
}

func (f *FILE_PATH) IsAbs() bool {
	return __filepath.IsAbs(f.path)
}

func (f *FILE_PATH) IsLocal() bool {
	return __filepath.IsLocal(f.path)
}

func (f *FILE_PATH) Join(elem ...string) FILE_PATH {
	args := make([]string, len(elem)+1)
	args[0] = f.path
	copy(args[1:], elem)
	return FILE_PATH{path: __filepath.Join(args...)}
}

func (f *FILE_PATH) Localize() (FILE_PATH, error) {
	var result FILE_PATH
	var err error
	if localPath, err := __filepath.Localize(f.path); err == nil {
		result = FILE_PATH{path: localPath}
	}
	return result, err
}

func (f *FILE_PATH) Match(pattern string) (bool, error) {
	return __filepath.Match(pattern, f.path)
}

func (f *FILE_PATH) Path(filePath ...string) string {
	if len(filePath) > 0 {
		f.path = filePath[0]
	}
	return f.path
}

func (f *FILE_PATH) Rel(basepath string) (FILE_PATH, error) {
	var result FILE_PATH
	var err error
	if relPath, err := __filepath.Rel(basepath, f.path); err == nil {
		result = FILE_PATH{path: relPath}
	}
	return result, err
}

func (f *FILE_PATH) Split() (string, string) {
	return __filepath.Split(f.path)
}

func (f *FILE_PATH) SplitList() []string {
	return __filepath.SplitList(f.path)
}

func (f *FILE_PATH) ToSlash() FILE_PATH {
	return FILE_PATH{path: __filepath.ToSlash(f.path)}
}

func (f *FILE_PATH) VolumeName() string {
	return __filepath.VolumeName(f.path)
}

func (f *FILE_PATH) Walk(fn __filepath.WalkFunc) error {
	return __filepath.Walk(f.path, fn)
}

func (f *FILE_PATH) WalkDir(fn fs.WalkDirFunc) error {
	return __filepath.WalkDir(f.path, fn)
}
