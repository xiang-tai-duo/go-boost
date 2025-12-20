// Package tusd
// File:        tusd.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/tusd/tusd.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: TUS resumable upload helpers.
// --------------------------------------------------------------------------------
package tusd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName,SpellCheckingInspection
type (
	ProgressHandler func(uploaded int64, total int64)

	SERVER_CAPABILITIES struct {
		ChecksumAlgorithms []string
		Extensions         []string
		MaxSize            int64
		ProtocolVersions   []string
	}

	TUSD_CLIENT struct {
		BaseUrl         *url.URL
		Capabilities    *SERVER_CAPABILITIES
		ChunkSize       int64
		Headers         map[string]string
		HttpClient      *http.Client
		Mutex           sync.Mutex
		ProgressHandler ProgressHandler
		RequestTimeout  time.Duration
		RetryInterval   time.Duration
		RetryMax        int
	}

	TUS_ERROR struct {
		Inner   error
		Message string
	}

	UPLOAD struct {
		Location      string
		Metadata      map[string]string
		Partial       bool
		RemoteOffset  int64
		RemoteSize    int64
		UploadExpired *time.Time
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	CONCAT_FINAL                  = "final"
	CONCAT_PARTIAL                = "partial"
	DEFAULT_CHUNK_SIZE            = 2 * 1024 * 1024
	DEFAULT_REQUEST_TIMEOUT       = 60 * time.Second
	DEFAULT_RETRY_INTERVAL        = 1 * time.Second
	DEFAULT_RETRY_MAX             = 3
	ERROR_MSG_BASE_URL_EMPTY      = "base URL cannot be empty"
	ERROR_MSG_CHUNK_DATA_EMPTY    = "chunk data cannot be empty"
	ERROR_MSG_CHUNK_SIZE_NEGATIVE = "chunk size must be non-negative"
	ERROR_MSG_FILE_PATH_EMPTY     = "file path cannot be empty"
	ERROR_MSG_HTTP_BODY_FMT       = "HTTP %d: %s"
	ERROR_MSG_KEY_HAS_SPACE       = "key %q contains spaces"
	ERROR_MSG_LACK_UPLOAD_OFFSET  = "lack of Upload-Offset required header in response"
	ERROR_MSG_METADATA_BAD_FORMAT = "metadata item %q has bad format"
	ERROR_MSG_NIL_POINTER         = "nil pointer"
	ERROR_MSG_OBTAIN_SERVER_CAP   = "cannot obtain server capabilities: %w"
	ERROR_MSG_OFFSET_NEGATIVE     = "offset cannot be negative"
	ERROR_MSG_PARSE_HEADER        = "cannot parse %s header %q: %w"
	ERROR_MSG_READER_NIL          = "reader cannot be nil"
	ERROR_MSG_REQUEST_TIMEOUT_GT0 = "request timeout must be greater than 0"
	ERROR_MSG_RESPONSE_NIL        = "response is nil"
	ERROR_MSG_RETRY_INTERVAL_NEG  = "retry interval cannot be negative"
	ERROR_MSG_RETRY_MAX_NEG       = "retry max cannot be negative"
	ERROR_MSG_SIZE_NEGATIVE       = "upload size is negative: %d"
	ERROR_MSG_SIZE_TOO_SMALL      = "upload size must be greater than 0"
	ERROR_MSG_UNEXPECTED_STATUS   = "unexpected status code %d when uploading chunk"
	ERROR_MSG_UPLOAD_URL_EMPTY    = "upload URL cannot be empty"
	EXTENSION_CREATION            = "creation"
	EXTENSION_CREATION_DEFER_LEN  = "creation-defer-length"
	EXTENSION_TERMINATION         = "termination"
	HEADER_CONTENT_LENGTH         = "Content-Length"
	HEADER_CONTENT_TYPE           = "Content-Type"
	HEADER_LOCATION               = "Location"
	HEADER_TUS_CHECKSUM_ALGORITHM = "Tus-Checksum-Algorithm"
	HEADER_TUS_EXTENSION          = "Tus-Extension"
	HEADER_TUS_MAX_SIZE           = "Tus-Max-Size"
	HEADER_TUS_RESUMABLE          = "Tus-Resumable"
	HEADER_TUS_VERSION            = "Tus-Version"
	HEADER_UPLOAD_CONCAT          = "Upload-Concat"
	HEADER_UPLOAD_DEFER_LENGTH    = "Upload-Defer-Length"
	HEADER_UPLOAD_EXPIRES         = "Upload-Expires"
	HEADER_UPLOAD_LENGTH          = "Upload-Length"
	HEADER_UPLOAD_METADATA        = "Upload-Metadata"
	HEADER_UPLOAD_OFFSET          = "Upload-Offset"
	METADATA_FILENAME             = "filename"
	MODULE_NAME_TUSD              = "tusd"
	OFFSET_OCTET_STREAM_TYPE      = "application/offset+octet-stream"
	RESPONSE_BODY_NO_CONTENT      = "<no body>"
	RESPONSE_BODY_READ_ERROR      = "<read error>"
	RESPONSE_BODY_READ_SIZE       = 256
	RESPONSE_PROTOCOL_VERSION_FMT = "request protocol version %q, server supported versions are %q"
	SIZE_UNKNOWN                  = -1
	STR_ONE                       = "1"
	TUS_PROTOCOL_VERSION          = "1.0.0"
)

var (
	ErrCannotUpload       = TUS_ERROR{Message: "can not upload"}
	ErrChecksumMismatch   = TUS_ERROR{Message: "checksum mismatch"}
	ErrOffsetsNotSynced   = TUS_ERROR{Message: "client stream and server offsets are not synced"}
	ErrProtocol           = TUS_ERROR{Message: "protocol error"}
	ErrUnexpectedResponse = TUS_ERROR{Message: "unexpected HTTP response code"}
	ErrUnsupportedFeature = TUS_ERROR{Message: "unsupported feature"}
	ErrUploadDoesNotExist = TUS_ERROR{Message: "upload does not exist"}
	ErrUploadTooLarge     = TUS_ERROR{Message: "upload is too large"}
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_TUSD, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_TUSD, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_TUSD, logger.SKIP_STACK_FRAMES_BASE)
}

