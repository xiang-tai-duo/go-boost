// Package tusdclient
// File:        client.go
// Url:         https://github.com/xiang-tai-duo/go-bootstrap/blob/master/tusd/client/client.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: TUSD client implementation for the tus resumable upload protocol.
// --------------------------------------------------------------------------------
package tusdclient

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	DEFAULT_CHUNK_SIZE       = 5 * 1024 * 1024
	DEFAULT_REQUEST_TIMEOUT  = 60 * time.Second
	DEFAULT_RETRY_INTERVAL   = 1 * time.Second
	DEFAULT_RETRY_MAX        = 3
	HEADER_CONTENT_TYPE      = "Content-Type"
	HEADER_LOCATION          = "Location"
	HEADER_TUS_RESUMABLE     = "Tus-Resumable"
	HEADER_UPLOAD_LENGTH     = "Upload-Length"
	HEADER_UPLOAD_METADATA   = "Upload-Metadata"
	HEADER_UPLOAD_OFFSET     = "Upload-Offset"
	OFFSET_OCTET_STREAM_TYPE = "application/offset+octet-stream"
	TUS_PROTOCOL_VERSION     = "1.0.0"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	PROGRESS_HANDLER func(uploaded int64, total int64)

	TUSD_CLIENT struct {
		baseUrl         string
		chunkSize       int64
		requestTimeout  time.Duration
		retryInterval   time.Duration
		retryMax        int
		headers         map[string]string
		httpClient      *http.Client
		progressHandler PROGRESS_HANDLER
		lock            sync.Mutex
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(baseUrl string) *TUSD_CLIENT {
	return &TUSD_CLIENT{
		baseUrl:        baseUrl,
		chunkSize:      DEFAULT_CHUNK_SIZE,
		requestTimeout: DEFAULT_REQUEST_TIMEOUT,
		retryInterval:  DEFAULT_RETRY_INTERVAL,
		retryMax:       DEFAULT_RETRY_MAX,
		headers:        make(map[string]string),
		httpClient: &http.Client{
			Timeout: DEFAULT_REQUEST_TIMEOUT,
		},
	}
}

//goland:noinspection GoUnusedExportedFunction
func CreateMetadata(metadata map[string]string) string {
	result := ""
	if len(metadata) > 0 {
		parts := make([]string, 0, len(metadata))
		for key, value := range metadata {
			encoded := base64.StdEncoding.EncodeToString([]byte(value))
			parts = append(parts, fmt.Sprintf("%s %s", key, encoded))
		}
		result = strings.Join(parts, ",")
	}
	return result
}

func (c *TUSD_CLIENT) CreateUpload(size int64, metadata map[string]string) (string, error) {
	result := ""
	err := error(nil)
	if size <= 0 {
		err = fmt.Errorf("upload size must be greater than 0")
	} else if c.baseUrl == "" {
		err = fmt.Errorf("base URL cannot be empty")
	} else {
		if request, requestErr := http.NewRequest(http.MethodPost, c.baseUrl, nil); requestErr == nil {
			request.Header.Set(HEADER_TUS_RESUMABLE, TUS_PROTOCOL_VERSION)
			request.Header.Set(HEADER_UPLOAD_LENGTH, strconv.FormatInt(size, 10))
			if metadataValue := CreateMetadata(metadata); metadataValue != "" {
				request.Header.Set(HEADER_UPLOAD_METADATA, metadataValue)
			}
			c.lock.Lock()
			for key, value := range c.headers {
				request.Header.Set(key, value)
			}
			c.lock.Unlock()
			if response, responseErr := c.httpClient.Do(request); responseErr == nil {
				_ = response.Body.Close()
				if response.StatusCode == http.StatusCreated {
					location := response.Header.Get(HEADER_LOCATION)
					if location != "" {
						result = c.resolveLocation(location)
					} else {
						err = fmt.Errorf("missing location header in response")
					}
				} else {
					err = fmt.Errorf("unexpected status code %d when creating upload", response.StatusCode)
				}
			} else {
				err = responseErr
			}
		} else {
			err = requestErr
		}
	}
	return result, err
}

func (c *TUSD_CLIENT) GetBaseUrl() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.baseUrl
}

func (c *TUSD_CLIENT) GetChunkSize() int64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.chunkSize
}

func (c *TUSD_CLIENT) GetHeader(key string) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.headers[key]
}

