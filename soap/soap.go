// Package soap
// File:        soap.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/soap/soap.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: SOAP provides functionality for SOAP web service calls with support for attachments and headers
// --------------------------------------------------------------------------------
package soap

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/xiang-tai-duo/go-boost/http"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
const (
	DEFAULT_SOAP_CONTENT_TYPE   = "text/xml;charset=UTF-8"
	SOAP_MULTIPART_CONTENT_TYPE = "multipart/related"
	DEFAULT_BOUNDARY_PREFIX     = "----=_Part_"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
var (
	//go:embed soap/*.xml
	SOAP_TEMPLATES             embed.FS
	envelopeTemplate           string
	envelopeWithHeaderTemplate string
)

type (
	SOAP struct {
		*http.HTTP
	}
)

//goland:noinspection GoUnusedExportedFunction
func New() *SOAP {
	return &SOAP{
		HTTP: http.New(),
	}
}

func (s *SOAP) CreateEnvelope(body string) (string, error) {
	result := ""
	err := error(nil)
	if envelopeTemplate == "" {
		content, err := SOAP_TEMPLATES.ReadFile("soap/xml/envelope.xml")
		if err == nil {
			envelopeTemplate = string(content)
		}
	}
	if err == nil {
		result = fmt.Sprintf(envelopeTemplate, body)
	}
	return result, err
}

func (s *SOAP) Invoke(url string, action string, body string) (string, int, error) {
	result := ""
	statusCode := 0
	err := error(nil)
	envelope, err := s.CreateEnvelope(body)
	if err == nil {
		result, statusCode, err = s.Do(action, url, DEFAULT_SOAP_CONTENT_TYPE, envelope)
	}
	return result, statusCode, err
}

func (s *SOAP) InvokeWithHeader(url string, action string, header string, body string) (string, int, error) {
	result := ""
	statusCode := 0
	err := error(nil)
	if envelopeWithHeaderTemplate == "" {
		content, err := SOAP_TEMPLATES.ReadFile("soap/xml/envelope_with_header.xml")
		if err == nil {
			envelopeWithHeaderTemplate = string(content)
		}
	}
	if err == nil {
		envelope := fmt.Sprintf(envelopeWithHeaderTemplate, header, body)
		result, statusCode, err = s.Do(action, url, DEFAULT_SOAP_CONTENT_TYPE, envelope)
	}
	return result, statusCode, err
}

func (s *SOAP) InvokeWithAttachment(url string, action string, body string, attachment []byte, fileName string, mimeType string) (string, int, error) {
	result := ""
	statusCode := 0
	err := error(nil)
	boundary := DEFAULT_BOUNDARY_PREFIX + fmt.Sprintf("%d_%d", time.Now().UnixNano(), len(attachment))
	envelope, err := s.CreateEnvelope(body)
	if err == nil {
		multipartBody := fmt.Sprintf("--%s\r\n", boundary)
		multipartBody += fmt.Sprintf("Content-Type: %s\r\n", DEFAULT_SOAP_CONTENT_TYPE)
		multipartBody += "Content-Transfer-Encoding: 8bit\r\n"
		multipartBody += "Content-ID: <soapPart>\r\n\r\n"
		multipartBody += envelope + "\r\n\r\n"
		multipartBody += fmt.Sprintf("--%s\r\n", boundary)
		multipartBody += fmt.Sprintf("Content-Type: %s\r\n", mimeType)
		multipartBody += "Content-Transfer-Encoding: binary\r\n"
		multipartBody += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName)
		multipartBody += fmt.Sprintf("Content-ID: <%s>\r\n\r\n", fileName)
		multipartBody += string(attachment) + "\r\n"
		multipartBody += fmt.Sprintf("--%s--\r\n", boundary)
		contentType := fmt.Sprintf("%s; boundary=%s; type=\"%s\"", SOAP_MULTIPART_CONTENT_TYPE, boundary, DEFAULT_SOAP_CONTENT_TYPE)
		result, statusCode, err = s.Do(action, url, contentType, multipartBody)
	}
	return result, statusCode, err
}

//goland:noinspection SpellCheckingInspection
func (s *SOAP) ExtractBody(soapResponse string) string {
	result := soapResponse
	start := strings.Index(soapResponse, "<soapenv:Body>")
	if start == -1 {
		start = strings.Index(soapResponse, "<Body>")
	}
	if start != -1 {
		start += strings.Index(soapResponse[start:], ">") + 1
		end := strings.Index(soapResponse, "</soapenv:Body>")
		if end == -1 {
			end = strings.Index(soapResponse, "</Body>")
		}
		if end != -1 {
			extracted := soapResponse[start:end]
			extracted = strings.TrimSpace(extracted)
			result = extracted
		}
	}
	return result
}
