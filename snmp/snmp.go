// Package snmp
// File:        snmp.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SNMP is a wrapper for SNMP v1/v2c/v3 operations, providing methods for GET/SET/WALK/BULK queries.
// --------------------------------------------------------------------------------
package snmp

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	SNMP_CLIENT struct {
		gosnmp.GoSNMP
	}
	SNMP_OPTION        func(*SNMP_CLIENT)
	SNMP_RESULT        = gosnmp.SnmpPDU
	SNMP_VERSION       = gosnmp.SnmpVersion
	SNMP_WALK_CALLBACK = func(pdu gosnmp.SnmpPDU) error
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,GoUnusedConst
const (
	AUTHENTICATION_MESSAGE_DIGEST_ALGORITHM       = gosnmp.MD5
	AUTHENTICATION_NO_AUTH                        = gosnmp.NoAuth
	AUTHENTICATION_SECURE_HASH_ALGORITHM          = gosnmp.SHA
	AUTHENTICATION_SECURE_HASH_ALGORITHM_256      = gosnmp.SHA256
	AUTHENTICATION_SECURE_HASH_ALGORITHM_384      = gosnmp.SHA384
	AUTHENTICATION_SECURE_HASH_ALGORITHM_512      = gosnmp.SHA512
	DEFAULT_COMMUNITY_STRING                      = "public"
	DEFAULT_MAXIMUM_REPETITION                    = 32
	DEFAULT_PORT_NUMBER                           = 161
	DEFAULT_RETRY_COUNT                           = 1
	DEFAULT_TIMEOUT_SECONDS                       = 5
	DEFAULT_TIMEOUT_DURATION                      = DEFAULT_TIMEOUT_SECONDS * time.Second
	PRIVACY_ADVANCED_ENCRYPTION_STANDARD          = gosnmp.AES
	PRIVACY_ADVANCED_ENCRYPTION_STANDARD_192      = gosnmp.AES192
	PRIVACY_ADVANCED_ENCRYPTION_STANDARD_192_C    = gosnmp.AES192C
	PRIVACY_ADVANCED_ENCRYPTION_STANDARD_256      = gosnmp.AES256
	PRIVACY_ADVANCED_ENCRYPTION_STANDARD_256_C    = gosnmp.AES256C
	PRIVACY_DATA_ENCRYPTION_STANDARD              = gosnmp.DES
	PRIVACY_NO_PRIVACY                            = gosnmp.NoPriv
	SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_1  = gosnmp.Version1
	SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_2C = gosnmp.Version2c
	SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_3  = gosnmp.Version3
	MODULE_NAME_SNMP                              = "snmp"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SNMP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SNMP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SNMP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func New(target string, options ...SNMP_OPTION) *SNMP_CLIENT {
	result := &SNMP_CLIENT{
		GoSNMP: gosnmp.GoSNMP{
			Community: DEFAULT_COMMUNITY_STRING,
			Port:      DEFAULT_PORT_NUMBER,
			Retries:   DEFAULT_RETRY_COUNT,
			Target:    target,
			Timeout:   DEFAULT_TIMEOUT_DURATION,
			Version:   SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_2C,
		},
	}
	for _, option := range options {
		option(result)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) BulkWalk(rootOID string, callback SNMP_WALK_CALLBACK) error {
	err := error(nil)
	if err = client.Connect(); err == nil {
		defer client.Close()
		err = client.GoSNMP.BulkWalk(rootOID, func(dataUnit gosnmp.SnmpPDU) error {
			return callback(dataUnit)
		})
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) BulkWalkAll(rootOID string) ([]SNMP_RESULT, error) {
	result := make([]SNMP_RESULT, 0)
	err := error(nil)
	err = client.BulkWalk(rootOID, func(pdu gosnmp.SnmpPDU) error {
		result = append(result, pdu)
		return nil
	})
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) Close() error {
	err := error(nil)
	if client.Conn != nil {
		err = client.Conn.Close()
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) Connect() error {
	err := error(nil)
	err = client.GoSNMP.Connect()
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) Get(oid string) (*SNMP_RESULT, error) {
	result := (*SNMP_RESULT)(nil)
	err := error(nil)
	__debug(fmt.Sprintf("SNMP GET - Target: %s, OID: %s", client.Target, oid))
	if err = client.Connect(); err == nil {
		defer client.Close()
		var packet *gosnmp.SnmpPacket
		if packet, err = client.GoSNMP.Get([]string{oid}); err == nil {
			if len(packet.Variables) > 0 {
				result = &packet.Variables[0]
				__debug(fmt.Sprintf("SNMP GET succeeded - Target: %s, OID: %s, Value: %v", client.Target, oid, result.Value))
			}
		} else {
			__debug(fmt.Sprintf("SNMP GET failed - Target: %s, OID: %s, Error: %v", client.Target, oid, err))
		}
	} else {
		__debug(fmt.Sprintf("SNMP connection failed - Target: %s, Error: %v", client.Target, err))
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) GetBulk(oids []string, nonRepeaters uint8, maxRepetitions uint32) ([]SNMP_RESULT, error) {
	result := make([]SNMP_RESULT, 0)
	err := error(nil)
	if err = client.Connect(); err == nil {
		defer client.Close()
		var packet *gosnmp.SnmpPacket
		if packet, err = client.GoSNMP.GetBulk(oids, nonRepeaters, maxRepetitions); err == nil {
			result = packet.Variables
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) GetMulti(oids []string) ([]SNMP_RESULT, error) {
	result := make([]SNMP_RESULT, 0)
	err := error(nil)
	if err = client.Connect(); err == nil {
		defer client.Close()
		var packet *gosnmp.SnmpPacket
		if packet, err = client.GoSNMP.Get(oids); err == nil {
			result = packet.Variables
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) GetNext(oid string) (*SNMP_RESULT, error) {
	result := (*SNMP_RESULT)(nil)
	err := error(nil)
	if err = client.Connect(); err == nil {
		defer client.Close()
		var packet *gosnmp.SnmpPacket
		if packet, err = client.GoSNMP.GetNext([]string{oid}); err == nil {
			if len(packet.Variables) > 0 {
				result = &packet.Variables[0]
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) GetTimeoutSeconds() int {
	result := 0
	if client.Timeout > 0 {
		result = int(client.Timeout / time.Second)
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func NewV1(target string, community string, options ...SNMP_OPTION) *SNMP_CLIENT {
	result := New(target, append(options, WithCommunity(community))...)
	result.Version = SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_1
	return result
}

//goland:noinspection GoUnusedExportedFunction
func NewV2C(target string, community string, options ...SNMP_OPTION) *SNMP_CLIENT {
	result := New(target, append(options, WithCommunity(community))...)
	result.Version = SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_2C
	return result
}

//goland:noinspection GoUnusedExportedFunction
func NewV3(target string, securityUserName string, authProtocol gosnmp.SnmpV3AuthProtocol, authPassword string, privateProtocol gosnmp.SnmpV3PrivProtocol, privatePassword string,
	securityLevel gosnmp.SnmpV3MsgFlags, options ...SNMP_OPTION) *SNMP_CLIENT {

	result := New(target, options...)
	result.Version = SIMPLE_NETWORK_MANAGEMENT_PROTOCOL_VERSION_3
	result.SecurityModel = gosnmp.UserSecurityModel
	result.MsgFlags = securityLevel
	result.SecurityParameters = &gosnmp.UsmSecurityParameters{
		UserName:                 securityUserName,
		AuthenticationProtocol:   authProtocol,
		AuthenticationPassphrase: authPassword,
		PrivacyProtocol:          privateProtocol,
		PrivacyPassphrase:        privatePassword,
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) Set(oid string, value interface{}, asn1Type gosnmp.Asn1BER) (*SNMP_RESULT, error) {
	result := (*SNMP_RESULT)(nil)
	err := error(nil)
	if err = client.Connect(); err == nil {
		defer client.Close()
		pdu := gosnmp.SnmpPDU{
			Name:  oid,
			Type:  asn1Type,
			Value: value,
		}
		var packet *gosnmp.SnmpPacket
		if packet, err = client.GoSNMP.Set([]gosnmp.SnmpPDU{pdu}); err == nil {
			if len(packet.Variables) > 0 {
				result = &packet.Variables[0]
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func (client *SNMP_CLIENT) SetTimeoutSeconds(seconds int) {
	if seconds > 0 {
		client.Timeout = time.Duration(seconds) * time.Second
	}
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) Walk(rootOID string, callback SNMP_WALK_CALLBACK) error {
	err := error(nil)
	if err = client.Connect(); err == nil {
		defer client.Close()
		err = client.GoSNMP.Walk(rootOID, func(dataUnit gosnmp.SnmpPDU) error {
			return callback(dataUnit)
		})
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult
func (client *SNMP_CLIENT) WalkAll(rootOID string) ([]SNMP_RESULT, error) {
	result := make([]SNMP_RESULT, 0)
	err := error(nil)
	__debug(fmt.Sprintf("SNMP WalkAll - Target: %s, RootOID: %s", client.Target, rootOID))
	if err = client.Connect(); err == nil {
		defer client.Close()
		err = client.GoSNMP.Walk(rootOID, func(dataUnit gosnmp.SnmpPDU) error {
			result = append(result, dataUnit)
			return nil
		})
		if err == nil {
			__debug(fmt.Sprintf("SNMP WalkAll succeeded - Target: %s, RootOID: %s, Result count: %d", client.Target, rootOID, len(result)))
		} else {
			__debug(fmt.Sprintf("SNMP WalkAll failed - Target: %s, RootOID: %s, Error: %v", client.Target, rootOID, err))
		}
	} else {
		__debug(fmt.Sprintf("SNMP WalkAll connection failed - Target: %s, Error: %v", client.Target, err))
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func WithCommunity(community string) SNMP_OPTION {
	return func(client *SNMP_CLIENT) {
		client.Community = community
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithContextEngineID(identifier string) SNMP_OPTION {
	return func(client *SNMP_CLIENT) {
		client.ContextEngineID = identifier
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithContextName(name string) SNMP_OPTION {
	return func(client *SNMP_CLIENT) {
		client.ContextName = name
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithTimeout(duration time.Duration) SNMP_OPTION {
	return func(client *SNMP_CLIENT) {
		client.Timeout = duration
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithTimeoutSeconds(seconds int) SNMP_OPTION {
	return func(client *SNMP_CLIENT) {
		client.SetTimeoutSeconds(seconds)
	}
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SNMP, logger.SKIP_STACK_FRAMES_BASE)
}
