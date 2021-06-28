package database

import (
	"database/sql"
	"log"

	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"
)

type Bank struct {
	DB *sql.DB
}

func NewBank(db *sql.DB) *Bank {

	stmt, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS
			bank ( rollno TEXT NOT NULL PRIMARY KEY UNIQUE,
				coin TEXT
			)
	`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec()

	return &Bank{
		DB: db,
	}

}

func (bank *Bank) OpenAccount(bankdata models.BankData) error {
	stmt, err := bank.DB.Prepare(`
			INSERT INTO bank 
				(rollno, coin) VALUES(?, ?)
	`)
	utils.PrintError(err)
	stmt.Exec(bankdata.Rollno, bankdata.Coin)
	if err != nil {
		return err
	}
	return nil
}
