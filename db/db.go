package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("could not load database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	if err = DB.Ping(); err != nil {
		panic(fmt.Sprintf("could not connect to database: %v", err))
	}

	createTables()
	migrateEventsTable()
}

func createTables() {

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,		
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic(fmt.Sprintf("could not create users table: %v", err))
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic(fmt.Sprintf("could not create events table: %v", err))
	}
}

func migrateEventsTable() {
	rows, err := DB.Query(`PRAGMA table_info(events)`)
	if err != nil {
		panic(fmt.Sprintf("could not inspect events table: %v", err))
	}
	defer rows.Close()

	hasTitle := false
	hasName := false

	for rows.Next() {
		var cid int
		var columnName string
		var dataType string
		var notNull int
		var defaultValue any
		var pk int

		if err := rows.Scan(&cid, &columnName, &dataType, &notNull, &defaultValue, &pk); err != nil {
			panic(fmt.Sprintf("could not scan events table info: %v", err))
		}

		if columnName == "title" {
			hasTitle = true
		}
		if columnName == "name" {
			hasName = true
		}
	}

	if err = rows.Err(); err != nil {
		panic(fmt.Sprintf("could not read events table info: %v", err))
	}

	if hasTitle && !hasName {
		if _, err = DB.Exec(`ALTER TABLE events RENAME COLUMN title TO name`); err != nil {
			panic(fmt.Sprintf("could not migrate events table from title to name: %v", err))
		}
	}
}
