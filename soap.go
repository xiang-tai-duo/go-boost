// Package boost
// File:        soap.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: SOAP provides functionality for SOAP web service calls with support for attachments and headers
package boost

import (
	"embed"
	"fmt"
	"strings"
	"time"
)

const (
	DEFAULT_SOAP_CONTENT_TYPE  = "text/xml;charset=UTF-8"
	SOAP_MULTIPART_CONTENT_TYPE = "multipart/related"
	DEFAULT_BOUNDARY_PREFIX     = "----=_Part_"
)

//go:embed soap/*.xml
var SOAP_TEMPLATES embed.FS

var (
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
	if envelopeTemplate == "" {
		content, err := SOAP_TEMPLATES.ReadFile("soap/envelope.xml")
		if err != nil {
			return "", err
		}
		envelopeTemplate = string(content)
	}
	return fmt.Sprintf(envelopeTemplate, body), nil
}

func (s *SOAP) Call(url string, action string, body string) (string, int, error) {
	envelope, err := s.CreateEnvelope(body)
	if err != nil {
		return "", 0, err
	}
	return s.Do("POST", url, DEFAULT_SOAP_CONTENT_TYPE, envelope)
}

func (s *SOAP) CallWithHeader(url string, action string, header string, body string) (string, int, error) {
	if envelopeWithHeaderTemplate == "" {
		content, err := SOAP_TEMPLATES.ReadFile("soap/envelope_with_header.xml")
		if err != nil {
			return "", 0, err
		}
		envelopeWithHeaderTemplate = string(content)
	}
	envelope := fmt.Sprintf(envelopeWithHeaderTemplate, header, body)
	return s.Do("POST", url, DEFAULT_SOAP_CONTENT_TYPE, envelope)
}

func (s *SOAP) CallWithAttachment(url string, action string, body string, attachment []byte, fileName string, mimeType string) (string, int, error) {
	boundary := DEFAULT_BOUNDARY_PREFIX + fmt.Sprintf("%d_%d", time.Now().UnixNano(), len(attachment))

	envelope, err := s.CreateEnvelope(body)
	if err != nil {
		return "", 0, err
	}

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

	return s.Do("POST", url, contentType, multipartBody)
}

func (s *SOAP) ExtractBody(soapResponse string) string {
	start := strings.Index(soapResponse, "<soapenv:Body>")
	if start == -1 {
		start = strings.Index(soapResponse, "<Body>")
	}
	if start == -1 {
		return soapResponse
	}
	start += strings.Index(soapResponse[start:], ">" ) + 1
	end := strings.Index(soapResponse, "</soapenv:Body>")
	if end == -1 {
		end = strings.Index(soapResponse, "</Body>")
	}
	if end == -1 {
		return soapResponse
	}
	extracted := soapResponse[start:end]
	extracted = strings.TrimSpace(extracted)
	return extracted
}