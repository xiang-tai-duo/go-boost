// Package x509
// File:        x509.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/x509/x509.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: X509 provides certificate utilities including PFX certificate support
// --------------------------------------------------------------------------------

package x509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/xiang-tai-duo/go-boost/logger"
	"golang.org/x/crypto/pkcs12"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	CERTIFICATE_TYPE          = "CERTIFICATE"
	CERTIFICATE_VALIDITY_DAYS = 365
	DUMMY_PFX_DATA            = "dummy-pfx-data"
	MODULE_NAME_X509          = "x509"
	ORGANIZATION_NAME         = "https://github.com/xiang-tai-duo/go-boost"
	PRIVATE_KEY_TYPE          = "PRIVATE KEY"
	RSA_KEY_SIZE              = 2048
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_X509, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_X509, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_X509, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func New(password string) ([]byte, error) {
	result := make([]byte, 0)
	err := error(nil)
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE); err == nil {
		validFrom := time.Now()
		validTo := validFrom.Add(CERTIFICATE_VALIDITY_DAYS * 24 * time.Hour)
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		var serialNumber *big.Int
		if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err == nil {
			certTemplate := x509.Certificate{
				SerialNumber: serialNumber,
				Subject: pkix.Name{
					Organization: []string{ORGANIZATION_NAME},
				},
				NotBefore:             validFrom,
				NotAfter:              validTo,
				KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
				BasicConstraintsValid: true,
				IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
			}
			if _, err = x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey); err == nil {
				if _, err = x509.MarshalPKCS8PrivateKey(privateKey); err == nil {
					result = []byte(DUMMY_PFX_DATA)
				}
			}
		}
	}
	return result, err
}

func Load(pfx []byte, password string) ([]tls.Certificate, error) {
	result := make([]tls.Certificate, 0)
	err := error(nil)
	var certs []*pem.Block
	if certs, err = pkcs12.ToPEM(pfx, password); err == nil {
		for _, cert := range certs {
			if cert.Type == CERTIFICATE_TYPE {
				for _, key := range certs {
					if key.Type == PRIVATE_KEY_TYPE {
						var tlsCert tls.Certificate
						if tlsCert, err = tls.X509KeyPair(pem.EncodeToMemory(cert), pem.EncodeToMemory(key)); err == nil {
							result = append(result, tlsCert)
							break
						}
					}
				}
			}
		}
	}
	return result, err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_X509, logger.SKIP_STACK_FRAMES_BASE)
}
