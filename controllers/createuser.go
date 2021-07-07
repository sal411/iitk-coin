package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sal411/iitk-coin/database"
	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"

	"golang.org/x/crypto/bcrypt"
)

var db = utils.ConnectDB()

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Page not found",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == "POST" {
		var newuser models.UserData
		err := json.NewDecoder(r.Body).Decode(&newuser)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(newuser)

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(newuser.Password), bcrypt.DefaultCost)
		utils.PrintError(err)

		item := database.NewUser(db)

		newUserData := models.UserData{
			Rollno:       newuser.Rollno,
			Name:         newuser.Name,
			Password:     string(hashed_password),
			Account_type: newuser.Account_type,
		}

		if newuser.Rollno == "" || newuser.Password == "" || newuser.Account_type == "" {
			var resp = map[string]interface{}{
				"status":  false,
				"message": "one of the field is empty",
			}
			json.NewEncoder(w).Encode(resp)
			return

		}

		item2 := database.NewBank(db)

		newBankData := models.BankData{
			Rollno: newuser.Rollno,
			Coins:  "0.00",
		}

		err_in_write := item.Add(newUserData)
		err1 := item2.OpenAccount(newBankData)
		if err_in_write != nil && err1 != nil {
			log.Printf("Body read error, %v", err_in_write)
			w.WriteHeader(500) // Return 500 Internal Server Error.

			var resp = map[string]interface{}{
				"status":  true,
				"message": "Roll no must be unique",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		var resp = map[string]interface{}{
			"status":  true,
			"message": "Congratulations! Your account has been successfully created",
		}
		json.NewEncoder(w).Encode(resp)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only Post Methods are supported, please try again")

	}

}
