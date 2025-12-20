// Package lpr
// File:        lpr.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/lpr/lpr.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: LPR provides functions to send print jobs to network printers using the LPR/LPD protocol (RFC 1179) and raw 9100 port.
// --------------------------------------------------------------------------------
package lpr

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xiang-tai-duo/go-boost/logger"
	"github.com/xiang-tai-duo/go-boost/spl"
	"github.com/xiang-tai-duo/go-boost/system"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	ACK_BUFFER_SIZE                          = 1
	CONTROL_BYTE_NULL                        = 0x00
	CONTROL_BYTE_RECEIVE_CONTROL             = 0x02
	CONTROL_BYTE_RECEIVE_DATA                = 0x03
	CONTROL_FILE_PREFIX                      = "cfA000"
	CONTROL_TAG_FILE_NAME                    = "N"
	CONTROL_TAG_HOST_NAME                    = "H"
	CONTROL_TAG_PRINT_DATA_FILE              = "f"
	CONTROL_TAG_UNLINK_DATA_FILE             = "U"
	CONTROL_TAG_USER_NAME                    = "P"
	CONNECT_RETRY_COUNT                      = 2
	CONNECT_RETRY_INTERVAL                   = 1 * time.Second
	CARRIAGE_RETURN                          = "\r"
	CARRIAGE_RETURN_LINE_FEED                = "\r\n"
	DATA_FILE_PREFIX                         = "dfA000"
	DEFAULT_JOB_NAME                         = "print_job"
	DEFAULT_PORT                             = 515
	DEFAULT_QUEUE                            = "lp"
	DEFAULT_TIMEOUT                          = 30 * time.Second
	FILE_PERMISSION                          = 0644
	HTTP_RESPONSE_PARTS_COUNT                = 3
	HTTP_RESPONSE_STATUS_PARTS_MIN           = 2
	HTTP_SCHEME                              = "http"
	HTTPS_SCHEME                             = "https"
	LINE_FEED                                = "\n"
	MAX_TRANSMIT_SIZE                        = 16384
	MODULE_NAME_LPR                          = "lpr"
	NETWORK_TCP                              = "tcp"
	NULL_TERMINATOR                          = "\x00"
	PROGRESS_LOG_CHUNK_INTERVAL              = 64
	PROXY_ENV_VAR_HTTP                       = "http_proxy"
	PROXY_ENV_VAR_HTTP_UPPER                 = "HTTP_PROXY"
	PROXY_ENV_VAR_HTTPS                      = "https_proxy"
	PROXY_ENV_VAR_HTTPS_UPPER                = "HTTPS_PROXY"
	PROXY_HEADER_PARSE_STATE_CR_AFTER_CRLF   = 2
	PROXY_HEADER_PARSE_STATE_CRLF_AFTER_CRLF = 3
	PROXY_HEADER_PARSE_STATE_DONE            = 4
	PROXY_HEADER_PARSE_STATE_INITIAL         = 0
	PROXY_HEADER_PARSE_STATE_LF_AFTER_CR     = 1
	PROXY_HTTP_CONNECT_FORMAT                = "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n"
	PROXY_HTTP_RESPONSE_OK                   = "200"
	PROXY_RESPONSE_BUFFER_SIZE               = 512
	RAW_PORT                                 = 9100
	SEND_RETRY_COUNT                         = 0
	SEND_RETRY_INTERVAL                      = 1 * time.Second
	SPACE                                    = " "
	TEMP_FILE_NAME_FORMAT                    = "lpr_%s_%s"
	URL_SCHEME_SEPARATOR                     = "://"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_LPR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_LPR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_LPR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_LPR, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult,SpellCheckingInspection
