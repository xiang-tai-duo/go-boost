// Package main
// File:        main.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ca/server/main.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: Certificate authority server helper.
// --------------------------------------------------------------------------------
//goland:noinspection GoSnakeCaseUsage
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "common"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage
type (
	CERTIFICATE_PATHS struct {
		CaDirectory            string
		CertificatePath        string
		CertificatesDirectory  string
		PrivateKeyPath         string
		SerialPath             string
		ServeCaCertificatePath string
		ServeCertificatePath   string
		ServeDirectory         string
		ServePrivateKeyPath    string
	}
	SAN_CONFIG struct {
		Dns []string `json:"dns"`
		Ips []string `json:"ip"`
	}
)

//goland:noinspection GoSnakeCaseUsage
const (
	DEFAULT_SAN_JSON_CONTENT    = "{}"
	LOG_PREFIX_DONE             = "[DONE] "
	LOG_PREFIX_FAIL             = "[FAIL] "
	LOG_PREFIX_INFO             = "[INFO] "
	LOG_PREFIX_WARN             = "[WARN] "
	MODULE_NAME_MAIN            = "ca.server.main"
	ORGANIZATION_UNIT_NAME      = "go-boost Certificate Authority"
	PRIVATE_KEY_TYPE            = "RSA PRIVATE KEY"
	SAN_JSON_FILE_NAME          = "san.json"
	SERIAL_INITIAL_VALUE        = "1001"
	SERVE_CERTIFICATE_FILE_NAME = "serve.crt"
	SERVE_COMMON_NAME           = "go-boost Serve"
	SERVE_DIRECTORY             = "serve"
	SERVE_PRIVATE_KEY_FILE_NAME = "serve.key"
	SERVE_SERIAL_NUMBER         = 1
	SERVE_VALIDITY_YEARS        = 3
	SUBJECT_COMMON_NAME         = "go-boost Root CA"
)

//go:embed san.json
var embeddedFiles embed.FS

//goland:noinspection GoUnhandledErrorResult
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, LOG_PREFIX_FAIL+err.Error())
		os.Exit(1)
	}
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_MAIN, logger.SKIP_STACK_FRAMES_BASE)
}