func (c *TUSD_CLIENT) GetOffset(uploadUrl string) (int64, int64, error) {
	result := int64(0)
	total := int64(0)
	err := error(nil)
	if uploadUrl == "" {
		err = fmt.Errorf("upload URL cannot be empty")
	} else {
		if request, requestErr := http.NewRequest(http.MethodHead, uploadUrl, nil); requestErr == nil {
			request.Header.Set(HEADER_TUS_RESUMABLE, TUS_PROTOCOL_VERSION)
			c.lock.Lock()
			for key, value := range c.headers {
				request.Header.Set(key, value)
			}
			c.lock.Unlock()
			if response, responseErr := c.httpClient.Do(request); responseErr == nil {
				_ = response.Body.Close()
				if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusNoContent {
					if offsetValue, parseErr := strconv.ParseInt(response.Header.Get(HEADER_UPLOAD_OFFSET), 10, 64); parseErr == nil {
						result = offsetValue
					} else {
						err = parseErr
					}
					if err == nil {
						if lengthValue, parseErr := strconv.ParseInt(response.Header.Get(HEADER_UPLOAD_LENGTH), 10, 64); parseErr == nil {
							total = lengthValue
						}
					}
				} else {
					err = fmt.Errorf("unexpected status code %d when getting offset", response.StatusCode)
				}
			} else {
				err = responseErr
			}
		} else {
			err = requestErr
		}
	}
	return result, total, err
}

func (c *TUSD_CLIENT) GetRequestTimeout() time.Duration {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.requestTimeout
}

func (c *TUSD_CLIENT) GetRetryInterval() time.Duration {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.retryInterval
}

func (c *TUSD_CLIENT) GetRetryMax() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.retryMax
}

func (c *TUSD_CLIENT) RemoveHeader(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.headers, key)
}

func (c *TUSD_CLIENT) SetBaseUrl(baseUrl string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.baseUrl = baseUrl
}

func (c *TUSD_CLIENT) SetChunkSize(chunkSize int64) error {
	err := error(nil)
	if chunkSize > 0 {
		c.lock.Lock()
		c.chunkSize = chunkSize
		c.lock.Unlock()
	} else {
		err = fmt.Errorf("chunk size must be greater than 0")
	}
	return err
}

func (c *TUSD_CLIENT) SetHeader(key string, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.headers[key] = value
}

func (c *TUSD_CLIENT) SetProgressHandler(handler PROGRESS_HANDLER) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.progressHandler = handler
}

func (c *TUSD_CLIENT) SetRequestTimeout(timeout time.Duration) error {
	err := error(nil)
	if timeout > 0 {
		c.lock.Lock()
		c.requestTimeout = timeout
		c.httpClient.Timeout = timeout
		c.lock.Unlock()
	} else {
		err = fmt.Errorf("request timeout must be greater than 0")
	}
	return err
}

func (c *TUSD_CLIENT) SetRetryInterval(interval time.Duration) error {
	err := error(nil)
	if interval >= 0 {
		c.lock.Lock()
		c.retryInterval = interval
		c.lock.Unlock()
	} else {
		err = fmt.Errorf("retry interval cannot be negative")
	}
	return err
}

func (c *TUSD_CLIENT) SetRetryMax(retryMax int) error {
	err := error(nil)
	if retryMax >= 0 {
		c.lock.Lock()
		c.retryMax = retryMax
		c.lock.Unlock()
	} else {
		err = fmt.Errorf("retry max cannot be negative")
	}
	return err
}

func (c *TUSD_CLIENT) Upload(filePath string, metadata map[string]string) (string, error) {
	result := ""
	err := error(nil)
	if filePath == "" {
		err = fmt.Errorf("file path cannot be empty")
	} else {
		if fileInfo, statErr := os.Stat(filePath); statErr == nil {
			size := fileInfo.Size()
			finalMetadata := make(map[string]string)
			for key, value := range metadata {
				finalMetadata[key] = value
			}
			if _, exists := finalMetadata["filename"]; !exists {
				finalMetadata["filename"] = filepath.Base(filePath)
			}
			if uploadUrl, createErr := c.CreateUpload(size, finalMetadata); createErr == nil {
				if uploadErr := c.UploadFile(uploadUrl, filePath); uploadErr == nil {
					result = uploadUrl
				} else {
					err = uploadErr
				}
			} else {
				err = createErr
			}
		} else {
			err = statErr
		}
	}
	return result, err
}

func (c *TUSD_CLIENT) UploadChunk(uploadUrl string, offset int64, data []byte) (int64, error) {
	result := int64(0)
	err := error(nil)
	if uploadUrl == "" {
		err = fmt.Errorf("upload URL cannot be empty")
	} else if offset < 0 {
		err = fmt.Errorf("offset cannot be negative")
	} else if len(data) == 0 {
		err = fmt.Errorf("chunk data cannot be empty")
	} else {
		if request, requestErr := http.NewRequest(http.MethodPatch, uploadUrl, strings.NewReader(string(data))); requestErr == nil {
			request.Header.Set(HEADER_TUS_RESUMABLE, TUS_PROTOCOL_VERSION)
			request.Header.Set(HEADER_UPLOAD_OFFSET, strconv.FormatInt(offset, 10))
			request.Header.Set(HEADER_CONTENT_TYPE, OFFSET_OCTET_STREAM_TYPE)
			request.ContentLength = int64(len(data))
			c.lock.Lock()
			for key, value := range c.headers {
				request.Header.Set(key, value)
			}
			c.lock.Unlock()
			if response, responseErr := c.httpClient.Do(request); responseErr == nil {
				_ = response.Body.Close()
				if response.StatusCode == http.StatusNoContent {
					if newOffset, parseErr := strconv.ParseInt(response.Header.Get(HEADER_UPLOAD_OFFSET), 10, 64); parseErr == nil {
						result = newOffset
					} else {
						err = parseErr
					}
				} else {
					err = fmt.Errorf("unexpected status code %d when uploading chunk", response.StatusCode)
				}
			} else {
				err = responseErr
			}
		} else {
			err = requestErr
		}
	}
	return result, err
}

