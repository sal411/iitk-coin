package database

import (
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

func GetUserFromRollNo(rollno string) (string, string, error) {

	var db = utils.ConnectDB()
	sqlStatement := `SELECT name,account_type FROM user WHERE rollno= $1;`
	row := db.QueryRow(sqlStatement, rollno)
	var userName string
	var userType string
	err := row.Scan(&userName, &userType)
	if err != nil {
		return "", "", err
	}

	return userName, userType, nil
}
