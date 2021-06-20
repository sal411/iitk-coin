package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"
)

type User struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) *User {

	stmt, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS 
				data (rollno TEXT NOT NULL PRIMARY KEY UNIQUE, 
				name TEXT,
				password TEXT ) 
	`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()

	return &User{
		DB: db,
	}
}

func (user *User) Add(userdata models.UserData) error {
	stmt, err := user.DB.Prepare(`
			INSERT INTO data 
				(rollno, name, password) VALUES(?, ?, ?)
	`)
	utils.PrintError(err)
	stmt.Exec(userdata.Rollno, userdata.Name, userdata.Password)
	if err != nil {
		return err
	}
	return nil

}
