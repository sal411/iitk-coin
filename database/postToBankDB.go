package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sal411/iitk-coin/utils"
)

var MaxCoins float64
var Options = sql.TxOptions{
	Isolation: sql.LevelSerializable,
}

// function to update coins in roll no
func WriteCoins(rollno string, numberOfCoins string, remarks string) (string, error) {

	db := utils.ConnectDB()
	godotenv.Load()
	MaxCoins, _ = strconv.ParseFloat(os.Getenv("MAXCOINS"), 32)

	coins_number, e := strconv.ParseFloat(numberOfCoins, 32)
	if e != nil {
		return "Coins not valid ", e
	}
	_, _, err := GetUserFromRollNo(rollno)
	if err != nil {
		return "User not present ", err
	}

	tx, _ := db.BeginTx(context.Background(), &Options)

	res, execErr := tx.Exec(`UPDATE bank SET coins = coins + ? WHERE rollno= ? AND coins + ?<= ?;`, coins_number, rollno, coins_number, MaxCoins)
	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return "", execErr
		}
		overflowError := errors.New("Balance cannot exceed " + fmt.Sprintf("%f", MaxCoins))
		return "", overflowError
	}
	_, err = tx.Exec(`INSERT INTO rewards (user,amount,remarks,time) VALUES (?,?,?,?)`, rollno, coins_number, remarks, time.Now())
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return "Some error occured in the transaction, please try again later ", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return "Coins added sucessfully ", nil

}

// function to transfer coins between two roll numbers
func TransferCoin(firstRollno string, secondRollno string, transferAmount float64) (float64, error) {

	db := utils.ConnectDB()

	if firstRollno == secondRollno {
		return 0, nil
	}
	_, _, err := GetUserFromRollNo(firstRollno)
	if err != nil {
		return 0, errors.New("user " + firstRollno + " not present ")
	}
	_, _, err = GetUserFromRollNo(secondRollno)
	if err != nil {
		return 0, errors.New("user " + secondRollno + " not present ")
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

	batch1 := firstRollno[0:2]
	batch2 := secondRollno[0:2]
	var taxRate float32 = 0.02
	if batch1 != batch2 {
		taxRate = 0.33
	}
	taxAmount := taxRate * float32(transferAmount)
	res, execErr := tx.Exec("UPDATE bank SET coins = coins - (?+?) WHERE rollno=? AND  coins - (?+?) >= 0 ", transferAmount, taxAmount, firstRollno, transferAmount, taxAmount)
	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return 0, err
		}

		balanceError := errors.New("not enough balance  ")
		return 0, balanceError

	}

	res, execErr = tx.Exec("UPDATE bank SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", transferAmount, secondRollno, transferAmount, MaxCoins)

	rowsAffected, _ = res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return 0, execErr
		}
		overflowError := errors.New("Balance cannot exceed " + fmt.Sprintf("%f", MaxCoins))
		return 0, overflowError
	}

	_, execErr = tx.Exec(`INSERT INTO transfers (TransferFrom,TransferTo,amount,tax,time) VALUES (?,?,?,?,?)`, firstRollno, secondRollno, transferAmount, taxAmount, time.Now())
	if execErr != nil {
		_ = tx.Rollback()
		return 0, execErr
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return float64(taxAmount), nil
}
