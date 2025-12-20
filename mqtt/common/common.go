// Package mqttcommon
// File:        common.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/mqtt/common/common.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: MQTT common types and constants
// --------------------------------------------------------------------------------
package mqttcommon

import (
	"time"
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	DEFAULT_MQTT_QOS = 0
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	MQTT_MESSAGE struct {
		Topic            string    `json:"topic"`
		Payload          string    `json:"payload"`
		Timestamp        time.Time `json:"timestamp"`
		QualityOfService int       `json:"qos"`
		Retained         bool      `json:"retained"`
		Duplicate        bool      `json:"duplicate"`
	}
)
