package models

type UserData struct {
	Name         string `json:"name"`
	Rollno       string `json:"rollno"`
	Password     string `json:"password"`
	Account_type string `json:"account_type"`
}
