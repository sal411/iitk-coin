package models

type TranfarCoin struct {
	Account_1_Rollno string `json:"firstrollno"`
	Account_2_Rollno string `json:"secondrollno"`
	Amount           int    `json:"amount"`
}
