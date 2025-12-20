// Package boost
// File:        sqlite.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: SQLite wrapper providing database connection management, query execution,
package boost

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

//goland:noinspection GoSnakeCaseUsage
type (
	SQLITE struct {
		Trace          bool
		SqliteFilePath string
		Database       *sql.DB
		PragmaKey      string
	}

	SQLITE_VALUE struct {
		Name  string
		Value interface{}
	}
)

func (sqlite *SQLITE) Close() {
	if sqlite != nil {
		if sqlite.Database != nil {
			_ = sqlite.Database.Close()
			sqlite.Database = nil
		}
		sqlite.SqliteFilePath = ""
	}
}

func (sqlite *SQLITE) Create(sqliteFilePath string) error {
	var err error
	var statErr error
	if _, statErr = os.Stat(sqliteFilePath); statErr == nil {
		err = fmt.Errorf("file already exists: %s", sqliteFilePath)
	} else if !os.IsNotExist(statErr) {
		err = fmt.Errorf("check file existence failed: %w", statErr)
	} else if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else {
		sqlite.Close()
		if sqlite.Database, err = sql.Open("sqlite3", sqliteFilePath); err != nil {
			err = fmt.Errorf("unable to open database: %w", err)
		} else {
			sqlite.Database.SetMaxOpenConns(1)
			sqlite.Database.SetMaxIdleConns(1)
			sqlite.Database.SetConnMaxLifetime(0)
			if sqlite.PragmaKey != "" {
				if _, err = sqlite.Database.Exec("PRAGMA key = ?", sqlite.PragmaKey); err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to set pragma key: %w", err)
				} else {
					var cipherVersion string
					if err = sqlite.Database.QueryRow("PRAGMA cipher_version").Scan(&cipherVersion); err != nil {
						sqlite.Close()
						err = fmt.Errorf("unable to check cipher version: %w", err)
					} else if cipherVersion == "" {
						sqlite.Close()
						err = errors.New("cipher version not found, encryption may not be available")
					}
				}
			} else {
				if _, err = sqlite.Database.Exec("PRAGMA synchronous = NORMAL"); err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to set synchronous mode: %w", err)
				}
			}

			if err == nil {
				if _, err = sqlite.Database.Exec("VACUUM"); err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to vacuum database: %w", err)
				} else {
					sqlite.SqliteFilePath = sqliteFilePath
				}
			}
		}
	}
	return err
}

func (sqlite *SQLITE) ExecNonQuery(query string) error {
	var err error
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.Database == nil {
		err = errors.New("database connection not open")
	} else if query == "" {
		err = errors.New("query cannot be empty")
	} else {
		if _, err = sqlite.Database.Exec(query); err != nil {
			err = fmt.Errorf("unable to execute query: %w", err)
		}
	}
	return err
}

func (sqlite *SQLITE) ExecuteQuery(query string) ([]SQLITE_VALUE, error) {
	var results []SQLITE_VALUE
	var err error
	var rows *sql.Rows
	var columns []string

	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
		Logger.Error(err)
	} else if sqlite.Database == nil {
		err = errors.New("database connection not open")
		Logger.Error(err)
	} else if query == "" {
		err = errors.New("query cannot be empty")
		Logger.Error(err)
	} else {
		if rows, err = sqlite.Database.Query(query); err != nil {
			err = fmt.Errorf("unable to execute query: %w", err)
			Logger.Error(err)
		} else {
			defer func() {
				_ = rows.Close()
			}()

			if columns, err = rows.Columns(); err != nil {
				err = fmt.Errorf("unable to get columns: %w", err)
				Logger.Error(err)
			} else {
				columnCount := len(columns)
				values := make([]interface{}, columnCount)
				valuesPtr := make([]interface{}, columnCount)

				results = make([]SQLITE_VALUE, 0)

				for rows.Next() {
					for i := range values {
						valuesPtr[i] = &values[i]
					}

					if err = rows.Scan(valuesPtr...); err != nil {
						err = fmt.Errorf("unable to scan row: %w", err)
						Logger.Error(err)
						break
					}

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
				}

				if err == nil {
					if err = rows.Err(); err != nil {
						err = fmt.Errorf("error during row iteration: %w", err)
						Logger.Error(err)
					}
				}
			}
		}
	}

	return results, err
}

func (sqlite *SQLITE) Open(sqliteFilePath string) error {
	var err error
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else {
		sqlite.Close()

		sqlite.SqliteFilePath = sqliteFilePath
		if sqlite.Database, err = sql.Open("sqlite3", sqlite.SqliteFilePath); err != nil {
			err = fmt.Errorf("unable to open %s: %w", sqliteFilePath, err)
		} else {
			sqlite.Database.SetMaxOpenConns(1)
			sqlite.Database.SetMaxIdleConns(1)
			sqlite.Database.SetConnMaxLifetime(0)

			if sqlite.PragmaKey != "" {
				if _, err = sqlite.Database.Exec("PRAGMA key = ?", sqlite.PragmaKey); err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to set pragma key: %w", err)
				} else {
					var keyStatus string
					if err = sqlite.Database.QueryRow("PRAGMA key").Scan(&keyStatus); err != nil {
						sqlite.Close()
						err = fmt.Errorf("unable to verify database connection: %w", err)
					}
				}
			}

			if err == nil {
				if _, err = sqlite.Database.Exec("PRAGMA journal_mode = WAL"); err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to set journal mode: %w", err)
				}
			}
		}
	}
	return err
}

func (sqlite *SQLITE) Parse(data interface{}) SQLITE_VALUE {
	return SQLITE_VALUE{Value: data}
}

func (sv SQLITE_VALUE) ToBool() bool {
	result := false
	if sv.Value != nil {
		switch val := sv.Value.(type) {
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

func (sv SQLITE_VALUE) ToFloat() float64 {
	result := 0.0
	if sv.Value != nil {
		switch val := sv.Value.(type) {
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

func (sv SQLITE_VALUE) ToInt() int {
	result := 0
	if sv.Value != nil {
		switch val := sv.Value.(type) {
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

func (sv SQLITE_VALUE) ToString() string {
	result := ""
	if sv.Value != nil {
		if s, ok := sv.Value.(string); ok {
			result = s
		}
	}
	return result
}

func (sv SQLITE_VALUE) Type() string {
	result := "nil"
	if sv.Value != nil {
		result = reflect.TypeOf(sv.Value).String()
	}
	return result
}

func init() {
}
