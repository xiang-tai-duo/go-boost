package http2

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	httplib "net/http"
	"net/url"

	"github.com/xiang-tai-duo/go-boost/ca"
	"github.com/xiang-tai-duo/go-boost/hash"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/system"
)

//goland:noinspection GoSnakeCaseUsage
type (
	DOWNLOAD_PROGRESS struct {
		Downloaded int64
		TotalSize  int64
	}

	DOWNLOAD_PROGRESS_HANDLER func(downloaded int64, total int64)
	DOWNLOAD_RETRY_HANDLER    func(attempt int, maxRetries int, chunkStart int64, chunkEnd int64, backoff time.Duration, err error)

	HTTP_STATUS_ERROR struct {
		Method     string
		RequestURL string
		StatusCode int
	}

	HTTP struct {
		allow_self_signed_certificates bool
		client                         *httplib.Client
		timeout                        time.Duration
	}

	PARALLELS_DOWNLOAD_CHUNK struct {
		EndByte   int64
		StartByte int64
	}
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection
const (
	ACCEPT_RANGES_HEADER              = "Accept-Ranges"
	BYTES_UNIT                        = "bytes"
	CONTENT_LENGTH_HEADER             = "Content-Length"
	CONTENT_TYPE_HEADER               = "Content-Type"
	DEFAULT_DOWNLOAD_PARALLEL         = 10
	DEFAULT_DOWNLOAD_PARALLELS_SIZE   = 1024 * 1024
	DEFAULT_DOWNLOAD_RETRIES          = 10
	DEFAULT_HTTP_TIMEOUT              = 30 * time.Second
	DIRECTORY_PERMISSION              = 0755
	DOWNLOAD_RETRY_INTERVAL           = 5 * time.Second
	ENTITY_TAG_HEADER                 = "ETag"
	ERR_INCOMPLETE_CHUNK              = "incomplete chunk: expected %d bytes, got %d"
	ERR_PARALLELS_SIZE_NOT_POSITIVE   = "parallels size must be positive"
	ERR_RANGE_NOT_SUPPORTED           = "range request not supported"
	FILE_PERMISSION                   = 0644
	HTTP_STATUS_OK_MAX                = 300
	HTTP_STATUS_OK_MIN                = 200
	HTTP_STATUS_PARTIAL_CONTENT       = 206
	HTTP_TEST_PATH                    = "/"
	HTTP_TEST_URL_FORMAT              = "%s%s%s%s"
	LAST_MODIFIED_HEADER              = "Last-Modified"
	MAX_DOWNLOAD_PARALLEL             = 400
	MAX_DOWNLOAD_RETRIES              = 999
	MIN_DOWNLOAD_CHUNK_SIZE           = 10 * 1024
	METHOD_DELETE                     = "DELETE"
	METHOD_GET                        = "GET"
	METHOD_POST                       = "POST"
	METHOD_PUT                        = "PUT"
	MODULE_NAME_HTTP2                 = "http2"
	NETWORK_ERROR_CONNECTION_REFUSED  = "connection refused"
	NETWORK_ERROR_HOST_UNREACHABLE    = "host is unreachable"
	NETWORK_ERROR_NO_ROUTE            = "no route to host"
	NETWORK_ERROR_NETWORK_UNREACHABLE = "network is unreachable"
	RANGE_HEADER                      = "Range"
	RANGE_PREFIX                      = "bytes="
	RANGE_PROBE_VALUE                 = "0-0"
	RANGE_SEPARATOR                   = "-"
	REQUEST_HASH_HEADER_NAME          = "X-Request-Hash"
	REQUEST_TIMESTAMP_HEADER_NAME     = "X-Request-Timestamp"
	REQUEST_UUID_HEADER_NAME          = "X-Request-UUID"
	SALT_ENDPOINT_PATH                = "/api/salt"
	SCHEME_HTTP                       = "http"
	SCHEME_HTTPS                      = "https"
	SCHEME_SEPARATOR                  = "://"
	TLS_ALERT_CERTIFICATE_REQUIRED    = "certificate required"
	UUID_BYTE_COUNT                   = 16
	UUID_FORMAT                       = "%x-%x-%x-%x-%x"
	UUID_VARIANT_MASK                 = 0x3f
	UUID_VARIANT_SET                  = 0x80
	UUID_VERSION_MASK                 = 0x0f
	UUID_VERSION_SET                  = 0x40
)

var (
	defaultAllowSelfSignedCertificates bool
	defaultClientCertificate           *tls.Certificate
	defaultRootCertificate             *x509.CertPool
	defaultTransportMutex              sync.RWMutex
	serverSaltCache                    map[string]string
	serverSaltFetched                  map[string]bool
	serverSaltMutex                    sync.Mutex
	systemProxyEnabled                 bool
	systemProxyURL                     *url.URL
)

func init() {
	defaultAllowSelfSignedCertificates = false
	defaultClientCertificate = loadClientCertificate()
	defaultRootCertificate = loadRootCertificate()
	serverSaltCache = make(map[string]string)
	serverSaltFetched = make(map[string]bool)
	systemProxyEnabled = false
	systemProxyURL = nil
}

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_HTTP2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_HTTP2, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_HTTP2, logger.SKIP_STACK_FRAMES_BASE)
}

func New() *HTTP {
	result := &HTTP{
		timeout:                        DEFAULT_HTTP_TIMEOUT,
		allow_self_signed_certificates: GetDefaultAllowSelfSignedCertificates(),
	}
	result.client = &httplib.Client{
		Timeout:   result.timeout,
		Transport: createTransport(result.allow_self_signed_certificates),
	}
	return result
}

func (h *HTTP) Delete(requestURL string) (string, int, error) {
	return h.Invoke(METHOD_DELETE, requestURL, "", "", nil)
}

//goland:noinspection GoUnhandledErrorResult
func (h *HTTP) Download(url string, filePath string) (*DOWNLOAD_PROGRESS, error) {
	__info(fmt.Sprintf("[Download] Starting SDK download: %s -> %s", url, filePath))
	result := &DOWNLOAD_PROGRESS{}
	err := error(nil)
	resp := (*httplib.Response)(nil)
	if resp, err = httplib.Get(url); err == nil {
		defer resp.Body.Close()
		body := make([]byte, 0)
		if body, err = io.ReadAll(resp.Body); err == nil {
			if err = os.WriteFile(filePath, body, FILE_PERMISSION); err == nil {
				result.Downloaded = int64(len(body))
				result.TotalSize = int64(len(body))
				__info(fmt.Sprintf("[Download] SDK download completed: %s -> %s, %d bytes", url, filePath, result.Downloaded))
			}
		}
	}
	if err != nil {
		__debug(fmt.Sprintf("[Download] SDK download failed: %v", err))
	}
	return result, err
}

func (h *HTTP) DownloadParallels(url string, filePath string) error {
	err := error(nil)
	err = h.DownloadParallelsEx(url, filePath, DEFAULT_DOWNLOAD_PARALLELS_SIZE, DEFAULT_DOWNLOAD_PARALLEL, nil, nil, MAX_DOWNLOAD_RETRIES, nil)
	return err
}

