// Package soap
// File:        soap.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/soap/soap.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SOAP provides functionality for SOAP web service calls with support for attachments, headers and basic authentication.
// --------------------------------------------------------------------------------
package soap

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
	"time"

	httplib "net/http"

	"github.com/xiang-tai-duo/go-boost/http2"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	SOAP struct {
		*http2.HTTP
		password string
		userName string
	}
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,GoUnusedConst
const (
	BODY_END_TAG                    = "</Body>"
	BODY_START_TAG                  = "<Body>"
	CONTENT_DISPOSITION_FORMAT      = "Content-Disposition: attachment; filename=\"%s\"\r\n"
	CONTENT_ID_FORMAT               = "Content-ID: <%s>\r\n\r\n"
	CONTENT_ID_SOAP_PART            = "Content-ID: <soapPart>\r\n\r\n"
	CONTENT_TRANSFER_ENCODING_8BIT  = "Content-Transfer-Encoding: 8bit\r\n"
	CONTENT_TRANSFER_ENCODING_BIN   = "Content-Transfer-Encoding: binary\r\n"
	CONTENT_TYPE_FORMAT             = "Content-Type: %s\r\n"
	CONTENT_TYPE_HEADER             = "Content-Type"
	CRLF                            = "\r\n"
	DEFAULT_BOUNDARY_PREFIX         = "----=_Part_"
	DEFAULT_SOAP_CONTENT_TYPE       = "text/xml;charset=UTF-8"
	ENVELOPE_TEMPLATE_PATH          = "xml/envelope.xml"
	ENVELOPE_WITH_HEADER_PATH       = "xml/envelope_with_header.xml"
	MODULE_NAME_SOAP                = "soap"
	MULTIPART_BOUNDARY_END_FORMAT   = "--%s--\r\n"
	MULTIPART_BOUNDARY_START_FORMAT = "--%s\r\n"
	MULTIPART_CONTENT_TYPE_FORMAT   = "%s; boundary=%s; type=\"%s\""
	SOAP_ACTION_HEADER              = "SOAPAction"
	SOAP_BODY_END_TAG               = "</soapenv:Body>"
	SOAP_BODY_START_TAG             = "<soapenv:Body>"
	SOAP_MULTIPART_CONTENT_TYPE     = "multipart/related"
	TAG_CLOSE                       = ">"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed xml/*.xml
	SOAP_TEMPLATES             embed.FS
	envelopeTemplate           string
	envelopeWithHeaderTemplate string
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SOAP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SOAP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SOAP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SOAP, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func New() *SOAP {
	__debug("Creating new SOAP instance")
	result := &SOAP{
		HTTP: http2.New(),
	}
	return result
}

func (s *SOAP) CreateEnvelope(body string) (string, error) {
	__debug(fmt.Sprintf("CreateEnvelope called, body length=%d", len(body)))
	result := ""
	err := error(nil)
	if envelopeTemplate == "" {
		var content []byte
		if content, err = SOAP_TEMPLATES.ReadFile(ENVELOPE_TEMPLATE_PATH); err == nil {
			envelopeTemplate = string(content)
			__debug(fmt.Sprintf("Envelope template loaded from %s, length=%d", ENVELOPE_TEMPLATE_PATH, len(envelopeTemplate)))
		} else {
			__error(fmt.Sprintf("Failed to load envelope template %s: %v", ENVELOPE_TEMPLATE_PATH, err))
		}
	}
	if err == nil {
		result = fmt.Sprintf(envelopeTemplate, body)
		__debug(fmt.Sprintf("Envelope created, total length=%d", len(result)))
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func (s *SOAP) ExtractBody(soapResponse string) string {
	__debug(fmt.Sprintf("ExtractBody called, response length=%d", len(soapResponse)))
	result := soapResponse
	start := strings.Index(soapResponse, SOAP_BODY_START_TAG)
	if start == -1 {
		start = strings.Index(soapResponse, BODY_START_TAG)
	}
	if start != -1 {
		start += strings.Index(soapResponse[start:], TAG_CLOSE) + 1
		end := strings.Index(soapResponse, SOAP_BODY_END_TAG)
		if end == -1 {
			end = strings.Index(soapResponse, BODY_END_TAG)
		}
		if end != -1 {
			result = strings.TrimSpace(soapResponse[start:end])
			__debug(fmt.Sprintf("Extracted body length=%d", len(result)))
		} else {
			__error("ExtractBody: body end tag not found")
		}
	} else {
		__error("ExtractBody: body start tag not found")
	}
	return result
}

func (s *SOAP) Invoke(url string, action string, body string) (string, int, error) {
	__debug(fmt.Sprintf("Invoke url=%s, action=%s, body length=%d", url, action, len(body)))
	result := ""
	statusCode := 0
	err := error(nil)
	var envelope string
	if envelope, err = s.CreateEnvelope(body); err == nil {
		result, statusCode, err = s.doRequest(url, action, DEFAULT_SOAP_CONTENT_TYPE, envelope)
	} else {
		__error(fmt.Sprintf("Invoke failed to create envelope, url=%s, action=%s, error=%v", url, action, err))
	}
	return result, statusCode, err
}

func (s *SOAP) InvokeRaw(url string, action string, envelope string) (string, int, error) {
	__debug(fmt.Sprintf("InvokeRaw url=%s, action=%s, envelope length=%d", url, action, len(envelope)))
	return s.doRequest(url, action, DEFAULT_SOAP_CONTENT_TYPE, envelope)
}

func (s *SOAP) InvokeWithAttachment(url string, action string, body string, attachment []byte, fileName string, mimeType string) (string, int, error) {
	__debug(fmt.Sprintf("InvokeWithAttachment url=%s, action=%s, fileName=%s, mimeType=%s, attachment size=%d", url, action, fileName, mimeType, len(attachment)))
	result := ""
	statusCode := 0
	err := error(nil)
	var envelope string
	if envelope, err = s.CreateEnvelope(body); err == nil {
		boundary := DEFAULT_BOUNDARY_PREFIX + fmt.Sprintf("%d_%d", time.Now().UnixNano(), len(attachment))
		__debug(fmt.Sprintf("InvokeWithAttachment boundary=%s", boundary))
		multipartBody := fmt.Sprintf(MULTIPART_BOUNDARY_START_FORMAT, boundary)
		multipartBody += fmt.Sprintf(CONTENT_TYPE_FORMAT, DEFAULT_SOAP_CONTENT_TYPE)
		multipartBody += CONTENT_TRANSFER_ENCODING_8BIT
		multipartBody += CONTENT_ID_SOAP_PART
		multipartBody += envelope + CRLF + CRLF
		multipartBody += fmt.Sprintf(MULTIPART_BOUNDARY_START_FORMAT, boundary)
		multipartBody += fmt.Sprintf(CONTENT_TYPE_FORMAT, mimeType)
		multipartBody += CONTENT_TRANSFER_ENCODING_BIN
		multipartBody += fmt.Sprintf(CONTENT_DISPOSITION_FORMAT, fileName)
		multipartBody += fmt.Sprintf(CONTENT_ID_FORMAT, fileName)
		multipartBody += string(attachment) + CRLF
		multipartBody += fmt.Sprintf(MULTIPART_BOUNDARY_END_FORMAT, boundary)
		contentType := fmt.Sprintf(MULTIPART_CONTENT_TYPE_FORMAT, SOAP_MULTIPART_CONTENT_TYPE, boundary, DEFAULT_SOAP_CONTENT_TYPE)
		result, statusCode, err = s.doRequest(url, action, contentType, multipartBody)
	} else {
		__error(fmt.Sprintf("InvokeWithAttachment failed to create envelope, url=%s, action=%s, error=%v", url, action, err))
	}
	return result, statusCode, err
}

func (s *SOAP) InvokeWithHeader(url string, action string, header string, body string) (string, int, error) {
	__debug(fmt.Sprintf("InvokeWithHeader url=%s, action=%s, header length=%d, body length=%d", url, action, len(header), len(body)))
	result := ""
	statusCode := 0
	err := error(nil)
	if envelopeWithHeaderTemplate == "" {
		var content []byte
		if content, err = SOAP_TEMPLATES.ReadFile(ENVELOPE_WITH_HEADER_PATH); err == nil {
			envelopeWithHeaderTemplate = string(content)
			__debug(fmt.Sprintf("Envelope-with-header template loaded from %s, length=%d", ENVELOPE_WITH_HEADER_PATH, len(envelopeWithHeaderTemplate)))
		} else {
			__error(fmt.Sprintf("Failed to load envelope-with-header template %s: %v", ENVELOPE_WITH_HEADER_PATH, err))
		}
	}
	if err == nil {
		envelope := fmt.Sprintf(envelopeWithHeaderTemplate, header, body)
		result, statusCode, err = s.doRequest(url, action, DEFAULT_SOAP_CONTENT_TYPE, envelope)
	}
	return result, statusCode, err
}

func (s *SOAP) SetBasicAuth(userName string, password string) {
	__debug(fmt.Sprintf("SetBasicAuth called, userName length=%d, password length=%d", len(userName), len(password)))
	s.userName = userName
	s.password = password
}

//goland:noinspection GoUnhandledErrorResult
func (s *SOAP) doRequest(url string, action string, contentType string, body string) (string, int, error) {
	__debug(fmt.Sprintf("doRequest url=%s, action=%s, contentType=%s, body length=%d", url, action, contentType, len(body)))
	result := ""
	statusCode := 0
	err := error(nil)
	var request *httplib.Request
	if request, err = httplib.NewRequest(httplib.MethodPost, url, bytes.NewBuffer([]byte(body))); err == nil {
		request.Header.Set(CONTENT_TYPE_HEADER, contentType)
		if action != "" {
			request.Header.Set(SOAP_ACTION_HEADER, action)
		}
		if s.userName != "" || s.password != "" {
			request.SetBasicAuth(s.userName, s.password)
		}
		var response *httplib.Response
		if response, err = s.HTTP.GetClient().Do(request); err == nil {
			defer func(response *httplib.Response) {
				_ = response.Body.Close()
			}(response)
			statusCode = response.StatusCode
			var responseBodyBytes []byte
			if responseBodyBytes, err = io.ReadAll(response.Body); err == nil {
				result = string(responseBodyBytes)
				__debug(fmt.Sprintf("doRequest succeeded, url=%s, statusCode=%d, response length=%d", url, statusCode, len(result)))
			} else {
				__error(fmt.Sprintf("doRequest failed to read response body, url=%s, statusCode=%d, error=%v", url, statusCode, err))
			}
		} else {
			__error(fmt.Sprintf("doRequest HTTP call failed, url=%s, action=%s, error=%v", url, action, err))
		}
	} else {
		__error(fmt.Sprintf("doRequest failed to build request, url=%s, action=%s, error=%v", url, action, err))
	}
	return result, statusCode, err
}
