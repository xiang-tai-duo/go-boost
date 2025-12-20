// Package common
// File:        issue_ip_address.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Issue IP address certificate helpers.
// --------------------------------------------------------------------------------
package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage
type CLIENT_PATH struct {
	CaCertificatePath     string
	CaDirectory           string
	CaPrivateKeyPath      string
	CaSerialPath          string
	ClientCaCertPath      string
	ClientCertificatePath string
	ClientDirectory       string
	ClientPrivateKeyPath  string
	OutputDirectory       string
	SanPath               string
}

type SAN struct {
	DNS []string `json:"dns"`
	IP  []string `json:"ip"`
}

//goland:noinspection GoSnakeCaseUsage
const (
	//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
	CLIENT_CERTIFICATE_FILE_NAME   = "client.crt"
	CLIENT_DIRECTORY               = "client"
	CLIENT_PRIVATE_KEY_FILE_NAME   = "client.key"
	CLIENT_VALIDITY_YEARS          = 100
	ERROR_MESSAGE_DECODE_CA_CERT   = "failed to decode CA certificate PEM block at %s"
	ERROR_MESSAGE_DECODE_CA_KEY    = "failed to decode CA private key PEM block at %s"
	ERROR_MESSAGE_MISSING_CA_FILE  = "Missing required CA file or directory: %s"
	ERROR_MESSAGE_NOT_RSA_KEY      = "CA private key is not an RSA private key"
	PEM_TYPE_RSA_PRIVATE_KEY       = "RSA PRIVATE KEY"
	SAN_FILE_NAME                  = "san.json"
	SERIAL_FALLBACK_VALUE          = int64(1000)
	SUBJECT_ORGANIZATIONAL_DEFAULT = "issueIPAddress"
	IP_PREFIX_192                  = "192."
	IP_PREFIX_10                   = "10."
	IP_PREFIX_172                  = "172."
	MODULE_NAME_ISSUE_IP_ADDRESS   = "ca.common.issue_ip_address"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_ISSUE_IP_ADDRESS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_ISSUE_IP_ADDRESS, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_ISSUE_IP_ADDRESS, logger.SKIP_STACK_FRAMES_BASE)
}

func BuildClientPath(executableDirectory, ipText string) CLIENT_PATH {
	var result CLIENT_PATH
	certificatesDirectory := filepath.Join(executableDirectory, CERTIFICATES_DIRECTORY)
	caDirectory := filepath.Join(certificatesDirectory, CA_DIRECTORY)
	outputDirectory := filepath.Join(certificatesDirectory, CLIENT_DIRECTORY, ipText)
	result.CaCertificatePath = filepath.Join(caDirectory, CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME)
	result.CaDirectory = caDirectory
	result.CaPrivateKeyPath = filepath.Join(caDirectory, CERTIFICATE_AUTHORITY_PRIVATE_KEY_FILE_NAME)
	result.CaSerialPath = filepath.Join(caDirectory, CERTIFICATE_AUTHORITY_SERIAL_FILE_NAME)
	result.ClientCaCertPath = filepath.Join(outputDirectory, CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME)
	result.ClientCertificatePath = filepath.Join(outputDirectory, CLIENT_CERTIFICATE_FILE_NAME)
	result.ClientDirectory = filepath.Join(certificatesDirectory, CLIENT_DIRECTORY)
	result.ClientPrivateKeyPath = filepath.Join(outputDirectory, CLIENT_PRIVATE_KEY_FILE_NAME)
	result.OutputDirectory = outputDirectory
	result.SanPath = filepath.Join(outputDirectory, SAN_FILE_NAME)
	return result
}

func CollectConsoleIPv4s(text string) []string {
	result := make([]string, 0)
	items := strings.FieldsFunc(text, func(r rune) bool {
		return r == ';'
	})
	for _, item := range items {
		ipText := strings.TrimSpace(item)
		if IsIPv4(ipText) {
			result = append(result, ipText)
		}
	}
	return result
}

