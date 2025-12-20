// Package logger
// File:        logger.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/logger/logger.go
// Author:      TRAE.AI
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
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/xiang-tai-duo/go-boost/debugger"
)

type (
	LOGGER struct {
		Blue                string
		BlueBackground      string
		BrightYellow        string
		ColorReset          string
		CurrentLevel        int
		Cyan                string
		CyanBackground      string
		Gray                string
		Green               string
		GreenBackground     string
		Magenta             string
		MagentaBackground   string
		MaxDailyLogSize     int64
		OutputMode          int
		Red                 string
		RedBackground       string
		White               string
		WhiteBackground     string
		Yellow              string
		YellowBackground    string
		currentLogDate      string
		executableBaseName  string
		executableDirectory string
		folderName          string
		isWindowsService    bool
		logger              *log.Logger
		loggerMutex         sync.Mutex
	}
)

//goland:noinspection GoSnakeCaseUsage
const (
	ANSI_BLUE                  = "\033[34m"
	ANSI_BLUE_BACKGROUND       = "\033[44m"
	ANSI_BRIGHT_YELLOW         = "\033[93m"
	ANSI_COLOR_PATTERN         = `\x1b\[[0-9;]*m`
	ANSI_COLOR_RESET           = "\033[0m"
	ANSI_CYAN                  = "\033[36m"
	ANSI_CYAN_BACKGROUND       = "\033[46m"
	ANSI_GRAY                  = "\033[37;2m"
	ANSI_GREEN                 = "\033[92m"
	ANSI_GREEN_BACKGROUND      = "\033[42m"
	ANSI_MAGENTA               = "\033[35m"
	ANSI_MAGENTA_BACKGROUND    = "\033[45m"
	ANSI_RED                   = "\033[31m"
	ANSI_RED_BACKGROUND        = "\033[41m"
	ANSI_WHITE                 = "\033[97m"
	ANSI_WHITE_BACKGROUND      = "\033[47m"
	ANSI_YELLOW                = "\033[33m"
	ANSI_YELLOW_BACKGROUND     = "\033[43m"
	BUFFER_SIZE                = 64
	CONSOLE_TEXT_FORMAT        = "[%s][%s::%s::%d][PID:%d][Goroutine:%d][%s%s%s]> %s%s%s"
	CR                         = "\r"
	DEFAULT_FOLDER_NAME        = "Logs"
	DOT                        = "."
	FILE_TEXT_FORMAT           = "[%s][%s::%s::%d][PID:%d][Goroutine:%d][%s]> %s"
	GOROUTINE_PREFIX           = "goroutine "
	GOROUTINE_PREFIX_LENGTH    = 9
	INITIAL_GOROUTINE_ID       = -1
	LEVEL_NAME_DEBUG           = "DEBUG"
	LEVEL_NAME_ERROR           = "ERROR"
	LEVEL_NAME_INFO            = "INFO"
	LEVEL_NAME_SECRET          = "SECRET"
	LEVEL_NAME_WARN            = "WARN"
	LF                         = "\n"
	LOG_EMPTY_PREFIX           = ""
	LOG_FILE_DATE_FORMAT       = "20060102"
	LOG_FILE_EXTENSION         = ".log"
	LOG_FILE_RETENTION_DAYS    = 100
	LOG_LEVEL_DEBUG            = 0
	LOG_LEVEL_ERROR            = 3
	LOG_LEVEL_INFO             = 1
	LOG_LEVEL_SECRET           = 4
	LOG_LEVEL_WARN             = 2
	MIN_SECRET_MATCH_LENGTH    = 2
	OS_WINDOWS                 = "windows"
	OUTPUT_MODE_BOTH           = 2
	OUTPUT_MODE_CONSOLE_ONLY   = 0
	OUTPUT_MODE_FILE_ONLY      = 1
	SCAN_TAIL_OFFSET           = 10
	SECRET_KEY                 = "X7K!2p9@A4z8#F1m5$C3b7%N6D2^q8&W0e9*R4t7(Y2u5)I8o1P3"
	SECRET_LOG_PATTERN         = `\[SECRET]]>\s*([A-Za-z0-9+/=]+)`
	SECRET_MASK                = "******"
	SKIP_STACK_FRAMES_BASE     = 1
	SKIP_STACK_FRAMES_OUTPUT   = 2
	TIMESTAMP_FORMAT           = "2006/01/02 15:04:05"
	DEFAULT_MAX_DAILY_LOG_SIZE = 1 * 1024 * 1024 * 1024
	MIN_FREE_DISK_SPACE        = 100 * 1024 * 1024
)

