package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"

	_ "github.com/mattn/go-sqlite3"
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

		item := NewUser(db)

		newUserData := models.UserData{
			Rollno:   newuser.Rollno,
			Name:     newuser.Name,
			Password: string(hashed_password),
		}

		err_in_write := item.Add(newUserData)
		if err_in_write != nil {
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
		db.Close()

	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Only Post Methods are supported, please try again")

	}

}
