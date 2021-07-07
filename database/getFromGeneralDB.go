package database

import (
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sal411/iitk-coin/utils"
)

func getItemFromId(item_id int) (float64, int, error) {
	var db = utils.ConnectDB()
	var cost float64
	var available int

	sqlStatement := `SELECT cost,available FROM items WHERE id= $1;`
	row := db.QueryRow(sqlStatement, strconv.Itoa(item_id))

	err := row.Scan(&cost, &available)
	if err != nil {
		return 0, 0, err
	}
	return cost, available, nil
}

func GetNumEvents(rollno string) (int, error) { // returns the number of awards given to a user
	var number int
	var db = utils.ConnectDB()

	sqlStatement := `SELECT COUNT(user)
	FROM rewards
	WHERE user = $1;`

	row := db.QueryRow(sqlStatement, rollno)
	err := row.Scan(&number)
	if err != nil {
		return 0, err
	}
	return number, nil
}
