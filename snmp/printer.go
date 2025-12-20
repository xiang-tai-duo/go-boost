// Package snmp
// File:        printer.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Printer provides SNMP-based printer status query functions (toner levels, error states, page counts, etc.).
// --------------------------------------------------------------------------------
package snmp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	PRINTER_ALERT struct {
		Description  string
		Severity     int
		SeverityName string
	}
	PRINTER_INFORMATION struct {
		Alerts           []PRINTER_ALERT
		DeviceStatus     int
		DeviceStatusName string
		Model            string
		Name             string
		Supplies         []PRINTER_SUPPLY
		TotalPages       int64
	}
	PRINTER_SUPPLY struct {
		Description     string
		Index           int
		Level           int
		MaximumCapacity int
		Type            int
		TypeName        string
	}
	SUPPLY_DATA struct {
		description     string
		index           int
		level           int
		maximumcapacity int
		supplytype      int
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	ALERT_KEYWORD_BIN                       = "bin"
	ALERT_KEYWORD_CASSETTE                  = "cassette"
	ALERT_KEYWORD_COVER                     = "cover"
	ALERT_KEYWORD_DOOR                      = "door"
	ALERT_KEYWORD_DRAWER                    = "drawer"
	ALERT_KEYWORD_DRUM                      = "drum"
	ALERT_KEYWORD_EMPTY                     = "empty"
	ALERT_KEYWORD_ERROR                     = "error"
	ALERT_KEYWORD_FAILURE                   = "failure"
	ALERT_KEYWORD_FAULT                     = "fault"
	ALERT_KEYWORD_FEED                      = "feed"
	ALERT_KEYWORD_FUSER                     = "fuser"
	ALERT_KEYWORD_IMAGING                   = "imaging"
	ALERT_KEYWORD_INK                       = "ink"
	ALERT_KEYWORD_JAM                       = "jam"
	ALERT_KEYWORD_LOW                       = "low"
	ALERT_KEYWORD_MEDIA                     = "media"
	ALERT_KEYWORD_MISMATCH                  = "mismatch"
	ALERT_KEYWORD_OUT                       = "out"
	ALERT_KEYWORD_OUTPUT                    = "output"
	ALERT_KEYWORD_PAPER                     = "paper"
	ALERT_KEYWORD_PHOTOCONDUCTOR            = "photoconductor"
	ALERT_KEYWORD_SIZE                      = "size"
	ALERT_KEYWORD_STACKER                   = "stacker"
	ALERT_KEYWORD_TONER                     = "toner"
	ALERT_KEYWORD_TRAY                      = "tray"
	ALERT_KEYWORD_TYPE                      = "type"
	ALERT_KEYWORD_WASTE                     = "waste"
	ALERT_SEVERITY_CRITICAL                 = 2
	ALERT_SEVERITY_OK                       = 0
	ALERT_SEVERITY_WARNING                  = 1
	DESCRIPTION_UNKNOW                      = "UNKNOWN"
	ERROR_MESSAGE_COVER_OPEN                = "COVER OPEN"
	ERROR_MESSAGE_DRUM_ABNORMAL             = "DRUM ABNORMAL"
	ERROR_MESSAGE_FUSER_ABNORMAL            = "FUSER ABNORMAL"
	ERROR_MESSAGE_INK_ABNORMAL              = "INK ABNORMAL"
	ERROR_MESSAGE_INK_LOW                   = "INK LOW"
	ERROR_MESSAGE_OUT_OF_INK                = "OUT OF INK"
	ERROR_MESSAGE_OUT_OF_PAPER              = "OUT OF PAPER"
	ERROR_MESSAGE_OUT_OF_TONER              = "OUT OF TONER"
	ERROR_MESSAGE_OUTPUT_ABNORMAL           = "OUTPUT ABNORMAL"
	ERROR_MESSAGE_PAPER_ABNORMAL            = "PAPER ABNORMAL"
	ERROR_MESSAGE_PAPER_JAM                 = "PAPER JAM"
	ERROR_MESSAGE_PAPER_MISMATCH            = "PAPER MISMATCH"
	ERROR_MESSAGE_PRINTER_ALERT_PREFIX      = "PRINTER ALERT: "
	ERROR_MESSAGE_PRINTER_FAULT_PREFIX      = "PRINTER FAULT: "
	ERROR_MESSAGE_PRINTER_HARDWARE_FAILURE  = "PRINTER HARDWARE FAILURE"
	ERROR_MESSAGE_PRINTER_WARNING_STATE     = "PRINTER IN WARNING STATE"
	ERROR_MESSAGE_SUPPLY_DEPLETED_PREFIX    = "SUPPLY DEPLETED: "
	ERROR_MESSAGE_TONER_ABNORMAL            = "TONER ABNORMAL"
	ERROR_MESSAGE_TONER_LOW                 = "TONER LOW"
	ERROR_MESSAGE_TRAY_ABNORMAL             = "TRAY ABNORMAL"
	ERROR_MESSAGE_WASTE_TONER_BOX_ABNORMAL  = "WASTE TONER BOX ABNORMAL"
	HOST_RESOURCES_DEVICE_DOWN              = 5
	HOST_RESOURCES_DEVICE_RUNNING           = 2
	HOST_RESOURCES_DEVICE_TESTING           = 4
	HOST_RESOURCES_DEVICE_WARNING           = 3
	OID_HOST_RESOURCES_DEVICE_STATUS        = ".1.3.6.1.2.1.25.3.2.1.5.1"
	OID_PRINTER_ALERT_DESCRIPTION           = ".1.3.6.1.2.1.43.18.1.1.4.1"
	OID_PRINTER_ALERT_SEVERITY_LEVEL        = ".1.3.6.1.2.1.43.18.1.1.7.1"
	OID_PRINTER_GENERAL_PRINTER_NAME        = ".1.3.6.1.2.1.43.5.1.1.16.1"
	OID_PRINTER_MARKER_COLORANT_VALUE       = ".1.3.6.1.2.1.43.12.1.1.4.1"
	OID_PRINTER_MARKER_COUNTER_LIFE         = ".1.3.6.1.2.1.43.10.2.1.4.1"
	OID_PRINTER_MARKER_LIFE_COUNT           = ".1.3.6.1.2.1.43.10.2.1.4.1"
	OID_PRINTER_MARKER_SUPPLIES_DESCRIPTION = ".1.3.6.1.2.1.43.11.1.1.6.1"
	OID_PRINTER_MARKER_SUPPLIES_LEVEL       = ".1.3.6.1.2.1.43.11.1.1.9.1"
	OID_PRINTER_MARKER_SUPPLIES_MAXIMUM     = ".1.3.6.1.2.1.43.11.1.1.8.1"
	OID_PRINTER_MARKER_SUPPLIES_TYPE        = ".1.3.6.1.2.1.43.11.1.1.5.1"
	OID_SYSTEM_DESCRIPTION                  = ".1.3.6.1.2.1.1.1.0"
	PARSE_LAST_INDEX_INVALID                = -1
	SUPPLY_LEVEL_EMPTY                      = 0
	SUPPLY_LEVEL_PERCENTAGE_CRITICAL        = 5
	SUPPLY_TYPE_DEVELOPER                   = 8
	SUPPLY_TYPE_INK                         = 5
	SUPPLY_TYPE_TONER                       = 3
	SUPPLY_TYPE_WASTE_TONER                 = 4
)

