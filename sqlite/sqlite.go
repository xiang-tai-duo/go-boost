// Package sqlite
// File:        sqlite.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/sqlite/sqlite.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: SQLite is a wrapper for SQLite instance operations, providing a set of methods for instance management and query execution.
// --------------------------------------------------------------------------------
package sqlite

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	SQLITE struct {
		Trace          bool
		PragmaKey      string
		instance       *sql.DB
		sqliteFilePath string
		mutex          sync.Mutex
		inTransaction  bool
	}

	SQLITE_VALUE struct {
		Name  string
		Value interface{}
	}

	SQLITE_EXEC_CALLBACK func(sql string, err error) bool
)

//goland:noinspection GoUnusedExportedFunction
func New() *SQLITE {
	return &SQLITE{}
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (sqlite *SQLITE) BeginTransaction() error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.instance == nil {
		err = errors.New("instance connection not open")
	} else {
		sqlite.mutex.Lock()
		defer func() {
			if err != nil {
				sqlite.mutex.Unlock()
			}
		}()
		if _, err = sqlite.instance.Exec("BEGIN TRANSACTION"); err == nil {
			sqlite.inTransaction = true
		} else {
			err = fmt.Errorf("unable to begin transaction: %w", err)
		}
	}
	return err
}

//goland:noinspection SqlNoDataSourceInspection,GoUnhandledErrorResult,SpellCheckingInspection
func (sqlite *SQLITE) DesensitizeColumn(tableName string, columnName string) error {
	err := error(nil)
	var rows *sql.Rows
	var rowID int64
	var value string
	var newHashValue string
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.instance == nil {
		err = errors.New("instance connection not open")
	} else if tableName == "" {
		err = errors.New("table name cannot be empty")
	} else if columnName == "" {
		err = errors.New("column name cannot be empty")
	} else {
		if err = sqlite.BeginTransaction(); err == nil {
			defer func() {
				if err != nil {
					if rollbackError := sqlite.Rollback(); rollbackError != nil {
						err = fmt.Errorf("%w (and rollback failed: %v)", err, rollbackError)
					}
				} else {
					if commitError := sqlite.Commit(); commitError != nil {
						err = fmt.Errorf("failed to commit transaction: %w", commitError)
					}
				}
			}()
			query := fmt.Sprintf("SELECT rowid, %s FROM %s", columnName, tableName)
			if rows, err = sqlite.instance.Query(query); err == nil {
				defer rows.Close()
				for rows.Next() {
					if err = rows.Scan(&rowID, &value); err == nil {
						if value != "" {
							if !sqlite.isMD5Hash(value) {
								newHashValue = sqlite.MD5(value)
								updateQuery := fmt.Sprintf("UPDATE %s SET %s = ? WHERE rowid = ?", tableName, columnName)
								if _, updateError := sqlite.instance.Exec(updateQuery, newHashValue, rowID); updateError != nil {
									err = fmt.Errorf("failed to update record with rowid %d: %w", rowID, updateError)
									break
								}
							}
						}
					} else {
						err = fmt.Errorf("failed to scan row: %w", err)
						break
					}
				}
				if err == nil {
					if rowsErr := rows.Err(); rowsErr != nil {
						err = fmt.Errorf("error during row iteration: %w", rowsErr)
					}
				}
			} else {
				err = fmt.Errorf("failed to query table %s: %w", tableName, err)
			}
		} else {
			err = fmt.Errorf("failed to begin transaction: %w", err)
		}
	}
	return err
}

