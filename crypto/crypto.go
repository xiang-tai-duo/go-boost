// Package crypto
// File:        crypto.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/crypto/crypto.go
// Author:      Vibe Coding
// Created:     2026/02/14 15:36:58
// Description: RSA encryption implementation
// --------------------------------------------------------------------------------
package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"time"
)

//goland:noinspection GoUnusedExportedFunction
func RSAEncrypt(data interface{}) error {
	// Generate RSA key pair with 4096 bits (highest strength)
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Convert input data to string
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case int:
		dataStr = strconv.Itoa(v)
	default:
		return fmt.Errorf("unsupported data type, only string and int are supported")
	}

	// Encrypt data with public key
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, []byte(dataStr))
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Base64 encode the encrypted data
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	fmt.Println("Encrypted and Base64 encoded data:")
	fmt.Println(encodedCiphertext)

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405000")

	// Save private key to file
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyFilename := fmt.Sprintf("private_key_%s.txt", timestamp)
	privateKeyFile, err := os.Create(privateKeyFilename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, privateKeyBlock); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	// Save public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyFilename := fmt.Sprintf("public_key_%s.txt", timestamp)
	publicKeyFile, err := os.Create(publicKeyFilename)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer publicKeyFile.Close()

	if err := pem.Encode(publicKeyFile, publicKeyBlock); err != nil {
		return fmt.Errorf("failed to write public key to file: %v", err)
	}

	fmt.Println("\nKeys saved to:")
	fmt.Printf("- %s (Private Key)\n", privateKeyFilename)
	fmt.Printf("- %s (Public Key)\n", publicKeyFilename)

	return nil
}
