// Package http
// File:        http.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/http/http.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: HTTP provides utility methods for HTTP/HTTPS/WS/WSS operations, including support for self-signed certificates.
// --------------------------------------------------------------------------------
package http

import (
	"crypto/tls"
	"io"
	httplib "net/http"
	"strings"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
const (
	DEFAULT_HTTP_TIMEOUT = 30 * time.Second
)

type (
	HTTP struct {
		client  *httplib.Client
		timeout time.Duration
	}
)

var (
	AllowSelfSignedCertificates = false
)

func New() *HTTP {
	timeout := DEFAULT_HTTP_TIMEOUT
	return &HTTP{
		timeout: timeout,
		client: &httplib.Client{
			Timeout: timeout,
			Transport: &httplib.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: AllowSelfSignedCertificates,
				},
			},
		},
	}
}

func (h *HTTP) Delete(requestURL string) (string, int, error) {
	return h.Do("DELETE", requestURL, "", "")
}

func (h *HTTP) Do(method string, requestURL string, contentType string, body string) (string, int, error) {
	result := ""
	statusCode := 0
	err := error(nil)
	var response *httplib.Response
	var request *httplib.Request
	var requestBody io.Reader
	if body != "" {
		requestBody = io.NopCloser(strings.NewReader(body))
	}
	if request, err = httplib.NewRequest(method, requestURL, requestBody); err == nil {
		if contentType != "" {
			request.Header.Set("Content-Type", contentType)
		}
		if response, err = h.client.Do(request); err == nil {
			defer func(response *httplib.Response) {
				_ = response.Body.Close()
			}(response)
			statusCode = response.StatusCode
			var responseBodyBytes []byte
			if responseBodyBytes, err = io.ReadAll(response.Body); err == nil {
				result = string(responseBodyBytes)
			}
		}
	}
	return result, statusCode, err
}

func (h *HTTP) Get(requestURL string) (string, int, error) {
	return h.Do("GET", requestURL, "", "")
}

func (h *HTTP) GetTimeout() time.Duration {
	return h.timeout
}

func (h *HTTP) Post(requestURL string, contentType string, body string) (string, int, error) {
	return h.Do("POST", requestURL, contentType, body)
}

func (h *HTTP) Put(requestURL string, contentType string, body string) (string, int, error) {
	return h.Do("PUT", requestURL, contentType, body)
}

func (h *HTTP) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
	h.client.Timeout = timeout
}

func (h *HTTP) UpdateClientTransport() error {
	result := error(nil)
	h.client.Transport = &httplib.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: AllowSelfSignedCertificates,
		},
	}
	return result
}

func SetAllowSelfSignedCertificates(allow bool) {
	AllowSelfSignedCertificates = allow
}