func (sqlite *SQLITE) Close() {
	if sqlite != nil {
		if sqlite.instance != nil {
			_ = sqlite.instance.Close()
			sqlite.instance = nil
		}
		sqlite.sqliteFilePath = ""
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (sqlite *SQLITE) Commit() error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.instance == nil {
		err = errors.New("instance connection not open")
	} else if !sqlite.inTransaction {
		err = errors.New("not in transaction")
	} else {
		if _, err = sqlite.instance.Exec("COMMIT"); err != nil {
			err = fmt.Errorf("unable to commit transaction: %w", err)
		}
		sqlite.inTransaction = false
		sqlite.mutex.Unlock()
	}

	return err
}

//goland:noinspection SqlNoDataSourceInspection,GoUnhandledErrorResult,SqlDialectInspection
func (sqlite *SQLITE) Create(sqliteFilePath string) error {
	err := error(nil)
	sqliteAbsoluteFilePath := ""
	if sqliteAbsoluteFilePath, err = sqlite.getAbsolutePath(sqliteFilePath); err == nil {
		var statError error
		if _, statError = os.Stat(sqliteAbsoluteFilePath); statError == nil {
			if sqlite == nil {
				err = errors.New("sqlite instance is nil")
			} else {
				sqlite.Close()
				if sqlite.instance, err = sql.Open("sqlite3", sqliteAbsoluteFilePath); err == nil {
					sqlite.instance.SetMaxOpenConns(1)
					sqlite.instance.SetMaxIdleConns(1)
					sqlite.instance.SetConnMaxLifetime(0)
					var schemaVersion int
					if err = sqlite.instance.QueryRow("PRAGMA schema_version").Scan(&schemaVersion); err != nil {
						sqlite.Close()
						err = fmt.Errorf("file is not a valid SQLite instance: %w", err)
						sqlite.handleCGOError(err)
					} else {
						sqlite.sqliteFilePath = sqliteAbsoluteFilePath
					}
				} else {
					err = fmt.Errorf("unable to open existing file as SQLite instance: %w", err)
					sqlite.handleCGOError(err)
				}
			}
		} else if !os.IsNotExist(statError) {
			err = fmt.Errorf("check file existence failed: %w", statError)
		} else if sqlite == nil {
			err = errors.New("sqlite instance is nil")
		} else {
			sqlite.Close()
			if sqlite.instance, err = sql.Open("sqlite3", sqliteAbsoluteFilePath); err == nil {
				sqlite.instance.SetMaxOpenConns(1)
				sqlite.instance.SetMaxIdleConns(1)
				sqlite.instance.SetConnMaxLifetime(0)
				if sqlite.PragmaKey == "" {
					if _, err = sqlite.instance.Exec("PRAGMA synchronous = NORMAL"); err != nil {
						sqlite.Close()
						err = fmt.Errorf("unable to set synchronous mode: %w", err)
						sqlite.handleCGOError(err)
					}
				} else {
					if _, err = sqlite.instance.Exec("PRAGMA key = ?", sqlite.PragmaKey); err == nil {
						var cipherVersion string
						if err = sqlite.instance.QueryRow("PRAGMA cipher_version").Scan(&cipherVersion); err == nil {
							if cipherVersion == "" {
								sqlite.Close()
								err = errors.New("cipher version not found, encryption may not be available")
							}
						} else {
							sqlite.Close()
							err = fmt.Errorf("unable to check cipher version: %w", err)
							sqlite.handleCGOError(err)
						}
					} else {
						sqlite.Close()
						err = fmt.Errorf("unable to set pragma key: %w", err)
						sqlite.handleCGOError(err)
					}
				}
				if err == nil {
					if _, err = sqlite.instance.Exec("VACUUM"); err == nil {
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

//goland:noinspection SqlNoDataSourceInspection
func (sqlite *SQLITE) CreateNew(sqliteFilePath string) error {
	sqliteAbsoluteFilePath, err := sqlite.getAbsolutePath(sqliteFilePath)
	if err == nil {
		var statError error
		if _, statError = os.Stat(sqliteAbsoluteFilePath); statError == nil {
			err = fmt.Errorf("file already exists: %s", sqliteAbsoluteFilePath)
		} else {
			err = sqlite.Create(sqliteAbsoluteFilePath)
		}
	}
	return err
}

func (sqlite *SQLITE) ExecNonQuery(query string) error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.instance == nil {
		err = errors.New("instance connection not open")
	} else if query == "" {
		err = errors.New("query cannot be empty")
	} else {
		sqlite.mutex.Lock()
		defer sqlite.mutex.Unlock()
		if _, err = sqlite.instance.Exec(query); err != nil {
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
	} else if sqlite.instance == nil {
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
					if _, err = sqlite.instance.Exec(query); err != nil {
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

func (sqlite *SQLITE) ExecuteQuery(query string) ([]SQLITE_VALUE, error) {
	err := error(nil)
	var results []SQLITE_VALUE
	var rows *sql.Rows
	var columns []string
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.instance == nil {
		err = errors.New("instance connection not open")
	} else if query == "" {
		err = errors.New("query cannot be empty")
	} else {
		sqlite.mutex.Lock()
		if rows, err = sqlite.instance.Query(query); err == nil {
			sqlite.mutex.Unlock()
			defer func() {
				_ = rows.Close()
			}()
			if columns, err = rows.Columns(); err == nil {
				columnCount := len(columns)
				values := make([]interface{}, columnCount)
				valuesPtr := make([]interface{}, columnCount)
				results = make([]SQLITE_VALUE, 0)
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

func (sqlite *SQLITE) MD5(value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}

//goland:noinspection SqlNoDataSourceInspection
func (sqlite *SQLITE) Open(sqliteFilePath string) error {
	sqliteAbsoluteFilePath, err := sqlite.getAbsolutePath(sqliteFilePath)
	if err == nil {
		var statError error
		if _, statError = os.Stat(sqliteAbsoluteFilePath); os.IsNotExist(statError) {
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

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (sqlite *SQLITE) Rollback() error {
	err := error(nil)
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.instance == nil {
		err = errors.New("instance connection not open")
	} else if !sqlite.inTransaction {
		err = errors.New("not in transaction")
	} else {
		if _, err = sqlite.instance.Exec("ROLLBACK"); err != nil {
			err = fmt.Errorf("unable to rollback transaction: %w", err)
		}
		sqlite.inTransaction = false
		sqlite.mutex.Unlock()
	}

	return err
}

func (sqlite *SQLITE) TEXT(str string) string {
	if sqlite != nil {
		str = strings.ReplaceAll(str, "'", "''")
		str = strings.ReplaceAll(str, "\\", "\\\\")
		str = strings.ReplaceAll(str, "\"", "\\\"")
		str = strings.ReplaceAll(str, "\n", "\\n")
		str = strings.ReplaceAll(str, "\r", "\\r")
		str = strings.ReplaceAll(str, "\t", "\\t")
	}
	return str
}

// 检查字符串是否为有效的 MD5 哈希值（32 位十六进制字符串）
func (sqlite *SQLITE) isMD5Hash(value string) bool {
	b := false
	if len(value) == 32 {
		b = true
		for _, char := range value {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				b = false
				break
			}
		}
	}
	return b
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
			result = strVal == "true" || strVal == "1" || strVal == "yes" || strVal == "on"
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
			if f, err := strconv.ParseFloat(val, 64); err == nil {
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
	result := "nil"
	if value.Value != nil {
		result = reflect.TypeOf(value.Value).String()
	}
	return result
}
