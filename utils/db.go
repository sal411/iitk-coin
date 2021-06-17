package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {

	db, err := sql.Open("sqlite3", "./users.db")
	PrintError(err)

	return db
}
