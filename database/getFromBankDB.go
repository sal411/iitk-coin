package database

import (
	"github.com/sal411/iitk-coin/utils"
)

func GetCoinsFromRollno(rollno string) (string, error) {
	db := utils.ConnectDB()
	var err error

	if AccountExists(rollno) {

		sqlStatement := `SELECT coin FROM bank WHERE rollno= $1;`
		row := db.QueryRow(sqlStatement, rollno)

		var coin string
		row.Scan(&coin)
		err = row.Scan(&coin)

		if err != nil {
			return "", err
		}

		return coin, nil
	}
	return "", err

}

func AccountExists(rollno string) bool {
	db := utils.ConnectDB()

	rows, _ := db.Query("SELECT coin FROM bank WHERE rollno = $1", rollno)

	var coins string
	err := rows.Scan(&coins)

	return err == nil
}