func New(baseUrl string) *TUSD_CLIENT {
	u, _ := url.Parse(baseUrl)
	return &TUSD_CLIENT{
		BaseUrl:      u,
		Capabilities: nil,
		ChunkSize:    DEFAULT_CHUNK_SIZE,
		Headers:      make(map[string]string),
		HttpClient: &http.Client{
			Timeout: DEFAULT_REQUEST_TIMEOUT,
		},
		Mutex:           sync.Mutex{},
		ProgressHandler: nil,
		RequestTimeout:  DEFAULT_REQUEST_TIMEOUT,
		RetryInterval:   DEFAULT_RETRY_INTERVAL,
		RetryMax:        DEFAULT_RETRY_MAX,
	}
}

func (c *TUSD_CLIENT) CreateUpload(u *UPLOAD, remoteSize int64, partial bool, meta map[string]string) (*http.Response, error) {
	err := error(nil)
	var response *http.Response

	if u == nil {
		err = fmt.Errorf(ERROR_MSG_NIL_POINTER)
	} else if remoteSize < 0 && remoteSize != SIZE_UNKNOWN {
		err = fmt.Errorf(ERROR_MSG_SIZE_NEGATIVE, remoteSize)
	} else if err = c.ensureExtension(EXTENSION_CREATION); err == nil {
		var req *http.Request
		if req, err = http.NewRequest(http.MethodPost, c.BaseUrl.String(), nil); err == nil {
			req.Header.Set(HEADER_CONTENT_LENGTH, strconv.FormatInt(0, 10))
			if partial {
				req.Header.Set(HEADER_UPLOAD_CONCAT, CONCAT_PARTIAL)
			}

			switch {
			case remoteSize == SIZE_UNKNOWN:
				if err = c.ensureExtension(EXTENSION_CREATION_DEFER_LEN); err == nil {
					req.Header.Set(HEADER_UPLOAD_DEFER_LENGTH, STR_ONE)
				}
			case remoteSize > 0:
				req.Header.Set(HEADER_UPLOAD_LENGTH, strconv.FormatInt(remoteSize, 10))
			}

			if err == nil && len(meta) > 0 {
				m := ""
				if m, err = EncodeMetadata(meta); err == nil {
					req.Header.Set(HEADER_UPLOAD_METADATA, m)
				}
			}

			if err == nil {
				if response, err = c.tusRequest(req); err == nil {
					defer response.Body.Close()
					switch response.StatusCode {
					case http.StatusCreated:
						u2 := UPLOAD{}
						u2.Location = response.Header.Get(HEADER_LOCATION)
						u2.Metadata = meta
						u2.Partial = partial
						u2.RemoteSize = remoteSize
						if v := response.Header.Get(HEADER_UPLOAD_EXPIRES); v != "" {
							t := time.Time{}
							if t, err = time.Parse(time.RFC1123, v); err == nil {
								u2.UploadExpired = &t
							} else {
								err = ErrProtocol.WithErr(fmt.Errorf(ERROR_MSG_PARSE_HEADER, HEADER_UPLOAD_EXPIRES, v, err))
							}
						}
						if err == nil {
							*u = u2
						}
					case http.StatusRequestEntityTooLarge:
						err = ErrUploadTooLarge.WithResponse(response)
					default:
						err = ErrUnexpectedResponse
					}
				}
			}
		}
	}

	return response, err
}