func (h *HTTP) DownloadParallelsEx(url string, filePath string, chunksSize int64, parallelsCount int, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string, maxDownloadRetries int, retryHandler DOWNLOAD_RETRY_HANDLER) error {
	__debug(fmt.Sprintf("[DownloadParallelsEx] Input parameters: url=%s, filePath=%s, chunksSize=%d, parallel=%d, hasProgressHandler=%v, hasHeaders=%v, maxDownloadRetries=%d, hasRetryHandler=%v, timeout=%v", url, filePath, chunksSize, parallelsCount, progressHandler != nil, len(headers) > 0, maxDownloadRetries, retryHandler != nil, h.timeout))
	err := error(nil)
	if maxDownloadRetries <= 0 {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Invalid maxDownloadRetries=%d, fallback to %d", maxDownloadRetries, MAX_DOWNLOAD_RETRIES))
		maxDownloadRetries = MAX_DOWNLOAD_RETRIES
	} else if maxDownloadRetries > MAX_DOWNLOAD_RETRIES {
		__debug(fmt.Sprintf("[DownloadParallelsEx] maxDownloadRetries=%d exceeds maximum=%d, clamped to maximum", maxDownloadRetries, MAX_DOWNLOAD_RETRIES))
		maxDownloadRetries = MAX_DOWNLOAD_RETRIES
	}
	if chunksSize > 0 {
		if chunksSize < MIN_DOWNLOAD_CHUNK_SIZE {
			__debug(fmt.Sprintf("[DownloadParallelsEx] chunksSize=%d is smaller than minimum=%d, using minimum chunk size", chunksSize, MIN_DOWNLOAD_CHUNK_SIZE))
			chunksSize = MIN_DOWNLOAD_CHUNK_SIZE
		}
		if parallelsCount <= 0 {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Invalid parallel=%d, fallback to 1", parallelsCount))
			parallelsCount = 1
		}
		if parallelsCount > MAX_DOWNLOAD_PARALLEL {
			__debug(fmt.Sprintf("[DownloadParallelsEx] parallel=%d is greater than maximum=%d, using maximum parallel", parallelsCount, MAX_DOWNLOAD_PARALLEL))
			parallelsCount = MAX_DOWNLOAD_PARALLEL
		}
		__debug(fmt.Sprintf("[DownloadParallelsEx] Getting remote file info: %s", url))
		fileSize := int64(0)
		supportRange := false
		if fileSize, supportRange, _, _, _, err = h.getURLFileInfo(url, headers); err == nil {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Remote file info: fileSize=%d, supportRange=%v", fileSize, supportRange))
			if fileSize > 0 && supportRange {
				__debug("[DownloadParallelsEx] Content-Length available and range supported, switching to chunked downloader")
				err = h.downloadParallelsWithRange(url, filePath, fileSize, chunksSize, parallelsCount, progressHandler, headers, maxDownloadRetries, retryHandler)
			} else if fileSize > 0 {
				__info(fmt.Sprintf("[DownloadParallelsEx] Content-Length=%d available but server does not support range, falling back to single-thread download", fileSize))
				err = h.downloadWithoutRange(url, filePath, progressHandler, headers, maxDownloadRetries, retryHandler)
			} else {
				__info("[DownloadParallelsEx] Content-Length unavailable, falling back to single-thread download from start to end")
				err = h.downloadWithoutRange(url, filePath, progressHandler, headers, maxDownloadRetries, retryHandler)
			}
		} else {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to get file info: %v", err))
		}
	} else {
		err = errors.New(ERR_PARALLELS_SIZE_NOT_POSITIVE)
		__debug(fmt.Sprintf("[DownloadParallelsEx] %v", err))
	}
	if err == nil {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Finished successfully: %s", filePath))
	} else {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Finished with error: %v", err))
	}
	return err
}

func (err *HTTP_STATUS_ERROR) Error() string {
	result := ""
	if err == nil {
		return result
	}
	statusText := httplib.StatusText(err.StatusCode)
	if statusText == "" {
		result = fmt.Sprintf("%s %s %d", err.Method, err.RequestURL, err.StatusCode)
	} else {
		result = fmt.Sprintf("%s %s %d (%s)", err.Method, err.RequestURL, err.StatusCode, statusText)
	}
	return result
}

func (h *HTTP) Get(requestURL string) (string, int, error) {
	return h.Invoke(METHOD_GET, requestURL, "", "", nil)
}

func (h *HTTP) GetAllowSelfSignedCertificates() bool {
	return h.allow_self_signed_certificates
}

func (h *HTTP) GetClient() *httplib.Client {
	return h.client
}

func GetDefaultAllowSelfSignedCertificates() bool {
	defaultTransportMutex.RLock()
	defer defaultTransportMutex.RUnlock()
	return defaultAllowSelfSignedCertificates
}

func GetHTTPStatusCode(err error) (int, bool) {
	result := 0
	ok := false
	statusError := (*HTTP_STATUS_ERROR)(nil)
	if errors.As(err, &statusError) {
		result = statusError.StatusCode
		ok = true
	}
	return result, ok
}

func (h *HTTP) GetTimeout() time.Duration {
	return h.timeout
}

