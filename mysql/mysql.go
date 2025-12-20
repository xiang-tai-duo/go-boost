// Package mysql
// File:        mysql.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/mysql/mysql.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: MySQL is a wrapper for MySQL instance operations, providing a set of methods for instance management and query execution.
// --------------------------------------------------------------------------------
package mysql

import (
	"crypto/md5"
	"crypto/sha3"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoNameStartsWithPackageName
type (
	MYSQL struct {
		Trace         bool
		Database      string
		Host          string
		InTransaction bool
		Password      string
		Port          int
		SqlDb         *sql.DB
		User          string
	}
	MYSQL_EXEC_CALLBACK func(sql string, err error) bool
	MYSQL_TX_CALLBACK   func(tx *MYSQL) error
	MYSQL_VALUE         struct {
		Name  string
		Value interface{}
	}
)

//goland:noinspection GoSnakeCaseUsage,SpellCheckingInspection,GoNameStartsWithPackageName,GoUnusedConst
const (
	MODULE_NAME_MYSQL                           = "mysql"
	MYSQL_CHARSET_DSN_FORMAT                    = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"
	MYSQL_CHARSET_DSN_NO_DB_FORMAT              = "%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local"
	MYSQL_CONNECTION_MAX_IDLE                   = 10
	MYSQL_CONNECTION_MAX_LIFETIME_HOUR          = 1
	MYSQL_CONNECTION_MAX_OPEN                   = 100
	MYSQL_DRIVER_NAME                           = "mysql"
	MYSQL_ERROR_BATCH_CANCELLED                 = "batch execution cancelled by callback"
	MYSQL_ERROR_BLOB_CHUNK_EMPTY                = "blob chunk cannot be empty"
	MYSQL_ERROR_BLOB_COLUMN_EMPTY               = "blob column name cannot be empty"
	MYSQL_ERROR_BLOB_COLUMN_INVALID             = "blob column name contains invalid characters"
	MYSQL_ERROR_BLOB_KEY_COLUMN_EMPTY           = "blob key column name cannot be empty"
	MYSQL_ERROR_BLOB_KEY_COLUMN_INVALID         = "blob key column name contains invalid characters"
	MYSQL_ERROR_BLOB_TABLE_EMPTY                = "blob table name cannot be empty"
	MYSQL_ERROR_BLOB_TABLE_INVALID              = "blob table name contains invalid characters"
	MYSQL_ERROR_COLUMN_EMPTY                    = "column name cannot be empty"
	MYSQL_ERROR_COLUMN_INVALID                  = "column name contains invalid characters"
	MYSQL_ERROR_COMMIT_FORMAT                   = "unable to commit transaction: %w"
	MYSQL_ERROR_CONNECT_FORMAT                  = "unable to connect to mysql server: %w"
	MYSQL_ERROR_CREATE_DATABASE_FORMAT          = "unable to create database: %w"
	MYSQL_ERROR_DANGEROUS_STATEMENT_FORMAT      = "dangerous SQL statement rejected by safe mode: %s"
	MYSQL_ERROR_DATABASE_EMPTY                  = "database cannot be empty"
	MYSQL_ERROR_DATABASE_INVALID                = "database name contains invalid characters"
	MYSQL_ERROR_GET_COLUMNS_FORMAT              = "unable to get columns: %w"
	MYSQL_ERROR_HOST_EMPTY                      = "host cannot be empty"
	MYSQL_ERROR_INSTANCE_NIL                    = "mysql instance is nil"
	MYSQL_ERROR_NOT_IN_TRANSACTION              = "not in transaction"
	MYSQL_ERROR_NOT_OPEN                        = "instance connection not open"
	MYSQL_ERROR_NO_STATEMENTS                   = "no valid sql statements found in file"
	MYSQL_ERROR_OPEN_FORMAT                     = "unable to open mysql connection: %w"
	MYSQL_ERROR_QUERY_EMPTY                     = "query cannot be empty"
	MYSQL_ERROR_QUERY_FORMAT                    = "unable to execute query: %w"
	MYSQL_ERROR_QUERY_INFORMATION_SCHEMA_FORMAT = "unable to query information_schema: %w"
	MYSQL_ERROR_QUERY_TABLE_FORMAT              = "failed to query table %s: %w"
	MYSQL_ERROR_READ_SQL_FILE_FORMAT            = "failed to read sql file: %w"
	MYSQL_ERROR_ROLLBACK_FORMAT                 = "unable to rollback transaction: %w"
	MYSQL_ERROR_ROLLBACK_NESTED_FORMAT          = "%w (and rollback failed: %v)"
	MYSQL_ERROR_ROW_ITERATION_FORMAT            = "error during row iteration: %w"
	MYSQL_ERROR_SCAN_ROW_FORMAT                 = "unable to scan row: %w"
	MYSQL_ERROR_SCAN_ROW_SIMPLE_FORMAT          = "failed to scan row: %w"
	MYSQL_ERROR_SQL_FILE_PATH_EMPTY             = "sql file path cannot be empty"
	MYSQL_ERROR_SQL_STATEMENTS_EMPTY            = "sql statements array cannot be empty"
	MYSQL_ERROR_START_TRANSACTION_FORMAT        = "unable to begin transaction: %w"
	MYSQL_ERROR_TABLE_EMPTY                     = "table name cannot be empty"
	MYSQL_ERROR_TABLE_INVALID                   = "table name contains invalid characters"
	MYSQL_ERROR_TRANSACTION_ALREADY_STARTED     = "transaction already started, nested transactions are not supported"
	MYSQL_ERROR_TRANSACTION_CALLBACK_NIL        = "transaction callback cannot be nil"
	MYSQL_ERROR_UPDATE_RECORD_FORMAT            = "failed to update record with id %d: %w"
	MYSQL_ERROR_USER_EMPTY                      = "user cannot be empty"
	MYSQL_IDENTIFIER_BACKTICK                   = "`"
	MYSQL_IDENTIFIER_BACKTICK_DOUBLE            = "``"
	MYSQL_IDENTIFIER_MAX_LENGTH                 = 64
	MYSQL_MD5_HASH_LENGTH                       = 32
	MYSQL_SHA3_HASH_LENGTH                      = 64
	MYSQL_SQL_APPEND_BLOB_CHUNK_FORMAT          = "UPDATE %s SET %s = CONCAT(IFNULL(%s, _binary ''), ?) WHERE %s = ?"
	MYSQL_SQL_BEGIN_TRANSACTION                 = "START TRANSACTION"
	MYSQL_SQL_COMMIT                            = "COMMIT"
	MYSQL_SQL_CREATE_DATABASE_FORMAT            = "CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	MYSQL_SQL_IS_DATABASE_EXISTS                = "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?"
	MYSQL_SQL_ROLLBACK                          = "ROLLBACK"
	MYSQL_SQL_SELECT_ID_COLUMN_FORMAT           = "SELECT id, %s FROM %s"
	MYSQL_SQL_UPDATE_COLUMN_BY_ID_FORMAT        = "UPDATE %s SET %s = ? WHERE id = ?"
	MYSQL_STRING_BOOLEAN_TRUE_ONE               = "1"
	MYSQL_STRING_BOOLEAN_TRUE_ON                = "on"
	MYSQL_STRING_BOOLEAN_TRUE_TRUE              = "true"
	MYSQL_STRING_BOOLEAN_TRUE_YES               = "yes"
	MYSQL_TYPE_NIL                              = "nil"
)

var (
	dangerousSqlKeywords = []string{"DROP", "TRUNCATE", "ALTER", "GRANT", "REVOKE", "RENAME", "SHUTDOWN", "CREATE", "REPLACE"}
	mysqlSafeMode        = true
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_MYSQL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_MYSQL, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_MYSQL, logger.SKIP_STACK_FRAMES_BASE)
}