func DecodeMetadata(raw string) (map[string]string, error) {
	result := make(map[string]string)
	err := error(nil)

	for _, item := range strings.Split(raw, ",") {
		if err == nil {
			kv := strings.SplitN(item, " ", 2)
			if len(kv) > 1 {
				val := make([]byte, 0)
				if val, err = base64.StdEncoding.DecodeString(kv[1]); err == nil {
					result[kv[0]] = string(val)
				}
			} else {
				err = fmt.Errorf(ERROR_MSG_METADATA_BAD_FORMAT, item)
			}
		}
	}

	return result, err
}

func (c *TUSD_CLIENT) DeleteUpload(u UPLOAD) (*http.Response, error) {
	err := error(nil)
	var response *http.Response

	if err = c.ensureExtension(EXTENSION_TERMINATION); err == nil {
		var loc *url.URL
		if loc, err = url.Parse(u.Location); err == nil {
			ref := c.BaseUrl.ResolveReference(loc).String()
			var req *http.Request
			if req, err = http.NewRequest(http.MethodDelete, ref, nil); err == nil {
				if response, err = c.tusRequest(req); err == nil {
					defer response.Body.Close()
					switch response.StatusCode {
					case http.StatusNoContent:
					case http.StatusNotFound, http.StatusGone, http.StatusForbidden:
						err = ErrUploadDoesNotExist.WithResponse(response)
					default:
						err = ErrUnexpectedResponse
					}
				}
			}
		}
	}

	return response, err
}

func EncodeMetadata(metadata map[string]string) (string, error) {
	result := ""
	err := error(nil)
	encoded := make([]string, 0)

	for k, v := range metadata {
		if err == nil {
			if !strings.Contains(k, " ") {
				encoded = append(encoded, fmt.Sprintf("%s %s", k, base64.StdEncoding.EncodeToString([]byte(v))))
			} else {
				err = fmt.Errorf(ERROR_MSG_KEY_HAS_SPACE, k)
			}
		}
	}

	if err == nil {
		result = strings.Join(encoded, ",")
	}

	return result, err
}

func (te TUS_ERROR) Error() string {
	result := ""
	if te.Inner == nil {
		result = te.Message
	} else {
		result = fmt.Sprintf("%s: %s", te.Message, te.Inner)
	}
	return result
}

func (c *TUSD_CLIENT) GetBaseUrl() string {
	result := ""
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.BaseUrl != nil {
		result = c.BaseUrl.String()
	}
	return result
}

func (c *TUSD_CLIENT) GetCapabilities() *SERVER_CAPABILITIES {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Capabilities
}

func (c *TUSD_CLIENT) GetChunkSize() int64 {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.ChunkSize
}

