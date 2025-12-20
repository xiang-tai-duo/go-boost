// Package license2
// File:        license2.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/license2/license2.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: License generation and validation helpers.
// --------------------------------------------------------------------------------
package license2

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_LICENSE2 = "license2"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_LICENSE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_LICENSE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_LICENSE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_LICENSE2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func GenerateLicense1(fingerprint string, extension string, publicKeyPEM string, useExtension bool) string {
	__debug(fmt.Sprintf("GenerateLicense1: start, fingerprint length=%d useExtension=%v", len(fingerprint), useExtension))
	result := ""
	data := []byte(fingerprint)
	if useExtension {
		payload := struct {
			Fingerprint string `json:"fingerprint"`
			Extension   string `json:"extension"`
		}{
			Fingerprint: fingerprint,
			Extension:   extension,
		}
		if jsonData, err := json.Marshal(payload); err == nil {
			data = jsonData
			__debug(fmt.Sprintf("GenerateLicense1: wrapped in JSON payload length=%d", len(data)))
		} else {
			__debug(fmt.Sprintf("GenerateLicense1: marshal payload failed: %v", err))
		}
	}
	if block, _ := pem.Decode([]byte(publicKeyPEM)); block != nil {
		__debug(fmt.Sprintf("GenerateLicense1: PEM block decoded type=%s", block.Type))
		if pubKey, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			if rsaPubKey, ok := pubKey.(*rsa.PublicKey); ok {
				__debug(fmt.Sprintf("GenerateLicense1: RSA public key size=%d bits", rsaPubKey.N.BitLen()))
				if encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, data, nil); err == nil {
					result = base64.StdEncoding.EncodeToString(encrypted)
					__debug(fmt.Sprintf("GenerateLicense1: success, result length=%d", len(result)))
				} else {
					__debug(fmt.Sprintf("GenerateLicense1: RSA encrypt failed: %v", err))
				}
			} else {
				__debug("GenerateLicense1: public key is not RSA")
			}
		} else {
			__debug(fmt.Sprintf("GenerateLicense1: ParsePKIXPublicKey failed: %v", err))
		}
	} else {
		__debug("GenerateLicense1: PEM decode failed")
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GenerateLicense2(fingerprint string, salt string) string {
	__debug(fmt.Sprintf("GenerateLicense2: start, fingerprint length=%d salt length=%d", len(fingerprint), len(salt)))
	data := []byte(fingerprint + salt)
	hash := sha256.Sum256(data)
	fullHash := hex.EncodeToString(hash[:])
	__debug(fmt.Sprintf("GenerateLicense2: SHA256 hash=%s", fullHash))
	parts := make([]string, 5)
	for i := 0; i < 5; i++ {
		parts[i] = strings.ToUpper(fullHash[i*5 : (i+1)*5])
	}
	result := strings.Join(parts, "-")
	__debug(fmt.Sprintf("GenerateLicense2: result=%s", result))
	return result
}

//goland:noinspection GoUnusedExportedFunction
func GetSHA3(input string) string {
	hash := sha3.Sum256([]byte(input))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

//goland:noinspection DuplicatedCode
func IsValidLicense(license string, privateKeyPEM string, fingerprint string) bool {
	__debug(fmt.Sprintf("IsValidLicense: start, license length=%d fingerprint length=%d", len(license), len(fingerprint)))
	result := false
	if license != "" {
		if encrypted, err := base64.StdEncoding.DecodeString(license); err == nil {
			__debug(fmt.Sprintf("IsValidLicense: base64 decoded length=%d", len(encrypted)))
			if block, _ := pem.Decode([]byte(privateKeyPEM)); block != nil {
				__debug(fmt.Sprintf("IsValidLicense: PEM block decoded type=%s", block.Type))
				if privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
					__debug("IsValidLicense: private key parsed")
					if decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encrypted, nil); err == nil {
						__debug(fmt.Sprintf("IsValidLicense: decrypted length=%d", len(decrypted)))
						result = strings.EqualFold(string(decrypted), fingerprint)
						__debug(fmt.Sprintf("IsValidLicense: RSA validation result=%v", result))
					} else {
						__debug(fmt.Sprintf("IsValidLicense: RSA decrypt failed: %v", err))
					}
				} else {
					__debug(fmt.Sprintf("IsValidLicense: ParsePKCS1PrivateKey failed: %v", err))
				}
			} else {
				__debug("IsValidLicense: PEM decode failed")
			}
		} else {
			__debug(fmt.Sprintf("IsValidLicense: base64 decode failed: %v (license is not base64, likely GenerateLicense2 format)", err))
		}
	} else {
		__debug("IsValidLicense: license is empty")
	}
	__debug(fmt.Sprintf("IsValidLicense: result=%v", result))
	return result
}

//goland:noinspection GoUnusedExportedFunction
func IsValidLicenseFile(privateKeyPEM string, fingerprint string, licenseFileName string, serialSalt string) bool {
	__debug(fmt.Sprintf("IsValidLicenseFile: start, licenseFileName=%s fingerprint length=%d", licenseFileName, len(fingerprint)))
	result := false
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		licenseFilePath := filepath.Join(exeDir, licenseFileName)
		__debug(fmt.Sprintf("IsValidLicenseFile: looking for license file at %s", licenseFilePath))
		if content, err := os.ReadFile(licenseFilePath); err == nil {
			result = validateLicenseContent(string(content), privateKeyPEM, fingerprint, serialSalt)
		} else {
			__debug(fmt.Sprintf("IsValidLicenseFile: failed to read license file from exe dir: %v", err))
			if cwd, err2 := os.Getwd(); err2 == nil {
				cwdFilePath := filepath.Join(cwd, licenseFileName)
				__debug(fmt.Sprintf("IsValidLicenseFile: trying working directory at %s", cwdFilePath))
				if content, err2 := os.ReadFile(cwdFilePath); err2 == nil {
					result = validateLicenseContent(string(content), privateKeyPEM, fingerprint, serialSalt)
				} else {
					__debug(fmt.Sprintf("IsValidLicenseFile: failed to read license file from working dir: %v", err2))
				}
			} else {
				__debug(fmt.Sprintf("IsValidLicenseFile: os.Getwd failed: %v", err2))
			}
		}
	} else {
		__debug(fmt.Sprintf("IsValidLicenseFile: os.Executable failed: %v", err))
	}
	__debug(fmt.Sprintf("IsValidLicenseFile: final result=%v", result))
	return result
}

func validateLicenseContent(content string, privateKeyPEM string, fingerprint string, serialSalt string) bool {
	result := false
	license := strings.TrimSpace(content)
	__debug(fmt.Sprintf("validateLicenseContent: license length=%d content=%s", len(license), license))
	if result = IsValidLicense(license, privateKeyPEM, fingerprint); !result {
		expected := GenerateLicense2(fingerprint, serialSalt)
		__debug(fmt.Sprintf("validateLicenseContent: RSA validation failed, trying GenerateLicense2 comparison: expected=%s actual=%s", expected, license))
		result = strings.EqualFold(license, expected)
		__debug(fmt.Sprintf("validateLicenseContent: GenerateLicense2 comparison result=%v", result))
	} else {
		__debug("validateLicenseContent: RSA validation succeeded")
	}
	return result
}
