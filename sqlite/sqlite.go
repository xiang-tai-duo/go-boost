// Package sqlite
// File:        sqlite.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/sqlite/sqlite.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: SQLite is a wrapper for SQLite instance operations, providing a set of methods for instance management and query execution.
// --------------------------------------------------------------------------------
package sqlite

import (
	"crypto/sha3"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mutecomm/go-sqlcipher/v4"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	SQLITE struct {
		PragmaKey      string
		Trace          bool
		inTransaction  bool
		mutex          sync.Mutex
		sqlDb          *sql.DB
		sqliteFilePath string
	}

	SQLITE_EXEC_CALLBACK func(sql string, err error) bool

	SQLITE_VALUE struct {
		Name  string
		Value interface{}
	}
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	BOOL_STRING_ON               = "on"
	BOOL_STRING_ONE              = "1"
	BOOL_STRING_TRUE             = "true"
	BOOL_STRING_YES              = "yes"
	CHAR_UNDERSCORE              = '_'
	DANGEROUS_STATEMENT_ERROR    = "dangerous SQL statement rejected by safe mode: %s"
	FLOAT_BIT_SIZE               = 64
	MODULE_NAME_SQLITE           = "sqlite"
	SHA3_HASH_LENGTH             = 64
	SQL_BEGIN_TRANSACTION        = "BEGIN TRANSACTION"
	SQL_COMMIT                   = "COMMIT"
	SQL_DRIVER_NAME              = "sqlite3"
	SQL_PRAGMA_CIPHER_VERSION    = "PRAGMA cipher_version"
	SQL_PRAGMA_SCHEMA_VERSION    = "PRAGMA schema_version"
	SQL_PRAGMA_SYNCHRONOUS       = "PRAGMA synchronous = NORMAL"
	SQL_ROLLBACK                 = "ROLLBACK"
	SQL_SELECT_ROW_ID_FORMAT     = "SELECT rowid, %s FROM %s"
	SQL_UPDATE_BY_ROW_ID_FORMAT  = "UPDATE %s SET %s = ? WHERE rowid = ?"
	SQL_VACUUM                   = "VACUUM"
	TYPE_NIL                     = "nil"
	VALUE_CONN_MAX_LIFETIME_ZERO = 0
	VALUE_MAX_CONNECTION         = 1
)