func (c *TUSD_CLIENT) GetHeader(key string) string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Headers[key]
}

func (c *TUSD_CLIENT) GetOffset(uploadUrl string) (int64, int64, error) {
	resultOffset := int64(0)
	resultTotal := int64(0)
	err := error(nil)

	if uploadUrl == "" {
		err = fmt.Errorf(ERROR_MSG_UPLOAD_URL_EMPTY)
	} else {
		u := UPLOAD{}
		if _, err = c.GetUpload(&u, uploadUrl); err == nil {
			resultOffset = u.RemoteOffset
			resultTotal = u.RemoteSize
		}
	}

	return resultOffset, resultTotal, err
}

func (c *TUSD_CLIENT) GetRequestTimeout() time.Duration {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.RequestTimeout
}

func (c *TUSD_CLIENT) GetRetryInterval() time.Duration {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.RetryInterval
}

func (c *TUSD_CLIENT) GetRetryMax() int {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.RetryMax
}

func (c *TUSD_CLIENT) GetUpload(u *UPLOAD, location string) (*http.Response, error) {
	err := error(nil)
	var response *http.Response

	if u == nil {
		err = fmt.Errorf(ERROR_MSG_NIL_POINTER)
	} else {
		var loc *url.URL
		if loc, err = url.Parse(location); err == nil {
			ref := c.BaseUrl.ResolveReference(loc).String()
			var req *http.Request
			if req, err = http.NewRequest(http.MethodHead, ref, nil); err == nil {
				if response, err = c.tusRequest(req); err == nil {
					defer response.Body.Close()
					switch response.StatusCode {
					case http.StatusOK:
						u2 := UPLOAD{}
						u2.Location = location
						u2.Partial = response.Header.Get(HEADER_UPLOAD_CONCAT) == CONCAT_PARTIAL

						uploadOffset := response.Header.Get(HEADER_UPLOAD_OFFSET)
						if uploadOffset == "" {
							if response.Header.Get(HEADER_UPLOAD_CONCAT) == CONCAT_FINAL {
								u2.RemoteOffset = SIZE_UNKNOWN
							} else {
								err = ErrProtocol.WithText(ERROR_MSG_LACK_UPLOAD_OFFSET)
							}
						} else if u2.RemoteOffset, err = strconv.ParseInt(uploadOffset, 10, 64); err != nil {
							err = ErrProtocol.WithErr(fmt.Errorf(ERROR_MSG_PARSE_HEADER, HEADER_UPLOAD_OFFSET, uploadOffset, err))
						}

						if err == nil {
							if v := response.Header.Get(HEADER_UPLOAD_LENGTH); v != "" {
								if u2.RemoteSize, err = strconv.ParseInt(v, 10, 64); err != nil {
									err = ErrProtocol.WithErr(fmt.Errorf(ERROR_MSG_PARSE_HEADER, HEADER_UPLOAD_LENGTH, v, err))
								}
							}
						}

						if err == nil {
							if v := response.Header.Get(HEADER_UPLOAD_METADATA); v != "" {
								if u2.Metadata, err = DecodeMetadata(v); err != nil {
									err = ErrProtocol.WithErr(fmt.Errorf(ERROR_MSG_PARSE_HEADER, HEADER_UPLOAD_METADATA, v, err))
								}
							}
						}

						if err == nil {
							*u = u2
						}
					case http.StatusNotFound, http.StatusGone, http.StatusForbidden:
						err = ErrUploadDoesNotExist.WithResponse(response)
					default:
						err = ErrUnexpectedResponse
					}
				}
			}
		}
	}

	return response, err
}

func (c *TUSD_CLIENT) HasExtension(extension string) bool {
	result := false
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.Capabilities != nil {
		for _, e := range c.Capabilities.Extensions {
			if e == extension {
				result = true
				break
			}
		}
	}
	return result
}

func (te TUS_ERROR) Is(e error) bool {
	result := false
	v := TUS_ERROR{}
	ok := errors.As(e, &v)
	if ok && v.Message == te.Message || errors.Is(te.Inner, e) {
		result = true
	}
	return result
}

func (c *TUSD_CLIENT) RemoveHeader(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	delete(c.Headers, key)
}

