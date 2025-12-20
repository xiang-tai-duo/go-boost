// --------------------------------------------------------------------------------
// File:        sqlite.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: SQLite wrapper providing database connection management, query execution,
//              and encryption support for SQLite databases.
// --------------------------------------------------------------------------------

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

// SQLITE provides a wrapper for SQLite database operations with encryption support.
type SQLITE struct {
	Trace          bool
	SqliteFilePath string
	Database       *sql.DB
	PragmaKey      string
}

func init() {

}

// Open opens an existing SQLite database file.
// sqliteFilePath: Path to the SQLite database file
// Returns: Error encountered during database opening
// Usage:
// sqlite := &SQLITE{}
// err := sqlite.Open("database.db")
func (sqlite *SQLITE) Open(sqliteFilePath string) (err error) {
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
				_, err = sqlite.Database.Exec("PRAGMA key = ?", sqlite.PragmaKey)
				if err != nil {
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
				_, err = sqlite.Database.Exec("PRAGMA journal_mode = WAL")
				if err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to set journal mode: %w", err)
				}
			}
		}
	}
	return err
}

// Close closes the SQLite database connection.
// Usage:
// sqlite.Close()
func (sqlite *SQLITE) Close() {
	if sqlite != nil {
		if sqlite.Database != nil {
			_ = sqlite.Database.Close()
			sqlite.Database = nil
		}
		sqlite.SqliteFilePath = ""
	}
}

