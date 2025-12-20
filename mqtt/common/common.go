// Package mqttcommon
// File:        common.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mqtt/common/common.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: MQTT common types and constants
// --------------------------------------------------------------------------------

package mqttcommon

import (
	"time"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	MQTT_MESSAGE struct {
		Duplicate        bool      `json:"duplicate"`
		Payload          string    `json:"payload"`
		QualityOfService int       `json:"qos"`
		Retained         bool      `json:"retained"`
		Timestamp        time.Time `json:"timestamp"`
		Topic            string    `json:"topic"`
	}
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	DEFAULT_MQTT_QUALITY_OF_SERVICE = 0
	MODULE_NAME_COMMON              = "mqtt.common"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_COMMON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_COMMON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_COMMON, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_COMMON, logger.SKIP_STACK_FRAMES_BASE)
}