func (c *TUSD_CLIENT) SetBaseUrl(baseUrl string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	u, _ := url.Parse(baseUrl)
	c.BaseUrl = u
}

func (c *TUSD_CLIENT) SetChunkSize(chunkSize int64) error {
	err := error(nil)
	if chunkSize >= 0 {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.ChunkSize = chunkSize
	} else {
		err = fmt.Errorf(ERROR_MSG_CHUNK_SIZE_NEGATIVE)
	}
	return err
}

func (c *TUSD_CLIENT) SetHeader(key string, value string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Headers[key] = value
}

func (c *TUSD_CLIENT) SetProgressHandler(handler ProgressHandler) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.ProgressHandler = handler
}

func (c *TUSD_CLIENT) SetRequestTimeout(timeout time.Duration) error {
	err := error(nil)
	if timeout > 0 {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.RequestTimeout = timeout
		c.HttpClient.Timeout = timeout
	} else {
		err = fmt.Errorf(ERROR_MSG_REQUEST_TIMEOUT_GT0)
	}
	return err
}

func (c *TUSD_CLIENT) SetRetryInterval(interval time.Duration) error {
	err := error(nil)
	if interval >= 0 {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.RetryInterval = interval
	} else {
		err = fmt.Errorf(ERROR_MSG_RETRY_INTERVAL_NEG)
	}
	return err
}

func (c *TUSD_CLIENT) SetRetryMax(retryMax int) error {
	err := error(nil)
	if retryMax >= 0 {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.RetryMax = retryMax
	} else {
		err = fmt.Errorf(ERROR_MSG_RETRY_MAX_NEG)
	}
	return err
}

func (te TUS_ERROR) Unwrap() error {
	return te.Inner
}

func (c *TUSD_CLIENT) UpdateCapabilities() (*http.Response, error) {
	err := error(nil)
	var response *http.Response
	var req *http.Request

	if req, err = http.NewRequest(http.MethodOptions, c.BaseUrl.String(), nil); err == nil {
		if response, err = c.tusRequest(req); err == nil {
			defer response.Body.Close()
			switch response.StatusCode {
			case http.StatusNoContent, http.StatusOK:
				caps := &SERVER_CAPABILITIES{}
				if v := response.Header.Get(HEADER_TUS_MAX_SIZE); v != "" {
					if caps.MaxSize, err = strconv.ParseInt(v, 10, 64); err != nil {
						err = ErrProtocol.WithErr(fmt.Errorf(ERROR_MSG_PARSE_HEADER, HEADER_TUS_MAX_SIZE, v, err))
					}
				}

				if err == nil {
					if v := response.Header.Get(HEADER_TUS_EXTENSION); v != "" {
						caps.Extensions = strings.Split(v, ",")
					}
				}

				if err == nil {
					if v := response.Header.Get(HEADER_TUS_VERSION); v != "" {
						caps.ProtocolVersions = strings.Split(v, ",")
					}
				}

				if err == nil {
					if v := response.Header.Get(HEADER_TUS_CHECKSUM_ALGORITHM); v != "" {
						caps.ChecksumAlgorithms = strings.Split(v, ",")
					}
				}

				if err == nil {
					c.Mutex.Lock()
					c.Capabilities = caps
					c.Mutex.Unlock()
				}
			default:
				err = ErrUnexpectedResponse
			}
		}
	}

	return response, err
}

func (c *TUSD_CLIENT) Upload(filePath string, metadata map[string]string) (string, error) {
	result := ""
	err := error(nil)
	var fileInfo os.FileInfo

	if filePath == "" {
		err = fmt.Errorf(ERROR_MSG_FILE_PATH_EMPTY)
	} else if fileInfo, err = os.Stat(filePath); err == nil {
		size := fileInfo.Size()

		finalMetadata := make(map[string]string)
		for key, value := range metadata {
			finalMetadata[key] = value
		}
		if _, exists := finalMetadata[METADATA_FILENAME]; !exists {
			finalMetadata[METADATA_FILENAME] = filepath.Base(filePath)
		}

		u := UPLOAD{}
		if _, err = c.CreateUpload(&u, size, false, finalMetadata); err == nil {
			if err = c.UploadFile(u.Location, filePath); err == nil {
				result = u.Location
			}
		}
	}

	return result, err
}