//goland:noinspection DuplicatedCode
func (h *HTTP) Invoke(method string, requestURL string, contentType string, body string, headers map[string]string) (string, int, error) {
	__debug(fmt.Sprintf("[HTTP] %s %s", method, requestURL))
	result := ""
	statusCode := 0
	err := error(nil)
	serverSalt := ""
	var parsedRequestURL *url.URL
	if parsedRequestURL, err = url.Parse(requestURL); err == nil && strings.EqualFold(parsedRequestURL.Scheme, SCHEME_HTTPS) {
		serverSalt, err = h.getServerSalt(requestURL)
	}
	if err == nil {
		var requestBody io.Reader
		if body != "" {
			requestBody = io.NopCloser(strings.NewReader(body))
		}
		var request *httplib.Request
		if request, err = httplib.NewRequest(method, requestURL, requestBody); err == nil {
			addDefaultHeaders(request)
			if contentType != "" {
				request.Header.Set(CONTENT_TYPE_HEADER, contentType)
			}
			for key, value := range headers {
				request.Header.Set(key, value)
			}
			request.Header.Set(REQUEST_HASH_HEADER_NAME, hash.SHA3(body+serverSalt))
			var response *httplib.Response
			if response, err = h.client.Do(request); err == nil {
				defer func(response *httplib.Response) {
					_ = response.Body.Close()
				}(response)
				statusCode = response.StatusCode
				__debug(fmt.Sprintf("[HTTP] %s %s -> Status: %d", method, requestURL, statusCode))
				responseBodyBytes := make([]byte, 0)
				if responseBodyBytes, err = io.ReadAll(response.Body); err == nil {
					result = string(responseBodyBytes)
				} else {
					__debug(fmt.Sprintf("[HTTP] Failed to read response body: %v", err))
				}
			} else {
				__debug(fmt.Sprintf("[HTTP] Request failed: %v", err))
			}
		} else {
			__debug(fmt.Sprintf("[HTTP] Failed to create request: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("[HTTP] Failed to get server salt: %v", err))
	}
	return result, statusCode, err
}

func IsHTTPProtocol(ip string, port int) error {
	return testProtocol(SCHEME_HTTP, ip, port)
}

func IsHTTPSProtocol(ip string, port int) error {
	return testProtocol(SCHEME_HTTPS, ip, port)
}

func (h *HTTP) Post(requestURL string, contentType string, body string) (string, int, error) {
	return h.Invoke(METHOD_POST, requestURL, contentType, body, nil)
}

func (h *HTTP) Put(requestURL string, contentType string, body string) (string, int, error) {
	return h.Invoke(METHOD_PUT, requestURL, contentType, body, nil)
}

func (h *HTTP) SetAllowSelfSignedCertificates(allow bool) {
	h.allow_self_signed_certificates = allow
	h.client.Transport = createTransport(h.allow_self_signed_certificates)
}

//goland:noinspection GoUnusedExportedFunction
func SetDefaultAllowSelfSignedCertificates(allow bool) {
	defaultTransportMutex.Lock()
	defaultAllowSelfSignedCertificates = allow
	defaultTransportMutex.Unlock()
}

func EnableSystemProxy() {
	proxyURL, err := loadSystemProxy()
	defaultTransportMutex.Lock()
	systemProxyEnabled = err == nil && proxyURL != nil
	if systemProxyEnabled {
		systemProxyURL = proxyURL
	} else {
		systemProxyURL = nil
	}
	defaultTransportMutex.Unlock()
}

func DisableSystemProxy() {
	defaultTransportMutex.Lock()
	systemProxyEnabled = false
	systemProxyURL = nil
	defaultTransportMutex.Unlock()
}

func IsSystemProxyEnabled() bool {
	defaultTransportMutex.RLock()
	defer defaultTransportMutex.RUnlock()
	return systemProxyEnabled
}

func (h *HTTP) SetTimeout(t time.Duration) {
	h.timeout = t
	h.client.Timeout = t
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_HTTP2, logger.SKIP_STACK_FRAMES_BASE)
}

func addDefaultHeaders(request *httplib.Request) {
	if request != nil {
		request.Header.Set(REQUEST_UUID_HEADER_NAME, newRequestUUID())
		request.Header.Set(REQUEST_TIMESTAMP_HEADER_NAME, strconv.FormatInt(time.Now().Unix(), 10))
	}
}

func applyCustomHeaders(request *httplib.Request, headers map[string]string) {
	if request != nil {
		for key, value := range headers {
			if strings.TrimSpace(key) != "" {
				request.Header.Set(key, value)
			}
		}
	}
}

func createTransport(allowSelfSignedCertificates bool) httplib.RoundTripper {
	defaultTransportMutex.RLock()
	clientCertificate := defaultClientCertificate
	rootCAs := defaultRootCertificate
	proxyEnabled := systemProxyEnabled
	proxyAddr := systemProxyURL
	defaultTransportMutex.RUnlock()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: allowSelfSignedCertificates,
		RootCAs:            rootCAs,
	}
	if clientCertificate != nil {
		tlsConfig.Certificates = []tls.Certificate{*clientCertificate}
	}
	result := &httplib.Transport{
		TLSClientConfig: tlsConfig,
	}
	if proxyEnabled && proxyAddr != nil {
		result.Proxy = httplib.ProxyURL(proxyAddr)
	} else if system.IsUnix() {
		result.Proxy = httplib.ProxyFromEnvironment
	}
	return result
}

func (h *HTTP) downloadParallelsChunkOnce(ctx context.Context, requestURL string, file *os.File, chunk PARALLELS_DOWNLOAD_CHUNK, progress *int64, totalSize int64, mutex *sync.Mutex, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string) error {
	chunkSize := chunk.EndByte - chunk.StartByte + 1
	__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk request start: range=%d-%d, chunkSize=%d", chunk.StartByte, chunk.EndByte, chunkSize))
	err := error(nil)
	var request *httplib.Request
	if request, err = httplib.NewRequest(METHOD_GET, requestURL, nil); err == nil {
		chunkCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()
		request = request.WithContext(chunkCtx)
		request.Close = true
		addDefaultHeaders(request)
		applyCustomHeaders(request, headers)
		rangeValue := RANGE_PREFIX + strconv.FormatInt(chunk.StartByte, 10) + RANGE_SEPARATOR + strconv.FormatInt(chunk.EndByte, 10)
		request.Header.Set(RANGE_HEADER, rangeValue)
		__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk HTTP request prepared: method=%s, url=%s, range=%s, close=%v", request.Method, request.URL.String(), rangeValue, request.Close))
		var response *httplib.Response
		__debug(fmt.Sprintf("[DownloadParallelsEx] Sending chunk HTTP request: range=%s", rangeValue))
		if response, err = h.client.Do(request); err == nil {
			defer func() {
				__debug(fmt.Sprintf("[DownloadParallelsEx] Closing response body: range=%s", rangeValue))
				_ = response.Body.Close()
			}()
			__debug(fmt.Sprintf("[DownloadParallelsEx] Range %d-%d response status: %d, contentLength=%d", chunk.StartByte, chunk.EndByte, response.StatusCode, response.ContentLength))
			if response.StatusCode == HTTP_STATUS_PARTIAL_CONTENT {
				__debug(fmt.Sprintf("[DownloadParallelsEx] Writing chunk response: range=%s", rangeValue))
				err = writeParallelsDownloadResponse(file, response.Body, chunk, progress, totalSize, mutex, progressHandler)
			} else {
				err = newHTTPStatusError(METHOD_GET, requestURL, response.StatusCode)
				__debug(fmt.Sprintf("[DownloadParallelsEx] Expected partial content (206), got %d: %v", response.StatusCode, err))
			}
		} else {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Request failed: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to create request: %v", err))
	}
	if err == nil {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk request finished: range=%d-%d", chunk.StartByte, chunk.EndByte))
	}
	return err
}

func (h *HTTP) downloadParallelsChunkWithRetry(ctx context.Context, requestURL string, file *os.File, chunk PARALLELS_DOWNLOAD_CHUNK, progress *int64, totalSize int64, mutex *sync.Mutex, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string, maxDownloadRetries int, retryHandler DOWNLOAD_RETRY_HANDLER) error {
	err := error(nil)
	chunkSize := chunk.EndByte - chunk.StartByte + 1
	__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk retry loop start: range=%d-%d, chunkSize=%d, maxRetries=%d", chunk.StartByte, chunk.EndByte, chunkSize, maxDownloadRetries))
	lastErr := error(nil)
	actualAttempts := 0
	for attempt := 0; attempt < maxDownloadRetries; attempt++ {
		if ctx.Err() != nil {
			err = ctx.Err()
			break
		}
		actualAttempts++
		__debug(fmt.Sprintf("[DownloadParallelsEx] Downloading range %d-%d, chunkSize=%d bytes (attempt %d/%d)", chunk.StartByte, chunk.EndByte, chunkSize, attempt+1, maxDownloadRetries))
		if err = h.downloadParallelsChunkOnce(ctx, requestURL, file, chunk, progress, totalSize, mutex, progressHandler, headers); err == nil {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk attempt succeeded: range=%d-%d, attempt=%d", chunk.StartByte, chunk.EndByte, attempt+1))
			break
		}
		lastErr = err
		__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk attempt failed: range=%d-%d, attempt=%d/%d, err=%v", chunk.StartByte, chunk.EndByte, actualAttempts, maxDownloadRetries, err))
		if ctx.Err() != nil {
			err = ctx.Err()
			__debug(fmt.Sprintf("[DownloadParallelsEx] Context cancelled after attempt %d: %v", actualAttempts, err))
			break
		}
		if errors.Is(err, context.Canceled) || isNonRetryableHTTPError(err) {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Non-retryable error: %v", err))
			break
		}
		if attempt < maxDownloadRetries-1 {
			backoff := DOWNLOAD_RETRY_INTERVAL
			if attempt > 0 {
				backoff = time.Duration(attempt+1) * DOWNLOAD_RETRY_INTERVAL
				if backoff > 30*time.Second {
					backoff = 30 * time.Second
				}
			}
			__warning(fmt.Sprintf("[DownloadParallelsEx] Attempt %d/%d failed, retrying in %v: %v", actualAttempts, maxDownloadRetries, backoff, err))
			__debug(fmt.Sprintf("[DownloadParallelsEx] Sleeping before retry: interval=%v", backoff))
			if retryHandler != nil {
				retryHandler(actualAttempts, maxDownloadRetries, chunk.StartByte, chunk.EndByte, backoff, err)
			}
			select {
			case <-ctx.Done():
				if err == nil {
					err = ctx.Err()
				}
				attempt = maxDownloadRetries
			case <-time.After(backoff):
			}
		}
	}
	if err == nil && lastErr != nil {
		err = fmt.Errorf("range %d-%d failed after %d/%d attempts: %w", chunk.StartByte, chunk.EndByte, actualAttempts, maxDownloadRetries, lastErr)
	}
	if err != nil {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Range %d-%d failed after %d/%d attempts: %v", chunk.StartByte, chunk.EndByte, actualAttempts, maxDownloadRetries, err))
	} else {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk retry loop finished successfully: range=%d-%d, attempts=%d", chunk.StartByte, chunk.EndByte, actualAttempts))
	}
	return err
}

func (h *HTTP) downloadParallelsChunks(requestURL string, file *os.File, chunks []PARALLELS_DOWNLOAD_CHUNK, totalSize int64, parallel int, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string, maxDownloadRetries int, retryHandler DOWNLOAD_RETRY_HANDLER) error {
	err := error(nil)
	__debug(fmt.Sprintf("[DownloadParallelsEx] Chunk scheduler start: chunkCount=%d, totalSize=%d, parallel=%d, hasProgressHandler=%v", len(chunks), totalSize, parallel, progressHandler != nil))
	ctx := context.Background()
	mutex := &sync.Mutex{}
	downloaded := int64(0)
	var firstErr error
	var errOnce sync.Once
	semaphore := make(chan struct{}, parallel)
	var wg sync.WaitGroup
	for i := range chunks {
		chunk := chunks[i]
		chunkIndex := i
		select {
		case semaphore <- struct{}{}:
		}
		wg.Add(1)
		__debug(fmt.Sprintf("[DownloadParallelsEx] Scheduling chunk: index=%d/%d, range=%d-%d", chunkIndex+1, len(chunks), chunk.StartByte, chunk.EndByte))
		go func(chunk PARALLELS_DOWNLOAD_CHUNK, idx int) {
			defer wg.Done()
			defer func() {
				<-semaphore
				__debug(fmt.Sprintf("[DownloadParallelsEx] Worker slot released: index=%d, range=%d-%d", idx, chunk.StartByte, chunk.EndByte))
			}()
			defer func() {
				if r := recover(); r != nil {
					panicErr := fmt.Errorf("chunk %d-%d panic: %v", chunk.StartByte, chunk.EndByte, r)
					__debug(fmt.Sprintf("[DownloadParallelsEx] %v", panicErr))
					errOnce.Do(func() { firstErr = panicErr })
				}
			}()
			__debug(fmt.Sprintf("[DownloadParallelsEx] Worker started: index=%d, range=%d-%d", idx, chunk.StartByte, chunk.EndByte))
			chunkErr := h.downloadParallelsChunkWithRetry(ctx, requestURL, file, chunk, &downloaded, totalSize, mutex, progressHandler, headers, maxDownloadRetries, retryHandler)
			if chunkErr != nil {
				__debug(fmt.Sprintf("[DownloadParallelsEx] Worker failed: index=%d, range=%d-%d, err=%v", idx, chunk.StartByte, chunk.EndByte, chunkErr))
				errOnce.Do(func() { firstErr = chunkErr })
			} else {
				__debug(fmt.Sprintf("[DownloadParallelsEx] Worker completed: index=%d, range=%d-%d", idx, chunk.StartByte, chunk.EndByte))
			}
		}(chunk, chunkIndex)
	}
	__debug("[DownloadParallelsEx] Waiting for all chunk workers")
	wg.Wait()
	__debug(fmt.Sprintf("[DownloadParallelsEx] All chunk workers finished: downloaded=%d/%d", downloaded, totalSize))
	if firstErr != nil {
		err = firstErr
		__debug(fmt.Sprintf("[DownloadParallelsEx] Scheduler captured error: %v", err))
	} else {
		__debug("[DownloadParallelsEx] Scheduler finished without chunk error")
	}
	return err
}

func (h *HTTP) downloadParallelsWithRange(requestURL string, filePath string, fileSize int64, parallelsSize int64, parallel int, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string, maxDownloadRetries int, retryHandler DOWNLOAD_RETRY_HANDLER) error {
	__info(fmt.Sprintf("[DownloadParallelsEx] Using ranged download, fileSize=%d", fileSize))
	__debug(fmt.Sprintf("[DownloadParallelsEx] Preparing ranged download: requestURL=%s, filePath=%s, fileSize=%d, chunksSize=%d, parallel=%d", requestURL, filePath, fileSize, parallelsSize, parallel))
	err := ensureDownloadDirectory(filePath)
	temporaryFilePath := getDownloadTemporaryFilePath(filePath)
	file := (*os.File)(nil)
	if err == nil {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Download directory ready for: %s", filePath))
		file, err = openTruncatedDownloadFile(temporaryFilePath, fileSize)
	}
	if err == nil {
		defer func() {
			if file != nil {
				__debug(fmt.Sprintf("[DownloadParallelsEx] Closing target file: %s", temporaryFilePath))
				_ = file.Close()
			}
		}()
		if progressHandler != nil {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Sending initial progress: downloaded=0, total=%d", fileSize))
			progressHandler(0, fileSize)
		}
		chunks := newParallelsDownloadChunks(fileSize, parallelsSize)
		__debug(fmt.Sprintf("[DownloadParallelsEx] Created %d chunks", len(chunks)))
		err = h.downloadParallelsChunks(requestURL, file, chunks, fileSize, parallel, progressHandler, headers, maxDownloadRetries, retryHandler)
		if err == nil {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Syncing target file: %s", temporaryFilePath))
			if syncErr := file.Sync(); syncErr != nil {
				err = syncErr
				__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to sync target file: %v", syncErr))
			} else if closeErr := file.Close(); closeErr != nil {
				file = nil
				err = closeErr
				__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to close target file: %v", closeErr))
			} else {
				file = nil
				if renameErr := replaceDownloadTargetFile(temporaryFilePath, filePath); renameErr != nil {
					err = renameErr
					__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to replace target file: %v", renameErr))
				} else {
					__info(fmt.Sprintf("[DownloadParallelsEx] Completed successfully: %s", filePath))
				}
			}
		}
	}
	if err != nil {
		removeDownloadTemporaryFile(temporaryFilePath)
		__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to prepare or complete ranged download: %v", err))
	}
	return err
}

func (h *HTTP) downloadWithoutRange(requestURL string, filePath string, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string, maxDownloadRetries int, retryHandler DOWNLOAD_RETRY_HANDLER) error {
	__info(fmt.Sprintf("[DownloadWithoutRange] Starting single-thread download: %s -> %s, maxRetries=%d", requestURL, filePath, maxDownloadRetries))
	err := ensureDownloadDirectory(filePath)
	if err == nil {
		for attempt := 0; attempt < maxDownloadRetries; attempt++ {
			__debug(fmt.Sprintf("[DownloadWithoutRange] Attempt %d/%d", attempt+1, maxDownloadRetries))
			if progressHandler != nil {
				progressHandler(0, 0)
			}
			err = h.downloadWithoutRangeOnce(requestURL, filePath, progressHandler, headers)
			if err == nil {
				__info(fmt.Sprintf("[DownloadWithoutRange] Completed successfully: %s", filePath))
				break
			}
			__debug(fmt.Sprintf("[DownloadWithoutRange] Attempt %d/%d failed: %v", attempt+1, maxDownloadRetries, err))
			if errors.Is(err, context.Canceled) || isNonRetryableHTTPError(err) {
				__debug(fmt.Sprintf("[DownloadWithoutRange] Non-retryable error, aborting: %v", err))
				break
			}
			if attempt < maxDownloadRetries-1 {
				__warning(fmt.Sprintf("[DownloadWithoutRange] Attempt %d failed, retrying after %v: %v", attempt+1, DOWNLOAD_RETRY_INTERVAL, err))
				if retryHandler != nil {
					retryHandler(attempt+1, maxDownloadRetries, 0, 0, DOWNLOAD_RETRY_INTERVAL, err)
				}
				time.Sleep(DOWNLOAD_RETRY_INTERVAL)
			}
		}
	}
	if err != nil {
		__debug(fmt.Sprintf("[DownloadWithoutRange] Failed after retries: %v", err))
	}
	return err
}

func (h *HTTP) downloadWithoutRangeOnce(requestURL string, filePath string, progressHandler DOWNLOAD_PROGRESS_HANDLER, headers map[string]string) (err error) {
	file := (*os.File)(nil)
	if file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FILE_PERMISSION); err == nil {
		defer func() {
			__debug(fmt.Sprintf("[DownloadWithoutRange] Closing target file: %s", filePath))
			_ = file.Close()
		}()
		request := (*httplib.Request)(nil)
		if request, err = httplib.NewRequest(METHOD_GET, requestURL, nil); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			request = request.WithContext(ctx)
			addDefaultHeaders(request)
			applyCustomHeaders(request, headers)
			__debug(fmt.Sprintf("[DownloadWithoutRange] Sending GET request: %s", requestURL))
			response := (*httplib.Response)(nil)
			if response, err = h.client.Do(request); err == nil {
				defer func() {
					__debug("[DownloadWithoutRange] Closing response body")
					_ = response.Body.Close()
				}()
				__debug(fmt.Sprintf("[DownloadWithoutRange] Response status: %d, contentLength=%d", response.StatusCode, response.ContentLength))
				if response.StatusCode >= HTTP_STATUS_OK_MIN && response.StatusCode < HTTP_STATUS_OK_MAX {
					totalSize := response.ContentLength
					downloaded := int64(0)
					buffer := make([]byte, DEFAULT_DOWNLOAD_PARALLELS_SIZE)
					lastProgress := time.Now()
					for {
						n, readErr := response.Body.Read(buffer)
						if n > 0 {
							if _, writeErr := file.Write(buffer[:n]); writeErr == nil {
								downloaded += int64(n)
								if progressHandler != nil && time.Since(lastProgress) >= time.Second {
									progressHandler(downloaded, totalSize)
									lastProgress = time.Now()
								}
							} else {
								err = writeErr
								__debug(fmt.Sprintf("[DownloadWithoutRange] Write failed: %v", writeErr))
								break
							}
						}
						if readErr == io.EOF {
							err = nil
							break
						} else if readErr != nil {
							err = readErr
							__debug(fmt.Sprintf("[DownloadWithoutRange] Read failed: %v", readErr))
							break
						}
					}
					if err == nil {
						if syncErr := file.Sync(); syncErr != nil {
							__warning(fmt.Sprintf("[DownloadWithoutRange] Sync failed: %v", syncErr))
						}
						if progressHandler != nil {
							progressHandler(downloaded, totalSize)
						}
						__debug(fmt.Sprintf("[DownloadWithoutRange] Downloaded %d bytes, expected=%d", downloaded, totalSize))
					}
				} else {
					err = newHTTPStatusError(METHOD_GET, requestURL, response.StatusCode)
					__debug(fmt.Sprintf("[DownloadWithoutRange] Unexpected status: %v", err))
				}
			} else {
				__debug(fmt.Sprintf("[DownloadWithoutRange] Request failed: %v", err))
			}
		} else {
			__debug(fmt.Sprintf("[DownloadWithoutRange] Failed to create request: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("[DownloadWithoutRange] Failed to open target file: %v", err))
	}
	return err
}

func ensureDownloadDirectory(filePath string) error {
	directory := filepath.Dir(filePath)
	err := error(nil)
	__debug(fmt.Sprintf("[DownloadParallelsEx] Checking download directory: %s", directory))
	if _, statErr := os.Stat(directory); os.IsNotExist(statErr) {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Creating directory: %s", directory))
		if err = os.MkdirAll(directory, DIRECTORY_PERMISSION); err != nil {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to create directory: %v", err))
		} else {
			__debug(fmt.Sprintf("[DownloadParallelsEx] Directory created: %s", directory))
		}
	} else if statErr != nil {
		err = statErr
		__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to stat directory: %v", statErr))
	} else {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Directory already exists: %s", directory))
	}
	return err
}

func findCertificateFile(fileName string) string {
	result := ""
	logger.Logger.Debug(fmt.Sprintf("[Certificate] Looking for certificate file: %s", fileName))
	var executableFilePath string
	var executableError error
	if executableFilePath, executableError = os.Executable(); executableError == nil {
		executableDirectory := filepath.Dir(executableFilePath)
		candidate := filepath.Join(executableDirectory, fileName)
		logger.Logger.Debug(fmt.Sprintf("[Certificate] Checking executable directory: %s", candidate))
		var statusError error
		if _, statusError = os.Stat(candidate); statusError == nil {
			result = candidate
			logger.Logger.Debug(fmt.Sprintf("[Certificate] Found in executable directory: %s", candidate))
		}
	}
	if result == "" {
		var currentWorkingDirectory string
		var workingDirectoryError error
		if currentWorkingDirectory, workingDirectoryError = os.Getwd(); workingDirectoryError == nil {
			candidate := filepath.Join(currentWorkingDirectory, fileName)
			logger.Logger.Debug(fmt.Sprintf("[Certificate] Checking working directory: %s", candidate))
			var statusError error
			if _, statusError = os.Stat(candidate); statusError == nil {
				result = candidate
				logger.Logger.Debug(fmt.Sprintf("[Certificate] Found in working directory: %s", candidate))
			}
		}
	}
	if result == "" {
		logger.Logger.Debug(fmt.Sprintf("[Certificate] File not found: %s", fileName))
	}
	return result
}

func getDownloadTemporaryFilePath(filePath string) string {
	return filePath + ".part"
}

func (h *HTTP) getServerSalt(requestURL string) (string, error) {
	__debug(fmt.Sprintf("[Salt] Getting server salt for: %s", requestURL))
	result := ""
	err := error(nil)
	var parsedURL *url.URL
	if parsedURL, err = url.Parse(requestURL); err == nil && parsedURL.Host != "" {
		if strings.EqualFold(parsedURL.Scheme, SCHEME_HTTPS) {
			cacheKey := parsedURL.Scheme + SCHEME_SEPARATOR + parsedURL.Host
			serverSaltMutex.Lock()
			if serverSaltFetched[cacheKey] {
				result = serverSaltCache[cacheKey]
				__debug(fmt.Sprintf("[Salt] Using cached salt for: %s", cacheKey))
				serverSaltMutex.Unlock()
			} else {
				serverSaltFetched[cacheKey] = true
				serverSaltMutex.Unlock()
				saltURL := cacheKey + SALT_ENDPOINT_PATH
				__debug(fmt.Sprintf("[Salt] Fetching salt from: %s", saltURL))
				var saltRequest *httplib.Request
				if saltRequest, err = httplib.NewRequest(METHOD_POST, saltURL, nil); err == nil {
					addDefaultHeaders(saltRequest)
					var saltResponse *httplib.Response
					if saltResponse, err = h.client.Do(saltRequest); err == nil {
						defer func(response *httplib.Response) {
							_ = response.Body.Close()
						}(saltResponse)
						__debug(fmt.Sprintf("[Salt] Response status: %d", saltResponse.StatusCode))
						if saltResponse.StatusCode >= HTTP_STATUS_OK_MIN && saltResponse.StatusCode < HTTP_STATUS_OK_MAX {
							responseBodyBytes := make([]byte, 0)
							if responseBodyBytes, err = io.ReadAll(saltResponse.Body); err == nil {
								payload := struct {
									Hash string `json:"hash"`
								}{}
								if err = json.Unmarshal(responseBodyBytes, &payload); err == nil {
									result = payload.Hash
									serverSaltMutex.Lock()
									serverSaltCache[cacheKey] = result
									serverSaltMutex.Unlock()
									__debug("[Salt] Server salt fetched successfully")
								} else {
									__debug(fmt.Sprintf("[Salt] Failed to parse response: %v", err))
								}
							} else {
								__debug(fmt.Sprintf("[Salt] Failed to read response: %v", err))
							}
						} else {
							__warning(fmt.Sprintf("[Salt] Unexpected status: %d", saltResponse.StatusCode))
						}
					} else {
						__debug(fmt.Sprintf("[Salt] Request failed: %v", err))
					}
				} else {
					__debug(fmt.Sprintf("[Salt] Failed to create request: %v", err))
				}
			}
		} else {
			err = fmt.Errorf("salt endpoint requires https scheme, got: %s", parsedURL.Scheme)
			__warning(fmt.Sprintf("[Salt] %v", err))
		}
	} else {
		__debug(fmt.Sprintf("[Salt] Failed to parse URL: %v", err))
	}
	return result, err
}

//goland:noinspection SpellCheckingInspection
func (h *HTTP) getURLFileInfo(requestURL string, headers map[string]string) (int64, bool, string, string, string, error) {
	__debug(fmt.Sprintf("[FileInfo] Getting file info: %s", requestURL))
	size := int64(0)
	supportRange := false
	entityTag := ""
	lastModified := ""
	contentType := ""
	err := error(nil)
	var headRequest *httplib.Request
	var response *httplib.Response
	for headAttempt := 0; headAttempt < MAX_DOWNLOAD_RETRIES; headAttempt++ {
		headRequest = nil
		response = nil
		__debug(fmt.Sprintf("[FileInfo] Creating HEAD request (attempt %d/%d): %s", headAttempt+1, MAX_DOWNLOAD_RETRIES, requestURL))
		if headRequest, err = httplib.NewRequest(httplib.MethodHead, requestURL, nil); err == nil {
			addDefaultHeaders(headRequest)
			applyCustomHeaders(headRequest, headers)
			__debug(fmt.Sprintf("[FileInfo] Sending HEAD request (attempt %d/%d): %s", headAttempt+1, MAX_DOWNLOAD_RETRIES, requestURL))
			response, err = h.client.Do(headRequest)
			if err == nil {
				break
			}
			__debug(fmt.Sprintf("[FileInfo] HEAD request attempt %d/%d failed: %v", headAttempt+1, MAX_DOWNLOAD_RETRIES, err))
			if isNonRetryableHTTPError(err) {
				__debug(fmt.Sprintf("[FileInfo] HEAD request non-retryable error, aborting: %v", err))
				break
			}
		} else {
			__debug(fmt.Sprintf("[FileInfo] Failed to create HEAD request: %v", err))
			break
		}
		if headAttempt < MAX_DOWNLOAD_RETRIES-1 {
			__debug(fmt.Sprintf("[FileInfo] Waiting %v before HEAD retry", DOWNLOAD_RETRY_INTERVAL))
			time.Sleep(DOWNLOAD_RETRY_INTERVAL)
		}
	}
	if err == nil && response != nil {
		defer func(response *httplib.Response) {
			__debug("[FileInfo] Closing HEAD response body")
			_ = response.Body.Close()
		}(response)
		__debug(fmt.Sprintf("[FileInfo] HEAD response status: %d", response.StatusCode))
		if response.StatusCode >= HTTP_STATUS_OK_MIN && response.StatusCode < HTTP_STATUS_OK_MAX {
			contentLength := response.Header.Get(CONTENT_LENGTH_HEADER)
			if contentLength != "" {
				size, _ = strconv.ParseInt(contentLength, 10, 64)
			}
			acceptRanges := response.Header.Get(ACCEPT_RANGES_HEADER)
			if acceptRanges == BYTES_UNIT {
				supportRange = true
			}
			entityTag = response.Header.Get(ENTITY_TAG_HEADER)
			lastModified = response.Header.Get(LAST_MODIFIED_HEADER)
			contentType = response.Header.Get(CONTENT_TYPE_HEADER)
			__debug(fmt.Sprintf("[FileInfo] HEAD headers: Content-Length=%s, Accept-Ranges=%s, ETag=%s, Last-Modified=%s, Content-Type=%s", contentLength, acceptRanges, entityTag, lastModified, contentType))
			__debug(fmt.Sprintf("[FileInfo] HEAD result: size=%d, range=%v", size, supportRange))
		} else if isNonRetryableHTTPStatus(response.StatusCode) {
			err = newHTTPStatusError(METHOD_GET, requestURL, response.StatusCode)
			__debug(fmt.Sprintf("[FileInfo] HEAD request failed: %v", err))
		} else {
			__debug(fmt.Sprintf("[FileInfo] HEAD returned retryable/unexpected status: %d", response.StatusCode))
		}
	}
	if !isNonRetryableHTTPError(err) && (err != nil || size == 0) {
		__debug(fmt.Sprintf("[FileInfo] Need GET fallback: err=%v, size=%d", err, size))
		size, supportRange, entityTag, lastModified, contentType, err = h.getURLFileInfoByGet(requestURL, size, supportRange, entityTag, lastModified, contentType, headers)
	}
	if size > 0 && !supportRange {
		__debug("[FileInfo] Probing range support")
		supportRange = h.probeRangeSupport(requestURL, headers)
	}
	__debug(fmt.Sprintf("[FileInfo] Final: size=%d, range=%v, etag=%s, lastModified=%s, contentType=%s, err=%v", size, supportRange, entityTag, lastModified, contentType, err))
	return size, supportRange, entityTag, lastModified, contentType, err
}

//goland:noinspection SpellCheckingInspection
func (h *HTTP) getURLFileInfoByGet(requestURL string, size int64, supportRange bool, entityTag string, lastModified string, contentType string, headers map[string]string) (int64, bool, string, string, string, error) {
	__debug("[FileInfo] Falling back to GET request")
	err := error(nil)
	var getRequest *httplib.Request
	var getResponse *httplib.Response
	__debug(fmt.Sprintf("[FileInfo] Creating GET fallback request: %s", requestURL))
	if getRequest, err = httplib.NewRequest(METHOD_GET, requestURL, nil); err == nil {
		addDefaultHeaders(getRequest)
		applyCustomHeaders(getRequest, headers)
		__debug(fmt.Sprintf("[FileInfo] Sending GET fallback request: %s", requestURL))
		getResponse, err = h.client.Do(getRequest)
	} else {
		__debug(fmt.Sprintf("[FileInfo] Failed to create GET request: %v", err))
	}
	if err == nil && getResponse != nil {
		defer func(getResponse *httplib.Response) {
			__debug("[FileInfo] Closing GET fallback response body")
			_ = getResponse.Body.Close()
		}(getResponse)
		__debug(fmt.Sprintf("[FileInfo] GET response status: %d", getResponse.StatusCode))
		if getResponse.StatusCode >= HTTP_STATUS_OK_MIN && getResponse.StatusCode < HTTP_STATUS_OK_MAX {
			contentLength := getResponse.Header.Get(CONTENT_LENGTH_HEADER)
			if contentLength != "" {
				size, _ = strconv.ParseInt(contentLength, 10, 64)
			}
			acceptRanges := getResponse.Header.Get(ACCEPT_RANGES_HEADER)
			if acceptRanges == BYTES_UNIT {
				supportRange = true
			}
			if entityTag == "" {
				entityTag = getResponse.Header.Get(ENTITY_TAG_HEADER)
			}
			if lastModified == "" {
				lastModified = getResponse.Header.Get(LAST_MODIFIED_HEADER)
			}
			if contentType == "" {
				contentType = getResponse.Header.Get(CONTENT_TYPE_HEADER)
			}
			__debug(fmt.Sprintf("[FileInfo] GET headers: Content-Length=%s, Accept-Ranges=%s, ETag=%s, Last-Modified=%s, Content-Type=%s", contentLength, acceptRanges, entityTag, lastModified, contentType))
			__debug(fmt.Sprintf("[FileInfo] GET result: size=%d, range=%v", size, supportRange))
		} else {
			err = newHTTPStatusError(METHOD_GET, requestURL, getResponse.StatusCode)
			__debug(fmt.Sprintf("[FileInfo] GET request failed: %v", err))
		}
	}
	__debug(fmt.Sprintf("[FileInfo] GET fallback final: size=%d, range=%v, err=%v", size, supportRange, err))
	return size, supportRange, entityTag, lastModified, contentType, err
}

func isNonRetryableHTTPError(err error) bool {
	result := false
	if statusCode, ok := GetHTTPStatusCode(err); ok {
		result = isNonRetryableHTTPStatus(statusCode)
	} else {
		result = IsNetworkUnreachableError(err)
	}
	return result
}

func IsNetworkUnreachableError(err error) bool {
	result := false
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			errMsg := strings.ToLower(err.Error())
			if strings.Contains(errMsg, NETWORK_ERROR_NO_ROUTE) ||
				strings.Contains(errMsg, NETWORK_ERROR_NETWORK_UNREACHABLE) ||
				strings.Contains(errMsg, NETWORK_ERROR_HOST_UNREACHABLE) ||
				strings.Contains(errMsg, NETWORK_ERROR_CONNECTION_REFUSED) {
				result = true
			}
		}
		var syscallErr syscall.Errno
		if !result && errors.As(err, &syscallErr) {
			if syscallErr == syscall.ENETUNREACH || syscallErr == syscall.EHOSTUNREACH || syscallErr == syscall.ECONNREFUSED {
				result = true
			}
		}
	}
	return result
}

func isNonRetryableHTTPStatus(statusCode int) bool {
	result := false
	if statusCode == httplib.StatusUnauthorized ||
		statusCode == httplib.StatusPaymentRequired ||
		statusCode == httplib.StatusForbidden ||
		statusCode == httplib.StatusNotFound {
		result = true
	}
	return result
}

func isProtocolTestSuccessStatus(statusCode int) bool {
	result := false
	if statusCode >= HTTP_STATUS_OK_MIN && statusCode < HTTP_STATUS_OK_MAX ||
		isNonRetryableHTTPStatus(statusCode) {
		result = true
	}
	return result
}

func isTLSCertificateRequiredError(err error) bool {
	return err != nil && strings.Contains(err.Error(), TLS_ALERT_CERTIFICATE_REQUIRED)
}

func loadClientCertificate() *tls.Certificate {
	var result *tls.Certificate
	logger.Logger.Debug("[Certificate] Loading client certificate...")
	certificatePath := findCertificateFile(ca.CLIENT_CERTIFICATE_FILE_NAME)
	privateKeyPath := findCertificateFile(ca.CLIENT_PRIVATE_KEY_FILE_NAME)
	if certificatePath != "" && privateKeyPath != "" {
		logger.Logger.Debug(fmt.Sprintf("[Certificate] Loading client certificate: %s, key: %s", certificatePath, privateKeyPath))
		var certificateError error
		var certificate tls.Certificate
		if certificate, certificateError = tls.LoadX509KeyPair(certificatePath, privateKeyPath); certificateError == nil {
			result = &certificate
			logger.Logger.Debug("[Certificate] Client certificate loaded successfully")
		} else {
			logger.Logger.Debug(fmt.Sprintf("[Certificate] Failed to load client certificate: %v", certificateError))
		}
	} else {
		logger.Logger.Debug(fmt.Sprintf("[Certificate] Client certificate or key not found - cert: '%s', key: '%s'", certificatePath, privateKeyPath))
	}
	return result
}

func loadRootCertificate() *x509.CertPool {
	var result *x509.CertPool
	logger.Logger.Debug("[Certificate] Loading root CA certificate...")
	certificatePath := findCertificateFile(ca.CERTIFICATE_AUTHORITY_CERTIFICATE_FILE_NAME)
	if certificatePath != "" {
		logger.Logger.Debug(fmt.Sprintf("[Certificate] Loading CA certificate from: %s", certificatePath))
		var certificateError error
		var certificateBytes []byte
		if certificateBytes, certificateError = os.ReadFile(certificatePath); certificateError == nil {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(certificateBytes) {
				result = pool
				logger.Logger.Debug("[Certificate] CA certificate pool loaded successfully")
			} else {
				logger.Logger.Debug("[Certificate] Failed to append CA certificate to pool")
			}
		} else {
			logger.Logger.Debug(fmt.Sprintf("[Certificate] Failed to read CA certificate file: %v", certificateError))
		}
	} else {
		logger.Logger.Debug("[Certificate] CA certificate file not found")
	}
	return result
}

func newHTTPStatusError(method string, requestURL string, statusCode int) error {
	return &HTTP_STATUS_ERROR{Method: method, RequestURL: requestURL, StatusCode: statusCode}
}

func newParallelsDownloadChunks(fileSize int64, parallelsSize int64) []PARALLELS_DOWNLOAD_CHUNK {
	__debug(fmt.Sprintf("[DownloadParallelsEx] Building chunks: fileSize=%d, chunksSize=%d", fileSize, parallelsSize))
	result := make([]PARALLELS_DOWNLOAD_CHUNK, 0)
	for startByte := int64(0); startByte < fileSize; startByte += parallelsSize {
		endByte := startByte + parallelsSize - 1
		if endByte >= fileSize {
			endByte = fileSize - 1
		}
		result = append(result, PARALLELS_DOWNLOAD_CHUNK{StartByte: startByte, EndByte: endByte})
	}
	if len(result) > 0 {
		firstChunk := result[0]
		lastChunk := result[len(result)-1]
		__debug(fmt.Sprintf("[DownloadParallelsEx] Chunks built: count=%d, first=%d-%d, last=%d-%d", len(result), firstChunk.StartByte, firstChunk.EndByte, lastChunk.StartByte, lastChunk.EndByte))
	} else {
		__debug("[DownloadParallelsEx] Chunks built: count=0")
	}
	return result
}

func newProtocolTestClient(scheme string, includeClientCertificate bool) *httplib.Client {
	result := &httplib.Client{
		Timeout: DEFAULT_HTTP_TIMEOUT,
	}
	if scheme == SCHEME_HTTPS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		if includeClientCertificate {
			defaultTransportMutex.RLock()
			clientCertificate := defaultClientCertificate
			defaultTransportMutex.RUnlock()
			if clientCertificate != nil {
				tlsConfig.Certificates = []tls.Certificate{*clientCertificate}
			}
		}
		result.Transport = &httplib.Transport{
			TLSClientConfig: tlsConfig,
		}
	}
	return result
}