func New(host string, port int, user string, password string, database string) *MYSQL {
	return &MYSQL{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
}

func (mysql *MYSQL) AppendBlobChunk(tableName string, blobColumn string, keyColumn string, keyValue interface{}, chunk []byte) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if tableName == "" {
		err = errors.New(MYSQL_ERROR_BLOB_TABLE_EMPTY)
	} else if blobColumn == "" {
		err = errors.New(MYSQL_ERROR_BLOB_COLUMN_EMPTY)
	} else if keyColumn == "" {
		err = errors.New(MYSQL_ERROR_BLOB_KEY_COLUMN_EMPTY)
	} else if len(chunk) == 0 {
		err = errors.New(MYSQL_ERROR_BLOB_CHUNK_EMPTY)
	} else if !mysql.IsValidIdentifier(tableName) {
		err = errors.New(MYSQL_ERROR_BLOB_TABLE_INVALID)
	} else if !mysql.IsValidIdentifier(blobColumn) {
		err = errors.New(MYSQL_ERROR_BLOB_COLUMN_INVALID)
	} else if !mysql.IsValidIdentifier(keyColumn) {
		err = errors.New(MYSQL_ERROR_BLOB_KEY_COLUMN_INVALID)
	} else {
		quotedTable := mysql.QuoteIdentifier(tableName)
		quotedBlob := mysql.QuoteIdentifier(blobColumn)
		quotedKey := mysql.QuoteIdentifier(keyColumn)
		query := fmt.Sprintf(MYSQL_SQL_APPEND_BLOB_CHUNK_FORMAT, quotedTable, quotedBlob, quotedBlob, quotedKey)
		if _, err = mysql.SqlDb.Exec(query, chunk, keyValue); err != nil {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_FORMAT, err)
		}
		mysql.LogSQLDebug(query, []interface{}{chunk, keyValue})
	}
	return err
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection
func (mysql *MYSQL) BeginTransaction() error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if mysql.InTransaction {
		err = errors.New(MYSQL_ERROR_TRANSACTION_ALREADY_STARTED)
	} else {
		_, err = mysql.SqlDb.Exec(MYSQL_SQL_BEGIN_TRANSACTION)
		mysql.LogSQLDebug(MYSQL_SQL_BEGIN_TRANSACTION, nil)
		if err == nil {
			mysql.InTransaction = true
		} else {
			err = fmt.Errorf(MYSQL_ERROR_START_TRANSACTION_FORMAT, err)
		}
	}
	return err
}

