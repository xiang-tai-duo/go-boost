// Package rsa2
// File:        rsa.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/rsa/rsa.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: RSA is a wrapper for RSA encryption and decryption with caller-supplied key pairs.
// --------------------------------------------------------------------------------
package rsa2

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage,GoUnusedConst
const (
	RSA_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwrS4aLF/V5gXfS97GJJ7
ESomOYtTK+Qa/OD1wBFZQcigh4BfdBBt0RQtWCi+Hh/TgurM7ZFWeU0sne+WivNh
QjR+I/cA5hQufraB1ydAUFqOxfXclyPtfqqbLt1Z8G4wz0B+3rFk5K/MHPRcEcQF
qCJ5w4WWTR3eLlJsKppY7XOOUupU9Wl2ATb/ChESKb+J4tLeUlY4b7yDfTMlK0Ge
OI1I7HZHnAnn3AsV0ilUgZMdwZOMWC5gzEZCuq/l99plxW74BI/W8snPQWgDlZ/x
Lbszf2E3+VSIdSycdMg3O8nWgFffWbEWrTv2y82KQ1XLSnMaQwIL4KymOJCDtghx
qZ1tXyqbjcONkMMgwBfWeoPJdhPk/e5x4+6TkIRja8D1vkXZPCs5bG8Nu1lB+pGQ
sb1uLl/uZgmsP9WXkJlK2tsYI2Yn8FnmReyAEDnkfXONbrdz7mgh1u+eFnmfptvz
GIIEnUf2DA636rbu2XST1/XGLY9R1d//HMwY7+V2R+d7XOoKr78DINbsIimoOCBI
Ibz1tSuF9admDKsdJm/6+fu8O+LRAeBF1KlFasr1kJ+fOlm8lTZ0MNsdIuv4ywQ+
EIrOZSbIisjLQlsryaVviRTuLRWe/LxxfQWySjSEaSP7Apmh7bh+ael7ChPL7oDK
utEoBR1bhyLBx9WvVZRwXZsCAwEAAQ==
-----END PUBLIC KEY-----`

	RSA_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

	MODULE_NAME_RSA          = "rsa"
	PEM_TYPE_PUBLIC_KEY      = "PUBLIC KEY"
	PEM_TYPE_RSA_PRIVATE_KEY = "RSA PRIVATE KEY"
	RSA_KEY_SIZE             = 4096
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_RSA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_RSA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_RSA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_RSA, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func Decrypt(privateKeyPEM string, value string) ([]byte, error) {
	var result []byte
	err := error(nil)
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		err = fmt.Errorf("failed to decode private key PEM block")
	} else {
		var privateKey *rsa.PrivateKey
		if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			var ciphertext []byte
			if ciphertext, err = base64.StdEncoding.DecodeString(value); err == nil {
				var plaintext []byte
				if plaintext, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil); err == nil {
					result = plaintext
				}
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func Encrypt(publicKeyPEM string, value string) (string, error) {
	result := ""
	err := error(nil)
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		err = fmt.Errorf("failed to decode public key PEM block")
	} else {
		var publicKeyInterface interface{}
		if publicKeyInterface, err = x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
			if !ok {
				err = fmt.Errorf("public key is not an RSA public key")
			} else {
				var ciphertext []byte
				if ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(value), nil); err == nil {
					result = base64.StdEncoding.EncodeToString(ciphertext)
				}
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedExportedFunction
func GenerateKeyPair() (string, string, error) {
	publicKeyBase64 := ""
	privateKeyBase64 := ""
	err := error(nil)
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE); err != nil {
		err = fmt.Errorf("failed to generate RSA private key: %w", err)
	} else {
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  PEM_TYPE_RSA_PRIVATE_KEY,
			Bytes: privateKeyBytes,
		})
		var publicKeyBytes []byte
		if publicKeyBytes, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey); err != nil {
			err = fmt.Errorf("failed to marshal RSA public key: %w", err)
		} else {
			publicKeyPEM := pem.EncodeToMemory(&pem.Block{
				Type:  PEM_TYPE_PUBLIC_KEY,
				Bytes: publicKeyBytes,
			})
			privateKeyBase64 = base64.StdEncoding.EncodeToString(privateKeyPEM)
			publicKeyBase64 = base64.StdEncoding.EncodeToString(publicKeyPEM)
		}
	}
	return publicKeyBase64, privateKeyBase64, err
}
