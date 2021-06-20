package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/sal411/iitk-coin/utils"
)

// function to update coins in roll no
func WriteCoins(rollno string, coins string) error {
	coinsInteger, e := strconv.Atoi(coins)
	if e != nil {
		return e
	}

	_, err := GetUserFromRollNo(rollno)

	if err != nil {
		return err
	}
	var success bool
	total_coins, err := GetCoinsFromRollno(rollno)
	if success {
		return err
	}

	total_coins = total_coins + coinsInteger
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
	db, _ := sql.Open("sqlite3", "./database/user.db")
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