var (
	ansiColorRegex = regexp.MustCompile(ANSI_COLOR_PATTERN)
	logLevelNames  = map[int]string{
		LOG_LEVEL_DEBUG:  LEVEL_NAME_DEBUG,
		LOG_LEVEL_INFO:   LEVEL_NAME_INFO,
		LOG_LEVEL_WARN:   LEVEL_NAME_WARN,
		LOG_LEVEL_ERROR:  LEVEL_NAME_ERROR,
		LOG_LEVEL_SECRET: LEVEL_NAME_SECRET,
	}
	processID        = os.Getpid()
	lastCleanupDate  = ""
	Logger           = &LOGGER{}
	compressWarnOnce sync.Once
)

func init() {
	Logger.Blue = ANSI_BLUE
	Logger.BlueBackground = ANSI_BLUE_BACKGROUND
	Logger.BrightYellow = ANSI_BRIGHT_YELLOW
	Logger.ColorReset = ANSI_COLOR_RESET
	Logger.CurrentLevel = LOG_LEVEL_INFO
	Logger.Cyan = ANSI_CYAN
	Logger.CyanBackground = ANSI_CYAN_BACKGROUND
	Logger.Gray = ANSI_GRAY
	Logger.Green = ANSI_GREEN
	Logger.GreenBackground = ANSI_GREEN_BACKGROUND
	Logger.Magenta = ANSI_MAGENTA
	Logger.MagentaBackground = ANSI_MAGENTA_BACKGROUND
	Logger.OutputMode = OUTPUT_MODE_BOTH
	Logger.Red = ANSI_RED
	Logger.RedBackground = ANSI_RED_BACKGROUND
	Logger.White = ANSI_WHITE
	Logger.WhiteBackground = ANSI_WHITE_BACKGROUND
	Logger.Yellow = ANSI_YELLOW
	Logger.YellowBackground = ANSI_YELLOW_BACKGROUND
	Logger.MaxDailyLogSize = DEFAULT_MAX_DAILY_LOG_SIZE
	Logger.isWindowsService = IsWindowsService()
	if executableFilePath, err := os.Executable(); err == nil {
		Logger.executableDirectory = filepath.Dir(executableFilePath)
		Logger.executableBaseName = filepath.Base(executableFilePath)
	}
	Logger.SetFolderName(DEFAULT_FOLDER_NAME)
	Logger.cleanupExpiredLogFiles()
}

func (logger *LOGGER) DebugEx(message string, moduleName string, skipStackFrames int) {
	logger.output(LOG_LEVEL_DEBUG, message, moduleName, SKIP_STACK_FRAMES_OUTPUT+skipStackFrames)
}

func (logger *LOGGER) Debug(message string) {
	logger.DebugEx(message, "", 1)
}

func (logger *LOGGER) DecryptSecretLogs(logFilePath string) ([]string, error) {
	result := make([]string, 0)
	err := error(nil)
	file := (*os.File)(nil)
	if file, err = os.Open(logFilePath); err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		content := make([]byte, 0)
		if content, err = io.ReadAll(file); err == nil {
			logPattern := regexp.MustCompile(SECRET_LOG_PATTERN)
			matches := logPattern.FindAllSubmatch(content, -1)
			for _, match := range matches {
				if len(match) >= MIN_SECRET_MATCH_LENGTH {
					encryptedBase64 := string(match[1])
					decryptedMsg := logger.decryptSecret(encryptedBase64)
					result = append(result, decryptedMsg)
				}
			}
		}
	}
	return result, err
}