func Send(ipAddress string, filePath string, jobName *string, queueName *string, timeout *time.Duration) (string, error) {
	err := error(nil)
	jcid := ""
	__debug(fmt.Sprintf("ipAddress=%s, filePath=%s", ipAddress, filePath))
	if filePath == "" {
		err = fmt.Errorf("file path is empty")
		__debug(fmt.Sprintf("validation failed: %v", err))
	} else {
		newPRNFilePath := filePath
		if jcid, newPRNFilePath, err = insertJCID(filePath); err == nil {
			filePath = newPRNFilePath
		} else {
			__debug(fmt.Sprintf("failed to insert JCID, continuing without it: %v", err))
			err = nil
		}
		fileName := filepath.Base(filePath)
		effectiveJobName := DEFAULT_JOB_NAME
		if jobName != nil && *jobName != "" {
			effectiveJobName = *jobName
		}
		effectiveQueue := DEFAULT_QUEUE
		if queueName != nil && *queueName != "" {
			effectiveQueue = *queueName
		}
		actualTimeout := DEFAULT_TIMEOUT
		if timeout != nil {
			actualTimeout = *timeout
		}
		dataFileName := DATA_FILE_PREFIX + uuid.New().String()
		controlContent := buildControlFileContent(dataFileName, fileName, effectiveJobName)
		address := fmt.Sprintf("%s:%d", ipAddress, DEFAULT_PORT)
		var connection net.Conn
		__debug(fmt.Sprintf("connecting to LPR: address=%s, timeout=%s", address, actualTimeout))
		if connection, err = dialWithRetry(NETWORK_TCP, address, actualTimeout); err == nil {
			__debug(fmt.Sprintf("connected to LPR: address=%s", address))
			defer connection.Close()
			if err = sendReceivePrinterJob(connection, effectiveQueue, actualTimeout); err == nil {
				if err = sendDataFileSubcommand(connection, dataFileName, filePath, actualTimeout); err == nil {
					if err = sendControlFileSubcommand(connection, dataFileName, controlContent, actualTimeout); err == nil {
						__debug("done: success")
					}
				}
			}
			if err != nil {
				__debug(fmt.Sprintf("LPR operation failed: %v", err))
			}
		} else {
			__debug(fmt.Sprintf("connect failed: address=%s, error=%v", address, err))
		}
	}
	__debug(fmt.Sprintf("result: jcid=%s, err=%v", jcid, err))
	return jcid, err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult,GoDfaConstantCondition
func SendReader(ipAddress string, reader io.Reader, queueName *string, timeout *time.Duration) error {
	err := error(nil)
	__debug(fmt.Sprintf("ipAddress=%s", ipAddress))
	effectiveQueue := DEFAULT_QUEUE
	if queueName != nil && *queueName != "" {
		effectiveQueue = *queueName
	}
	actualTimeout := DEFAULT_TIMEOUT
	if timeout != nil {
		actualTimeout = *timeout
	}
	dataFileName := DATA_FILE_PREFIX + uuid.New().String()
	controlContent := buildControlFileContent(dataFileName, "", DEFAULT_JOB_NAME)
	address := fmt.Sprintf("%s:%d", ipAddress, DEFAULT_PORT)
	var connection net.Conn
	__debug(fmt.Sprintf("connecting to LPR (reader): address=%s, timeout=%s", address, actualTimeout))
	if connection, err = dialWithRetry(NETWORK_TCP, address, actualTimeout); err == nil {
		__debug(fmt.Sprintf("connected to LPR: address=%s", address))
		defer connection.Close()
		if err = sendReceivePrinterJob(connection, effectiveQueue, actualTimeout); err == nil {
			attemptCount := 0
			for {
				__debug(fmt.Sprintf("sending data file (reader): filename=%s, timeout=%s, attempt=%d/%d", dataFileName, actualTimeout, attemptCount, SEND_RETRY_COUNT))
				if err = sendDataFileSubcommandFromReader(connection, dataFileName, reader, actualTimeout); err == nil {
					break
				}
				if err == nil || attemptCount >= SEND_RETRY_COUNT {
					__debug(fmt.Sprintf("send failed and retry limit reached: address=%s, retry=%d/%d, error=%v", address, attemptCount, SEND_RETRY_COUNT, err))
					break
				}
				attemptCount++
				__debug(fmt.Sprintf("send failed, retrying: address=%s, retry=%d/%d, interval=%s, error=%v", address, attemptCount, SEND_RETRY_COUNT, SEND_RETRY_INTERVAL, err))
				time.Sleep(SEND_RETRY_INTERVAL)
				if connection, err = dialWithRetry(NETWORK_TCP, address, actualTimeout); err == nil {
					__debug(fmt.Sprintf("reconnected to LPR: address=%s", address))
					if err = sendReceivePrinterJob(connection, effectiveQueue, actualTimeout); err != nil {
						connection.Close()
						break
					}
				} else {
					__debug(fmt.Sprintf("reconnect failed: address=%s, error=%v", address, err))
					break
				}
			}
			if err == nil {
				if err = sendControlFileSubcommand(connection, dataFileName, controlContent, actualTimeout); err == nil {
					__debug("done: success")
				}
			}
		}
		if err != nil {
			__debug(fmt.Sprintf("LPR operation failed: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("connect failed: address=%s, error=%v", address, err))
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult,GoDfaConstantCondition
func SendReaderTo9100Port(ipAddress string, reader io.Reader) error {
	err := error(nil)
	__debug(fmt.Sprintf("ipAddress=%s", ipAddress))
	address := fmt.Sprintf("%s:%d", ipAddress, RAW_PORT)
	var connection net.Conn
	__debug(fmt.Sprintf("connecting to raw port: address=%s, timeout=%s", address, DEFAULT_TIMEOUT))
	if connection, err = dialWithRetry(NETWORK_TCP, address, DEFAULT_TIMEOUT); err == nil {
		__debug(fmt.Sprintf("connected to raw port: address=%s", address))
		defer connection.Close()
		buffer := make([]byte, MAX_TRANSMIT_SIZE)
		totalBytes := int64(0)
		chunkCount := 0
		for {
			var n int
			if n, err = reader.Read(buffer); err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
			chunkCount++
			attemptCount := 0
			for {
				__debug(fmt.Sprintf("sending data chunk: bytes=%d, total=%d, attempt=%d/%d", n, totalBytes+int64(n), attemptCount, SEND_RETRY_COUNT))
				if _, err = connection.Write(buffer[:n]); err == nil {
					totalBytes += int64(n)
					if chunkCount%PROGRESS_LOG_CHUNK_INTERVAL == 0 {
						__debug(fmt.Sprintf("send progress: total_bytes=%d, chunks=%d", totalBytes, chunkCount))
					}
					break
				}
				if attemptCount >= SEND_RETRY_COUNT {
					__debug(fmt.Sprintf("send failed and retry limit reached: address=%s, retry=%d/%d, error=%v", address, attemptCount, SEND_RETRY_COUNT, err))
					break
				}
				attemptCount++
				__debug(fmt.Sprintf("send failed, retrying: address=%s, retry=%d/%d, interval=%s, error=%v", address, attemptCount, SEND_RETRY_COUNT, SEND_RETRY_INTERVAL, err))
				time.Sleep(SEND_RETRY_INTERVAL)
			}
			if err != nil {
				break
			}
		}
		if err == nil {
			__debug(fmt.Sprintf("done: success, total_bytes=%d, chunks=%d", totalBytes, chunkCount))
		} else {
			__debug(fmt.Sprintf("raw port send failed: %v", err))
		}
	} else {
		__debug(fmt.Sprintf("connect failed: address=%s, error=%v", address, err))
	}
	return err
}

//goland:noinspection GoUnusedExportedFunction,GoUnhandledErrorResult,SpellCheckingInspection
func SendTo9100Port(ipAddress string, filePath string) (string, error) {
	err := error(nil)
	jcid := ""
	__debug(fmt.Sprintf("ipAddress=%s, filePath=%s", ipAddress, filePath))
	newPRNFilePath := filePath
	if jcid, newPRNFilePath, err = insertJCID(filePath); err == nil {
		var file *os.File
		if file, err = os.Open(newPRNFilePath); err == nil {
			defer file.Close()
			err = SendReaderTo9100Port(ipAddress, file)
		}
	} else {
		__debug(fmt.Sprintf("failed to insert JCID, continuing without it: %v", err))
		err = nil
		var file *os.File
		if file, err = os.Open(filePath); err == nil {
			defer file.Close()
			err = SendReaderTo9100Port(ipAddress, file)
		}
	}
	return jcid, err
}

func buildControlFileContent(dataFileName string, fileName string, jobName string) string {
	result := ""
	hostName, _ := os.Hostname()
	result += CONTROL_TAG_HOST_NAME + hostName + LINE_FEED
	result += CONTROL_TAG_USER_NAME + CONTROL_FILE_PREFIX + LINE_FEED
	result += CONTROL_TAG_FILE_NAME + jobName + LINE_FEED
	result += CONTROL_TAG_PRINT_DATA_FILE + dataFileName + LINE_FEED
	result += CONTROL_TAG_FILE_NAME + fileName + LINE_FEED
	result += CONTROL_TAG_UNLINK_DATA_FILE + dataFileName + LINE_FEED
	__debug(fmt.Sprintf("done: length=%d", len(result)))
	return result
}

func getProxyURL() string {
	result := ""
	if system.IsUnix() {
		proxyURL := os.Getenv(PROXY_ENV_VAR_HTTP)
		if proxyURL == "" {
			proxyURL = os.Getenv(PROXY_ENV_VAR_HTTP_UPPER)
		}
		if proxyURL == "" {
			proxyURL = os.Getenv(PROXY_ENV_VAR_HTTPS)
		}
		if proxyURL == "" {
			proxyURL = os.Getenv(PROXY_ENV_VAR_HTTPS_UPPER)
		}
		result = proxyURL
	}
	return result
}

//goland:noinspection SpellCheckingInspection
func insertJCID(filePath string) (string, string, error) {
	err := error(nil)
	jcid := ""
	resultPath := filePath
	var info os.FileInfo
	if info, err = os.Stat(filePath); err == nil {
		dir := filepath.Dir(filePath)
		base := filepath.Base(filePath)
		ext := filepath.Ext(base)
		jcid = strings.ReplaceAll(uuid.New().String(), "-", "")
		newName := fmt.Sprintf(TEMP_FILE_NAME_FORMAT, jcid, base)
		resultPath = filepath.Join(dir, newName)
		if err = copyFile(filePath, resultPath); err == nil {
			if err = spl.SetJCID(resultPath, jcid); err == nil {
				__debug(fmt.Sprintf("JCID inserted: jcid=%s, newFile=%s", jcid, resultPath))
			}
		}
		if err != nil {
			resultPath = filePath
		}
		_ = ext
		_ = info
	}
	return jcid, resultPath, err
}

func copyFile(src string, dst string) error {
	err := error(nil)
	var data []byte
	if data, err = os.ReadFile(src); err == nil {
		err = os.WriteFile(dst, data, FILE_PERMISSION)
	}
	return err
}

func sendControlFileSubcommand(connection net.Conn, dataFileName string, controlContent string, timeout time.Duration) error {
	err := error(nil)
	controlBytes := []byte(controlContent)
	command := fmt.Sprintf("%c%d %s%s\n", CONTROL_BYTE_RECEIVE_CONTROL, len(controlBytes), CONTROL_FILE_PREFIX, dataFileName)
	__debug(fmt.Sprintf("sending control file subcommand: length=%d", len(controlBytes)))
	if err = sendWithTimeout(connection, []byte(command), timeout); err == nil {
		if err = readAckWithTimeout(connection, timeout); err == nil {
			chunks := 0
			totalWritten := 0
			for i := 0; i < len(controlBytes); i += MAX_TRANSMIT_SIZE {
				end := i + MAX_TRANSMIT_SIZE
				if end > len(controlBytes) {
					end = len(controlBytes)
				}
				chunk := controlBytes[i:end]
				if _, err = connection.Write(chunk); err != nil {
					break
				}
				totalWritten += len(chunk)
				chunks++
			}
			if err == nil {
				nullByte := []byte{CONTROL_BYTE_NULL}
				if _, err = connection.Write(nullByte); err == nil {
					if err = readAckWithTimeout(connection, timeout); err == nil {
						__debug(fmt.Sprintf("control file sent: total=%d bytes, chunks=%d", totalWritten, chunks))
					}
				}
			}
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult,GoDfaConstantCondition
func sendDataFileSubcommand(connection net.Conn, dataFileName string, filePath string, timeout time.Duration) error {
	err := error(nil)
	var file *os.File
	if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		var fileInfo os.FileInfo
		if fileInfo, err = file.Stat(); err == nil {
			fileSize := fileInfo.Size()
			command := fmt.Sprintf("%c%d %s\n", CONTROL_BYTE_RECEIVE_DATA, fileSize, dataFileName)
			__debug(fmt.Sprintf("sending data file subcommand: size=%d, file=%s", fileSize, dataFileName))
			if err = sendWithTimeout(connection, []byte(command), timeout); err == nil {
				if err = readAckWithTimeout(connection, timeout); err == nil {
					buffer := make([]byte, MAX_TRANSMIT_SIZE)
					totalWritten := int64(0)
					chunks := 0
					for {
						var n int
						if n, err = file.Read(buffer); err != nil {
							if err == io.EOF {
								err = nil
							}
							break
						}
						attemptCount := 0
						for {
							if _, err = connection.Write(buffer[:n]); err == nil {
								totalWritten += int64(n)
								chunks++
								if chunks%PROGRESS_LOG_CHUNK_INTERVAL == 0 {
									__debug(fmt.Sprintf("data file progress: total=%d/%d bytes, chunks=%d", totalWritten, fileSize, chunks))
								}
								break
							}
							if err == nil || attemptCount >= SEND_RETRY_COUNT {
								__debug(fmt.Sprintf("data file send failed: total=%d/%d, retry=%d/%d, error=%v", totalWritten, fileSize, attemptCount, SEND_RETRY_COUNT, err))
								break
							}
							attemptCount++
							__debug(fmt.Sprintf("data file send failed, retrying: total=%d/%d, retry=%d/%d, interval=%s, error=%v", totalWritten, fileSize, attemptCount, SEND_RETRY_COUNT, SEND_RETRY_INTERVAL, err))
							time.Sleep(SEND_RETRY_INTERVAL)
						}
						if err != nil {
							break
						}
					}
					if err == nil {
						nullByte := []byte{CONTROL_BYTE_NULL}
						if _, err = connection.Write(nullByte); err == nil {
							if err = readAckWithTimeout(connection, timeout); err == nil {
								__debug(fmt.Sprintf("data file sent: total=%d/%d bytes, chunks=%d", totalWritten, fileSize, chunks))
							}
						}
					}
				}
			}
		}
	}
	return err
}

func sendDataFileSubcommandFromReader(connection net.Conn, dataFileName string, reader io.Reader, timeout time.Duration) error {
	err := error(nil)
	command := fmt.Sprintf("%c0 %s\n", CONTROL_BYTE_RECEIVE_DATA, dataFileName)
	__debug(fmt.Sprintf("sending data file subcommand (reader, size unknown): file=%s", dataFileName))
	if err = sendWithTimeout(connection, []byte(command), timeout); err == nil {
		if err = readAckWithTimeout(connection, timeout); err == nil {
			buffer := make([]byte, MAX_TRANSMIT_SIZE)
			totalWritten := int64(0)
			chunks := 0
			for {
				var n int
				if n, err = reader.Read(buffer); err != nil {
					if err == io.EOF {
						err = nil
					}
					break
				}
				chunks++
				__debug(fmt.Sprintf("sending data chunk: bytes=%d, total=%d", n, totalWritten+int64(n)))
				if _, err = connection.Write(buffer[:n]); err == nil {
					totalWritten += int64(n)
				} else {
					break
				}
			}
			if err == nil {
				nullByte := []byte{CONTROL_BYTE_NULL}
				if _, err = connection.Write(nullByte); err == nil {
					if err = readAckWithTimeout(connection, timeout); err == nil {
						__debug(fmt.Sprintf("data file sent (reader): total=%d bytes, chunks=%d", totalWritten, chunks))
					}
				}
			}
		}
	}
	return err
}

func sendReceivePrinterJob(connection net.Conn, queueName string, timeout time.Duration) error {
	err := error(nil)
	command := fmt.Sprintf("%c%s\n", CONTROL_BYTE_RECEIVE_CONTROL, queueName)
	__debug(fmt.Sprintf("sending receive printer job: queue=%s", queueName))
	if err = sendWithTimeout(connection, []byte(command), timeout); err == nil {
		err = readAckWithTimeout(connection, timeout)
	}
	return err
}

func sendWithTimeout(connection net.Conn, data []byte, timeout time.Duration) error {
	err := error(nil)
	if err = connection.SetWriteDeadline(time.Now().Add(timeout)); err == nil {
		_, err = connection.Write(data)
		if err != nil {
			__debug(fmt.Sprintf("write failed: error=%v", err))
		}
	}
	return err
}

func readAckWithTimeout(connection net.Conn, timeout time.Duration) error {
	err := error(nil)
	if err = connection.SetReadDeadline(time.Now().Add(timeout)); err == nil {
		buffer := make([]byte, ACK_BUFFER_SIZE)
		var n int
		if n, err = connection.Read(buffer); err == nil {
			if n == 1 && buffer[0] == 0 {
				__debug("ACK received")
			} else {
				err = fmt.Errorf("unexpected ACK response: %d (expected 0x00)", buffer[0])
			}
		} else {
			__debug(fmt.Sprintf("read ACK failed: error=%v", err))
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult,GoDfaConstantCondition
func dialWithProxy(proxyURL string, targetAddress string, timeout time.Duration) (net.Conn, error) {
	result := net.Conn(nil)
	err := error(nil)
	var proxyAddr string
	var parsedURL *url.URL
	if parsedURL, err = url.Parse(proxyURL); err == nil {
		if parsedURL.Scheme == "" {
			proxyURL = HTTP_SCHEME + URL_SCHEME_SEPARATOR + proxyURL
			if parsedURL, err = url.Parse(proxyURL); err != nil {
				__debug(fmt.Sprintf("proxy URL re-parse failed: %s, error=%v", proxyURL, err))
			}
		}
		if err == nil {
			if parsedURL.Scheme == HTTP_SCHEME || parsedURL.Scheme == HTTPS_SCHEME {
				proxyAddr = parsedURL.Host
			}
			if proxyAddr != "" {
				var conn net.Conn
				if conn, err = net.DialTimeout(NETWORK_TCP, proxyAddr, timeout); err == nil {
					closed := false
					defer func() {
						if !closed {
							conn.Close()
						}
					}()
					__debug(fmt.Sprintf("proxy connected: proxy=%s", proxyAddr))
					connectReq := fmt.Sprintf(PROXY_HTTP_CONNECT_FORMAT, targetAddress, targetAddress)
					if err = conn.SetWriteDeadline(time.Now().Add(timeout)); err == nil {
						if _, err = conn.Write([]byte(connectReq)); err == nil {
							__debug(fmt.Sprintf("proxy CONNECT request sent: target=%s", targetAddress))
							if err = conn.SetReadDeadline(time.Now().Add(timeout)); err == nil {
								response := make([]byte, 0, PROXY_RESPONSE_BUFFER_SIZE)
								buffer := make([]byte, PROXY_RESPONSE_BUFFER_SIZE)
								state := PROXY_HEADER_PARSE_STATE_INITIAL
								for state < PROXY_HEADER_PARSE_STATE_DONE {
									var n int
									if n, err = conn.Read(buffer); err != nil {
										break
									}
									for i := 0; i < n && state < PROXY_HEADER_PARSE_STATE_DONE; i++ {
										b := buffer[i]
										response = append(response, b)
										if state == PROXY_HEADER_PARSE_STATE_INITIAL {
											if b == CARRIAGE_RETURN[0] {
												state = PROXY_HEADER_PARSE_STATE_LF_AFTER_CR
											}
										} else if state == PROXY_HEADER_PARSE_STATE_LF_AFTER_CR {
											if b == LINE_FEED[0] {
												state = PROXY_HEADER_PARSE_STATE_CR_AFTER_CRLF
											} else if b == CARRIAGE_RETURN[0] {
												state = PROXY_HEADER_PARSE_STATE_LF_AFTER_CR
											} else {
												state = PROXY_HEADER_PARSE_STATE_INITIAL
											}
										} else if state == PROXY_HEADER_PARSE_STATE_CR_AFTER_CRLF {
											if b == CARRIAGE_RETURN[0] {
												state = PROXY_HEADER_PARSE_STATE_CRLF_AFTER_CRLF
											} else {
												state = PROXY_HEADER_PARSE_STATE_INITIAL
											}
										} else if state == PROXY_HEADER_PARSE_STATE_CRLF_AFTER_CRLF {
											if b == LINE_FEED[0] {
												state = PROXY_HEADER_PARSE_STATE_DONE
											} else if b == CARRIAGE_RETURN[0] {
												state = PROXY_HEADER_PARSE_STATE_LF_AFTER_CR
											} else {
												state = PROXY_HEADER_PARSE_STATE_INITIAL
											}
										}
									}
								}
								if err == nil {
									statusLine := string(response)
									if idx := strings.Index(statusLine, CARRIAGE_RETURN_LINE_FEED); idx != -1 {
										statusLine = statusLine[:idx]
									}
									parts := strings.SplitN(statusLine, SPACE, HTTP_RESPONSE_PARTS_COUNT)
									if len(parts) >= HTTP_RESPONSE_STATUS_PARTS_MIN && parts[1] == PROXY_HTTP_RESPONSE_OK {
										conn.SetReadDeadline(time.Time{})
										conn.SetWriteDeadline(time.Time{})
										closed = true
										result = conn
										__debug(fmt.Sprintf("proxy tunnel established: target=%s, response=%s", targetAddress, statusLine))
									} else {
										err = fmt.Errorf("proxy tunnel failed: %s", statusLine)
									}
								} else {
									__debug(fmt.Sprintf("proxy response read failed: target=%s, error=%v", targetAddress, err))
								}
							}
						}
					}
				} else {
					__debug(fmt.Sprintf("proxy dial failed: proxy=%s, error=%v", proxyAddr, err))
				}
			} else {
				err = fmt.Errorf("proxy URL host is empty: %s", proxyURL)
			}
		} else {
			__debug(fmt.Sprintf("proxy URL parse failed: %s, error=%v", proxyURL, err))
		}
	}
	return result, err
}

//goland:noinspection GoUnhandledErrorResult
func dialWithRetry(network string, address string, timeout time.Duration) (net.Conn, error) {
	result := net.Conn(nil)
	err := error(nil)
	attemptCount := 0
	proxyURL := getProxyURL()
	for {
		if proxyURL != "" {
			__debug(fmt.Sprintf("using proxy for connection: proxy=%s, target=%s", proxyURL, address))
			if result, err = dialWithProxy(proxyURL, address, timeout); err == nil {
				if attemptCount > 0 {
					__debug(fmt.Sprintf("connect retry succeeded via proxy: address=%s, retry=%d/%d", address, attemptCount, CONNECT_RETRY_COUNT))
				}
				break
			}
			__debug(fmt.Sprintf("proxy dial failed: address=%s, error=%v", address, err))
		} else {
			__debug(fmt.Sprintf("dialing TCP: network=%s, address=%s, timeout=%s, attempt=%d/%d", network, address, timeout, attemptCount, CONNECT_RETRY_COUNT))
			if result, err = net.DialTimeout(network, address, timeout); err == nil {
				if attemptCount > 0 {
					__debug(fmt.Sprintf("connect retry succeeded: address=%s, retry=%d/%d", address, attemptCount, CONNECT_RETRY_COUNT))
				}
				break
			}
		}
		if attemptCount >= CONNECT_RETRY_COUNT {
			__debug(fmt.Sprintf("connect failed and retry limit reached: address=%s, retry=%d/%d, error=%v", address, attemptCount, CONNECT_RETRY_COUNT, err))
			break
		}
		attemptCount++
		__debug(fmt.Sprintf("connect failed, retrying: address=%s, retry=%d/%d, interval=%s, error=%v", address, attemptCount, CONNECT_RETRY_COUNT, CONNECT_RETRY_INTERVAL, err))
		time.Sleep(CONNECT_RETRY_INTERVAL)
	}
	return result, err
}
