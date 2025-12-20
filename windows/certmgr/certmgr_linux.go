//go:build linux

// File:        certmgr_linux.go
// Description: Linux stub for certmgr package (no-op)
// --------------------------------------------------------------------------------

package certmgr

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst,SpellCheckingInspection
const (
	CERT_STORE_ADD_REPLACE_EXISTING       = 3
	ERROR_ADD_CERTIFICATE_CONTEXT_FAILED  = "CertAddCertificateContextToStore failed: %d"
	ERROR_CERTIFICATE_DATA_EMPTY          = "certificate data is empty"
	ERROR_CERTIFICATE_FILE_NOT_EXISTS     = "certificate file does not exist"
	ERROR_CERTIFICATE_PATH_EMPTY          = "certificate path is empty"
	ERROR_CERTIFICATE_STORE_EMPTY         = "certificate store is empty"
	ERROR_CERT_CLOSE_STORE_FAILED         = "CertCloseStore failed: %d"
	ERROR_CERT_OPEN_STORE_FAILED          = "CertOpenStore failed: %d"
	ERROR_CRYPT_MSG_CLOSE_FAILED          = "CryptMsgClose failed: %d"
	ERROR_CRYPT_QUERY_FILE_FAILED         = "CryptQueryObject for file failed: %d"
	ERROR_CRYPT_QUERY_FROM_DATA_FAILED    = "CryptQueryObject for data failed: %d"
	ERROR_INVALID_CERTIFICATE_CONTEXT     = "invalid certificate context"
	ERROR_STORE_NAME_EMPTY                = "store name is empty"
	MODULE_NAME_CERTMGR                   = "windows.certmgr"
	ROOT_STORE_NAME                       = "ROOT"
)

func CString(s string) *uint16 {
	return nil
}

func FreeCString(ptr *uint16) {
}

func AddCertificateBytesToStore(data []byte, storeName string) error {
	return nil
}

func AddCertificateToStore(certificatePath string, storeName string) error {
	return nil
}

func AddRootCertificate(certificatePath string) error {
	return nil
}

func AddRootCertificateBytes(data []byte) error {
	return nil
}

func GetLastError() uint32 {
	return 0
}
