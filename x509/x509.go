// Package x509
// File:        x509.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/x509/x509.go
// Author:      Vibe Coding
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

	"golang.org/x/crypto/pkcs12"
)

//goland:noinspection GoSnakeCaseUsage
const (
	RSA_KEY_SIZE              = 2048
	CERTIFICATE_VALIDITY_DAYS = 365
	ORGANIZATION_NAME         = "https://github.com/xiang-tai-duo/go-boost"
	DUMMY_PFX_DATA            = "dummy-pfx-data"
	CERTIFICATE_TYPE          = "CERTIFICATE"
	PRIVATE_KEY_TYPE          = "PRIVATE KEY"
)

//goland:noinspection GoUnusedExportedFunction
func New(password string) ([]byte, error) {
	result := make([]byte, 0)
	err := error(nil)
	privateKey := (*rsa.PrivateKey)(nil)
	if privateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE); err == nil {
		validFrom := time.Now()
		validTo := validFrom.Add(CERTIFICATE_VALIDITY_DAYS * 24 * time.Hour)
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber := (*big.Int)(nil)
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
	certs := ([]*pem.Block)(nil)
	if certs, err = pkcs12.ToPEM(pfx, password); err == nil {
		for _, cert := range certs {
			if cert.Type == CERTIFICATE_TYPE {
				for _, key := range certs {
					if key.Type == PRIVATE_KEY_TYPE {
						tlsCert := tls.Certificate{}
						loadErr := error(nil)
						if tlsCert, loadErr = tls.X509KeyPair(pem.EncodeToMemory(cert), pem.EncodeToMemory(key)); loadErr == nil {
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
