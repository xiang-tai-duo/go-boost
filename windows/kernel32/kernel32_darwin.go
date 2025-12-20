//go:build darwin

package kernel32

import "fmt"

func CheckRemoteDebuggerPresent() bool {
	return false
}

func GetLongPathNameW(path string) (string, error) {
	return path, nil
}

func GetNativeSystemInfo() string {
	return ""
}

func GetWindowsDirectoryW() (string, error) {
	return "", fmt.Errorf("GetWindowsDirectoryW not implemented on macOS")
}

func IsDebuggerPresent() bool {
	return false
}

func SetFileCompression(filePath string, enable bool) error {
	return nil
}
