//go:build darwin

// File:        mpdefender_darwin.go
// Description: Darwin stub for defender package (no-op)
// --------------------------------------------------------------------------------

package defender

const (
	ERROR_MPDEFENDER_ADD_FAILED    = "failed to add Windows Defender exclusion"
	ERROR_MPDEFENDER_REMOVE_FAILED = "failed to remove Windows Defender exclusion"
	ERROR_MPDEFENDER_INVALID_PATH  = "path cannot be empty"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_MPDEFENDER = "windows.defender.mpdefender"
)

func AddFileToWhiteList(filePath string) error {
	return nil
}

func AddFolderToWhiteList(folderPath string) error {
	return nil
}

func RemoveFileFromWhiteList(filePath string) error {
	return nil
}

func RemoveFolderFromWhiteList(folderPath string) error {
	return nil
}
