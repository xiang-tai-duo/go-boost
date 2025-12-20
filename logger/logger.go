// Package logger
// File:        logger.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/logger/logger.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: Logger is a utility for recording and managing application performance metrics and logs.
// --------------------------------------------------------------------------------
package logger

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/xiang-tai-duo/go-boost/debugger"
)

//goland:noinspection GoSnakeCaseUsage
const (
	SECRET_KEY       = "X7K!2p9@A4z8#F1m5$C3b7%N6D2^q8&W0e9*R4t7(Y2u5)I8o1P3"
	LOG_LEVEL_DEBUG  = 0
	LOG_LEVEL_INFO   = 1
	LOG_LEVEL_WARN   = 2
	LOG_LEVEL_ERROR  = 3
	LOG_LEVEL_SECRET = 4
)

//goland:noinspection GoSnakeCaseUsage
const (
	OUTPUT_MODE_CONSOLE_ONLY = 0
	OUTPUT_MODE_FILE_ONLY    = 1
	OUTPUT_MODE_BOTH         = 2
)

type (
	LOGGER struct {
		logger            *log.Logger
		mutex             sync.Mutex
		FolderName        string
		ColorReset        string
		Red               string
		Green             string
		Yellow            string
		Blue              string
		Magenta           string
		Cyan              string
		White             string
		Gray              string
		RedBackground     string
		GreenBackground   string
		YellowBackground  string
		BlueBackground    string
		MagentaBackground string
		CyanBackground    string
		WhiteBackground   string
		CurrentLevel      int
		OutputMode        int
	}
)

var (
	logLevelNames = map[int]string{
		LOG_LEVEL_DEBUG:  "DEBUG",
		LOG_LEVEL_INFO:   "INFO",
		LOG_LEVEL_WARN:   "WARN",
		LOG_LEVEL_ERROR:  "ERROR",
		LOG_LEVEL_SECRET: "SECRET",
	}
	Logger = &LOGGER{}
)

func init() {
	Logger.FolderName = "Logs"
	Logger.ColorReset = "\033[0m"
	Logger.Red = "\033[31m"
	Logger.Green = "\033[92m"
	Logger.Yellow = "\033[33m"
	Logger.Blue = "\033[34m"
	Logger.Magenta = "\033[35m"
	Logger.Cyan = "\033[36m"
	Logger.White = "\033[97m"
	Logger.Gray = "\033[37;2m"
	Logger.RedBackground = "\033[41m"
	Logger.GreenBackground = "\033[42m"
	Logger.YellowBackground = "\033[43m"
	Logger.BlueBackground = "\033[44m"
	Logger.MagentaBackground = "\033[45m"
	Logger.CyanBackground = "\033[46m"
	Logger.WhiteBackground = "\033[47m"
	Logger.CurrentLevel = LOG_LEVEL_INFO
	Logger.OutputMode = OUTPUT_MODE_BOTH
}

func (logger *LOGGER) Debug(message string) {
	logger.output(LOG_LEVEL_DEBUG, message, 2)
}

func (logger *LOGGER) DecryptSecretLogs(logFilePath string) ([]string, error) {
	var decryptedLogs []string
	var resultErr error

	file, err := os.Open(logFilePath)
	if err != nil {
		resultErr = err
	} else {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		content, err := io.ReadAll(file)
		if err != nil {
			resultErr = err
		} else {
			logPattern := regexp.MustCompile(`\[SECRET]]>\s*([A-Za-z0-9+/=]+)`)
			matches := logPattern.FindAllSubmatch(content, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					encryptedBase64 := string(match[1])
					decryptedMsg := logger.decryptSecret(encryptedBase64)
					decryptedLogs = append(decryptedLogs, decryptedMsg)
				}
			}
		}
	}

	return decryptedLogs, resultErr
}

func (logger *LOGGER) Error(message interface{}, skipStackFrames ...int) {
	var logMessage string
	var skip bool
	var skipCount int
	if len(skipStackFrames) > 0 {
		skipCount = skipStackFrames[0]
	}
	switch v := message.(type) {
	case error:
		if v != nil {
			var syscallErrno syscall.Errno
			if errors.As(v, &syscallErrno) && syscallErrno == 0 {
				skip = true
			}
			var exitError *exec.ExitError
			if errors.As(v, &exitError) {
				skip = true
			}
			if !skip {
				logMessage = v.Error()
			}
		}
	case string:
		logMessage = v
	default:
		logMessage = fmt.Sprint(v)
	}
	if logMessage != "" {
		logger.output(LOG_LEVEL_ERROR, logMessage, 2+skipCount)
	}
}

func (logger *LOGGER) GetGoroutineID() int {
	var buffer [64]byte
	var goroutineID int
	goroutineID = -1
	length := runtime.Stack(buffer[:], false)
	for i := 0; i < length-10; i++ {
		if string(buffer[i:i+9]) == "goroutine " {
			j := i + 9
			for ; j < length && buffer[j] >= '0' && buffer[j] <= '9'; j++ {
				goroutineID = goroutineID*10 + int(buffer[j]-'0')
			}
			break
		}
	}
	return goroutineID
}