func (mysql *MYSQL) Close() {
	if mysql != nil {
		if mysql.SqlDb != nil {
			_ = mysql.SqlDb.Close()
			mysql.SqlDb = nil
		}
		mysql.Host = ""
		mysql.Port = 0
		mysql.User = ""
		mysql.Password = ""
		mysql.Database = ""
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection,DuplicatedCode
func (mysql *MYSQL) Commit() error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if !mysql.InTransaction {
		err = errors.New(MYSQL_ERROR_NOT_IN_TRANSACTION)
	} else {
		_, err = mysql.SqlDb.Exec(MYSQL_SQL_COMMIT)
		mysql.LogSQLDebug(MYSQL_SQL_COMMIT, nil)
		if err == nil {
			mysql.InTransaction = false
		} else {
			err = fmt.Errorf(MYSQL_ERROR_COMMIT_FORMAT, err)
		}
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult,DuplicatedCode
func (mysql *MYSQL) ConvertTextToMD5Text(tableName string, columnName string) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if tableName == "" {
		err = errors.New(MYSQL_ERROR_TABLE_EMPTY)
	} else if columnName == "" {
		err = errors.New(MYSQL_ERROR_COLUMN_EMPTY)
	} else if !mysql.IsValidIdentifier(tableName) {
		err = errors.New(MYSQL_ERROR_TABLE_INVALID)
	} else if !mysql.IsValidIdentifier(columnName) {
		err = errors.New(MYSQL_ERROR_COLUMN_INVALID)
	} else if err = mysql.BeginTransaction(); err == nil {
		defer func() {
			if err == nil {
				if err = mysql.Commit(); err != nil {
					err = fmt.Errorf(MYSQL_ERROR_COMMIT_FORMAT, err)
				}
			} else {
				originalErr := err
				if rollbackErr := mysql.Rollback(); rollbackErr != nil {
					err = fmt.Errorf(MYSQL_ERROR_ROLLBACK_NESTED_FORMAT, originalErr, rollbackErr)
				}
			}
		}()
		quotedTable := mysql.QuoteIdentifier(tableName)
		quotedColumn := mysql.QuoteIdentifier(columnName)
		query := fmt.Sprintf(MYSQL_SQL_SELECT_ID_COLUMN_FORMAT, quotedColumn, quotedTable)
		var rows *sql.Rows
		rows, err = mysql.SqlDb.Query(query)
		mysql.LogSQLDebug(query, nil)
		if err == nil {
			defer rows.Close()
			var id int64
			var value string
			for rows.Next() {
				if err = rows.Scan(&id, &value); err == nil {
					if value != "" && !mysql.IsMD5Hash(value) {
						newHashValue := mysql.MD5(value)
						updateQuery := fmt.Sprintf(MYSQL_SQL_UPDATE_COLUMN_BY_ID_FORMAT, quotedTable, quotedColumn)
						_, execErr := mysql.SqlDb.Exec(updateQuery, newHashValue, id)
						mysql.LogSQLDebug(updateQuery, []interface{}{newHashValue, id})
						if execErr != nil {
							err = fmt.Errorf(MYSQL_ERROR_UPDATE_RECORD_FORMAT, id, execErr)
							break
						}
					}
				} else {
					err = fmt.Errorf(MYSQL_ERROR_SCAN_ROW_SIMPLE_FORMAT, err)
					break
				}
			}
			if err == nil {
				if err = rows.Err(); err != nil {
					err = fmt.Errorf(MYSQL_ERROR_ROW_ITERATION_FORMAT, err)
				}
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_TABLE_FORMAT, tableName, err)
		}
	} else {
		err = fmt.Errorf(MYSQL_ERROR_START_TRANSACTION_FORMAT, err)
	}
	return err
}

//goland:noinspection GoUnhandledErrorResult,DuplicatedCode
func (mysql *MYSQL) ConvertTextToSHA3Text(tableName string, columnName string) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if tableName == "" {
		err = errors.New(MYSQL_ERROR_TABLE_EMPTY)
	} else if columnName == "" {
		err = errors.New(MYSQL_ERROR_COLUMN_EMPTY)
	} else if !mysql.IsValidIdentifier(tableName) {
		err = errors.New(MYSQL_ERROR_TABLE_INVALID)
	} else if !mysql.IsValidIdentifier(columnName) {
		err = errors.New(MYSQL_ERROR_COLUMN_INVALID)
	} else if err = mysql.BeginTransaction(); err == nil {
		defer func() {
			if err == nil {
				if err = mysql.Commit(); err != nil {
					err = fmt.Errorf(MYSQL_ERROR_COMMIT_FORMAT, err)
				}
			} else {
				originalErr := err
				if rollbackErr := mysql.Rollback(); rollbackErr != nil {
					err = fmt.Errorf(MYSQL_ERROR_ROLLBACK_NESTED_FORMAT, originalErr, rollbackErr)
				}
			}
		}()
		quotedTable := mysql.QuoteIdentifier(tableName)
		quotedColumn := mysql.QuoteIdentifier(columnName)
		query := fmt.Sprintf(MYSQL_SQL_SELECT_ID_COLUMN_FORMAT, quotedColumn, quotedTable)
		var rows *sql.Rows
		rows, err = mysql.SqlDb.Query(query)
		mysql.LogSQLDebug(query, nil)
		if err == nil {
			defer rows.Close()
			var id int64
			var value string
			for rows.Next() {
				if err = rows.Scan(&id, &value); err == nil {
					if value != "" && !mysql.IsSHA3Hash(value) {
						newHashValue := mysql.SHA3(value)
						updateQuery := fmt.Sprintf(MYSQL_SQL_UPDATE_COLUMN_BY_ID_FORMAT, quotedTable, quotedColumn)
						_, execErr := mysql.SqlDb.Exec(updateQuery, newHashValue, id)
						mysql.LogSQLDebug(updateQuery, []interface{}{newHashValue, id})
						if execErr != nil {
							err = fmt.Errorf(MYSQL_ERROR_UPDATE_RECORD_FORMAT, id, execErr)
							break
						}
					}
				} else {
					err = fmt.Errorf(MYSQL_ERROR_SCAN_ROW_SIMPLE_FORMAT, err)
					break
				}
			}
			if err == nil {
				if err = rows.Err(); err != nil {
					err = fmt.Errorf(MYSQL_ERROR_ROW_ITERATION_FORMAT, err)
				}
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_TABLE_FORMAT, tableName, err)
		}
	} else {
		err = fmt.Errorf(MYSQL_ERROR_START_TRANSACTION_FORMAT, err)
	}
	return err
}

func (mysql *MYSQL) Create(host string, port int, user string, password string, database string) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if host == "" {
		err = errors.New(MYSQL_ERROR_HOST_EMPTY)
	} else if user == "" {
		err = errors.New(MYSQL_ERROR_USER_EMPTY)
	} else if database == "" {
		err = errors.New(MYSQL_ERROR_DATABASE_EMPTY)
	} else {
		mysql.Close()
		dataSourceName := fmt.Sprintf(MYSQL_CHARSET_DSN_FORMAT, user, password, host, port, database)
		if mysql.SqlDb, err = sql.Open(MYSQL_DRIVER_NAME, dataSourceName); err == nil {
			mysql.SqlDb.SetMaxOpenConns(MYSQL_CONNECTION_MAX_OPEN)
			mysql.SqlDb.SetMaxIdleConns(MYSQL_CONNECTION_MAX_IDLE)
			mysql.SqlDb.SetConnMaxLifetime(time.Duration(MYSQL_CONNECTION_MAX_LIFETIME_HOUR) * time.Hour)
			if err = mysql.SqlDb.Ping(); err == nil {
				mysql.Host = host
				mysql.Port = port
				mysql.User = user
				mysql.Password = password
				mysql.Database = database
			} else {
				mysql.Close()
				err = fmt.Errorf(MYSQL_ERROR_CONNECT_FORMAT, err)
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_OPEN_FORMAT, err)
		}
	}
	return err
}

