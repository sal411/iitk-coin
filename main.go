package main

import (
	"database/sql"
	"fmt"
	"strconv"

    _ "github.com/mattn/go-sqlite3"
)

/*
CREATE TABLE IF NOT EXIST "users" (
	"roll no"	INTEGER NOT NULL UNIQUE,
	"coins"	INTEGER,
	PRIMARY KEY("roll no")
);
*/

func main() {
	db, _ := sql.Open("sqlite3", "./users.db")

	// statement, _ := db.Prepare(`
	// 	INSERT INTO users (roll no, coins) VALUES(? , ?)
	// `)

	// statement.Exec("190951", 14)

	statement, _ := db.Prepare(
		"CREATE TABLE IF NOT EXIST users (
			roll no	INTEGER PRIMARY KEY NOT NULL UNIQUE,
			coins INTEGER)"
	)

	statement.Exec()

	fmt.Println("hello")
}
