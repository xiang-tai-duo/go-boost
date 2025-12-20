// --------------------------------------------------------------------------------
// File:        http.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
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

// AllowSelfSignedCertificates Global flag to allow self-signed certificates
var AllowSelfSignedCertificates = false

// DefaultHTTPTimeout Default timeout for HTTP requests
const DefaultHTTPTimeout = 30 * time.Second

// HTTP provides utility methods for HTTP/HTTPS/WS/WSS operations.
type HTTP struct {
	client  *http.Client
	timeout time.Duration
}

// Delete sends an HTTP DELETE request to the specified URL.
// requestURL: URL to send the request to
// Returns: response body as string, status code, error if any occurred
// Usage:
// body, statusCode, err := httpClient.Delete("https://api.example.com/resource")
// returns response body, status code 204, nil on success
func (h *HTTP) Delete(requestURL string) (string, int, error) {
	return h.Do("DELETE", requestURL, "", "")
}

// Do sends an HTTP request with the specified method to the given URL.
// method: HTTP method to use (e.g., "GET", "POST", "PUT", "DELETE", etc.)
// requestURL: URL to send the request to
// contentType: Content-Type header value, leave empty for no Content-Type
// body: Request body as string, leave empty for no body
// Returns: response body as string, status code, error if any occurred
// Usage:
// body, statusCode, err := httpClient.Do("PATCH", "https://api.example.com/resource", "application/json", `{"key": "value"}`)
// returns response body, status code 200, nil on success
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

// Get sends an HTTP GET request to the specified URL.
// requestURL: URL to send the request to
// Returns: response body as string, status code, error if any occurred
// Usage:
// body, statusCode, err := httpClient.Get("https://api.example.com/resource")
// returns response body, status code 200, nil on success
func (h *HTTP) Get(requestURL string) (string, int, error) {
	return h.Do("GET", requestURL, "", "")
}

// GetTimeout returns the current timeout for HTTP requests.
// Returns: current timeout duration
// Usage:
// timeout := httpClient.GetTimeout()
// returns current timeout duration
func (h *HTTP) GetTimeout() time.Duration {
	return h.timeout
}

// NewHTTP creates a new HTTP instance with default configurations.
// Returns: HTTP instance
// Usage:
// httpClient := NewHTTP()
// returns new HTTP instance with default configurations
func NewHTTP() *HTTP {
	timeout := DefaultHTTPTimeout
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

// Post sends an HTTP POST request to the specified URL with the given body.
// requestURL: URL to send the request to
// contentType: Content-Type header value
// body: Request body as string
// Returns: response body as string, status code, error if any occurred
// Usage:
// body, statusCode, err := httpClient.Post("https://api.example.com/resource", "application/json", `{"key": "value"}`)
// returns response body, status code 201, nil on success
func (h *HTTP) Post(requestURL string, contentType string, body string) (string, int, error) {
	return h.Do("POST", requestURL, contentType, body)
}

// Put sends an HTTP PUT request to the specified URL with the given body.
// requestURL: URL to send the request to
// contentType: Content-Type header value
// body: Request body as string
// Returns: response body as string, status code, error if any occurred
// Usage:
// body, statusCode, err := httpClient.Put("https://api.example.com/resource", "application/json", `{"key": "value"}`)
// returns response body, status code 200, nil on success
func (h *HTTP) Put(requestURL string, contentType string, body string) (string, int, error) {
	return h.Do("PUT", requestURL, contentType, body)
}

// SetAllowSelfSignedCertificates enables or disables the allowance of self-signed certificates globally.
// allow: Boolean indicating whether to allow self-signed certificates
// Usage:
// SetAllowSelfSignedCertificates(true) // Enable self-signed certificates
// returns nothing
func SetAllowSelfSignedCertificates(allow bool) {
	AllowSelfSignedCertificates = allow
	// Update existing HTTP clients
}

// SetTimeout sets the timeout for HTTP requests.
// timeout: Timeout duration to set
// Usage:
// httpClient.SetTimeout(60 * time.Second)
// returns nothing
func (h *HTTP) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
	h.client.Timeout = timeout
}

// UpdateClientTransport updates the client transport to reflect the current self-signed certificate setting.
// Returns: error if any occurred
// Usage:
// err := httpClient.UpdateClientTransport()
// returns nil on success
func (h *HTTP) UpdateClientTransport() error {
	var err error
	h.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: AllowSelfSignedCertificates,
		},
	}
	return err
}
