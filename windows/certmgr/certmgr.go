//go:build windows

// Package certmgr
// File:        certmgr.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/windows/certmgr/certmgr.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: certmgr.exe equivalent wrapper for Windows certificate store API functions
// --------------------------------------------------------------------------------
package certmgr

/*
#cgo LDFLAGS: -lcrypt32
#include <windows.h>
#include <wincrypt.h>
#include <stdlib.h>

static HCERTSTORE goCertOpenLocalMachineStore(LPCWSTR storeName) {
    return CertOpenStore(CERT_STORE_PROV_SYSTEM_W, 0, 0, CERT_SYSTEM_STORE_LOCAL_MACHINE, storeName);
}

static BOOL goCryptQueryCertificateFile(LPCWSTR fileName, HCERTSTORE *store, HCRYPTMSG *message, const CERT_CONTEXT **context) {
    DWORD encodingType = 0;
    DWORD contentType = 0;
    DWORD formatType = 0;
    const void *queryContext = NULL;
    BOOL result = CryptQueryObject(
        CERT_QUERY_OBJECT_FILE,
        fileName,
        CERT_QUERY_CONTENT_FLAG_CERT |
        CERT_QUERY_CONTENT_FLAG_SERIALIZED_CERT |
        CERT_QUERY_CONTENT_FLAG_PKCS7_SIGNED |
        CERT_QUERY_CONTENT_FLAG_PKCS7_SIGNED_EMBED,
        CERT_QUERY_FORMAT_FLAG_ALL,
        0,
        &encodingType,
        &contentType,
        &formatType,
        store,
        message,
        &queryContext
    );
    if (result) {
        *context = (PCCERT_CONTEXT)queryContext;
    }
    return result;
}

static BOOL goCryptQueryCertificateData(BYTE *data, DWORD dataSize, HCERTSTORE *store, HCRYPTMSG *message, const CERT_CONTEXT **context) {
    DWORD encodingType = 0;
    DWORD contentType = 0;
    DWORD formatType = 0;
    const void *queryContext = NULL;
    CRYPT_DATA_BLOB blob;
    blob.cbData = dataSize;
    blob.pbData = data;
    BOOL result = CryptQueryObject(
        CERT_QUERY_OBJECT_BLOB,
        &blob,
        CERT_QUERY_CONTENT_FLAG_CERT |
        CERT_QUERY_CONTENT_FLAG_SERIALIZED_CERT |
        CERT_QUERY_CONTENT_FLAG_PKCS7_SIGNED |
        CERT_QUERY_CONTENT_FLAG_PKCS7_SIGNED_EMBED,
        CERT_QUERY_FORMAT_FLAG_ALL,
        0,
        &encodingType,
        &contentType,
        &formatType,
        store,
        message,
        &queryContext
    );
    if (result) {
        *context = (PCCERT_CONTEXT)queryContext;
    }
    return result;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/xiang-tai-duo/go-boost/logger"
)

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
	ERROR_CRYPT_QUERY_OBJECT_FAILED       = "CryptQueryObject failed: %d"
	ERROR_CSTRING_FAILED                  = "CString failed"
	ERROR_FREE_CERTIFICATE_CONTEXT_FAILED = "CertFreeCertificateContext failed: %d"
	MODULE_NAME_CERTMGR                   = "windows.certmgr"
	ROOT_STORE_NAME                       = "root"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_CERTMGR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_CERTMGR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_CERTMGR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection DuplicatedCode,GoUnusedExportedFunction
func CString(s string) *C.wchar_t {
	result := (*C.wchar_t)(nil)
	if utf16, err := syscall.UTF16FromString(s); err == nil {
		if ptr := C.malloc(C.size_t((len(utf16) + 1) * int(unsafe.Sizeof(uint16(0))))); ptr != nil {
			src := utf16
			dst := unsafe.Slice((*uint16)(ptr), len(utf16)+1)
			copy(dst, src)
			dst[len(utf16)] = 0
			result = (*C.wchar_t)(ptr)
		}
	}
	return result
}

//goland:noinspection GoUnusedExportedFunction
func FreeCString(ptr *C.wchar_t) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func AddCertificateBytesToStore(data []byte, storeName string) error {
	result := error(nil)
	err := error(nil)
	if len(data) == 0 {
		err = errors.New(ERROR_CERTIFICATE_DATA_EMPTY)
	} else {
		err = addCertificateToStore(storeName, func(targetStore C.HCERTSTORE) error {
			return addCertificateDataToStore(data, targetStore)
		})
	}
	result = err
	return result
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func AddCertificateToStore(certificatePath string, storeName string) error {
	result := error(nil)
	err := error(nil)
	if certificatePath == "" {
		err = errors.New(ERROR_CERTIFICATE_PATH_EMPTY)
	} else if _, statErr := os.Stat(certificatePath); statErr != nil {
		err = errors.New(ERROR_CERTIFICATE_FILE_NOT_EXISTS)
	} else {
		err = addCertificateToStore(storeName, func(targetStore C.HCERTSTORE) error {
			return addCertificateFileToStore(certificatePath, targetStore)
		})
	}
	result = err
	return result
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func AddRootCertificate(certificatePath string) error {
	result := AddCertificateToStore(certificatePath, ROOT_STORE_NAME)
	return result
}

//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func AddRootCertificateBytes(data []byte) error {
	result := AddCertificateBytesToStore(data, ROOT_STORE_NAME)
	return result
}

//goland:noinspection GoSnakeCaseUsage
func GetLastError() uint32 {
	return uint32(C.GetLastError())
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_CERTMGR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoSnakeCaseUsage
func addCertificateContextToStore(context *C.CERT_CONTEXT, store C.HCERTSTORE) error {
	result := error(nil)
	if C.CertAddCertificateContextToStore(store, context, C.DWORD(CERT_STORE_ADD_REPLACE_EXISTING), nil) == 0 {
		result = fmt.Errorf(ERROR_ADD_CERTIFICATE_CONTEXT_FAILED, GetLastError())
	}
	return result
}

//goland:noinspection GoSnakeCaseUsage
func addCertificateDataToStore(data []byte, targetStore C.HCERTSTORE) error {
	result := error(nil)
	var sourceStore C.HCERTSTORE
	var message C.HCRYPTMSG
	var context *C.CERT_CONTEXT
	if C.goCryptQueryCertificateData((*C.BYTE)(unsafe.Pointer(&data[0])), C.DWORD(len(data)), &sourceStore, &message, &context) == 0 {
		result = fmt.Errorf(ERROR_CRYPT_QUERY_OBJECT_FAILED, GetLastError())
	} else {
		result = addQueriedCertificateToStore(sourceStore, message, context, targetStore)
	}
	return result
}

//goland:noinspection GoSnakeCaseUsage
func addCertificateFileToStore(certificatePath string, targetStore C.HCERTSTORE) error {
	result := error(nil)
	if certificatePtr := CString(certificatePath); certificatePtr == nil {
		result = errors.New(ERROR_CSTRING_FAILED)
	} else {
		defer FreeCString(certificatePtr)
		var sourceStore C.HCERTSTORE
		var message C.HCRYPTMSG
		var context *C.CERT_CONTEXT
		if C.goCryptQueryCertificateFile((*C.WCHAR)(unsafe.Pointer(certificatePtr)), &sourceStore, &message, &context) == 0 {
			result = fmt.Errorf(ERROR_CRYPT_QUERY_OBJECT_FAILED, GetLastError())
		} else {
			result = addQueriedCertificateToStore(sourceStore, message, context, targetStore)
		}
	}
	return result
}

//goland:noinspection GoSnakeCaseUsage
func addCertificateStoreToStore(sourceStore C.HCERTSTORE, targetStore C.HCERTSTORE) error {
	result := error(nil)
	context := (*C.CERT_CONTEXT)(nil)
	count := 0
	for {
		context = C.CertEnumCertificatesInStore(sourceStore, context)
		if context == nil {
			break
		}
		count++
		if result = addCertificateContextToStore(context, targetStore); result != nil {
			break
		}
	}
	if result == nil && count == 0 {
		result = errors.New(ERROR_CERTIFICATE_STORE_EMPTY)
	}
	return result
}

//goland:noinspection GoSnakeCaseUsage
func addCertificateToStore(storeName string, addCertificate func(C.HCERTSTORE) error) error {
	result := error(nil)
	if storeNamePtr := CString(storeName); storeNamePtr == nil {
		result = errors.New(ERROR_CSTRING_FAILED)
	} else {
		defer FreeCString(storeNamePtr)
		var targetStore C.HCERTSTORE
		if targetStore = C.goCertOpenLocalMachineStore((*C.WCHAR)(unsafe.Pointer(storeNamePtr))); targetStore == nil {
			result = fmt.Errorf(ERROR_CERT_OPEN_STORE_FAILED, GetLastError())
		} else {
			defer closeStore(targetStore, &result)
			result = addCertificate(targetStore)
		}
	}
	return result
}

//goland:noinspection GoSnakeCaseUsage
func addQueriedCertificateToStore(sourceStore C.HCERTSTORE, message C.HCRYPTMSG, context *C.CERT_CONTEXT, targetStore C.HCERTSTORE) error {
	result := error(nil)
	defer closeMessage(message, &result)
	defer closeStore(sourceStore, &result)
	if context != nil {
		result = addCertificateContextToStore(context, targetStore)
		defer freeCertificateContext(context, &result)
	} else if sourceStore != nil {
		result = addCertificateStoreToStore(sourceStore, targetStore)
	} else {
		result = errors.New(ERROR_CERTIFICATE_STORE_EMPTY)
	}
	return result
}

//goland:noinspection GoSnakeCaseUsage
func closeMessage(message C.HCRYPTMSG, err *error) {
	if message != nil {
		if C.CryptMsgClose(message) == 0 && *err == nil {
			*err = fmt.Errorf(ERROR_CRYPT_MSG_CLOSE_FAILED, GetLastError())
		}
	}
}

//goland:noinspection GoSnakeCaseUsage
func closeStore(store C.HCERTSTORE, err *error) {
	if store != nil {
		if C.CertCloseStore(store, 0) == 0 && *err == nil {
			*err = fmt.Errorf(ERROR_CERT_CLOSE_STORE_FAILED, GetLastError())
		}
	}
}

//goland:noinspection GoSnakeCaseUsage
func freeCertificateContext(context *C.CERT_CONTEXT, err *error) {
	if context != nil {
		if C.CertFreeCertificateContext(context) == 0 && *err == nil {
			*err = fmt.Errorf(ERROR_FREE_CERTIFICATE_CONTEXT_FAILED, GetLastError())
		}
	}
}