func (mysql *MYSQL) CreateDatabase(host string, port int, user string, password string, database string) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if host == "" {
		err = errors.New(MYSQL_ERROR_HOST_EMPTY)
	} else if user == "" {
		err = errors.New(MYSQL_ERROR_USER_EMPTY)
	} else if database == "" {
		err = errors.New(MYSQL_ERROR_DATABASE_EMPTY)
	} else if !mysql.IsValidIdentifier(database) {
		err = errors.New(MYSQL_ERROR_DATABASE_INVALID)
	} else {
		mysql.Close()
		dataSourceName := fmt.Sprintf(MYSQL_CHARSET_DSN_NO_DB_FORMAT, user, password, host, port)
		if mysql.SqlDb, err = sql.Open(MYSQL_DRIVER_NAME, dataSourceName); err == nil {
			mysql.SqlDb.SetMaxOpenConns(MYSQL_CONNECTION_MAX_OPEN)
			mysql.SqlDb.SetMaxIdleConns(MYSQL_CONNECTION_MAX_IDLE)
			mysql.SqlDb.SetConnMaxLifetime(time.Duration(MYSQL_CONNECTION_MAX_LIFETIME_HOUR) * time.Hour)
			if err = mysql.SqlDb.Ping(); err == nil {
				createQuery := fmt.Sprintf(MYSQL_SQL_CREATE_DATABASE_FORMAT, mysql.QuoteIdentifier(database))
				_, err = mysql.SqlDb.Exec(createQuery)
				mysql.LogSQLDebug(createQuery, nil)
				if err == nil {
					mysql.Host = host
					mysql.Port = port
					mysql.User = user
					mysql.Password = password
					mysql.Database = database
					_ = mysql.SqlDb.Close()
					mysql.SqlDb = nil
					err = mysql.Create(host, port, user, password, database)
				} else {
					mysql.Close()
					err = fmt.Errorf(MYSQL_ERROR_CREATE_DATABASE_FORMAT, err)
				}
			} else {
				mysql.Close()
				err = fmt.Errorf(MYSQL_ERROR_CONNECT_FORMAT, err)
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_OPEN_FORMAT, err)
		}
	}
	return err
}