var (
	sqliteSafeMode             = true
	sqliteDangerousSQLKeywords = []string{"DROP", "TRUNCATE", "ALTER", "GRANT", "REVOKE", "RENAME", "SHUTDOWN", "CREATE", "REPLACE"}
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_SQLITE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_SQLITE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_SQLITE, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedExportedFunction
func New() *SQLITE {
	return &SQLITE{}
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (sqlite *SQLITE) BeginTransaction() error {
	err := sqlite.beginTransactionLocked()
	return err
}

func (sqlite *SQLITE) Close() {
	if sqlite != nil {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		if sqlite.sqlDb != nil {
			if sqlite.inTransaction {
				_, _ = sqlite.sqlDb.Exec(SQL_ROLLBACK)
				sqlite.inTransaction = false
			}
			_ = sqlite.sqlDb.Close()
			sqlite.sqlDb = nil
		}
		sqlite.sqliteFilePath = ""
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (sqlite *SQLITE) Commit() error {
	err := sqlite.commitLocked()
	return err
}

//goland:noinspection SqlNoDataSourceInspection,GoUnhandledErrorResult,SpellCheckingInspection
func (sqlite *SQLITE) ConvertTextToSHA3Text(tableName string, columnName string) error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if tableName == "" {
		err = errors.New("table name cannot be empty")
	} else if columnName == "" {
		err = errors.New("column name cannot be empty")
	} else if !sqlite.isSQLiteIdentifier(tableName) {
		err = fmt.Errorf("invalid table name: %s", tableName)
	} else if !sqlite.isSQLiteIdentifier(columnName) {
		err = fmt.Errorf("invalid column name: %s", columnName)
	} else {
		err = sqlite.ExecuteInTransaction(func(tx *sql.Tx) error {
			transactionErr := error(nil)
			query := fmt.Sprintf(SQL_SELECT_ROW_ID_FORMAT, columnName, tableName)
			var rows *sql.Rows
			if rows, transactionErr = tx.Query(query); transactionErr == nil {
				defer rows.Close()
				var rowID int64
				var value string
				for rows.Next() {
					if transactionErr = rows.Scan(&rowID, &value); transactionErr == nil {
						if value != "" && !sqlite.isSHA3Hash(value) {
							newHashValue := sqlite.SHA3(value)
							updateQuery := fmt.Sprintf(SQL_UPDATE_BY_ROW_ID_FORMAT, tableName, columnName)
							if _, transactionErr = tx.Exec(updateQuery, newHashValue, rowID); transactionErr != nil {
								transactionErr = fmt.Errorf("failed to update record with row id %d: %w", rowID, transactionErr)
								break
							}
						}
					} else {
						transactionErr = fmt.Errorf("failed to scan row: %w", transactionErr)
						break
					}
				}
				if transactionErr == nil {
					if transactionErr = rows.Err(); transactionErr != nil {
						transactionErr = fmt.Errorf("error during row iteration: %w", transactionErr)
					}
				}
			} else {
				transactionErr = fmt.Errorf("failed to query table %s: %w", tableName, transactionErr)
			}
			return transactionErr
		})
	}
	return err
}

//goland:noinspection SqlNoDataSourceInspection,GoUnhandledErrorResult,SqlDialectInspection
func (sqlite *SQLITE) Create(sqliteFilePath string) error {
	err := error(nil)
	sqliteAbsoluteFilePath := ""
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqliteAbsoluteFilePath, err = sqlite.getAbsolutePath(sqliteFilePath); err == nil {
		dataSourceName := sqlite.buildDataSourceName(sqliteAbsoluteFilePath)
		if _, err = os.Stat(sqliteAbsoluteFilePath); err == nil {
			sqlite.Close()
			if sqlite.sqlDb, err = sql.Open(SQL_DRIVER_NAME, dataSourceName); err == nil {
				sqlite.sqlDb.SetMaxOpenConns(VALUE_MAX_CONNECTION)
				sqlite.sqlDb.SetMaxIdleConns(VALUE_MAX_CONNECTION)
				sqlite.sqlDb.SetConnMaxLifetime(VALUE_CONN_MAX_LIFETIME_ZERO)
				var schemaVersion int
				if err = sqlite.sqlDb.QueryRow(SQL_PRAGMA_SCHEMA_VERSION).Scan(&schemaVersion); err == nil {
					sqlite.sqliteFilePath = sqliteAbsoluteFilePath
				} else {
					sqlite.Close()
					err = fmt.Errorf("file is not a valid SQLite instance: %w", err)
					sqlite.handleCGOError(err)
				}
			} else {
				err = fmt.Errorf("unable to open existing file as SQLite instance: %w", err)
				sqlite.handleCGOError(err)
			}
		} else if os.IsNotExist(err) {
			sqlite.Close()
			if sqlite.sqlDb, err = sql.Open(SQL_DRIVER_NAME, dataSourceName); err == nil {
				sqlite.sqlDb.SetMaxOpenConns(VALUE_MAX_CONNECTION)
				sqlite.sqlDb.SetMaxIdleConns(VALUE_MAX_CONNECTION)
				sqlite.sqlDb.SetConnMaxLifetime(VALUE_CONN_MAX_LIFETIME_ZERO)
				if sqlite.PragmaKey == "" {
					if _, err = sqlite.sqlDb.Exec(SQL_PRAGMA_SYNCHRONOUS); err != nil {
						sqlite.Close()
						err = fmt.Errorf("unable to set synchronous mode: %w", err)
						sqlite.handleCGOError(err)
					}
				} else {
					var cipherVersion string
					if err = sqlite.sqlDb.QueryRow(SQL_PRAGMA_CIPHER_VERSION).Scan(&cipherVersion); err == nil {
						if cipherVersion == "" {
							sqlite.Close()
							err = errors.New("cipher version not found, encryption may not be available")
						}
					} else {
						sqlite.Close()
						err = fmt.Errorf("unable to check cipher version: %w", err)
						sqlite.handleCGOError(err)
					}
				}
				if err == nil {
					if _, err = sqlite.sqlDb.Exec(SQL_VACUUM); err == nil {
						sqlite.sqliteFilePath = sqliteAbsoluteFilePath
					} else {
						sqlite.Close()
						err = fmt.Errorf("unable to vacuum instance: %w", err)
					}
				}
			} else {
				err = fmt.Errorf("unable to open instance: %w", err)
				sqlite.handleCGOError(err)
			}
		}
	}
	return err
}

//goland:noinspection SqlNoDataSourceInspection,GoUnhandledErrorResult
func (sqlite *SQLITE) CreateNew(sqliteFilePath string) error {
	err := error(nil)
	sqliteAbsoluteFilePath := ""
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqliteAbsoluteFilePath, err = sqlite.getAbsolutePath(sqliteFilePath); err == nil {
		os.Remove(sqliteAbsoluteFilePath)
		if _, err = os.Stat(sqliteAbsoluteFilePath); os.IsNotExist(err) {
			err = sqlite.Create(sqliteAbsoluteFilePath)
		} else {
			err = fmt.Errorf("cannot create: %s", sqliteAbsoluteFilePath)
		}
	}
	return err
}

func (sqlite *SQLITE) Exec(query string, args ...interface{}) error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if query == "" {
		err = errors.New("query cannot be empty")
	} else {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		if _, err = sqlite.sqlDb.Exec(query, args...); err != nil {
			err = fmt.Errorf("unable to execute query: %w", err)
		}
	}
	return err
}

//goland:noinspection SqlDialectInspection
func (sqlite *SQLITE) ExecNonQueries(queries []string, callback SQLITE_EXEC_CALLBACK) error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if queries == nil || len(queries) == 0 {
		err = errors.New("sql statements array cannot be empty")
	} else {
		if err = sqlite.BeginTransaction(); err == nil {
			defer func() {
				if err != nil {
					if rollbackError := sqlite.Rollback(); rollbackError != nil {
						err = fmt.Errorf("%w (and rollback failed: %v)", err, rollbackError)
					}
				}
			}()
			for _, query := range queries {
				if query != "" {
					if sqliteSafeMode && isDangerousStatement(query) {
						err = fmt.Errorf(DANGEROUS_STATEMENT_ERROR, query)
						break
					}
					if _, err = sqlite.sqlDb.Exec(query); err != nil {
						err = fmt.Errorf("unable to execute query: %w", err)
					}
					if callback != nil && !callback(query, err) {
						err = errors.New("batch execution cancelled by callback")
						break
					}
				}
			}
			if err == nil {
				if err = sqlite.Commit(); err != nil {
					err = fmt.Errorf("unable to commit transaction: %w", err)
				}
			}
		} else {
			err = fmt.Errorf("unable to begin transaction: %w", err)
		}
	}
	return err
}

func (sqlite *SQLITE) ExecNonQuery(query string) error {
	return sqlite.Exec(query)
}

func (sqlite *SQLITE) ExecuteInTransaction(callback func(*sql.Tx) error) error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if callback == nil {
		err = errors.New("transaction callback cannot be nil")
	} else {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		var tx *sql.Tx
		if tx, err = sqlite.sqlDb.Begin(); err == nil {
			if err = callback(tx); err == nil {
				err = tx.Commit()
				if err != nil {
					err = fmt.Errorf("unable to commit transaction: %w", err)
				}
			} else {
				originalErr := err
				if rollbackError := tx.Rollback(); rollbackError != nil {
					err = fmt.Errorf("%w (and rollback failed: %v)", originalErr, rollbackError)
				}
			}
		} else {
			err = fmt.Errorf("unable to begin transaction: %w", err)
		}
	}
	return err
}