func (c *TUSD_CLIENT) UploadFile(uploadUrl string, filePath string) error {
	err := error(nil)
	if uploadUrl == "" {
		err = fmt.Errorf("upload URL cannot be empty")
	} else if filePath == "" {
		err = fmt.Errorf("file path cannot be empty")
	} else {
		if file, openErr := os.Open(filePath); openErr == nil {
			if currentOffset, total, offsetErr := c.GetOffset(uploadUrl); offsetErr == nil {
				if total <= 0 {
					if fileInfo, statErr := file.Stat(); statErr == nil {
						total = fileInfo.Size()
					}
				}
				if _, seekErr := file.Seek(currentOffset, io.SeekStart); seekErr == nil {
					c.lock.Lock()
					chunkSize := c.chunkSize
					handler := c.progressHandler
					c.lock.Unlock()
					buffer := make([]byte, chunkSize)
					done := false
					for !done && err == nil {
						if readCount, readErr := file.Read(buffer); readCount > 0 {
							chunkData := buffer[:readCount]
							attempt := 0
							chunkUploaded := false
							for !chunkUploaded && err == nil {
								if newOffset, uploadErr := c.UploadChunk(uploadUrl, currentOffset, chunkData); uploadErr == nil {
									currentOffset = newOffset
									chunkUploaded = true
									if handler != nil {
										handler(currentOffset, total)
									}
								} else {
									attempt++
									if attempt > c.retryMax {
										err = uploadErr
									} else {
										time.Sleep(c.retryInterval)
									}
								}
							}
							if readErr == io.EOF {
								done = true
							}
						} else if readErr == io.EOF {
							done = true
						} else if readErr != nil {
							err = readErr
						} else {
							done = true
						}
					}
				} else {
					err = seekErr
				}
			} else {
				err = offsetErr
			}
			_ = file.Close()
		} else {
			err = openErr
		}
	}
	return err
}

func (c *TUSD_CLIENT) UploadFromReader(uploadUrl string, reader io.Reader, total int64) error {
	err := error(nil)
	if uploadUrl == "" {
		err = fmt.Errorf("upload URL cannot be empty")
	} else if reader == nil {
		err = fmt.Errorf("reader cannot be nil")
	} else {
		if currentOffset, _, offsetErr := c.GetOffset(uploadUrl); offsetErr == nil {
			c.lock.Lock()
			chunkSize := c.chunkSize
			handler := c.progressHandler
			c.lock.Unlock()
			buffer := make([]byte, chunkSize)
			done := false
			for !done && err == nil {
				if readCount, readErr := reader.Read(buffer); readCount > 0 {
					chunkData := buffer[:readCount]
					attempt := 0
					chunkUploaded := false
					for !chunkUploaded && err == nil {
						if newOffset, uploadErr := c.UploadChunk(uploadUrl, currentOffset, chunkData); uploadErr == nil {
							currentOffset = newOffset
							chunkUploaded = true
							if handler != nil {
								handler(currentOffset, total)
							}
						} else {
							attempt++
							if attempt > c.retryMax {
								err = uploadErr
							} else {
								time.Sleep(c.retryInterval)
							}
						}
					}
					if readErr == io.EOF {
						done = true
					}
				} else if readErr == io.EOF {
					done = true
				} else if readErr != nil {
					err = readErr
				} else {
					done = true
				}
			}
		} else {
			err = offsetErr
		}
	}
	return err
}

func (c *TUSD_CLIENT) resolveLocation(location string) string {
	result := location
	if !strings.HasPrefix(location, "http://") && !strings.HasPrefix(location, "https://") {
		baseUrl := strings.TrimRight(c.baseUrl, "/")
		if strings.HasPrefix(location, "/") {
			if schemeIndex := strings.Index(baseUrl, "://"); schemeIndex >= 0 {
				if pathIndex := strings.Index(baseUrl[schemeIndex+3:], "/"); pathIndex >= 0 {
					result = baseUrl[:schemeIndex+3+pathIndex] + location
				} else {
					result = baseUrl + location
				}
			} else {
				result = baseUrl + location
			}
		} else {
			result = baseUrl + "/" + location
		}
	}
	return result
}