func (c *TUSD_CLIENT) UploadChunk(uploadUrl string, offset int64, data []byte) (int64, error) {
	result := int64(0)
	err := error(nil)
	var req *http.Request
	var resp *http.Response

	if uploadUrl == "" {
		err = fmt.Errorf(ERROR_MSG_UPLOAD_URL_EMPTY)
	} else if offset < 0 {
		err = fmt.Errorf(ERROR_MSG_OFFSET_NEGATIVE)
	} else if len(data) == 0 {
		err = fmt.Errorf(ERROR_MSG_CHUNK_DATA_EMPTY)
	} else if req, err = http.NewRequest(http.MethodPatch, uploadUrl, strings.NewReader(string(data))); err == nil {
		req.Header.Set(HEADER_TUS_RESUMABLE, TUS_PROTOCOL_VERSION)
		req.Header.Set(HEADER_UPLOAD_OFFSET, strconv.FormatInt(offset, 10))
		req.Header.Set(HEADER_CONTENT_TYPE, OFFSET_OCTET_STREAM_TYPE)
		req.ContentLength = int64(len(data))
		if resp, err = c.tusRequest(req); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNoContent {
				uploadOffset := resp.Header.Get(HEADER_UPLOAD_OFFSET)
				if result, err = strconv.ParseInt(uploadOffset, 10, 64); err != nil {
					err = fmt.Errorf(ERROR_MSG_PARSE_HEADER, HEADER_UPLOAD_OFFSET, uploadOffset, err)
				}
			} else {
				err = fmt.Errorf(ERROR_MSG_UNEXPECTED_STATUS, resp.StatusCode)
			}
		}
	}

	return result, err
}

func (c *TUSD_CLIENT) UploadFile(uploadUrl string, filePath string) error {
	err := error(nil)
	var file *os.File
	var currentOffset int64
	var total int64
	var chunkSize int64
	var handler ProgressHandler

	if uploadUrl == "" {
		err = fmt.Errorf(ERROR_MSG_UPLOAD_URL_EMPTY)
	} else if filePath == "" {
		err = fmt.Errorf(ERROR_MSG_FILE_PATH_EMPTY)
	} else if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		if currentOffset, total, err = c.GetOffset(uploadUrl); err == nil {
			if total <= 0 {
				var fileInfo os.FileInfo
				if fileInfo, err = file.Stat(); err == nil {
					total = fileInfo.Size()
				}
			}

			if err == nil {
				if _, err = file.Seek(currentOffset, io.SeekStart); err == nil {
					c.Mutex.Lock()
					chunkSize = c.ChunkSize
					handler = c.ProgressHandler
					c.Mutex.Unlock()

					buffer := make([]byte, chunkSize)
					done := false
					for !done && err == nil {
						readCount, readErr := file.Read(buffer)
						if readCount > 0 {
							chunkData := buffer[:readCount]
							attempt := 0
							chunkUploaded := false
							for !chunkUploaded && err == nil {
								newOffset := int64(0)
								if newOffset, err = c.UploadChunk(uploadUrl, currentOffset, chunkData); err == nil {
									currentOffset = newOffset
									chunkUploaded = true
									if handler != nil {
										handler(currentOffset, total)
									}
								} else {
									attempt++
									if attempt <= c.RetryMax {
										time.Sleep(c.RetryInterval)
										err = nil
									}
								}
							}
						}

						if readErr == io.EOF {
							done = true
						} else if readErr != nil {
							err = readErr
						} else if readCount == 0 {
							done = true
						}
					}
				}
			}
		}
	}

	return err
}