func newRequestUUID() string {
	result := ""
	bytesValue := make([]byte, UUID_BYTE_COUNT)
	var uuidError error
	if _, uuidError = rand.Read(bytesValue); uuidError == nil {
		bytesValue[6] = (bytesValue[6] & UUID_VERSION_MASK) | UUID_VERSION_SET
		bytesValue[8] = (bytesValue[8] & UUID_VARIANT_MASK) | UUID_VARIANT_SET
		result = fmt.Sprintf(UUID_FORMAT, bytesValue[0:4], bytesValue[4:6], bytesValue[6:8], bytesValue[8:10], bytesValue[10:])
	} else {
		result = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return result
}

func openTruncatedDownloadFile(filePath string, fileSize int64) (*os.File, error) {
	__debug(fmt.Sprintf("[DownloadParallelsEx] Opening file: %s", filePath))
	result, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, FILE_PERMISSION)
	if err == nil && fileSize > 0 {
		__debug(fmt.Sprintf("[DownloadParallelsEx] File opened, truncating file to: %d", fileSize))
		if err = result.Truncate(fileSize); err != nil {
			_ = result.Close()
			result = nil
			__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to truncate file: %v", err))
		} else {
			__debug(fmt.Sprintf("[DownloadParallelsEx] File truncated successfully: %s, size=%d", filePath, fileSize))
		}
	} else if err != nil {
		__debug(fmt.Sprintf("[DownloadParallelsEx] Failed to open file: %v", err))
	} else {
		__debug(fmt.Sprintf("[DownloadParallelsEx] File opened without truncate because fileSize=%d", fileSize))
	}
	return result, err
}

