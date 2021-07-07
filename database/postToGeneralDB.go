package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sal411/iitk-coin/utils"
)

var MinEvents int

var Option = sql.TxOptions{
	Isolation: sql.LevelSerializable,
}

func RedeemCoinsDb(roll_no string, item_id int) (float64, error) {

	db := utils.ConnectDB()
	godotenv.Load()

	MinEvents, _ = strconv.Atoi(os.Getenv("MINEVENTS"))

	numEvents, _ := GetNumEvents(roll_no)
	if numEvents < MinEvents {
		return 0, errors.New("You need to participate in at least " + strconv.Itoa(MinEvents) + " events to clam a reward ")
	}
	cost, available, err := getItemFromId(item_id)
	if err != nil {
		return 0, err
	}
	if available == 0 {
		return 0, errors.New("item not available, please select another item ")
	}

	_, _, err = GetUserFromRollNo(roll_no)
	if err != nil {
		return 0, errors.New("user " + roll_no + " not present ")
	}
	var options = sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	tx, err := db.BeginTx(context.Background(), &options)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return 0, err
	}
	res, err := tx.Exec(`UPDATE bank SET coins = coins - ? WHERE rollno= ? AND coins - ? >=0 `, cost, roll_no, cost)
	rowsAffected, _ := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, errors.New("insufficient coins to claim this item ")
	}
	res, err = tx.Exec(`UPDATE items SET available = available -1 WHERE id = ? `, item_id)
	rowsAffected, _ = res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, errors.New("error occured while transaction please try later ")
	}
	_, err = tx.Exec(`INSERT INTO redeems (user,item,time) VALUES (?,?,?)`, roll_no, item_id, time.Now())
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	coins, _ := GetCoinsFromRollno(roll_no)

	return coins, err
}

func WriteItems(item_id int, cost string, number int) (string, error) { // cpnvert this into a transaction

	db := utils.ConnectDB()
	var err error
	cost_number, e := strconv.ParseFloat(cost, 32)
	if e != nil {
		return "Coins not valid ", e
	}

	tx, _ := db.BeginTx(context.Background(), &Option)

	_, err = tx.Exec(`INSERT INTO items (id,cost,available) VALUES (?,?,?) ON CONFLICT(id) DO UPDATE SET available = available + ? ;`, item_id, cost_number, number, number)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return "Some error occured in the transaction, please try again later ", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}
	return "success", e

}
