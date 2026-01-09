// Package boost
// File:        soap.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/soap.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: SOAP provides functionality for SOAP web service calls with support for attachments and headers
// --------------------------------------------------------------------------------
package boost

import (
	"embed"
	"fmt"
	"strings"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
const (
	DEFAULT_SOAP_CONTENT_TYPE   = "text/xml;charset=UTF-8"
	SOAP_MULTIPART_CONTENT_TYPE = "multipart/related"
	DEFAULT_BOUNDARY_PREFIX     = "----=_Part_"
)

//goland:noinspection GoSnakeCaseUsage
var (
	//go:embed soap/*.xml
	SOAP_TEMPLATES             embed.FS
	envelopeTemplate           string
	envelopeWithHeaderTemplate string
)

type (
	SOAP struct {
		*HTTP
	}
)

func NewSOAP() *SOAP {
	return &SOAP{
		HTTP: NewHTTP(),
	}
}

func (s *SOAP) CreateEnvelope(body string) (string, error) {
	var result string
	var err error

	if envelopeTemplate == "" {
		var content []byte
		content, err = SOAP_TEMPLATES.ReadFile("soap/envelope.xml")
		if err == nil {
			envelopeTemplate = string(content)
		}
	}
	if err == nil {
		result = fmt.Sprintf(envelopeTemplate, body)
	}
	return result, err
}

func (s *SOAP) Call(url string, action string, body string) (string, int, error) {
	var result string
	var statusCode int
	var err error

	var envelope string
	envelope, err = s.CreateEnvelope(body)
	if err == nil {
		result, statusCode, err = s.Do(action, url, DEFAULT_SOAP_CONTENT_TYPE, envelope)
	}
	return result, statusCode, err
}

func (s *SOAP) CallWithHeader(url string, action string, header string, body string) (string, int, error) {
	var result string
	var statusCode int
	var err error

	if envelopeWithHeaderTemplate == "" {
		var content []byte
		content, err = SOAP_TEMPLATES.ReadFile("soap/envelope_with_header.xml")
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

func (s *SOAP) CallWithAttachment(url string, action string, body string, attachment []byte, fileName string, mimeType string) (string, int, error) {
	var result string
	var statusCode int
	var err error

	boundary := DEFAULT_BOUNDARY_PREFIX + fmt.Sprintf("%d_%d", time.Now().UnixNano(), len(attachment))

	var envelope string
	envelope, err = s.CreateEnvelope(body)
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

		result, statusCode, err = s.Do("POST", url, contentType, multipartBody)
	}
	return result, statusCode, err
}

func (s *SOAP) ExtractBody(soapResponse string) string {
	var result string = soapResponse

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