func (logger *LOGGER) ErrorEx(message interface{}, moduleName string, skipStackFrames int) {
	logMessage := ""
	skip := false
	switch value := message.(type) {
	case error:
		if value != nil {
			syscallErrorNumber := syscall.Errno(0)
			if errors.As(value, &syscallErrorNumber) && syscallErrorNumber == 0 {
				skip = true
			}
			exitError := (*exec.ExitError)(nil)
			if errors.As(value, &exitError) {
				skip = true
			}
			if !skip {
				logMessage = value.Error()
			}
		}
	case string:
		logMessage = value
	default:
		logMessage = fmt.Sprint(value)
	}
	if logMessage != "" {
		logger.output(LOG_LEVEL_ERROR, logMessage, moduleName, SKIP_STACK_FRAMES_OUTPUT+skipStackFrames)
	}
}

func (logger *LOGGER) Error(message interface{}) {
	logger.ErrorEx(message, "", SKIP_STACK_FRAMES_BASE)
}

func (logger *LOGGER) GetGoroutineID() int {
	result := INITIAL_GOROUTINE_ID
	buffer := [BUFFER_SIZE]byte{}
	length := runtime.Stack(buffer[:], false)
	for index := 0; index < length-SCAN_TAIL_OFFSET; index++ {
		if string(buffer[index:index+GOROUTINE_PREFIX_LENGTH]) == GOROUTINE_PREFIX {
			digitIndex := index + GOROUTINE_PREFIX_LENGTH
			for ; digitIndex < length && buffer[digitIndex] >= '0' && buffer[digitIndex] <= '9'; digitIndex++ {
				result = result*10 + int(buffer[digitIndex]-'0')
			}
			break
		}
	}
	return result
}

func (logger *LOGGER) GetMaxDailyLogSize() int64 {
	return logger.MaxDailyLogSize
}

func (logger *LOGGER) SetFolderName(name string) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		trimmed = DEFAULT_FOLDER_NAME
	}
	base := filepath.Base(trimmed)
	if base == "" || base == "." || base == string(filepath.Separator) {
		base = DEFAULT_FOLDER_NAME
	}
	logger.folderName = base
}

func (logger *LOGGER) GetFolderName() string {
	if logger.folderName == "" {
		return DEFAULT_FOLDER_NAME
	}
	return logger.folderName
}

func (logger *LOGGER) GetLogFilePath(packageName string) (string, error) {
	result := ""
	err := error(nil)
	today := time.Now().Format(LOG_FILE_DATE_FORMAT)
	if logger.currentLogDate != today {
		logger.currentLogDate = today
		logger.cleanupExpiredLogFiles()
	}
	logRootDirectory := ""
	if logRootDirectory, err = logger.getLogRootDirectory(); err == nil && logRootDirectory != "" {
		if mkErr := os.MkdirAll(logRootDirectory, os.ModePerm); mkErr == nil {
			logDailyDirectory := filepath.Join(logRootDirectory, today)
			if mkErr = os.MkdirAll(logDailyDirectory, os.ModePerm); mkErr == nil {
				logFileName := today
				if packageName != "" {
					logFileName = packageName
				} else if logger.executableBaseName != "" {
					logFileName = logger.executableBaseName
				}
				result = filepath.Join(logDailyDirectory, logFileName) + LOG_FILE_EXTENSION
			} else {
				err = mkErr
			}
		} else {
			err = mkErr
		}
	}
	return result, err
}

func (logger *LOGGER) InfoEx(message string, moduleName string, skipStackFrames int) {
	logger.output(LOG_LEVEL_INFO, message, moduleName, SKIP_STACK_FRAMES_OUTPUT+skipStackFrames)
}

func (logger *LOGGER) Info(message string) {
	logger.InfoEx(message, "", SKIP_STACK_FRAMES_BASE)
}

func (logger *LOGGER) SecretEx(message string, moduleName string, skipStackFrames int) {
	logger.output(LOG_LEVEL_SECRET, message, moduleName, SKIP_STACK_FRAMES_OUTPUT+skipStackFrames)
}

func (logger *LOGGER) Secret(message string) {
	logger.SecretEx(message, "", SKIP_STACK_FRAMES_BASE)
}

func (logger *LOGGER) SetMaxDailyLogSize(size int64) {
	logger.MaxDailyLogSize = size
}

