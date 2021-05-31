package user

import (
	"database/sql"
	"log"
)

type User struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) *User {

	stmt, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS 
				data (rollno INTEGER NOT NULL PRIMARY KEY UNIQUE, 
				name TEXT ) 
	`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()

	return &User{
		DB: db,
	}
}

func (user *User) Add(userdata UserData) {
	stmt, err := user.DB.Prepare(`
			INSERT INTO data 
				(rollno, name) VALUES(?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec(userdata.Rollno, userdata.Name)

}
