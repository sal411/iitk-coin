package database

import (
	"strconv"

	"github.com/sal411/iitk-coin/utils"
)

func GetCoinsFromRollno(rollno string) (int, error) {
	db := utils.ConnectDB()
	var err error

	if AccountExists(rollno) {
		integerRollNo, _ := strconv.Atoi(rollno)

		sqlStatement := `SELECT coin FROM bank WHERE rollno= $1;`
		row := db.QueryRow(sqlStatement, integerRollNo)

		var coin int
		row.Scan(&coin)
		err = row.Scan(&coin)

		if err != nil {
			return 0, err
		}

		return coin, nil
	}
	return 0, err

}

func AccountExists(rollno string) bool {
	db := utils.ConnectDB()

	rows, _ := db.Query("SELECT coin FROM bank WHERE rollno = $1", rollno)

	var coins int
	err := rows.Scan(&coins)

	return err == nil
}
