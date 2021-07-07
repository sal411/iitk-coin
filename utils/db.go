package utils

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var MaxCoins float64
var MinEvents int

func ConnectDB() *sql.DB {

	godotenv.Load()
	MaxCoins, _ = strconv.ParseFloat(os.Getenv("MAXCOINS"), 32)
	MinEvents, _ = strconv.Atoi(os.Getenv("MINEVENTS"))

	db, err := sql.Open("sqlite3", "./users.db")
	PrintError(err)

	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS data (
    name         TEXT,
    rollno       TEXT PRIMARY KEY,
    password     TEXT,
    account_type TEXT
);
`)
	statement.Exec()

	statement, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS transfers (
		TransferFrom TEXT          NOT NULL,
		TransferTo   TEXT          NOT NULL,
		amount       DOUBLE (7, 2) NOT NULL,
		tax          DOUBLE (7, 2) NOT NULL,
		time         DATETIME
	);	
	`)
	statement.Exec()

	statement, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS rewards (
		user    TEXT     NOT NULL,
		amount  INTEGER  NOT NULL,
		remarks TEXT,
		time    DATETIME
	);
		
	`)
	statement.Exec()

	statement, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS redeems (
		user TEXT     NOT NULL,
		item          REFERENCES items (id),
		time DATETIME
	);
	`)
	statement.Exec()

	statement, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS  items (
		id        INTEGER        PRIMARY KEY,
		cost      DECIMAL (7, 2),
		available INTEGER        NOT NULL
	);
	
	`)
	statement.Exec()

	statement, _ = db.Prepare(`CREATE TABLE IF NOT EXISTS bank (
		rollno TEXT           PRIMARY KEY,
		coins  DECIMAL (7, 2) 
	);	
	`)
	statement.Exec()

	if err != nil {
		return db
	}
	return nil
}