func (logger *LOGGER) WarningEx(message string, moduleName string, skipStackFrames int) {
	logger.output(LOG_LEVEL_WARN, message, moduleName, SKIP_STACK_FRAMES_OUTPUT+skipStackFrames)
}

func (logger *LOGGER) Warning(message string) {
	logger.WarningEx(message, "", SKIP_STACK_FRAMES_OUTPUT)
}

func (logger *LOGGER) cleanupExpiredLogFiles() {
	result := true
	err := error(nil)
	today := time.Now().Format(LOG_FILE_DATE_FORMAT)
	if lastCleanupDate == today {
		result = false
	}
	if result {
		lastCleanupDate = today
	}
	logRootDirectory := ""
	if result {
		logRootDirectory, err = logger.getLogRootDirectory()
		if err != nil || logRootDirectory == "" {
			result = false
		}
	}
	entries := make([]os.DirEntry, 0)
	if result {
		entries, err = os.ReadDir(logRootDirectory)
		if err != nil {
			result = false
		}
	}
	logDirectoryNames := make([]string, 0)
	if result {
		for _, entry := range entries {
			if entry.IsDir() {
				_, err = time.Parse(LOG_FILE_DATE_FORMAT, entry.Name())
				if err == nil {
					logDirectoryNames = append(logDirectoryNames, entry.Name())
				}
			}
		}
		sort.Sort(sort.Reverse(sort.StringSlice(logDirectoryNames)))
		for index, logDirectoryName := range logDirectoryNames {
			if index >= LOG_FILE_RETENTION_DAYS {
				_ = os.RemoveAll(filepath.Join(logRootDirectory, logDirectoryName))
			}
		}
	}
}