// Create creates a new SQLite database file.
// sqliteFilePath: Path to the SQLite database file to create
// Returns: Error encountered during database creation
// Usage:
//
//	sqlite := &SQLITE{
//	    PragmaKey: "encryption_key",
//	}
//
// err := sqlite.Create("new_database.db")
func (sqlite *SQLITE) Create(sqliteFilePath string) (err error) {
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
				_, err = sqlite.Database.Exec("PRAGMA key = ?", sqlite.PragmaKey)
				if err != nil {
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
				_, err = sqlite.Database.Exec("PRAGMA synchronous = NORMAL")
				if err != nil {
					sqlite.Close()
					err = fmt.Errorf("unable to set synchronous mode: %w", err)
				}
			}

			if err == nil {
				_, err = sqlite.Database.Exec("VACUUM")
				if err != nil {
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

// ExecNonQuery executes a SQL query that does not return rows.
// query: SQL query to execute
// Returns: Error encountered during query execution
// Usage:
// err := sqlite.ExecNonQuery("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
func (sqlite *SQLITE) ExecNonQuery(query string) (err error) {
	if sqlite == nil {
		err = errors.New("sqlite instance is nil")
	} else if sqlite.Database == nil {
		err = errors.New("database connection not open")
	} else if query == "" {
		err = errors.New("query cannot be empty")
	} else {
		_, err = sqlite.Database.Exec(query)
		if err != nil {
			err = fmt.Errorf("unable to execute query: %w", err)
		}
	}
	return err
}

// ExecuteQuery executes a SQL query that returns rows.
// query: SQL query to execute
// Returns: Slice of SQLITE_VALUE containing query results, or error encountered during query execution
// Usage:
//	for _, field := range result {
//	    fmt.Printf("Name: %s, Value: %v, Type: %s\n", field.Name, field.Value, field.Type())
//	    if field.Name == "name" {
//	        name := field.ToString()
//	    }
//	    if field.Name == "age" {
//	        age := field.ToInt()
//	    }
//	}
// result, err := sqlite.ExecuteQuery("SELECT * FROM users WHERE id = 1")

func (sqlite *SQLITE) ExecuteQuery(query string) (results []SQLITE_VALUE, err error) {
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
		var rows *sql.Rows
		if rows, err = sqlite.Database.Query(query); err != nil {
			err = fmt.Errorf("unable to execute query: %w", err)
			Logger.Error(err)
		} else {
			defer func(rows *sql.Rows) {
				_ = rows.Close()
			}(rows)

			// Get column names
			var columns []string
			if columns, err = rows.Columns(); err != nil {
				err = fmt.Errorf("unable to get columns: %w", err)
				Logger.Error(err)
			} else {
				columnCount := len(columns)
				values := make([]interface{}, columnCount)
				valuesPtr := make([]interface{}, columnCount)

				// Initialize results slice
				results = make([]SQLITE_VALUE, 0)

				for rows.Next() {
					// Prepare the pointers for scanning
					for i := range values {
						valuesPtr[i] = &values[i]
					}

					// Scan the row into the values slice
					if err = rows.Scan(valuesPtr...); err != nil {
						err = fmt.Errorf("unable to scan row: %w", err)
						Logger.Error(err)
						break
					}

					// Process each value and create SQLITE_VALUE instances
					for i := 0; i < columnCount; i++ {
						value := values[i]
						if value != nil {
							// Check if it's a pointer and dereference if needed
							if valType := reflect.TypeOf(value); valType.Kind() == reflect.Ptr {
								value = reflect.ValueOf(value).Elem().Interface()
							}
						}
						// Create SQLITE_VALUE with column name and value
						results = append(results, SQLITE_VALUE{
							Name:  columns[i],
							Value: value,
						})
					}
				}

				// Check for any errors during row iteration
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

// SQLITE_VALUE represents a parsed SQLite value with conversion methods
// This struct provides type-safe conversion methods for SQLite database values
// Usage:
// result, _ := SQLite.ExecuteQuery("SELECT * FROM users WHERE id = 1")
//
//	if len(result) > 0 {
//		for _, row := range result {
//		    for _, field := range row {
//		        if field.name == "name" {
//		            name := field.ToString()
//		        }
//		        if field.name == "age" {
//		            age := field.ToInt()
//		        }
//		    }
//	}
//
//goland:noinspection GoSnakeCaseUsage
type SQLITE_VALUE struct {
	Name  string
	Value interface{}
}

// Parse creates a new SQLITE_VALUE instance for type conversion
// data: The interface{} value to wrap for conversion
// Returns: SQLITE_VALUE instance with conversion methods
// Usage:
// parsed := SQLite.Parse(someData)
func (sqlite *SQLITE) Parse(data interface{}) SQLITE_VALUE {
	return SQLITE_VALUE{Value: data}
}

// ToString converts the wrapped value to string, returning empty string if conversion fails or value is nil
// Returns: String representation of the value, or empty string if conversion fails
// Usage:
// result := SQLite.Parse("hello").ToString()
// return: "hello"
//
// result := SQLite.Parse(nil).ToString()
// return: ""
func (sv SQLITE_VALUE) ToString() string {
	result := ""
	if sv.Value != nil {
		if s, ok := sv.Value.(string); ok {
			result = s
		}
	}
	return result
}

// ToInt converts the wrapped value to integer, returning 0 if conversion fails or value is nil
// Returns: Integer representation of the value, or 0 if conversion fails
// Usage:
// result := SQLite.Parse(42).ToInt()
// return: 42
//
// result := SQLite.Parse(3.14).ToInt()
// return: 3
//
// result := SQLite.Parse("123").ToInt()
// return: 123
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

// ToFloat converts the wrapped value to float64, returning 0.0 if conversion fails or value is nil
// Returns: Float64 representation of the value, or 0.0 if conversion fails
// Usage:
// result := SQLite.Parse(3.14).ToFloat()
// return: 3.14
//
// result := SQLite.Parse(42).ToFloat()
// return: 42.0
//
// result := SQLite.Parse("123.45").ToFloat()
// return: 123.45
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

// ToBool converts the wrapped value to bool, returning false if conversion fails or value is nil
// Returns: Bool representation of the value, or false if conversion fails
// Usage:
// result := SQLite.Parse(1).ToBool()
// return: true
//
// result := SQLite.Parse("true").ToBool()
// return: true
//
// result := SQLite.Parse(nil).ToBool()
// return: false
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

// Type returns the type name of the wrapped value
// Returns: String representation of the value's type, or "nil" if value is nil
// Usage:
// result := SQLite.Parse("hello").Type()
// return: "string"
//
// result := SQLite.Parse(42).Type()
// return: "int"
//
// result := SQLite.Parse(nil).Type()
// return: "nil"
func (sv SQLITE_VALUE) Type() string {
	result := "nil"
	if sv.Value != nil {
		result = reflect.TypeOf(sv.Value).String()
	}
	return result
}
