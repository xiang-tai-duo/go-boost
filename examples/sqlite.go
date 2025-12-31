// --------------------------------------------------------------------------------
// File:        sqlite.go
// Author:      TRAE AI
// Created:     12/30/2025 11:03:46
// Description: Example for SQLITE utility functions
// --------------------------------------------------------------------------------

package main

import (
	"fmt"

	. "github.com/xiang-tai-duo/go-boost"
)

func main() {

	// Create a new SQLITE instance
	sqlite := &SQLITE{}

	// Set Trace to true for debugging
	sqlite.Trace = true

	// Database file path
	dbPath := "./test.db"

	// Example: Create a new database
	fmt.Println("--- Creating Database ---")
	err := sqlite.Create(dbPath)
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		return
	}
	fmt.Printf("Database created successfully: %s\n", dbPath)

	// Example: Create a table
	fmt.Println("\n--- Creating Table ---")
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE,
		age INTEGER,
		active BOOLEAN DEFAULT true
	);
	`
	err = sqlite.ExecNonQuery(createTableQuery)
	if err != nil {
		fmt.Printf("Error creating table: %v\n", err)
		return
	}
	fmt.Println("Table 'users' created successfully")

	// Example: Insert data
	fmt.Println("\n--- Inserting Data ---")
	insertQueries := []string{
		`INSERT INTO users (name, email, age) VALUES ('John Doe', 'john@example.com', 30);`,
		`INSERT INTO users (name, email, age, active) VALUES ('Jane Smith', 'jane@example.com', 25, true);`,
		`INSERT INTO users (name, email, age, active) VALUES ('Bob Johnson', 'bob@example.com', 35, false);`,
	}
	for _, query := range insertQueries {
		err = sqlite.ExecNonQuery(query)
		if err != nil {
			fmt.Printf("Error inserting data: %v\n", err)
			return
		}
	}
	fmt.Println("3 users inserted successfully")

	// Example: Query data
	fmt.Println("\n--- Querying Data ---")
	selectQuery := "SELECT * FROM users;"
	results, err := sqlite.ExecuteQuery(selectQuery)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Display query results
	fmt.Printf("Query returned %d values\n", len(results))
	if len(results) > 0 {
		// Group results by row (since each value is a separate SQLITE_VALUE)
		// Assuming all rows have the same number of columns
		columns := []string{"id", "name", "email", "age", "active"}
		rows := len(results) / len(columns)
		for i := 0; i < rows; i++ {
			fmt.Printf("Row %d:\n", i+1)
			for j := 0; j < len(columns); j++ {
				index := i*len(columns) + j
				if index < len(results) {
					val := results[index]
					fmt.Printf("  %s: %v (Type: %s)\n", val.Name, val.Value, val.Type())

					// Example: Using SQLITE_VALUE conversion methods
					if val.Name == "name" {
						fmt.Printf("    ToString(): %s\n", val.ToString())
					}
					if val.Name == "age" {
						fmt.Printf("    ToInt(): %d\n", val.ToInt())
						fmt.Printf("    ToFloat(): %.2f\n", val.ToFloat())
					}
					if val.Name == "active" {
						fmt.Printf("    ToBool(): %v\n", val.ToBool())
					}
				}
			}
		}
	}

	// Example: Update data
	fmt.Println("\n--- Updating Data ---")
	updateQuery := "UPDATE users SET age = 31 WHERE name = 'John Doe';"
	err = sqlite.ExecNonQuery(updateQuery)
	if err != nil {
		fmt.Printf("Error updating data: %v\n", err)
		return
	}
	fmt.Println("User 'John Doe' updated successfully")

	// Example: Delete data
	fmt.Println("\n--- Deleting Data ---")
	deleteQuery := "DELETE FROM users WHERE name = 'Bob Johnson';"
	err = sqlite.ExecNonQuery(deleteQuery)
	if err != nil {
		fmt.Printf("Error deleting data: %v\n", err)
		return
	}
	fmt.Println("User 'Bob Johnson' deleted successfully")

	// Example: Query data after update/delete
	fmt.Println("\n--- Querying Updated Data ---")
	results, err = sqlite.ExecuteQuery(selectQuery)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return
	}
	fmt.Printf("Query returned %d values\n", len(results))
	if len(results) > 0 {
		columns := []string{"id", "name", "email", "age", "active"}
		rows := len(results) / len(columns)
		for i := 0; i < rows; i++ {
			fmt.Printf("Row %d:\n", i+1)
			for j := 0; j < len(columns); j++ {
				index := i*len(columns) + j
				if index < len(results) {
					val := results[index]
					fmt.Printf("  %s: %v\n", val.Name, val.Value)
				}
			}
		}
	}

	// Example: Parse data
	fmt.Println("\n--- Parsing Data ---")
	intVal := sqlite.Parse(42)
	stringVal := sqlite.Parse("test")
	boolVal := sqlite.Parse(true)
	floatVal := sqlite.Parse(3.14)
	fmt.Printf("Parse(42): %v (Type: %s)\n", intVal.Value, intVal.Type())
	fmt.Printf("Parse('test'): %v (Type: %s)\n", stringVal.Value, stringVal.Type())
	fmt.Printf("Parse(true): %v (Type: %s)\n", boolVal.Value, boolVal.Type())
	fmt.Printf("Parse(3.14): %v (Type: %s)\n", floatVal.Value, floatVal.Type())

	// Example: Close database
	fmt.Println("\n--- Closing Database ---")
	sqlite.Close()
	fmt.Println("Database closed successfully")

	// Example: Reopen database
	fmt.Println("\n--- Reopening Database ---")
	err = sqlite.Open(dbPath)
	if err != nil {
		fmt.Printf("Error reopening database: %v\n", err)
		return
	}
	fmt.Printf("Database reopened successfully: %s\n", dbPath)

	// Example: Drop table and clean up
	fmt.Println("\n--- Cleaning Up ---")
	dropTableQuery := "DROP TABLE IF EXISTS users;"
	err = sqlite.ExecNonQuery(dropTableQuery)
	if err != nil {
		fmt.Printf("Error dropping table: %v\n", err)
		return
	}
	fmt.Println("Table 'users' dropped successfully")

	// Close database again
	sqlite.Close()
	fmt.Println("Database closed successfully")
	fmt.Println("\n--- SQLITE Examples Complete ---")
}
