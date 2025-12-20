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
