// --------------------------------------------------------------------------------
// File:        logger.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: Logger is a utility for recording and managing application performance metrics and logs.
// --------------------------------------------------------------------------------

package boost

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
)

//goland:noinspection GoSnakeCaseUsage
const (
	// SECRET_KEY is a 32-byte AES-256 key for encrypting secret logs
	SECRET_KEY = "12345678901234567890123456789012"
)

// Log level constants
const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelSecret
)

// LOGGER provides utility methods for recording and managing application logs.
type LOGGER struct {
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
}

// Log level name mapping
var logLevelNames = map[int]string{
	LogLevelDebug:  "DEBUG",
	LogLevelInfo:   "INFO",
	LogLevelWarn:   "WARN",
	LogLevelError:  "ERROR",
	LogLevelSecret: "SECRET",
}

// Logger Global logger instance
var (
	Logger = &LOGGER{}
)

func init() {
	Logger.FolderName = "Logs"
	Logger.ColorReset = "\033[0m"
	Logger.Red = "\033[31m"
	Logger.Green = "\033[92m" // bright green
	Logger.Yellow = "\033[33m"
	Logger.Blue = "\033[34m"
	Logger.Magenta = "\033[35m"
	Logger.Cyan = "\033[36m"
	Logger.White = "\033[97m"  // bright white
	Logger.Gray = "\033[37;2m" // light gray
	Logger.RedBackground = "\033[41m"
	Logger.GreenBackground = "\033[42m"
	Logger.YellowBackground = "\033[43m"
	Logger.BlueBackground = "\033[44m"
	Logger.MagentaBackground = "\033[45m"
	Logger.CyanBackground = "\033[46m"
	Logger.WhiteBackground = "\033[47m"
	Logger.CurrentLevel = LogLevelInfo // Default log level is Info
}

// Debug logs a debug message.
// s: Message to log
// Usage: Logger.Debug("Debug message")
func (logger *LOGGER) Debug(s string) {
	logger.output(LogLevelDebug, s, 2)
}

// DecryptSecretLogs scans log files and decrypts SECRET logs
// logFilePath: Path to log file
// returns: List of decrypted SECRET logs
func (logger *LOGGER) DecryptSecretLogs(logFilePath string) ([]string, error) {
	// Open log file
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Define regular expression to match SECRET log lines
	// Match format: [time][filename::function::line][Goroutine:ID][SECRET]> base64_encrypted_string
	logPattern := regexp.MustCompile(`\[SECRET]]>\s*([A-Za-z0-9+/=]+)`)

	// Find all matching log lines
	matches := logPattern.FindAllSubmatch(content, -1)

	// Decrypt each matching SECRET log
	var decryptedLogs []string
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// Extract base64 encrypted string
		encryptedBase64 := string(match[1])

		// Decrypt string
		decryptedMsg := logger.decryptSecret(encryptedBase64)
		decryptedLogs = append(decryptedLogs, decryptedMsg)
	}

	return decryptedLogs, nil
}

// Error logs an error message with stack trace information (Error level).
// msg: Error message or error object to log
// skipStackFrames: Number of stack frames to skip when logging
// Usage:
// Logger.Error("An error occurred")
// Logger.Error(errors.New("An error occurred"))
func (logger *LOGGER) Error(msg interface{}, skipStackFrames ...int) {
	var message string
	var skip bool
	var skipCount int

	if len(skipStackFrames) > 0 {
		skipCount = skipStackFrames[0]
	}

	switch v := msg.(type) {
	case error:
		if v != nil {
			// Check if it's a syscall.Errno with value 0, skip logging
			var syscallErrno syscall.Errno
			if errors.As(v, &syscallErrno) && syscallErrno == 0 {
				skip = true
			}
			// Check if it's an exec.ExitError, skip logging
			var exitError *exec.ExitError
			if errors.As(v, &exitError) {
				skip = true
			}
			if !skip {
				message = v.Error()
			}
		}
	case string:
		message = v
	default:
		message = fmt.Sprint(v)
	}

	if message != "" {
		logger.output(LogLevelError, message, 2+skipCount)
	}
}

// Info logs an info message.
// s: Message to log
// Usage: Logger.Info("Info message")
func (logger *LOGGER) Info(s string) {
	logger.output(LogLevelInfo, s, 2)
}

// Secret logs a secret message.
// s: Message to log
// Usage: Logger.Secret("Secret message")
func (logger *LOGGER) Secret(s string) {
	logger.output(LogLevelSecret, s, 2)
}

// Warning logs a warning message.
// s: Message to log
// Usage: Logger.Warning("Warning message")
func (logger *LOGGER) Warning(s string) {
	logger.output(LogLevelWarn, s, 2)
}