func CollectIPv4s(text string) []string {
	result := make([]string, 0)
	items := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\r' || r == '\n'
	})
	for _, item := range items {
		ipText := strings.TrimSpace(item)
		if IsIPv4(ipText) {
			result = append(result, ipText)
		}
	}
	return result
}

func CopyFile(sourcePath string, destinationPath string) error {
	err := error(nil)
	data := make([]byte, 0)
	if data, err = os.ReadFile(sourcePath); err != nil {
		log.Printf("Failed to read source file: %v", err)
	} else if err = os.WriteFile(destinationPath, data, FILE_PERMISSION_PUBLIC); err != nil {
		log.Printf("Failed to write destination file: %v", err)
	}
	return err
}

func CreateClientCertificate(subjectName string, serialNumber int64, clientPrivateKey *rsa.PrivateKey, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey, ipAddresses []net.IP) ([]byte, error) {
	result := make([]byte, 0)
	err := error(nil)
	validFrom := time.Now()
	validTo := validFrom.AddDate(CLIENT_VALIDITY_YEARS, 0, 0)
	certificateTemplate := x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject: pkix.Name{
			Country:            []string{COUNTRY_NAME},
			Organization:       []string{ORGANIZATION_NAME},
			OrganizationalUnit: []string{SUBJECT_ORGANIZATIONAL_DEFAULT},
			CommonName:         subjectName,
		},
		NotBefore:             validFrom,
		NotAfter:              validTo,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		IPAddresses:           ipAddresses,
	}
	certificateDer := make([]byte, 0)
	if certificateDer, err = x509.CreateCertificate(rand.Reader, &certificateTemplate, caCertificate, &clientPrivateKey.PublicKey, caPrivateKey); err != nil {
		log.Printf("Failed to create client certificate: %v", err)
	} else {
		result = certificateDer
	}
	return result, err
}

func EnsureCaExists(path CLIENT_PATH) error {
	err := error(nil)
	checkList := []string{path.CaDirectory, path.CaCertificatePath, path.CaPrivateKeyPath, path.CaSerialPath}
	for _, item := range checkList {
		if _, err = os.Stat(item); err != nil {
			if os.IsNotExist(err) {
				err = fmt.Errorf(ERROR_MESSAGE_MISSING_CA_FILE, item)
			}
			log.Printf("CA file or directory check failed: %v", err)
			break
		}
	}
	return err
}

func IsIPv4(value string) bool {
	result := false
	if ip := net.ParseIP(value); ip != nil && ip.To4() != nil && ip.String() == value {
		result = true
	}
	return result
}