func (sqlite *SQLITE) ExecuteQuery(query string) ([]SQLITE_VALUE, error) {
	return sqlite.Query(query)
}

//goland:noinspection GoUnusedExportedFunction
func GetSafeMode() bool {
	return sqliteSafeMode
}

//goland:noinspection SqlNoDataSourceInspection
func (sqlite *SQLITE) Open(sqliteFilePath string) error {
	err := error(nil)
	sqliteAbsoluteFilePath := ""
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqliteAbsoluteFilePath, err = sqlite.getAbsolutePath(sqliteFilePath); err == nil {
		if _, err = os.Stat(sqliteAbsoluteFilePath); os.IsNotExist(err) {
			err = fmt.Errorf("file not found: %s", sqliteAbsoluteFilePath)
		} else {
			err = sqlite.Create(sqliteAbsoluteFilePath)
		}
	}
	return err
}

func (sqlite *SQLITE) Parse(data interface{}) SQLITE_VALUE {
	return SQLITE_VALUE{Value: data}
}

func (sqlite *SQLITE) Query(query string, args ...interface{}) ([]SQLITE_VALUE, error) {
	err := error(nil)
	results := make([]SQLITE_VALUE, 0)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if query == "" {
		err = errors.New("query cannot be empty")
	} else {
		sqlite.mutex.Lock()
		var rows *sql.Rows
		if rows, err = sqlite.sqlDb.Query(query, args...); err == nil {
			sqlite.mutex.Unlock()
			defer func() {
				_ = rows.Close()
			}()
			var columns []string
			if columns, err = rows.Columns(); err == nil {
				columnCount := len(columns)
				values := make([]interface{}, columnCount)
				valuesPtr := make([]interface{}, columnCount)
				for rows.Next() {
					for i := range values {
						valuesPtr[i] = &values[i]
					}
					if err = rows.Scan(valuesPtr...); err == nil {
						for i := 0; i < columnCount; i++ {
							value := values[i]
							if value != nil {
								if valType := reflect.TypeOf(value); valType.Kind() == reflect.Ptr {
									value = reflect.ValueOf(value).Elem().Interface()
								}
							}
							results = append(results, SQLITE_VALUE{
								Name:  columns[i],
								Value: value,
							})
						}
					} else {
						err = fmt.Errorf("unable to scan row: %w", err)
						break
					}
				}
				if err == nil {
					if err = rows.Err(); err != nil {
						err = fmt.Errorf("error during row iteration: %w", err)
					}
				}
			} else {
				err = fmt.Errorf("unable to get columns: %w", err)
			}
		} else {
			sqlite.mutex.Unlock()
			err = fmt.Errorf("unable to execute query: %w", err)
		}
	}
	return results, err
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (sqlite *SQLITE) Rollback() error {
	err := sqlite.rollbackLocked()
	return err
}

func (sqlite *SQLITE) SHA3(value string) string {
	hash := sha3.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

//goland:noinspection GoUnusedExportedFunction
func SetSafeMode(enabled bool) {
	sqliteSafeMode = enabled
}

func (sqlite *SQLITE) TEXT(str string) string {
	result := "NULL"
	if sqlite != nil && !strings.ContainsRune(str, 0) {
		result = "'" + strings.ReplaceAll(str, "'", "''") + "'"
	}
	return result
}

func (value SQLITE_VALUE) ToBool() bool {
	result := false
	if value.Value != nil {
		switch val := value.Value.(type) {
		case bool:
			result = val
		case int:
			result = val != 0
		case int64:
			result = val != 0
		case float32:
			result = val != 0
		case float64:
			result = val != 0
		case string:
			strVal := strings.ToLower(val)
			result = strVal == BOOL_STRING_TRUE || strVal == BOOL_STRING_ONE || strVal == BOOL_STRING_YES || strVal == BOOL_STRING_ON
		}
	}
	return result
}

func (value SQLITE_VALUE) ToFloat() float64 {
	result := 0.0
	if value.Value != nil {
		switch val := value.Value.(type) {
		case int:
			result = float64(val)
		case int64:
			result = float64(val)
		case float32:
			result = float64(val)
		case float64:
			result = val
		case string:
			if f, err := strconv.ParseFloat(val, FLOAT_BIT_SIZE); err == nil {
				result = f
			}
		}
	}
	return result
}

func (value SQLITE_VALUE) ToInt() int {
	result := 0
	if value.Value != nil {
		switch val := value.Value.(type) {
		case int:
			result = val
		case int64:
			result = int(val)
		case float32:
			result = int(val)
		case float64:
			result = int(val)
		case string:
			if i, err := strconv.Atoi(val); err == nil {
				result = i
			}
		}
	}
	return result
}

func (value SQLITE_VALUE) ToString() string {
	result := ""
	if value.Value != nil {
		if s, ok := value.Value.(string); ok {
			result = s
		}
	}
	return result
}

func (value SQLITE_VALUE) Type() string {
	result := TYPE_NIL
	if value.Value != nil {
		result = reflect.TypeOf(value.Value).String()
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_SQLITE, logger.SKIP_STACK_FRAMES_BASE)
}

func (sqlite *SQLITE) beginTransactionInternal() error {
	err := error(nil)
	if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if sqlite.inTransaction {
		err = errors.New("already in transaction")
	} else if _, err = sqlite.sqlDb.Exec(SQL_BEGIN_TRANSACTION); err == nil {
		sqlite.inTransaction = true
	} else {
		err = fmt.Errorf("unable to begin transaction: %w", err)
	}
	return err
}

func (sqlite *SQLITE) beginTransactionLocked() error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		err = sqlite.beginTransactionInternal()
	}
	return err
}

