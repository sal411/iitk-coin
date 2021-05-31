/*
@author Vatsal Chaudhary
@date 31/5/21

 properties of user are defined in user.go
-Add
	this takes a object UserData and adds it to the database file

main.go
	established a connection with database
	passed an instance, with will give access to methods of user.go
*/

package main

import (
	"database/sql"
	"iitk-coins/iitk-coins/user"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}

	item := user.NewUser(db)

	newUserData := user.UserData{
		Rollno: 190951,
		Name:   "Vatsal Chaudhary",
	}

	item.Add(newUserData)

	db.Close()
}