func (mysql *MYSQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	err := error(nil)
	if err = mysql.IsValidateConnectionAndQuery(query); err == nil {
		result, err = mysql.SqlDb.Exec(query, args...)
		mysql.LogSQLDebug(query, args)
		if err != nil {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_FORMAT, err)
		}
	}
	return result, err
}

//goland:noinspection SqlDialectInspection
func (mysql *MYSQL) ExecNonQueries(queries []string, callback MYSQL_EXEC_CALLBACK) error {
	err := error(nil)
	if err = mysql.IsValidateConnection(); err == nil {
		if queries == nil || len(queries) == 0 {
			err = errors.New(MYSQL_ERROR_SQL_STATEMENTS_EMPTY)
		} else if err = mysql.BeginTransaction(); err == nil {
			defer func() {
				if err != nil {
					if rollbackErr := mysql.Rollback(); rollbackErr != nil {
						err = fmt.Errorf(MYSQL_ERROR_ROLLBACK_NESTED_FORMAT, err, rollbackErr)
					}
				}
			}()
			for _, query := range queries {
				if query != "" {
					if mysqlSafeMode && isDangerousStatement(query) {
						err = fmt.Errorf(MYSQL_ERROR_DANGEROUS_STATEMENT_FORMAT, query)
						break
					}
					_, execErr := mysql.SqlDb.Exec(query)
					mysql.LogSQLDebug(query, nil)
					if execErr != nil {
						err = fmt.Errorf(MYSQL_ERROR_QUERY_FORMAT, execErr)
						break
					}
					if callback != nil && !callback(query, err) {
						err = errors.New(MYSQL_ERROR_BATCH_CANCELLED)
						break
					}
				}
			}
			if err == nil {
				if err = mysql.Commit(); err != nil {
					err = fmt.Errorf(MYSQL_ERROR_COMMIT_FORMAT, err)
				}
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_START_TRANSACTION_FORMAT, err)
		}
	}
	return err
}

func (mysql *MYSQL) ExecNonQuery(query string, args ...interface{}) error {
	err := error(nil)
	if err = mysql.IsValidateConnectionAndQuery(query); err == nil {
		_, err = mysql.SqlDb.Exec(query, args...)
		mysql.LogSQLDebug(query, args)
		if err != nil {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_FORMAT, err)
		}
	} else {
		mysql.LogSQLDebug(query, args)
	}
	return err
}

//goland:noinspection DuplicatedCode
func (mysql *MYSQL) ExecuteQuery(query string, args ...interface{}) ([]MYSQL_VALUE, error) {
	results := make([]MYSQL_VALUE, 0)
	err := error(nil)
	if err = mysql.IsValidateConnectionAndQuery(query); err == nil {
		var rows *sql.Rows
		rows, err = mysql.SqlDb.Query(query, args...)
		mysql.LogSQLDebug(query, args)
		if err == nil {
			defer func() {
				_ = rows.Close()
			}()
			columns := make([]string, 0)
			if columns, err = rows.Columns(); err == nil {
				columnCount := len(columns)
				values := make([]interface{}, columnCount)
				valuesPointer := make([]interface{}, columnCount)
				for rows.Next() {
					for i := range values {
						valuesPointer[i] = &values[i]
					}
					if err = rows.Scan(valuesPointer...); err == nil {
						for i := 0; i < columnCount; i++ {
							value := values[i]
							if value != nil {
								if valueType := reflect.TypeOf(value); valueType.Kind() == reflect.Ptr {
									value = reflect.ValueOf(value).Elem().Interface()
								}
							}
							results = append(results, MYSQL_VALUE{
								Name:  columns[i],
								Value: value,
							})
						}
					} else {
						err = fmt.Errorf(MYSQL_ERROR_SCAN_ROW_FORMAT, err)
						break
					}
				}
				if err == nil {
					if err = rows.Err(); err != nil {
						err = fmt.Errorf(MYSQL_ERROR_ROW_ITERATION_FORMAT, err)
					}
				}
			} else {
				err = fmt.Errorf(MYSQL_ERROR_GET_COLUMNS_FORMAT, err)
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_FORMAT, err)
		}
	} else {
		mysql.LogSQLDebug(query, args)
	}
	return results, err
}

func (mysql *MYSQL) ExecuteQueryRows(query string, args ...interface{}) ([][]MYSQL_VALUE, error) {
	results := make([][]MYSQL_VALUE, 0)
	err := error(nil)
	if err = mysql.IsValidateConnectionAndQuery(query); err == nil {
		var rows *sql.Rows
		rows, err = mysql.SqlDb.Query(query, args...)
		mysql.LogSQLDebug(query, args)
		if err == nil {
			defer func() {
				_ = rows.Close()
			}()
			var columns []string
			if columns, err = rows.Columns(); err == nil {
				columnCount := len(columns)
				values := make([]interface{}, columnCount)
				valuesPointer := make([]interface{}, columnCount)
				for rows.Next() {
					for i := range values {
						valuesPointer[i] = &values[i]
					}
					if err = rows.Scan(valuesPointer...); err == nil {
						row := make([]MYSQL_VALUE, columnCount)
						for i := 0; i < columnCount; i++ {
							value := values[i]
							if value != nil {
								if valueType := reflect.TypeOf(value); valueType.Kind() == reflect.Ptr {
									value = reflect.ValueOf(value).Elem().Interface()
								}
							}
							row[i] = MYSQL_VALUE{
								Name:  columns[i],
								Value: value,
							}
						}
						results = append(results, row)
					} else {
						err = fmt.Errorf(MYSQL_ERROR_SCAN_ROW_FORMAT, err)
						break
					}
				}
				if err == nil {
					if err = rows.Err(); err != nil {
						err = fmt.Errorf(MYSQL_ERROR_ROW_ITERATION_FORMAT, err)
					}
				}
			} else {
				err = fmt.Errorf(MYSQL_ERROR_GET_COLUMNS_FORMAT, err)
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_QUERY_FORMAT, err)
		}
	} else {
		mysql.LogSQLDebug(query, args)
	}
	return results, err
}