func (sqlite *SQLITE) buildDataSourceName(sqliteAbsoluteFilePath string) string {
	result := sqliteAbsoluteFilePath
	if sqlite != nil && sqlite.PragmaKey != "" && !strings.ContainsRune(sqlite.PragmaKey, 0) {
		pragmaKeyHex := hex.EncodeToString([]byte(sqlite.PragmaKey))
		pragmaKeyValue := fmt.Sprintf("x'%s'", pragmaKeyHex)
		result = fmt.Sprintf("%s?_pragma_key=%s", sqliteAbsoluteFilePath, url.QueryEscape(pragmaKeyValue))
	}
	return result
}

func (sqlite *SQLITE) commitInternal() error {
	err := error(nil)
	if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if !sqlite.inTransaction {
		err = errors.New("not in transaction")
	} else {
		if _, err = sqlite.sqlDb.Exec(SQL_COMMIT); err != nil {
			err = fmt.Errorf("unable to commit transaction: %w", err)
		}
		sqlite.inTransaction = false
	}
	return err
}

func (sqlite *SQLITE) commitLocked() error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		err = sqlite.commitInternal()
	}
	return err
}

func (sqlite *SQLITE) getAbsolutePath(sqliteFilePath string) (string, error) {
	err := errors.New("file path cannot be empty")
	absoluteFilePath := ""
	if sqliteFilePath != "" {
		absoluteFilePath, err = filepath.Abs(sqliteFilePath)
	}
	return absoluteFilePath, err
}

