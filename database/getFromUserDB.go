package database

import (
	"database/sql"

	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sal411/iitk-coin/utils"
)

func GetHashedPassword(rollno string) string {

	var db = utils.ConnectDB()

	integerRollNo, _ := strconv.Atoi(rollno)

	sqlStatement := `SELECT password FROM data WHERE rollno= $1;`
	row := db.QueryRow(sqlStatement, integerRollNo)

	var hashedPassword string
	row.Scan(&hashedPassword)
	//fmt.Println("hey getting hashed password")
	//fmt.Println(hashedPassword)
	return (hashedPassword)

}

func GetUserFromRollNo(rollno string) (*sql.Row, error) {

	sqlStatement := `SELECT name FROM user WHERE rollno= $1;`
	db := utils.ConnectDB()
	row := db.QueryRow(sqlStatement, rollno)
	err := row.Scan(&rollno)
	//fmt.Println(hashed_password)
	if err != nil {
		return nil, err
	}
	return row, nil
}