func (h *HTTP) probeRangeSupport(requestURL string, headers map[string]string) bool {
	__debug(fmt.Sprintf("[Range] Probing range support for: %s", requestURL))
	result := false
	err := error(nil)
	var request *httplib.Request
	if request, err = httplib.NewRequest(METHOD_GET, requestURL, nil); err == nil {
		addDefaultHeaders(request)
		applyCustomHeaders(request, headers)
		request.Header.Set(RANGE_HEADER, RANGE_PREFIX+RANGE_PROBE_VALUE)
		__debug(fmt.Sprintf("[Range] Probe request prepared: range=%s", RANGE_PREFIX+RANGE_PROBE_VALUE))
		var response *httplib.Response
		__debug("[Range] Sending range probe request")
		if response, err = h.client.Do(request); err == nil {
			defer func(response *httplib.Response) {
				__debug("[Range] Closing probe response body")
				_ = response.Body.Close()
			}(response)
			__debug(fmt.Sprintf("[Range] Probe response status: %d, contentLength=%d", response.StatusCode, response.ContentLength))
			if response.StatusCode == HTTP_STATUS_PARTIAL_CONTENT {
				result = true
				__debug("[Range] Range support confirmed")
			} else {
				__debug("[Range] Range not supported")
			}
		} else {
			__debug(fmt.Sprintf("[Range] Probe request failed: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("[Range] Failed to create probe request: %v", err))
	}
	__debug(fmt.Sprintf("[Range] Probe final: supportRange=%v", result))
	return result
}

func removeDownloadTemporaryFile(filePath string) {
	if filePath != "" {
		if err := os.Remove(filePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			__warning(fmt.Sprintf("[DownloadParallelsEx] Failed to remove temporary file: %s, error=%v", filePath, err))
		}
	}
}

func replaceDownloadTargetFile(temporaryFilePath string, filePath string) error {
	err := os.Remove(filePath)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		err = os.Rename(temporaryFilePath, filePath)
	}
	return err
}

//goland:noinspection DuplicatedCode
func testProtocol(scheme string, ip string, port int) error {
	__info(fmt.Sprintf("[Protocol] Testing %s protocol for %s:%d", scheme, ip, port))
	err := error(nil)
	requestURL := fmt.Sprintf(HTTP_TEST_URL_FORMAT, scheme, SCHEME_SEPARATOR, net.JoinHostPort(ip, strconv.Itoa(port)), HTTP_TEST_PATH)
	__debug(fmt.Sprintf("[Protocol] Test URL: %s", requestURL))
	err = testProtocolWithClient(newProtocolTestClient(scheme, false), requestURL)
	if scheme == SCHEME_HTTPS && isTLSCertificateRequiredError(err) {
		__info("[Protocol] Certificate required, retrying with client certificate")
		err = testProtocolWithClient(newProtocolTestClient(scheme, true), requestURL)
	}
	if err == nil {
		__info(fmt.Sprintf("[Protocol] %s protocol test succeeded", scheme))
	} else {
		__debug(fmt.Sprintf("[Protocol] %s protocol test failed: %v", scheme, err))
	}
	return err
}

func testProtocolWithClient(httpClient *httplib.Client, requestURL string) error {
	__debug(fmt.Sprintf("[Protocol] Testing with client: %s", requestURL))
	err := error(nil)
	var request *httplib.Request
	if request, err = httplib.NewRequest(METHOD_GET, requestURL, nil); err == nil {
		addDefaultHeaders(request)
		var response *httplib.Response
		if response, err = httpClient.Do(request); err == nil {
			defer func(response *httplib.Response) {
				_ = response.Body.Close()
			}(response)
			__debug(fmt.Sprintf("[Protocol] Response status: %d", response.StatusCode))
			if !isProtocolTestSuccessStatus(response.StatusCode) {
				err = newHTTPStatusError(METHOD_GET, requestURL, response.StatusCode)
				__debug(fmt.Sprintf("[Protocol] Unexpected status: %v", err))
			}
		} else {
			__debug(fmt.Sprintf("[Protocol] Request failed: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("[Protocol] Failed to create request: %v", err))
	}
	return err
}

func writeParallelsDownloadResponse(file *os.File, body io.Reader, chunk PARALLELS_DOWNLOAD_CHUNK, progress *int64, totalSize int64, mutex *sync.Mutex, progressHandler DOWNLOAD_PROGRESS_HANDLER) error {
	err := error(nil)
	expectedBytes := chunk.EndByte - chunk.StartByte + 1
	writeOffset := chunk.StartByte
	remaining := expectedBytes
	buf := make([]byte, 32*1024)
	var toRead int64
	var n int
	var readErr error
	var writeErr error
	__debug(fmt.Sprintf("[DownloadParallelsEx] Start streaming chunk: range=%d-%d, expectedBytes=%d", chunk.StartByte, chunk.EndByte, expectedBytes))
	for remaining > 0 {
		toRead = int64(len(buf))
		if remaining < toRead {
			toRead = remaining
		}
		n, readErr = body.Read(buf[:toRead])
		if n > 0 {
			mutex.Lock()
			_, writeErr = file.WriteAt(buf[:n], writeOffset)
			if writeErr == nil {
				writeOffset += int64(n)
				remaining -= int64(n)
				*progress += int64(n)
				if progressHandler != nil {
					progressHandler(*progress, totalSize)
				}
			} else {
				err = writeErr
				__debug(fmt.Sprintf("[DownloadParallelsEx] Write failed: %v", writeErr))
			}
			mutex.Unlock()
			if err != nil {
				break
			}
		}
		if readErr == io.EOF {
			break
		} else if readErr != nil {
			err = readErr
			__debug(fmt.Sprintf("[DownloadParallelsEx] Read failed: %v", readErr))
			break
		}
	}
	if err == nil && remaining > 0 {
		err = fmt.Errorf(ERR_INCOMPLETE_CHUNK, expectedBytes, expectedBytes-remaining)
		__warning(fmt.Sprintf("[DownloadParallelsEx] %v", err))
	}
	return err
}
