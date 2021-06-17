package utils

import (
	_ "database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func GetHashedPassword(rollno string) string {

	var db = ConnectDB()

	integerRollNo, _ := strconv.Atoi(rollno)

	sqlStatement := `SELECT password FROM data WHERE rollno= $1;`
	row := db.QueryRow(sqlStatement, integerRollNo)

	var hashedPassword string
	row.Scan(&hashedPassword)
	//fmt.Println("hey getting hashed password")
	//fmt.Println(hashedPassword)
	return (hashedPassword)

}