// internal AES decryption function using fixed key
func (logger *LOGGER) decryptSecret(s string) string {
	// Use constant key
	key := []byte(SECRET_KEY) // 32-byte AES-256 key

	// Decode base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return s // Return original string if decoding fails
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return s
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return s
	}

	// Separate nonce and ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return s
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt message
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return s
	}

	// Return decrypted original string
	return string(plaintext)
}

// internal AES encryption function using fixed key
func (logger *LOGGER) encryptSecret(s string) string {
	// Use constant key
	key := []byte(SECRET_KEY) // 32-byte AES-256 key

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return s // Return original string if encryption fails
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return s
	}

	// Generate random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return s
	}

	// Encrypt message
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(s), nil)

	// Return base64 encoded encryption result
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// getGoroutineID returns the current goroutine ID.
func (logger *LOGGER) getGoroutineID() int {
	var buf [64]byte
	var goroutineID int
	goroutineID = -1
	n := runtime.Stack(buf[:], false)
	for i := 0; i < n-10; i++ {
		if string(buf[i:i+9]) == "goroutine " {
			j := i + 9
			for ; j < n && buf[j] >= '0' && buf[j] <= '9'; j++ {
				goroutineID = goroutineID*10 + int(buf[j]-'0')
			}
			break
		}
	}
	return goroutineID
}

// output logs a message to console and file with specified log level.
// level: Log level of the message
// s: Message to log
// skipStackFrames: Number of stack frames to skip when logging
func (logger *LOGGER) output(level int, s string, skipStackFrames ...int) {
	if s != "" && level >= logger.CurrentLevel {
		logger.mutex.Lock()
		var skipCount int
		skipCount = 1
		if len(skipStackFrames) > 0 {
			skipCount = skipStackFrames[0]
		}
		pc, filename, line, _ := runtime.Caller(skipCount)
		parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
		funcName := parts[len(parts)-1]

		// Get log level name
		levelName := logLevelNames[level]

		// Base information uses default color
		baseColor := logger.White

		// Set log level text color based on log level
		levelTextColor := baseColor
		// Set log message color based on log level
		messageColor := baseColor
		// Log message display content
		logMessage := s

		switch level {
		case LogLevelDebug:
			levelTextColor = logger.Gray // light gray for debug tag
			messageColor = logger.Gray   // light gray for debug message
		case LogLevelInfo:
			levelTextColor = logger.Green // bright green for info tag
			messageColor = logger.White   // bright white for info message
		case LogLevelWarn:
			levelTextColor = logger.Yellow
			messageColor = logger.Yellow
		case LogLevelError:
			levelTextColor = logger.Red
			messageColor = logger.Red
		case LogLevelSecret:
			levelTextColor = logger.Magenta // Pink color
			messageColor = logger.Magenta   // SECRET log's ****** shows as pink
			logMessage = "******"           // SECRET logs always show as ****** in console
		}

		// Construct console output text with different colors for different parts
		consoleText := fmt.Sprintf(
			"[%s][%s::%s::%d][Goroutine:%d][%s%s%s]> %s%s%s",
			time.Now().Format("2006/01/02 15:04:05"),
			filepath.Base(filename),
			funcName,
			line,
			logger.getGoroutineID(),
			levelTextColor,
			levelName,
			logger.ColorReset,
			messageColor,
			logMessage,
			logger.ColorReset,
		)

		// Construct file output text
		fileText := fmt.Sprintf(
			"[%s][%s::%s::%d][Goroutine:%d][%s]> %s",
			time.Now().Format("2006/01/02 15:04:05"),
			filepath.Base(filename),
			funcName,
			line,
			logger.getGoroutineID(),
			levelName,
			s, // Keep original message in file
		)

		// If it's a SECRET level log, encrypt the message before writing to file
		if level == LogLevelSecret {
			// Encrypt original message
			encryptedMsg := logger.encryptSecret(s)
			// Reconstruct fileText with encrypted message
			fileText = fmt.Sprintf(
				"[%s][%s::%s::%d][Goroutine:%d][%s]> %s",
				time.Now().Format("2006/01/02 15:04:05"),
				filepath.Base(filename),
				funcName,
				line,
				logger.getGoroutineID(),
				levelName,
				encryptedMsg, // Write encrypted message to file
			)
		}

		if runtime.GOOS == "windows" && !Debugger.IsPresent() {
			consoleText = fileText
		}
		fmt.Println(consoleText)
		logger.writeFile(fileText)
		logger.mutex.Unlock()
	}
}

// writeFile writes the log text to a file.
// logText: Text to write to log file
func (logger *LOGGER) writeFile(logText string) {
	err := os.MkdirAll(logger.FolderName, os.ModePerm)
	if err == nil {
		logFilePath := path.Join(logger.FolderName, time.Now().Format("20060102")) + ".log"
		file, err := os.OpenFile(
			logFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			fileLogger := log.New(file, "", log.Ldate|log.Ltime)
			fileLogger.Println(logText)
		}
	}
}
