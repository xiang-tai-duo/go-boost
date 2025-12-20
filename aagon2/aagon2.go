package aagon2

import (
	"golang.org/x/crypto/argon2"

	"github.com/denisbrodbeck/machineid"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,SpellCheckingInspection,GoUnusedConst
const (
	AAGON2_AES_KEY_LENGTH            = 32
	AAGON2_ITERATIONS                = 3
	AAGON2_MEMORY                    = 64 * 1024
	AAGON2_PROTECTED_ID_SALT         = "go-boost"
	AAGON2_PARALLELISM               = 4
	AAGON2_INTERNAL_SALT             = "go-boost.aagon2.secret.salt.v1"
	AAGON2_HARDWARE_FINGERPRINT_SALT = "go-boost.aagon2.hardware.fingerprint.salt.v1"
	UNKNOWN_MACHINE_ID               = "UNKNOWN_MACHINE_ID"
	MODULE_NAME_AAGON2               = "aagon2"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_AAGON2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_AAGON2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_AAGON2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func DeriveKey() []byte {
	result := make([]byte, 0)
	err := error(nil)
	machineID := ""
	if machineID, err = machineid.ID(); err != nil {
		machineID, err = machineid.ProtectedID(AAGON2_PROTECTED_ID_SALT)
		if err != nil || machineID == "" {
			machineID = UNKNOWN_MACHINE_ID
		}
	}
	password := []byte(machineID + AAGON2_HARDWARE_FINGERPRINT_SALT + AAGON2_INTERNAL_SALT)
	result = argon2.IDKey(password, []byte(AAGON2_INTERNAL_SALT), AAGON2_ITERATIONS, AAGON2_MEMORY, AAGON2_PARALLELISM, AAGON2_AES_KEY_LENGTH)
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_AAGON2, logger.SKIP_STACK_FRAMES_BASE)
}
