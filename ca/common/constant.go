// Package common
// File:        constant.go
// Author:      TRAE.AI
// Created:     2026/06/09 00:00:00
// Description: Common certificate constants shared by CA client and server.
// --------------------------------------------------------------------------------
package common

//goland:noinspection GoSnakeCaseUsage
const (
	CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME = "ca.crt"
	CERTIFICATE_AUTHORITY_PRIVATE_KEY_FILE_NAME = "ca.key"
	CERTIFICATE_AUTHORITY_SERIAL_FILE_NAME      = "ca.serial"
	CERTIFICATE_AUTHORITY_VALIDITY_YEARS        = 100
	CERTIFICATE_TYPE                            = "CERTIFICATE"
	CERTIFICATES_DIRECTORY                      = "certificates"
	CA_DIRECTORY                                = "ca"
	COUNTRY_NAME                                = "CN"
	DIRECTORY_PERMISSION                        = 0o755
	FILE_PERMISSION_PRIVATE                     = 0o600
	FILE_PERMISSION_PUBLIC                      = 0o644
	ORGANIZATION_NAME                           = "go-boost CA"
	RSA_KEY_SIZE                                = 4096
)
