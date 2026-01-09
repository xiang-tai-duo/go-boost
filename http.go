// Package boost
// File:        http.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/http.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: HTTP provides utility methods for HTTP/HTTPS/WS/WSS operations,
//              including support for self-signed certificates.
// --------------------------------------------------------------------------------
package boost

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DEFAULT_HTTP_TIMEOUT = 30 * time.Second
)

type (
	HTTP struct {
		client  *http.Client
		timeout time.Duration
	}
)

var (
	AllowSelfSignedCertificates = false
)

func (h *HTTP) Delete(requestURL string) (string, int, error) {
	return h.Do("DELETE", requestURL, "", "")
}

func (h *HTTP) Do(method string, requestURL string, contentType string, body string) (string, int, error) {
	var err error
	var responseBody string
	var statusCode int
	var response *http.Response
	var request *http.Request
	var requestBody io.Reader
	if body != "" {
		requestBody = io.NopCloser(strings.NewReader(body))
	}
	if request, err = http.NewRequest(method, requestURL, requestBody); err == nil {
		if contentType != "" {
			request.Header.Set("Content-Type", contentType)
		}
		if response, err = h.client.Do(request); err == nil {
			defer func(response *http.Response) {
				_ = response.Body.Close()
			}(response)
			statusCode = response.StatusCode
			var responseBodyBytes []byte
			if responseBodyBytes, err = io.ReadAll(response.Body); err == nil {
				responseBody = string(responseBodyBytes)
			}
		}
	}
	return responseBody, statusCode, err
}

func (h *HTTP) Get(requestURL string) (string, int, error) {
	return h.Do("GET", requestURL, "", "")
}

func (h *HTTP) GetTimeout() time.Duration {
	return h.timeout
}

func NewHTTP() *HTTP {
	timeout := DEFAULT_HTTP_TIMEOUT
	return &HTTP{
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: AllowSelfSignedCertificates,
				},
			},
		},
	}
}

func (h *HTTP) Post(requestURL string, contentType string, body string) (string, int, error) {
	return h.Do("POST", requestURL, contentType, body)
}

func (h *HTTP) Put(requestURL string, contentType string, body string) (string, int, error) {
	return h.Do("PUT", requestURL, contentType, body)
}

func SetAllowSelfSignedCertificates(allow bool) {
	AllowSelfSignedCertificates = allow
}

func (h *HTTP) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
	h.client.Timeout = timeout
}

func (h *HTTP) UpdateClientTransport() error {
	var err error
	h.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: AllowSelfSignedCertificates,
		},
	}
	return err
}