//goland:noinspection GoUnusedExportedFunction
func AlertSeverityName(severity int) string {
	result := ""
	switch severity {
	case ALERT_SEVERITY_OK:
		result = "OK"
	case ALERT_SEVERITY_WARNING:
		result = "Warning"
	case ALERT_SEVERITY_CRITICAL:
		result = "Critical"
	default:
		result = fmt.Sprintf("Unknown(%d)", severity)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func addWalkData(
	supplyMap map[int]*SUPPLY_DATA,
	client *SNMP_CLIENT,
	oid string,
	setter func(data *SUPPLY_DATA, value interface{}),
) {
	var walkResultList []SNMP_RESULT
	var err error
	if walkResultList, err = client.WalkAll(oid); err == nil {
		for _, walkResult := range walkResultList {
			var index int
			index = parseLastIndex(walkResult.Name)
			if index >= 0 {
				if _, ok := supplyMap[index]; !ok {
					supplyMap[index] = &SUPPLY_DATA{index: index}
				}
				setter(supplyMap[index], walkResult.Value)
			}
		}
	}
}

//goland:noinspection GoUnusedExportedFunction
func DeviceStatusName(status int) string {
	result := ""
	switch status {
	case HOST_RESOURCES_DEVICE_RUNNING:
		result = "Running"
	case HOST_RESOURCES_DEVICE_WARNING:
		result = "Warning"
	case HOST_RESOURCES_DEVICE_TESTING:
		result = "Testing"
	case HOST_RESOURCES_DEVICE_DOWN:
		result = "Down"
	default:
		result = fmt.Sprintf("Unknown(%d)", status)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetAlerts() ([]PRINTER_ALERT, error) {
	result := make([]PRINTER_ALERT, 0)
	err := error(nil)
	logger.Logger.Debug(fmt.Sprintf("Fetching alert information - Target: %s", client.Target))

	var walkResultList []SNMP_RESULT
	if walkResultList, err = client.WalkAll(OID_PRINTER_ALERT_DESCRIPTION); err == nil {
		type alertPair struct {
			description string
			index       int
			severity    int
		}
		alertMap := make(map[int]*alertPair)
		for _, walkResult := range walkResultList {
			var index int
			index = parseLastIndex(walkResult.Name)
			if index >= 0 {
				if _, ok := alertMap[index]; !ok {
					alertMap[index] = &alertPair{index: index}
				}
				if value, ok := walkResult.Value.([]byte); ok {
					alertMap[index].description = string(value)
				}
			}
		}
		var severityList []SNMP_RESULT
		if severityList, err = client.WalkAll(OID_PRINTER_ALERT_SEVERITY_LEVEL); err == nil {
			for _, walkResult := range severityList {
				var index int
				index = parseLastIndex(walkResult.Name)
				if index >= 0 {
					if _, ok := alertMap[index]; !ok {
						alertMap[index] = &alertPair{index: index}
					}
					if value, ok := walkResult.Value.(int); ok {
						alertMap[index].severity = value
					}
				}
			}
		}
		keys := make([]int, 0, len(alertMap))
		for key := range alertMap {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		for _, key := range keys {
			value := alertMap[key]
			result = append(result, PRINTER_ALERT{
				Description:  value.description,
				Severity:     value.severity,
				SeverityName: AlertSeverityName(value.severity),
			})
		}
		logger.Logger.Debug(fmt.Sprintf("Successfully fetched alert information - Target: %s, alert count: %d", client.Target, len(result)))
	} else {
		logger.Logger.Debug(fmt.Sprintf("Failed to fetch alert information - Target: %s, Error: %v", client.Target, err))
	}

	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetDeviceStatus() (int, string, error) {
	status := 0
	statusName := ""
	err := error(nil)
	var response *SNMP_RESULT
	if response, err = client.Get(OID_HOST_RESOURCES_DEVICE_STATUS); err == nil {
		if response != nil {
			if value, ok := response.Value.(int); ok {
				status = value
				statusName = DeviceStatusName(status)
			}
		}
	}
	return status, statusName, err
}

//goland:noinspection GoUnusedExportedFunction
func GetFriendlyErrorMessages(information *PRINTER_INFORMATION) []string {
	result := make([]string, 0)
	if information != nil {
		logger.Logger.Debug(fmt.Sprintf("Parsing printer status - DeviceStatus: %d (%s), alert count: %d, supply count: %d",
			information.DeviceStatus, information.DeviceStatusName, len(information.Alerts), len(information.Supplies)))
		if information.DeviceStatus == HOST_RESOURCES_DEVICE_DOWN {
			result = append(result, ERROR_MESSAGE_PRINTER_HARDWARE_FAILURE)
		} else if information.DeviceStatus == HOST_RESOURCES_DEVICE_WARNING {
			result = append(result, ERROR_MESSAGE_PRINTER_WARNING_STATE)
		}
		for _, alert := range information.Alerts {
			logger.Logger.Debug(fmt.Sprintf("Processing alert - Severity: %d (%s), Description: %s", alert.Severity, alert.SeverityName, alert.Description))
			description := strings.ToLower(alert.Description)
			if strings.Contains(description, ALERT_KEYWORD_PAPER) || strings.Contains(description, ALERT_KEYWORD_MEDIA) {
				if strings.Contains(description, ALERT_KEYWORD_JAM) || strings.Contains(description, ALERT_KEYWORD_FUSER) {
					result = append(result, ERROR_MESSAGE_PAPER_JAM)
				} else if strings.Contains(description, ALERT_KEYWORD_EMPTY) || strings.Contains(description, ALERT_KEYWORD_OUT) || strings.Contains(description, ALERT_KEYWORD_LOW) || strings.Contains(description, ALERT_KEYWORD_FEED) {
					result = append(result, ERROR_MESSAGE_OUT_OF_PAPER)
				} else if strings.Contains(description, ALERT_KEYWORD_MISMATCH) || strings.Contains(description, ALERT_KEYWORD_SIZE) || strings.Contains(description, ALERT_KEYWORD_TYPE) {
					result = append(result, ERROR_MESSAGE_PAPER_MISMATCH)
				} else {
					result = append(result, ERROR_MESSAGE_PAPER_ABNORMAL)
				}
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_TONER) {
				if strings.Contains(description, ALERT_KEYWORD_EMPTY) || strings.Contains(description, ALERT_KEYWORD_OUT) {
					result = append(result, ERROR_MESSAGE_OUT_OF_TONER)
				} else if strings.Contains(description, ALERT_KEYWORD_LOW) {
					result = append(result, ERROR_MESSAGE_TONER_LOW)
				} else {
					result = append(result, ERROR_MESSAGE_TONER_ABNORMAL)
				}
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_INK) {
				if strings.Contains(description, ALERT_KEYWORD_EMPTY) || strings.Contains(description, ALERT_KEYWORD_OUT) {
					result = append(result, ERROR_MESSAGE_OUT_OF_INK)
				} else if strings.Contains(description, ALERT_KEYWORD_LOW) {
					result = append(result, ERROR_MESSAGE_INK_LOW)
				} else {
					result = append(result, ERROR_MESSAGE_INK_ABNORMAL)
				}
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_DRUM) || strings.Contains(description, ALERT_KEYWORD_IMAGING) || strings.Contains(description, ALERT_KEYWORD_PHOTOCONDUCTOR) {
				result = append(result, ERROR_MESSAGE_DRUM_ABNORMAL)
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_COVER) || strings.Contains(description, ALERT_KEYWORD_DOOR) {
				result = append(result, ERROR_MESSAGE_COVER_OPEN)
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_TRAY) || strings.Contains(description, ALERT_KEYWORD_CASSETTE) || strings.Contains(description, ALERT_KEYWORD_DRAWER) {
				result = append(result, ERROR_MESSAGE_TRAY_ABNORMAL)
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_FUSER) {
				result = append(result, ERROR_MESSAGE_FUSER_ABNORMAL)
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_WASTE) {
				result = append(result, ERROR_MESSAGE_WASTE_TONER_BOX_ABNORMAL)
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_OUTPUT) || strings.Contains(description, ALERT_KEYWORD_STACKER) || strings.Contains(description, ALERT_KEYWORD_BIN) {
				result = append(result, ERROR_MESSAGE_OUTPUT_ABNORMAL)
				continue
			}
			if strings.Contains(description, ALERT_KEYWORD_ERROR) || strings.Contains(description, ALERT_KEYWORD_FAILURE) || strings.Contains(description, ALERT_KEYWORD_FAULT) {
				result = append(result, ERROR_MESSAGE_PRINTER_FAULT_PREFIX+alert.Description)
				continue
			}
			if alert.Severity >= ALERT_SEVERITY_CRITICAL {
				alertDescription := alert.Description
				if alertDescription == "" {
					alertDescription = DESCRIPTION_UNKNOW
				}
				result = append(result, ERROR_MESSAGE_PRINTER_ALERT_PREFIX+alertDescription)
			}
		}
		for _, supply := range information.Supplies {
			logger.Logger.Debug(fmt.Sprintf("Processing supply - Type: %d (%s), Description: %s, Level: %d, MaximumCapacity: %d", supply.Type, supply.TypeName, supply.Description, supply.Level, supply.MaximumCapacity))
			if supply.MaximumCapacity > 0 {
				percentage := (supply.Level * 100) / supply.MaximumCapacity
				if percentage <= SUPPLY_LEVEL_PERCENTAGE_CRITICAL && supply.Level == SUPPLY_LEVEL_EMPTY {
					switch supply.Type {
					case SUPPLY_TYPE_TONER:
						result = append(result, ERROR_MESSAGE_OUT_OF_TONER)
					case SUPPLY_TYPE_INK:
						result = append(result, ERROR_MESSAGE_OUT_OF_INK)
					default:
						if supply.Description != "" {
							result = append(result, ERROR_MESSAGE_SUPPLY_DEPLETED_PREFIX+supply.Description)
						}
					}
				}
			}
		}
		seen := make(map[string]bool)
		unique := make([]string, 0, len(result))
		for _, message := range result {
			if !seen[message] {
				seen[message] = true
				unique = append(unique, message)
			}
		}
		result = unique
		logger.Logger.Debug(fmt.Sprintf("Parsing completed - error messages: %v", result))
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetModel() (string, error) {
	result := ""
	err := error(nil)
	var response *SNMP_RESULT
	if response, err = client.Get(OID_SYSTEM_DESCRIPTION); err == nil {
		if response != nil {
			if value, ok := response.Value.([]byte); ok {
				result = string(value)
			} else if value, ok := response.Value.(string); ok {
				result = value
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetName() (string, error) {
	result := ""
	err := error(nil)
	var response *SNMP_RESULT
	if response, err = client.Get(OID_PRINTER_GENERAL_PRINTER_NAME); err == nil {
		if response != nil {
			if value, ok := response.Value.([]byte); ok {
				result = string(value)
			} else if value, ok := response.Value.(string); ok {
				result = value
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetPrinterInformation() (*PRINTER_INFORMATION, error) {
	result := &PRINTER_INFORMATION{}
	err := error(nil)
	logger.Logger.Debug(fmt.Sprintf("Fetching printer information - Target: %s", client.Target))
	if result.Model, err = client.GetModel(); err == nil {
		logger.Logger.Debug(fmt.Sprintf("Successfully fetched printer model - Target: %s, Model: %s", client.Target, result.Model))
		if result.Name, err = client.GetName(); err == nil {
			logger.Logger.Debug(fmt.Sprintf("Successfully fetched printer name - Target: %s, Name: %s", client.Target, result.Name))
			if result.DeviceStatus, result.DeviceStatusName, err = client.GetDeviceStatus(); err == nil {
				logger.Logger.Debug(fmt.Sprintf("Successfully fetched device status - Target: %s, Status: %s", client.Target, result.DeviceStatusName))
				if result.Supplies, err = client.GetSupplies(); err == nil {
					logger.Logger.Debug(fmt.Sprintf("Successfully fetched supply information - Target: %s, supply count: %d", client.Target, len(result.Supplies)))
					if result.Alerts, err = client.GetAlerts(); err == nil {
						logger.Logger.Debug(fmt.Sprintf("Successfully fetched alert information - Target: %s, alert count: %d", client.Target, len(result.Alerts)))
						result.TotalPages, err = client.GetTotalPages()
						if err == nil {
							logger.Logger.Debug(fmt.Sprintf("Successfully fetched total page count - Target: %s, total page count: %d", client.Target, result.TotalPages))
						}
					} else {
						logger.Logger.Debug(fmt.Sprintf("Failed to fetch alert information - Target: %s, Error: %v", client.Target, err))
					}
				} else {
					logger.Logger.Debug(fmt.Sprintf("Failed to fetch supply information - Target: %s, Error: %v", client.Target, err))
				}
			} else {
				logger.Logger.Debug(fmt.Sprintf("Failed to fetch device status - Target: %s, Error: %v", client.Target, err))
			}
		} else {
			logger.Logger.Debug(fmt.Sprintf("Failed to fetch printer name - Target: %s, Error: %v", client.Target, err))
		}
	} else {
		logger.Logger.Debug(fmt.Sprintf("Failed to fetch printer model - Target: %s, Error: %v", client.Target, err))
	}

	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetSupplies() ([]PRINTER_SUPPLY, error) {
	result := make([]PRINTER_SUPPLY, 0)
	err := error(nil)
	logger.Logger.Debug(fmt.Sprintf("Fetching supply information - Target: %s", client.Target))

	var walkResultList []SNMP_RESULT
	if walkResultList, err = client.WalkAll(OID_PRINTER_MARKER_SUPPLIES_LEVEL); err == nil {
		supplyMap := make(map[int]*SUPPLY_DATA)
		for _, walkResult := range walkResultList {
			var index int
			index = parseLastIndex(walkResult.Name)
			if index >= 0 {
				if _, ok := supplyMap[index]; !ok {
					supplyMap[index] = &SUPPLY_DATA{index: index}
				}
				if value, ok := walkResult.Value.(int); ok {
					supplyMap[index].level = value
				}
			}
		}
		addWalkData(supplyMap, client, OID_PRINTER_MARKER_SUPPLIES_DESCRIPTION, func(data *SUPPLY_DATA, value interface{}) {
			if bytes, ok := value.([]byte); ok {
				data.description = string(bytes)
			}
		})
		addWalkData(supplyMap, client, OID_PRINTER_MARKER_SUPPLIES_MAXIMUM, func(data *SUPPLY_DATA, value interface{}) {
			if intValue, ok := value.(int); ok {
				data.maximumcapacity = intValue
			}
		})
		addWalkData(supplyMap, client, OID_PRINTER_MARKER_SUPPLIES_TYPE, func(data *SUPPLY_DATA, value interface{}) {
			if intValue, ok := value.(int); ok {
				data.supplytype = intValue
			}
		})
		keys := make([]int, 0, len(supplyMap))
		for key := range supplyMap {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		for _, key := range keys {
			value := supplyMap[key]
			result = append(result, PRINTER_SUPPLY{
				Description:     value.description,
				Index:           value.index,
				Level:           value.level,
				MaximumCapacity: value.maximumcapacity,
				Type:            value.supplytype,
				TypeName:        SupplyTypeName(value.supplytype),
			})
		}
		logger.Logger.Debug(fmt.Sprintf("Successfully fetched supply information - Target: %s, supply count: %d", client.Target, len(result)))
	} else {
		logger.Logger.Debug(fmt.Sprintf("Failed to fetch supply information - Target: %s, Error: %v", client.Target, err))
	}

	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetTotalPages() (int64, error) {
	result := int64(0)
	err := error(nil)
	var walkResultList []SNMP_RESULT
	if walkResultList, err = client.WalkAll(OID_PRINTER_MARKER_LIFE_COUNT); err == nil {
		for _, walkResult := range walkResultList {
			switch value := walkResult.Value.(type) {
			case int:
				result += int64(value)
			case int64:
				result += value
			case uint:
				result += int64(value)
			case uint64:
				result += int64(value)
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func IsPrinterOnline(client *SNMP_CLIENT) (bool, error) {
	result := false
	err := error(nil)
	var status int
	if status, _, err = client.GetDeviceStatus(); err == nil {
		if status == HOST_RESOURCES_DEVICE_RUNNING {
			result = true
		}
	}
	return result, err
}

func parseLastIndex(oid string) int {
	result := PARSE_LAST_INDEX_INVALID
	index := strings.LastIndex(oid, ".")
	if index >= 0 {
		if value, err := strconv.Atoi(oid[index+1:]); err == nil {
			result = value
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func SupplyTypeName(supplyType int) string {
	result := ""
	switch supplyType {
	case SUPPLY_TYPE_TONER:
		result = "Toner"
	case SUPPLY_TYPE_WASTE_TONER:
		result = "WasteToner"
	case SUPPLY_TYPE_INK:
		result = "Ink"
	case SUPPLY_TYPE_DEVELOPER:
		result = "Developer"
	default:
		result = fmt.Sprintf("Other(%d)", supplyType)
	}
	return result
}