//goland:noinspection SpellCheckingInspection
func (sqlite *SQLITE) handleCGOError(err error) {
	if strings.Contains(err.Error(), "Binary was compiled with 'CGO_ENABLED=0'") {
		fmt.Println("[ERROR] SQLite error: Binary was compiled with 'CGO_ENABLED=0'. go-sqlite3 requires cgo to work.")
		fmt.Println("[ERROR] Please set the CGO_ENABLED environment variable to 1 before compiling/running this application.")
		fmt.Println("[ERROR] Command examples:")
		fmt.Println("[ERROR]   --- TEMPORARY (RESETS AFTER RESTART) ---")
		fmt.Println("[ERROR]   Windows (Command Prompt): set CGO_ENABLED=1")
		fmt.Println("[ERROR]   Windows (PowerShell): $env:CGO_ENABLED=1")
		fmt.Println("[ERROR]   Mac/Linux (bash/zsh): export CGO_ENABLED=1")
		fmt.Println("[ERROR]")
		fmt.Println("[ERROR]   --- PERMANENT (PERSISTS AFTER RESTART) ---")
		fmt.Println("[ERROR]   Windows (Control Panel): Add 'CGO_ENABLED=1' to System Environment Variables")
		fmt.Println("[ERROR]   Windows (Command Line - User Level): setx CGO_ENABLED 1")
		fmt.Println("[ERROR]   Windows (Command Line - System Level, requires admin): setx CGO_ENABLED 1 /M")
		fmt.Println("[ERROR]   Mac/Linux (bash): Add 'export CGO_ENABLED=1' to ~/.bashrc")
		fmt.Println("[ERROR]   Mac/Linux (zsh): Add 'export CGO_ENABLED=1' to ~/.zshrc")
		fmt.Println("[ERROR]   After permanent setting, restart your terminal or run 'source ~/.bashrc' (or equivalent)")
	}
}

func isDangerousStatement(query string) bool {
	stripped := query
	for {
		stripped = strings.TrimLeft(stripped, " \t\r\n")
		if strings.HasPrefix(stripped, "--") {
			if idx := strings.IndexAny(stripped, "\r\n"); idx >= 0 {
				stripped = stripped[idx:]
				continue
			}
			return false
		}
		break
	}
	if stripped == "" {
		return false
	}
	end := len(stripped)
	for i, ch := range stripped {
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' || ch == '(' || ch == ';' {
			end = i
			break
		}
	}
	firstWord := strings.ToUpper(stripped[:end])
	for _, kw := range sqliteDangerousSQLKeywords {
		if firstWord == kw {
			return true
		}
	}
	return false
}

func (sqlite *SQLITE) isSHA3Hash(value string) bool {
	result := false
	if len(value) == SHA3_HASH_LENGTH {
		result = true
		for _, char := range value {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				result = false
				break
			}
		}
	}
	return result
}

func (sqlite *SQLITE) isSQLiteIdentifier(value string) bool {
	result := false
	if value != "" {
		result = true
		for index, char := range value {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == CHAR_UNDERSCORE || (index > 0 && char >= '0' && char <= '9')) {
				result = false
				break
			}
		}
	}
	return result
}

func (sqlite *SQLITE) rollbackInternal() error {
	err := error(nil)
	if sqlite.sqlDb == nil {
		err = errors.New("instance connection not open")
	} else if !sqlite.inTransaction {
		err = errors.New("not in transaction")
	} else {
		if _, err = sqlite.sqlDb.Exec(SQL_ROLLBACK); err != nil {
			err = fmt.Errorf("unable to rollback transaction: %w", err)
		}
		sqlite.inTransaction = false
	}
	return err
}

func (sqlite *SQLITE) rollbackLocked() error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		err = sqlite.rollbackInternal()
	}
	return err
}