//goland:noinspection GoUnusedExportedFunction
func GetSafeMode() bool {
	return mysqlSafeMode
}

//goland:noinspection GoUnhandledErrorResult
func (mysql *MYSQL) Hello(host string, port int, user string, password string) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if host == "" {
		err = errors.New(MYSQL_ERROR_HOST_EMPTY)
	} else if user == "" {
		err = errors.New(MYSQL_ERROR_USER_EMPTY)
	} else {
		dataSourceName := fmt.Sprintf(MYSQL_CHARSET_DSN_NO_DB_FORMAT, user, password, host, port)
		var probeDb *sql.DB
		if probeDb, err = sql.Open(MYSQL_DRIVER_NAME, dataSourceName); err == nil {
			defer probeDb.Close()
			if err = probeDb.Ping(); err != nil {
				err = fmt.Errorf(MYSQL_ERROR_CONNECT_FORMAT, err)
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_OPEN_FORMAT, err)
		}
	}
	return err
}

func (mysql *MYSQL) ImportSQLFile(sqlFilePath string) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if sqlFilePath == "" {
		err = errors.New(MYSQL_ERROR_SQL_FILE_PATH_EMPTY)
	} else {
		var absoluteFilePath string
		if absoluteFilePath, err = filepath.Abs(sqlFilePath); err == nil {
			var content []byte
			if content, err = os.ReadFile(absoluteFilePath); err == nil {
				sqlContent := string(content)
				statements := mysql.ParseSQLStatements(sqlContent)
				if len(statements) == 0 {
					err = errors.New(MYSQL_ERROR_NO_STATEMENTS)
				} else {
					err = mysql.ExecNonQueries(statements, nil)
				}
			} else {
				err = fmt.Errorf(MYSQL_ERROR_READ_SQL_FILE_FORMAT, err)
			}
		}
	}
	return err
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection,GoUnhandledErrorResult
func (mysql *MYSQL) IsDatabaseExists(host string, port int, user string, password string, database string) (bool, error) {
	result := false
	err := error(nil)
	if database == "" {
		err = errors.New(MYSQL_ERROR_DATABASE_EMPTY)
	} else if err = mysql.Hello(host, port, user, password); err == nil {
		dataSourceName := fmt.Sprintf(MYSQL_CHARSET_DSN_NO_DB_FORMAT, user, password, host, port)
		var probeDb *sql.DB
		if probeDb, err = sql.Open(MYSQL_DRIVER_NAME, dataSourceName); err == nil {
			defer probeDb.Close()
			row := probeDb.QueryRow(MYSQL_SQL_IS_DATABASE_EXISTS, database)
			if err = row.Scan(new("")); err == nil {
				result = true
			} else if errors.Is(err, sql.ErrNoRows) {
				err = nil
			} else {
				err = fmt.Errorf(MYSQL_ERROR_QUERY_INFORMATION_SCHEMA_FORMAT, err)
			}
		} else {
			err = fmt.Errorf(MYSQL_ERROR_OPEN_FORMAT, err)
		}
	}
	return result, err
}

