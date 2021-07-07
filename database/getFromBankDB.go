package database

import (
	"github.com/sal411/iitk-coin/utils"
)

func GetCoinsFromRollno(rollno string) (float64, error) {
	db := utils.ConnectDB()

	statement, _ :=
		db.Prepare("CREATE TABLE IF NOT EXISTS bank (rollno TEXT PRIMARY KEY ,coins INT)")
	statement.Exec()

	sqlStatement := `SELECT coins FROM bank WHERE rollno= $1;`
	row := db.QueryRow(sqlStatement, rollno)

	var coins float64
	err := row.Scan(&coins)

	if err != nil {
		return 0, err
	}
	return coins, nil

}

func AccountExists(rollno string) bool {
	db := utils.ConnectDB()

	rows, _ := db.Query("SELECT coin FROM bank WHERE rollno = $1", rollno)

	var coins string
	err := rows.Scan(&coins)

	return err == nil
}
