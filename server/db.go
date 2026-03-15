package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	_ "github.com/mattn/go-sqlite3" // blank import the sqlite3 driver to register it
)

func connectToDatabase(path string) (*sql.DB, error) {
	fmt.Println("Attempting to execute database connection")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		fmt.Println("Connection secured, but pinging database failed - closing connection")
		_ = db.Close()
		return nil, err
	}
	fmt.Println("Connected successfully to SQLite database")
	return db, nil
}

// create tables if they don't exist
func initSchema(db *sql.DB) {
	schema, _ := os.ReadFile("server/database/setup.sql")
	for _, stmt := range strings.Split(string(schema), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			db.Exec(stmt)
		}
	}
}