func buildCertificatePaths(executableDirectory string) CERTIFICATE_PATHS {
	result := CERTIFICATE_PATHS{}
	certificatesDirectory := filepath.Join(executableDirectory, CERTIFICATES_DIRECTORY)
	caDirectory := filepath.Join(certificatesDirectory, CA_DIRECTORY)
	serveDirectory := filepath.Join(certificatesDirectory, SERVE_DIRECTORY)
	result = CERTIFICATE_PATHS{
		CaDirectory:            caDirectory,
		CertificatePath:        filepath.Join(caDirectory, CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME),
		CertificatesDirectory:  filepath.Join(executableDirectory, CERTIFICATES_DIRECTORY),
		PrivateKeyPath:         filepath.Join(caDirectory, CERTIFICATE_AUTHORITY_PRIVATE_KEY_FILE_NAME),
		SerialPath:             filepath.Join(caDirectory, CERTIFICATE_AUTHORITY_SERIAL_FILE_NAME),
		ServeCaCertificatePath: filepath.Join(serveDirectory, CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME),
		ServeCertificatePath:   filepath.Join(serveDirectory, SERVE_CERTIFICATE_FILE_NAME),
		ServeDirectory:         serveDirectory,
		ServePrivateKeyPath:    filepath.Join(serveDirectory, SERVE_PRIVATE_KEY_FILE_NAME),
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func copyFile(sourcePath, targetPath string) error {
	err := error(nil)
	data := make([]byte, 0)
	if data, err = os.ReadFile(sourcePath); err != nil {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to read file: %s, error: %v\n", sourcePath, err)
	} else {
		if err = os.WriteFile(targetPath, data, FILE_PERMISSION_PUBLIC); err != nil {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write file: %s, error: %v\n", targetPath, err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func createCertificate(privateKey *rsa.PrivateKey) ([]byte, *x509.Certificate, error) {
	result := make([]byte, 0)
	err := error(nil)
	resultCertificate := (*x509.Certificate)(nil)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	var serialNumber *big.Int
	if serialNumber, err = rand.Int(rand.Reader, serialNumberLimit); err == nil {
		validFrom := time.Now()
		validTo := validFrom.AddDate(CERTIFICATE_AUTHORITY_VALIDITY_YEARS, 0, 0)
		certificateTemplate := x509.Certificate{
			SerialNumber: serialNumber,
			Subject: pkix.Name{
				Country:            []string{COUNTRY_NAME},
				Organization:       []string{ORGANIZATION_NAME},
				OrganizationalUnit: []string{ORGANIZATION_UNIT_NAME},
				CommonName:         SUBJECT_COMMON_NAME,
			},
			NotBefore:             validFrom,
			NotAfter:              validTo,
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
			BasicConstraintsValid: true,
			IsCA:                  true,
			MaxPathLen:            1,
			MaxPathLenZero:        false,
		}
		certificateDer, certErr := x509.CreateCertificate(rand.Reader, &certificateTemplate, &certificateTemplate, &privateKey.PublicKey, privateKey)
		if certErr != nil {
			err = certErr
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to create certificate, error: %v\n", err)
		} else {
			result = certificateDer
			resultCertificate = &certificateTemplate
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to generate serial number, error: %v\n", err)
	}
	return result, resultCertificate, err
}

//goland:noinspection GoUnhandledErrorResult
func createServeCertificate(privateKey *rsa.PrivateKey, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey, sanConfig SAN_CONFIG) ([]byte, error) {
	result := make([]byte, 0)
	err := error(nil)
	validFrom := time.Now()
	validTo := validFrom.AddDate(SERVE_VALIDITY_YEARS, 0, 0)
	certificateTemplate := x509.Certificate{
		SerialNumber: big.NewInt(SERVE_SERIAL_NUMBER),
		Subject: pkix.Name{
			Country:            []string{COUNTRY_NAME},
			Organization:       []string{ORGANIZATION_NAME},
			OrganizationalUnit: []string{ORGANIZATION_UNIT_NAME},
			CommonName:         SERVE_COMMON_NAME,
		},
		NotBefore:             validFrom,
		NotAfter:              validTo,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              normalizeDNSNames(sanConfig.Dns),
		IPAddresses:           parseIPAddresses(sanConfig.Ips),
	}

	certificateDer := make([]byte, 0)
	if certificateDer, err = x509.CreateCertificate(rand.Reader, &certificateTemplate, caCertificate, &privateKey.PublicKey, caPrivateKey); err == nil {
		result = certificateDer
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to create serve certificate, error: %v\n", err)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func ensureSanConfig(outputDirectory string) (string, error) {
	result := ""
	err := error(nil)
	executablePath := filepath.Join(outputDirectory, SAN_JSON_FILE_NAME)
	workingPath := filepath.Join(getWorkingDirectory(), SAN_JSON_FILE_NAME)
	_, err = os.Stat(executablePath)
	if err == nil {
		result = executablePath
	} else if os.IsNotExist(err) {
		_, workingStatErr := os.Stat(workingPath)
		if workingStatErr == nil {
			result = workingPath
		} else if os.IsNotExist(workingStatErr) {
			result = executablePath
			fmt.Printf(LOG_PREFIX_INFO+"Writing default SAN config to: %s\n", result)
			if err = writeDefaultSanConfig(result); err != nil {
				fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write default SAN config, error: %v\n", err)
			}
		} else {
			err = workingStatErr
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to stat working path: %s, error: %v\n", workingPath, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to stat executable path: %s, error: %v\n", executablePath, err)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func generateRootCertificate(paths CERTIFICATE_PATHS) (*rsa.PrivateKey, *x509.Certificate, error) {
	resultPrivateKey := (*rsa.PrivateKey)(nil)
	resultCertificate := (*x509.Certificate)(nil)
	err := error(nil)
	fmt.Printf(LOG_PREFIX_INFO+"Generating RSA private key (%d bits)...\n", RSA_KEY_SIZE)
	if resultPrivateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE); err == nil {
		fmt.Println(LOG_PREFIX_INFO + "Generating self-signed CA root certificate...")
		cert := (*x509.Certificate)(nil)
		certificateDer := make([]byte, 0)
		if certificateDer, cert, err = createCertificate(resultPrivateKey); err == nil {
			resultCertificate = cert
			fmt.Printf(LOG_PREFIX_INFO+"Saving CA certificate to: %s\n", paths.CertificatePath)
			if err = writeCertificateFile(paths.CertificatePath, certificateDer); err != nil {
				fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write certificate file: %s, error: %v\n", paths.CertificatePath, err)
			} else {
				fmt.Printf(LOG_PREFIX_INFO+"Saving CA private key to: %s\n", paths.PrivateKeyPath)
				if err = writePrivateKeyFile(paths.PrivateKeyPath, resultPrivateKey); err != nil {
					fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write private key file: %s, error: %v\n", paths.PrivateKeyPath, err)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to create certificate, error: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to generate RSA private key, error: %v\n", err)
	}
	return resultPrivateKey, resultCertificate, err
}

//goland:noinspection GoUnhandledErrorResult
func generateServeCertificate(paths CERTIFICATE_PATHS, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey, sanConfig SAN_CONFIG) error {
	err := error(nil)
	fmt.Println(LOG_PREFIX_INFO + "Generating CA-signed serve certificate...")
	var servePrivateKey *rsa.PrivateKey
	if servePrivateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE); err == nil {
		var serveCertificateDer []byte
		if serveCertificateDer, err = createServeCertificate(servePrivateKey, caCertificate, caPrivateKey, sanConfig); err == nil {
			fmt.Printf(LOG_PREFIX_INFO+"Saving serve certificate to: %s\n", paths.ServeCertificatePath)
			if err = writeCertificateFile(paths.ServeCertificatePath, serveCertificateDer); err != nil {
				fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write serve certificate file: %s, error: %v\n", paths.ServeCertificatePath, err)
			} else {
				fmt.Printf(LOG_PREFIX_INFO+"Saving serve private key to: %s\n", paths.ServePrivateKeyPath)
				if err = writePrivateKeyFile(paths.ServePrivateKeyPath, servePrivateKey); err != nil {
					fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write serve private key file: %s, error: %v\n", paths.ServePrivateKeyPath, err)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to create serve certificate, error: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to generate serve RSA private key, error: %v\n", err)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func getWorkingDirectory() string {
	result := ""
	var err error
	if result, err = os.Getwd(); err != nil {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to get working directory, error: %v\n", err)
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func loadSanConfig(path string) (SAN_CONFIG, error) {
	result := SAN_CONFIG{}
	err := error(nil)
	var data []byte
	if data, err = os.ReadFile(path); err == nil {
		if err = json.Unmarshal(data, &result); err != nil {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to unmarshal SAN config, error: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to read SAN config file: %s, error: %v\n", path, err)
	}
	return result, err
}

func normalizeDNSNames(values []string) []string {
	result := make([]string, 0)
	seen := make(map[string]bool)
	for _, value := range values {
		name := strings.TrimSpace(value)
		if name != "" && !seen[name] {
			result = append(result, name)
			seen[name] = true
		}
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func parseIPAddresses(values []string) []net.IP {
	result := make([]net.IP, 0)
	seen := make(map[string]bool)
	for _, value := range values {
		ipText := strings.TrimSpace(value)
		ip := net.ParseIP(ipText)
		if ip != nil && !seen[ip.String()] {
			result = append(result, ip)
			seen[ip.String()] = true
		}
	}
	if addresses, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addresses {
			var ip net.IP
			switch value := address.(type) {
			case *net.IPNet:
				ip = value.IP
			case *net.IPAddr:
				ip = value.IP
			}
			if ip != nil && !seen[ip.String()] {
				result = append(result, ip)
				seen[ip.String()] = true
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_WARN+"Failed to get interface addresses, error: %v\n", err)
	}
	return result
}

//goland:noinspection GoUnhandledErrorResult
func prepareCertificateDirectories(paths CERTIFICATE_PATHS) error {
	err := error(nil)
	statErr := error(nil)
	if _, statErr = os.Stat(paths.ServeDirectory); statErr == nil {
		err = fmt.Errorf("serve directory already exists: %s", paths.ServeDirectory)
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Serve directory already exists: %s\n", paths.ServeDirectory)
	} else if !os.IsNotExist(statErr) {
		err = statErr
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to stat serve directory: %s, error: %v\n", paths.ServeDirectory, err)
	}
	if err == nil {
		fmt.Printf(LOG_PREFIX_INFO+"Creating CA directory: %s\n", paths.CaDirectory)
		if err = os.MkdirAll(paths.CaDirectory, DIRECTORY_PERMISSION); err != nil {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to create CA directory: %s, error: %v\n", paths.CaDirectory, err)
		} else {
			fmt.Printf(LOG_PREFIX_INFO+"Creating serve directory: %s\n", paths.ServeDirectory)
			if err = os.MkdirAll(paths.ServeDirectory, DIRECTORY_PERMISSION); err != nil {
				fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to create serve directory: %s, error: %v\n", paths.ServeDirectory, err)
			}
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func prepareSanConfig(outputDirectory string) (SAN_CONFIG, error) {
	result := SAN_CONFIG{}
	err := error(nil)
	var sanPath string
	if sanPath, err = ensureSanConfig(outputDirectory); err == nil {
		fmt.Printf(LOG_PREFIX_INFO+"Loading SAN config from: %s\n", sanPath)
		if result, err = loadSanConfig(sanPath); err != nil {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to load SAN config, error: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to ensure SAN config, error: %v\n", err)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func resolveOutputDirectory() (string, error) {
	result := ""
	err := error(nil)
	var executablePath string
	if executablePath, err = os.Executable(); err == nil {
		result = filepath.Dir(executablePath)
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to get executable path, error: %v\n", err)
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func run() error {
	err := error(nil)
	var outputDirectory string
	if outputDirectory, err = resolveOutputDirectory(); err == nil {
		paths := buildCertificatePaths(outputDirectory)
		if err = prepareCertificateDirectories(paths); err != nil {
			fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to prepare certificate directories, error: %v\n", err)
		} else {
			var sanConfig SAN_CONFIG
			if sanConfig, err = prepareSanConfig(outputDirectory); err == nil {
				var caPrivateKey *rsa.PrivateKey
				var caCertificate *x509.Certificate
				if caPrivateKey, caCertificate, err = generateRootCertificate(paths); err == nil {
					if err = writeSerialFile(paths.SerialPath); err != nil {
						fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write serial file, error: %v\n", err)
					} else {
						if err = generateServeCertificate(paths, caCertificate, caPrivateKey, sanConfig); err != nil {
							fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to generate serve certificate, error: %v\n", err)
						} else {
							fmt.Printf(LOG_PREFIX_INFO+"Copying CA certificate to serve directory: %s\n", paths.ServeCaCertificatePath)
							if err = copyFile(paths.CertificatePath, paths.ServeCaCertificatePath); err != nil {
								fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to copy CA certificate, error: %v\n", err)
							} else {
								fmt.Println(LOG_PREFIX_DONE + "CA root and serve certificate generated successfully.")
								fmt.Println(LOG_PREFIX_INFO + "Please keep ca.key in a safe place.")
							}
						}
					}
				} else {
					fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to generate root certificate, error: %v\n", err)
				}
			} else {
				fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to prepare SAN config, error: %v\n", err)
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to resolve output directory, error: %v\n", err)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func writeCertificateFile(path string, certificateDer []byte) error {
	err := error(nil)
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  CERTIFICATE_TYPE,
		Bytes: certificateDer,
	})
	if err = os.WriteFile(path, pemBytes, FILE_PERMISSION_PUBLIC); err != nil {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write certificate file: %s, error: %v\n", path, err)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func writeDefaultSanConfig(path string) error {
	err := error(nil)
	var data []byte
	var readErr error
	if data, readErr = embeddedFiles.ReadFile(SAN_JSON_FILE_NAME); readErr != nil {
		data = []byte(DEFAULT_SAN_JSON_CONTENT)
	}
	if err = os.WriteFile(path, data, FILE_PERMISSION_PUBLIC); err != nil {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write default SAN config: %s, error: %v\n", path, err)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func writePrivateKeyFile(path string, privateKey *rsa.PrivateKey) error {
	err := error(nil)
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  PRIVATE_KEY_TYPE,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err = os.WriteFile(path, pemBytes, FILE_PERMISSION_PRIVATE); err != nil {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write private key file: %s, error: %v\n", path, err)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult
func writeSerialFile(path string) error {
	err := error(nil)
	if err = os.WriteFile(path, []byte(SERIAL_INITIAL_VALUE+"\r\n"), FILE_PERMISSION_PRIVATE); err != nil {
		fmt.Fprintf(os.Stderr, LOG_PREFIX_FAIL+"Failed to write serial file: %s, error: %v\n", path, err)
	}
	return err
}
