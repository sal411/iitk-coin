package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/sal411/iitk-coin/utils"
)

// function to update coins in roll no
func WriteCoins(rollno string, coins string) error {

	_, err := GetUserFromRollNo(rollno)

	if err != nil {
		return err
	}

	total_coins, err := GetCoinsFromRollno(rollno)
	if err != nil {
		return err
	}

	total_coins_f, _ := strconv.ParseFloat(total_coins, 64)
	coins_f, _ := strconv.ParseFloat(coins, 64)

	total_coins = fmt.Sprintf("%f", total_coins_f+coins_f)

	db := utils.ConnectDB()
	statement, _ :=
		db.Prepare(`UPDATE bank SET coins = $1 WHERE rollno= $2;`)
	_, err = statement.Exec(total_coins, rollno)
	if err != nil {
		return err
	}

	return nil

}

// function to transfer coins between two roll numbers
func TransferCoin(firstRollno string, secondRollno string, transferAmount int) error {
	if firstRollno == secondRollno {
		return nil
	}
	db, _ := sql.Open("sqlite3", "./user.db")
	var options = sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	tx, err := db.BeginTx(context.Background(), &options)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return err
	}

	res, execErr := tx.Exec("UPDATE bank SET coins = coins - ? WHERE rollno=? AND coins - ? >= 0", transferAmount, firstRollno, transferAmount)

	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return err
		}

		balanceError := errors.New("not enough balance ")
		return balanceError

	}

	res, execErr = tx.Exec("UPDATE bank SET coins = coins + ? WHERE rollno=? ", transferAmount, secondRollno)

	rowsAffected, _ = res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