//goland:noinspection DuplicatedCode
func (mysql *MYSQL) IsMD5Hash(value string) bool {
	result := false
	if len(value) == MYSQL_MD5_HASH_LENGTH {
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

//goland:noinspection DuplicatedCode
func (mysql *MYSQL) IsSHA3Hash(value string) bool {
	result := false
	if len(value) == MYSQL_SHA3_HASH_LENGTH {
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

func (mysql *MYSQL) IsValidIdentifier(name string) bool {
	result := false
	if name != "" && len(name) <= MYSQL_IDENTIFIER_MAX_LENGTH {
		result = true
		for i, ch := range name {
			isLetter := (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
			isDigit := ch >= '0' && ch <= '9'
			isUnderscore := ch == '_' || ch == '$'
			if !(isLetter || isDigit || isUnderscore) {
				result = false
				break
			}
			if i == 0 && isDigit {
				result = false
				break
			}
		}
	}
	return result
}

func (mysql *MYSQL) IsValidateColumnIdentifier(columnName string, columnInvalidErr string) error {
	err := error(nil)
	if columnName != "" {
		if !mysql.IsValidIdentifier(columnName) {
			err = errors.New(columnInvalidErr)
		}
	}
	return err
}

func (mysql *MYSQL) IsValidateConnection() error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	}
	return err
}

func (mysql *MYSQL) IsValidateConnectionAndQuery(query string) error {
	err := mysql.IsValidateConnection()
	if err == nil {
		err = mysql.IsValidateQuery(query)
	}
	return err
}

func (mysql *MYSQL) IsValidateConnectionAndTableName(tableName string, tableEmptyErr string) error {
	err := mysql.IsValidateConnection()
	if err == nil {
		err = mysql.IsValidateTableName(tableName, tableEmptyErr)
	}
	return err
}

func (mysql *MYSQL) IsValidateQuery(query string) error {
	err := error(nil)
	if query == "" {
		err = errors.New(MYSQL_ERROR_QUERY_EMPTY)
	}
	return err
}

func (mysql *MYSQL) IsValidateTableAndColumn(tableName string, columnName string, tableEmptyErr string, columnEmptyErr string, tableInvalidErr string, columnInvalidErr string) error {
	_ = columnEmptyErr
	err := mysql.IsValidateConnection()
	if err == nil {
		err = mysql.IsValidateTableName(tableName, tableEmptyErr)
	}
	if err == nil {
		err = mysql.IsValidateTableIdentifier(tableName, tableInvalidErr)
	}
	if err == nil {
		err = mysql.IsValidateColumnIdentifier(columnName, columnInvalidErr)
	}
	return err
}

func (mysql *MYSQL) IsValidateTableIdentifier(tableName string, tableInvalidErr string) error {
	err := error(nil)
	if !mysql.IsValidIdentifier(tableName) {
		err = errors.New(tableInvalidErr)
	}
	return err
}

func (mysql *MYSQL) IsValidateTableName(tableName string, tableEmptyErr string) error {
	err := error(nil)
	if tableName == "" {
		err = errors.New(tableEmptyErr)
	}
	return err
}

func (mysql *MYSQL) LogSQLDebug(query string, args []interface{}) {
	__debug(fmt.Sprintf("SQL: %s, args: %v", query, args))
}

func (mysql *MYSQL) MD5(value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}

func (mysql *MYSQL) Open(host string, port int, user string, password string, database string) error {
	return mysql.Create(host, port, user, password, database)
}

func (mysql *MYSQL) Parse(data interface{}) MYSQL_VALUE {
	return MYSQL_VALUE{Value: data}
}

func (mysql *MYSQL) ParseSQLStatements(sqlContent string) []string {
	statements := make([]string, 0)
	lines := strings.Split(sqlContent, "\n")
	var currentStatement strings.Builder
	inMultiLineComment := false
	inString := false
	stringDelimiter := byte(0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		i := 0
		for i < len(line) {
			if !inString && !inMultiLineComment && i+1 < len(line) && line[i] == '-' && line[i+1] == '-' {
				break
			}
			if !inString && !inMultiLineComment && i+1 < len(line) && line[i] == '/' && line[i+1] == '*' {
				inMultiLineComment = true
				i += 2
				continue
			}
			if inMultiLineComment && i+1 < len(line) && line[i] == '*' && line[i+1] == '/' {
				inMultiLineComment = false
				i += 2
				continue
			}
			if inMultiLineComment {
				i++
				continue
			}
			if (line[i] == '\'' || line[i] == '"') && (i == 0 || line[i-1] != '\\') {
				if !inString {
					inString = true
					stringDelimiter = line[i]
				} else if line[i] == stringDelimiter {
					inString = false
				}
			}
			if !inString && line[i] == ';' {
				currentStatement.WriteString(line[:i])
				stmt := strings.TrimSpace(currentStatement.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				currentStatement.Reset()
				line = strings.TrimSpace(line[i+1:])
				i = 0
				continue
			}
			i++
		}

		if line != "" {
			currentStatement.WriteString(line)
			currentStatement.WriteString(" ")
		}
	}

	remaining := strings.TrimSpace(currentStatement.String())
	if remaining != "" {
		statements = append(statements, remaining)
	}

	return statements
}

func (mysql *MYSQL) QuoteIdentifier(name string) string {
	return MYSQL_IDENTIFIER_BACKTICK + strings.ReplaceAll(name, MYSQL_IDENTIFIER_BACKTICK, MYSQL_IDENTIFIER_BACKTICK_DOUBLE) + MYSQL_IDENTIFIER_BACKTICK
}

//goland:noinspection SqlNoDataSourceInspection,SqlDialectInspection,DuplicatedCode
func (mysql *MYSQL) Rollback() error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if !mysql.InTransaction {
		err = errors.New(MYSQL_ERROR_NOT_IN_TRANSACTION)
	} else {
		_, err = mysql.SqlDb.Exec(MYSQL_SQL_ROLLBACK)
		mysql.LogSQLDebug(MYSQL_SQL_ROLLBACK, nil)
		if err == nil {
			mysql.InTransaction = false
		} else {
			err = fmt.Errorf(MYSQL_ERROR_ROLLBACK_FORMAT, err)
		}
	}
	return err
}

func (mysql *MYSQL) SHA3(value string) string {
	hash := sha3.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

//goland:noinspection GoUnusedExportedFunction
func SetSafeMode(enabled bool) {
	mysqlSafeMode = enabled
}

func (mysql *MYSQL) TEXT(str string) string {
	if mysql != nil {
		str = strings.ReplaceAll(str, "\\", "\\\\")
		str = strings.ReplaceAll(str, "'", "\\'")
		str = strings.ReplaceAll(str, "\"", "\\\"")
		str = strings.ReplaceAll(str, "\n", "\\n")
		str = strings.ReplaceAll(str, "\r", "\\r")
		str = strings.ReplaceAll(str, "\t", "\\t")
		str = strings.ReplaceAll(str, "\x00", "\\0")
		str = strings.ReplaceAll(str, "\x1a", "\\Z")
	}
	return str
}

func (value MYSQL_VALUE) ToBool() bool {
	result := false
	if value.Value != nil {
		switch val := value.Value.(type) {
		case bool:
			result = val
		case int:
			if val != 0 {
				result = true
			}
		case int64:
			if val != 0 {
				result = true
			}
		case float32:
			if val != 0 {
				result = true
			}
		case float64:
			if val != 0 {
				result = true
			}
		case string:
			strVal := strings.ToLower(val)
			if strVal == MYSQL_STRING_BOOLEAN_TRUE_TRUE || strVal == MYSQL_STRING_BOOLEAN_TRUE_ONE || strVal == MYSQL_STRING_BOOLEAN_TRUE_YES || strVal == MYSQL_STRING_BOOLEAN_TRUE_ON {
				result = true
			}
		case []byte:
			strVal := strings.ToLower(string(val))
			if strVal == MYSQL_STRING_BOOLEAN_TRUE_TRUE || strVal == MYSQL_STRING_BOOLEAN_TRUE_ONE || strVal == MYSQL_STRING_BOOLEAN_TRUE_YES || strVal == MYSQL_STRING_BOOLEAN_TRUE_ON {
				result = true
			}
		}
	}
	return result
}

func (value MYSQL_VALUE) ToBytes() []byte {
	result := make([]byte, 0)
	if value.Value != nil {
		switch val := value.Value.(type) {
		case []byte:
			result = val
		case string:
			result = []byte(val)
		}
	}
	return result
}

func (value MYSQL_VALUE) ToFloat() float64 {
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
		case []byte:
			if f, err := strconv.ParseFloat(string(val), 64); err == nil {
				result = f
			}
		}
	}
	return result
}

func (value MYSQL_VALUE) ToInt() int {
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
		case []byte:
			if i, err := strconv.Atoi(string(val)); err == nil {
				result = i
			}
		}
	}
	return result
}