func (logger *LOGGER) decryptSecret(message string) string {
	result := message
	err := error(nil)
	key := []byte(SECRET_KEY)
	ciphertext := make([]byte, 0)
	if ciphertext, err = base64.StdEncoding.DecodeString(message); err == nil {
		block := cipher.Block(nil)
		if block, err = aes.NewCipher(key); err == nil {
			aesGCM := cipher.AEAD(nil)
			if aesGCM, err = cipher.NewGCM(block); err == nil {
				nonceSize := aesGCM.NonceSize()
				if len(ciphertext) >= nonceSize {
					nonce := ciphertext[:nonceSize]
					ciphertext = ciphertext[nonceSize:]
					plaintext := make([]byte, 0)
					if plaintext, err = aesGCM.Open(nil, nonce, ciphertext, nil); err == nil {
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
	err := error(nil)
	key := []byte(SECRET_KEY)
	block := cipher.Block(nil)
	if block, err = aes.NewCipher(key); err == nil {
		aesGCM := cipher.AEAD(nil)
		if aesGCM, err = cipher.NewGCM(block); err == nil {
			nonce := make([]byte, aesGCM.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err == nil {
				ciphertext := aesGCM.Seal(nonce, nonce, []byte(message), nil)
				result = base64.StdEncoding.EncodeToString(ciphertext)
			}
		}
	}
	return result
}

func (logger *LOGGER) getLogRootDirectory() (string, error) {
	result := ""
	err := error(nil)
	if logger.isWindowsService {
		result = logger.executableDirectory
	} else {
		workingDirectory := ""
		workingDirectory, err = os.Getwd()
		if err == nil {
			result = workingDirectory
		}
	}
	if err == nil && result != "" {
		result = filepath.Join(result, logger.GetFolderName())
	}
	return result, err
}

func (logger *LOGGER) output(level int, message string, packageName string, skipStackFrames ...int) {
	if message != "" && level >= logger.CurrentLevel {
		for {
			isSuffixTrimmed := false
			if strings.HasSuffix(message, LF) {
				message = strings.TrimSuffix(message, LF)
				isSuffixTrimmed = true
			}
			if strings.HasSuffix(message, CR) {
				message = strings.TrimSuffix(message, CR)
				isSuffixTrimmed = true
			}
			if !isSuffixTrimmed {
				break
			}
		}
		logger.loggerMutex.Lock()
		skipCount := SKIP_STACK_FRAMES_BASE
		if len(skipStackFrames) > 0 {
			skipCount = skipStackFrames[0]
		}
		programCounter, filename, line, _ := runtime.Caller(skipCount)
		parts := strings.Split(runtime.FuncForPC(programCounter).Name(), DOT)
		functionName := parts[len(parts)-1]
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
			levelTextColor = logger.BrightYellow
			messageColor = logger.Yellow
		case LOG_LEVEL_ERROR:
			levelTextColor = logger.Red
			messageColor = logger.Red
		case LOG_LEVEL_SECRET:
			levelTextColor = logger.Magenta
			messageColor = logger.Magenta
			logMessage = SECRET_MASK
		}
		consoleText := fmt.Sprintf(
			CONSOLE_TEXT_FORMAT,
			time.Now().Format(TIMESTAMP_FORMAT),
			filepath.Base(filename),
			functionName,
			line,
			processID,
			logger.GetGoroutineID(),
			levelTextColor,
			levelName,
			logger.ColorReset,
			messageColor,
			logMessage,
			logger.ColorReset,
		)
		strippedMessage := logger.stripANSIColorCodes(message)
		fileText := fmt.Sprintf(
			FILE_TEXT_FORMAT,
			time.Now().Format(TIMESTAMP_FORMAT),
			filepath.Base(filename),
			functionName,
			line,
			processID,
			logger.GetGoroutineID(),
			levelName,
			strippedMessage,
		)
		if level == LOG_LEVEL_SECRET {
			encryptedMsg := logger.encryptSecret(strippedMessage)
			fileText = fmt.Sprintf(
				FILE_TEXT_FORMAT,
				time.Now().Format(TIMESTAMP_FORMAT),
				filepath.Base(filename),
				functionName,
				line,
				processID,
				logger.GetGoroutineID(),
				levelName,
				encryptedMsg,
			)
		}
		if runtime.GOOS == OS_WINDOWS && !debugger.IsPresent() {
			consoleText = fileText
		}
		switch logger.OutputMode {
		case OUTPUT_MODE_CONSOLE_ONLY:
			fmt.Println(consoleText)
		case OUTPUT_MODE_FILE_ONLY:
			logger.writeFile(fileText, packageName)
		case OUTPUT_MODE_BOTH:
			fmt.Println(consoleText)
			logger.writeFile(fileText, packageName)
		}
		logger.loggerMutex.Unlock()
	}
}

func (logger *LOGGER) stripANSIColorCodes(message string) string {
	return ansiColorRegex.ReplaceAllString(message, "")
}

//goland:noinspection GoBoolExpressions
func (logger *LOGGER) writeFile(logText string, packageName string) {
	result := true
	err := error(nil)
	logFilePath := ""
	if logFilePath, err = logger.GetLogFilePath(packageName); err != nil || logFilePath == "" {
		result = false
	}
	if result {
		freeBytes := uint64(0)
		if freeBytes, err = getDiskFreeSpace(filepath.Dir(logFilePath)); err == nil && freeBytes >= MIN_FREE_DISK_SPACE {
			file := (*os.File)(nil)
			isLogFileExists := true
			if _, err = os.Stat(logFilePath); err != nil {
				if os.IsNotExist(err) {
					isLogFileExists = false
				}
			}
			if isLogFileExists {
				file, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND, os.ModePerm)
			} else {
				if file, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm); err == nil {
					_ = file.Close()
					if compressErr := SetFileCompression(logFilePath, true); compressErr != nil {
						compressWarnOnce.Do(func() {
							log.Printf("SetFileCompression failed: %v", compressErr)
						})
					}
					file, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND, os.ModePerm)
				}
			}
			if err == nil && file != nil {
				defer file.Close()
				isFileTooLarge := false
				var stat os.FileInfo
				if stat, err = file.Stat(); err == nil && stat.Size() >= logger.MaxDailyLogSize {
					isFileTooLarge = true
				}
				if !isFileTooLarge {
					fileLogger := log.New(file, LOG_EMPTY_PREFIX, 0)
					fileLogger.Print(logText + CR)
				}
			}
		}
	}
}