func (c *TUSD_CLIENT) UploadFromReader(uploadUrl string, reader io.Reader, total int64) error {
	err := error(nil)
	var currentOffset int64
	var chunkSize int64
	var handler ProgressHandler

	if uploadUrl == "" {
		err = fmt.Errorf(ERROR_MSG_UPLOAD_URL_EMPTY)
	} else if reader == nil {
		err = fmt.Errorf(ERROR_MSG_READER_NIL)
	} else if currentOffset, _, err = c.GetOffset(uploadUrl); err == nil {
		c.Mutex.Lock()
		chunkSize = c.ChunkSize
		handler = c.ProgressHandler
		c.Mutex.Unlock()

		buffer := make([]byte, chunkSize)
		done := false
		for !done && err == nil {
			readCount, readErr := reader.Read(buffer)
			if readCount > 0 {
				chunkData := buffer[:readCount]
				attempt := 0
				chunkUploaded := false
				for !chunkUploaded && err == nil {
					newOffset := int64(0)
					if newOffset, err = c.UploadChunk(uploadUrl, currentOffset, chunkData); err == nil {
						currentOffset = newOffset
						chunkUploaded = true
						if handler != nil {
							handler(currentOffset, total)
						}
					} else {
						attempt++
						if attempt <= c.RetryMax {
							time.Sleep(c.RetryInterval)
							err = nil
						}
					}
				}
			}

			if readErr == io.EOF {
				done = true
			} else if readErr != nil {
				err = readErr
			} else if readCount == 0 {
				done = true
			}
		}
	}

	return err
}

func (te TUS_ERROR) WithErr(err error) TUS_ERROR {
	te.Inner = err
	return te
}

func (te TUS_ERROR) WithResponse(r *http.Response) TUS_ERROR {
	err := error(nil)
	if r == nil {
		te.Inner = errors.New(ERROR_MSG_RESPONSE_NIL)
	} else {
		b := make([]byte, RESPONSE_BODY_READ_SIZE)
		l := 0
		if l, err = io.ReadFull(r.Body, b); err == nil || err == io.EOF {
			if l > 0 {
				te.Inner = fmt.Errorf(ERROR_MSG_HTTP_BODY_FMT, r.StatusCode, b[:l])
			} else {
				te.Inner = fmt.Errorf(ERROR_MSG_HTTP_BODY_FMT, r.StatusCode, RESPONSE_BODY_NO_CONTENT)
			}
		} else {
			te.Inner = fmt.Errorf(ERROR_MSG_HTTP_BODY_FMT, r.StatusCode, RESPONSE_BODY_READ_ERROR)
		}
	}
	return te
}

func (te TUS_ERROR) WithText(s string) TUS_ERROR {
	te.Inner = errors.New(s)
	return te
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_TUSD, logger.SKIP_STACK_FRAMES_BASE)
}

func (c *TUSD_CLIENT) ensureExtension(extension string) error {
	err := error(nil)

	if c.Capabilities == nil {
		if _, err = c.UpdateCapabilities(); err != nil {
			err = fmt.Errorf(ERROR_MSG_OBTAIN_SERVER_CAP, err)
		}
	}

	if err == nil {
		found := false
		for _, e := range c.Capabilities.Extensions {
			if extension == e {
				found = true
				break
			}
		}
		if !found {
			err = ErrUnsupportedFeature.WithText(extension)
		}
	}

	return err
}

func (c *TUSD_CLIENT) tusRequest(req *http.Request) (*http.Response, error) {
	err := error(nil)
	var response *http.Response

	if req.Method != http.MethodOptions && req.Header.Get(HEADER_TUS_RESUMABLE) == "" {
		req.Header.Set(HEADER_TUS_RESUMABLE, TUS_PROTOCOL_VERSION)
	}

	c.Mutex.Lock()
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	c.Mutex.Unlock()

	if response, err = c.HttpClient.Do(req); err == nil {
		if response.StatusCode == http.StatusPreconditionFailed {
			versions := response.Header.Get(HEADER_TUS_VERSION)
			err = ErrProtocol.WithText(fmt.Sprintf(RESPONSE_PROTOCOL_VERSION_FMT, TUS_PROTOCOL_VERSION, versions))
		} else if response.Body != nil {
			bodyBytes := make([]byte, 0)
			if bodyBytes, err = io.ReadAll(response.Body); err == nil {
				response.Body.Close()
				response.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			}
		}
	}

	return response, err
}