func (value MYSQL_VALUE) ToString() string {
	result := ""
	if value.Value != nil {
		switch val := value.Value.(type) {
		case string:
			result = val
		case []byte:
			result = string(val)
		case bool:
			result = strconv.FormatBool(val)
		case int:
			result = strconv.Itoa(val)
		case int64:
			result = strconv.FormatInt(val, 10)
		case float32:
			result = strconv.FormatFloat(float64(val), 'f', -1, 32)
		case float64:
			result = strconv.FormatFloat(val, 'f', -1, 64)
		case time.Time:
			result = val.Format(time.RFC3339)
		default:
			result = fmt.Sprintf("%v", val)
		}
	}
	return result
}

func (mysql *MYSQL) Transaction(callback MYSQL_TX_CALLBACK) error {
	err := error(nil)
	if mysql == nil {
		err = errors.New(MYSQL_ERROR_INSTANCE_NIL)
	} else if mysql.SqlDb == nil {
		err = errors.New(MYSQL_ERROR_NOT_OPEN)
	} else if callback == nil {
		err = errors.New(MYSQL_ERROR_TRANSACTION_CALLBACK_NIL)
	} else if mysql.InTransaction {
		err = errors.New(MYSQL_ERROR_TRANSACTION_ALREADY_STARTED)
	} else if err = mysql.BeginTransaction(); err == nil {
		if err = callback(mysql); err == nil {
			if err = mysql.Commit(); err != nil {
				err = fmt.Errorf(MYSQL_ERROR_COMMIT_FORMAT, err)
			}
		} else {
			originalErr := err
			if rollbackErr := mysql.Rollback(); rollbackErr != nil {
				err = fmt.Errorf(MYSQL_ERROR_ROLLBACK_NESTED_FORMAT, originalErr, rollbackErr)
			}
		}
	} else {
		err = fmt.Errorf(MYSQL_ERROR_START_TRANSACTION_FORMAT, err)
	}
	return err
}

func (value MYSQL_VALUE) Type() string {
	result := MYSQL_TYPE_NIL
	if value.Value != nil {
		result = reflect.TypeOf(value.Value).String()
	}
	return result
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_MYSQL, logger.SKIP_STACK_FRAMES_BASE)
}

func isDangerousStatement(query string) bool {
	result := false
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
	if stripped != "" {
		end := len(stripped)
		for i, ch := range stripped {
			if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' || ch == '(' || ch == ';' {
				end = i
				break
			}
		}
		firstWord := strings.ToUpper(stripped[:end])
		for _, keyword := range dangerousSqlKeywords {
			if firstWord == keyword {
				result = true
				break
			}
		}
	}
	return result
}