func (logger *LOGGER) Info(message string) {
	logger.output(LOG_LEVEL_INFO, message, 2)
}

func (logger *LOGGER) Secret(message string) {
	logger.output(LOG_LEVEL_SECRET, message, 2)
}

func (logger *LOGGER) Warning(message string) {
	logger.output(LOG_LEVEL_WARN, message, 2)
}

func (logger *LOGGER) decryptSecret(message string) string {
	result := message
	key := []byte(SECRET_KEY)
	ciphertext, err := base64.StdEncoding.DecodeString(message)
	if err == nil {
		block, err := aes.NewCipher(key)
		if err == nil {
			aesGCM, err := cipher.NewGCM(block)
			if err == nil {
				nonceSize := aesGCM.NonceSize()
				if len(ciphertext) >= nonceSize {
					nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
					plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
					if err == nil {
						result = string(plaintext)
					}
				}
			}
		}
	}
	return result
}

func (logger *LOGGER) encryptSecret(message string) string {
	result := message
	key := []byte(SECRET_KEY)
	block, err := aes.NewCipher(key)
	if err == nil {
		aesGCM, err := cipher.NewGCM(block)
		if err == nil {
			nonce := make([]byte, aesGCM.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
				ciphertext := aesGCM.Seal(nonce, nonce, []byte(message), nil)
				result = base64.StdEncoding.EncodeToString(ciphertext)
			}
		}
	}
	return result
}

func (logger *LOGGER) output(level int, message string, skipStackFrames ...int) {
	if message != "" && level >= logger.CurrentLevel {
		logger.mutex.Lock()
		var skipCount int
		skipCount = 1
		if len(skipStackFrames) > 0 {
			skipCount = skipStackFrames[0]
		}
		pc, filename, line, _ := runtime.Caller(skipCount)
		parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		funcName := parts[len(parts)-1]
		levelName := logLevelNames[level]
		baseColor := logger.White
		levelTextColor := baseColor
		messageColor := baseColor
		logMessage := message
		switch level {
		case LOG_LEVEL_DEBUG:
			levelTextColor = logger.Gray
			messageColor = logger.Gray
		case LOG_LEVEL_INFO:
			levelTextColor = logger.Green
			messageColor = logger.White
		case LOG_LEVEL_WARN:
			levelTextColor = logger.Yellow
			messageColor = logger.Yellow
		case LOG_LEVEL_ERROR:
			levelTextColor = logger.Red
			messageColor = logger.Red
		case LOG_LEVEL_SECRET:
			levelTextColor = logger.Magenta
			messageColor = logger.Magenta
			logMessage = "******"
		}
		consoleText := fmt.Sprintf(
			"[%s][%s::%s::%d][Goroutine:%d][%s%s%s]> %s%s%s",
			time.Now().Format("2006/01/02 15:04:05"),
			filepath.Base(filename),
			funcName,
			line,
			logger.GetGoroutineID(),
			levelTextColor,
			levelName,
			logger.ColorReset,
			messageColor,
			logMessage,
			logger.ColorReset,
		)
		fileText := fmt.Sprintf(
			"[%s][%s::%s::%d][Goroutine:%d][%s]> %s",
			time.Now().Format("2006/01/02 15:04:05"),
			filepath.Base(filename),
			funcName,
			line,
			logger.GetGoroutineID(),
			levelName,
			message,
		)
		if level == LOG_LEVEL_SECRET {
			encryptedMsg := logger.encryptSecret(message)
			fileText = fmt.Sprintf(
				"[%s][%s::%s::%d][Goroutine:%d][%s]> %s",
				time.Now().Format("2006/01/02 15:04:05"),
				filepath.Base(filename),
				funcName,
				line,
				logger.GetGoroutineID(),
				levelName,
				encryptedMsg,
			)
		}
		if runtime.GOOS == "windows" && !debugger.IsPresent() {
			consoleText = fileText
		}

		switch logger.OutputMode {
		case OUTPUT_MODE_CONSOLE_ONLY:
			fmt.Println(consoleText)
		case OUTPUT_MODE_FILE_ONLY:
			logger.writeFile(fileText)
		case OUTPUT_MODE_BOTH:
			fmt.Println(consoleText)
			logger.writeFile(fileText)
		}

		logger.mutex.Unlock()
	}
}

func (logger *LOGGER) writeFile(logText string) {
	var err error
	if err = os.MkdirAll(logger.FolderName, os.ModePerm); err == nil {
		logFilePath := path.Join(logger.FolderName, time.Now().Format("20060102")) + ".log"
		var file *os.File
		if file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm); err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			fileLogger := log.New(file, "", log.Ldate|log.Ltime)
			fileLogger.Println(logText)
		}
	}
}