func IssueClientCertificate(ipTexts []string) (string, error) {
	result := ""
	err := error(nil)
	executableDirectory := ""
	if executableDirectory, err = ResolveExecutableDirectory(); err != nil {
		log.Printf("Failed to resolve executable directory: %v", err)
	} else {
		ipText := ipTexts[0]
		path := BuildClientPath(executableDirectory, ipText)
		if err = EnsureCaExists(path); err != nil {
			log.Printf("CA does not exist: %v", err)
		} else {
			var statErr error
			if _, statErr = os.Stat(path.OutputDirectory); statErr == nil {
				if err = os.RemoveAll(path.OutputDirectory); err != nil {
					log.Printf("Failed to remove existing output directory: %v", err)
				}
			} else if !os.IsNotExist(statErr) {
				err = statErr
			}
			if err == nil {
				if err = os.MkdirAll(path.OutputDirectory, DIRECTORY_PERMISSION); err != nil {
					log.Printf("Failed to create output directory: %v", err)
				} else {
					var caCertificate *x509.Certificate
					var caPrivateKey *rsa.PrivateKey
					if caCertificate, caPrivateKey, err = LoadCa(path.CaCertificatePath, path.CaPrivateKeyPath); err != nil {
						log.Printf("Failed to load CA: %v", err)
					} else {
						var clientPrivateKey *rsa.PrivateKey
						if clientPrivateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE); err != nil {
							log.Printf("Failed to generate client private key: %v", err)
						} else {
							serialNumber := ReadNextSerial(path.CaSerialPath)
							ipAddresses := make([]net.IP, 0)
							for _, item := range ipTexts {
								ipAddresses = append(ipAddresses, net.ParseIP(item))
							}
							var certificateDer []byte
							if certificateDer, err = CreateClientCertificate(ipText, serialNumber, clientPrivateKey, caCertificate, caPrivateKey, ipAddresses); err != nil {
								log.Printf("Failed to create client certificate: %v", err)
							} else {
								if err = WriteCertificateFile(path.ClientCertificatePath, certificateDer); err != nil {
									log.Printf("Failed to write certificate file: %v", err)
								} else {
									if err = WritePrivateKeyFile(path.ClientPrivateKeyPath, clientPrivateKey); err != nil {
										log.Printf("Failed to write private key file: %v", err)
									} else {
										if err = WriteNextSerial(path.CaSerialPath, serialNumber+1); err != nil {
											log.Printf("Failed to write next serial: %v", err)
										} else {
											if err = CopyFile(path.CaCertificatePath, path.ClientCaCertPath); err != nil {
												log.Printf("Failed to copy CA certificate to client directory: %v", err)
											} else {
												if err = WriteSanFile(path.SanPath, ipTexts); err != nil {
													log.Printf("Failed to write SAN file: %v", err)
												} else {
													result = path.OutputDirectory
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return result, err
}

func LoadCa(caCertificatePath, caPrivateKeyPath string) (*x509.Certificate, *rsa.PrivateKey, error) {
	var resultCertificate *x509.Certificate
	var resultPrivateKey *rsa.PrivateKey
	err := error(nil)
	certificatePemBytes := make([]byte, 0)
	if certificatePemBytes, err = os.ReadFile(caCertificatePath); err == nil {
		certificateBlock, _ := pem.Decode(certificatePemBytes)
		if certificateBlock == nil || certificateBlock.Type != CERTIFICATE_TYPE {
			err = fmt.Errorf(ERROR_MESSAGE_DECODE_CA_CERT, caCertificatePath)
			log.Printf("%v", err)
		} else {
			var caCertificate *x509.Certificate
			if caCertificate, err = x509.ParseCertificate(certificateBlock.Bytes); err != nil {
				log.Printf("Failed to parse CA certificate: %v", err)
			} else {
				var privateKeyPemBytes []byte
				if privateKeyPemBytes, err = os.ReadFile(caPrivateKeyPath); err != nil {
					log.Printf("Failed to read CA private key: %v", err)
				} else {
					privateKeyBlock, _ := pem.Decode(privateKeyPemBytes)
					if privateKeyBlock == nil {
						err = fmt.Errorf(ERROR_MESSAGE_DECODE_CA_KEY, caPrivateKeyPath)
						log.Printf("%v", err)
					} else {
						var caPrivateKey *rsa.PrivateKey
						if caPrivateKey, err = ParseRsaPrivateKey(privateKeyBlock); err != nil {
							log.Printf("Failed to parse RSA private key: %v", err)
						} else {
							resultCertificate = caCertificate
							resultPrivateKey = caPrivateKey
						}
					}
				}
			}
		}
	} else {
		log.Printf("Failed to read CA certificate: %v", err)
	}
	return resultCertificate, resultPrivateKey, err
}

func ParseRsaPrivateKey(privateKeyBlock *pem.Block) (*rsa.PrivateKey, error) {
	var result *rsa.PrivateKey
	err := error(nil)
	if privateKeyBlock.Type == PEM_TYPE_RSA_PRIVATE_KEY {
		if result, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes); err != nil {
			log.Printf("Failed to parse PKCS1 private key: %v", err)
		}
	} else {
		var parsed any
		if parsed, err = x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes); err != nil {
			log.Printf("Failed to parse PKCS8 private key: %v", err)
		} else {
			if rsaKey, ok := parsed.(*rsa.PrivateKey); ok {
				result = rsaKey
			} else {
				err = fmt.Errorf(ERROR_MESSAGE_NOT_RSA_KEY)
				log.Printf("%v", err)
			}
		}
	}
	return result, err
}

func ReadNextSerial(serialPath string) int64 {
	result := SERIAL_FALLBACK_VALUE
	err := error(nil)
	data := make([]byte, 0)
	if data, err = os.ReadFile(serialPath); err == nil {
		text := strings.TrimSpace(string(data))
		if parsed, parseError := strconv.ParseInt(text, 10, 64); parseError == nil && parsed > 0 {
			result = parsed
		}
	}
	return result
}

func ResolveExecutableDirectory() (string, error) {
	result := ""
	err := error(nil)
	executablePath := ""
	if executablePath, err = os.Executable(); err != nil {
		log.Printf("Failed to get executable path: %v", err)
	} else {
		result = filepath.Dir(executablePath)
	}
	return result, err
}

func SelectDefaultIPv4() string {
	result := ""
	err := error(nil)
	candidates := make([]string, 0)
	addresses := make([]net.Addr, 0)
	if addresses, err = net.InterfaceAddrs(); err == nil {
		for _, address := range addresses {
			ip := net.IP(nil)
			switch value := address.(type) {
			case *net.IPNet:
				ip = value.IP
			case *net.IPAddr:
				ip = value.IP
			}
			if ip != nil {
				ipv4 := ip.To4()
				if ipv4 != nil && !ipv4.IsLoopback() {
					candidates = append(candidates, ipv4.String())
				}
			}
		}
	} else {
		log.Printf("Failed to get interface addresses: %v", err)
	}
	prefixes := []string{IP_PREFIX_192, IP_PREFIX_10, IP_PREFIX_172}
	for _, prefix := range prefixes {
		for _, candidate := range candidates {
			if strings.HasPrefix(candidate, prefix) {
				result = candidate
				break
			}
		}
		if result != "" {
			break
		}
	}
	if result == "" && len(candidates) > 0 {
		result = candidates[0]
	}
	return result
}

func WriteCertificateFile(path string, certificateDer []byte) error {
	err := error(nil)
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  CERTIFICATE_TYPE,
		Bytes: certificateDer,
	})
	if err = os.WriteFile(path, pemBytes, FILE_PERMISSION_PUBLIC); err != nil {
		log.Printf("Failed to write certificate file: %v", err)
	}
	return err
}

func WriteNextSerial(serialPath string, nextSerial int64) error {
	err := error(nil)
	text := strconv.FormatInt(nextSerial, 10) + "\r\n"
	if err = os.WriteFile(serialPath, []byte(text), FILE_PERMISSION_PRIVATE); err != nil {
		log.Printf("Failed to write serial file: %v", err)
	}
	return err
}

func WritePrivateKeyFile(path string, privateKey *rsa.PrivateKey) error {
	err := error(nil)
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  PEM_TYPE_RSA_PRIVATE_KEY,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err = os.WriteFile(path, pemBytes, FILE_PERMISSION_PRIVATE); err != nil {
		log.Printf("Failed to write private key file: %v", err)
	}
	return err
}

func WriteSanFile(path string, ipTexts []string) error {
	err := error(nil)
	data := SAN{
		DNS: []string{"localhost"},
		IP:  ipTexts,
	}
	jsonBytes := make([]byte, 0)
	if jsonBytes, err = json.MarshalIndent(data, "", "  "); err != nil {
		log.Printf("Failed to marshal SAN file: %v", err)
	} else if err = os.WriteFile(path, append(jsonBytes, '\n'), FILE_PERMISSION_PUBLIC); err != nil {
		log.Printf("Failed to write SAN file: %v", err)
	}
	return err
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_ISSUE_IP_ADDRESS, logger.SKIP_STACK_FRAMES_BASE)
}
