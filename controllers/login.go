package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"
	"golang.org/x/crypto/bcrypt"
)

func setCookie(w http.ResponseWriter, token string, time time.Time) {

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time,
	})

}

func FindUser(userRollno string, userPassword string, w http.ResponseWriter) map[string]interface{} {

	hashedPassword := utils.GetHashedPassword(userRollno)

	if hashedPassword == "" {
		var resp = map[string]interface{}{
			"status":  false,
			"message": "User with given RollNo does not exist, Please try to SignUp",
		}

		return resp
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword)); err != nil {
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Wrong Password",
		}
		return resp
	}

	expirationTime, token, err := utils.GenerateToken(userRollno)
	if err != nil {

		var resp = map[string]interface{}{
			"status":  false,
			"message": "Could Not generate JWT token",
			"error":   []byte(err.Error()),
		}
		return resp
	}

	var resp = map[string]interface{}{
		"status":    true,
		"message":   "Password was correct!, You are logged in to localhost:8080 ",
		"token":     token,
		"rollno":    userRollno,
		"expiresAt": expirationTime,
	}

	setCookie(w, token, expirationTime)

	return resp

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == "POST" {

		var newuser models.UserData
		err := json.NewDecoder(r.Body).Decode(&newuser)
		if err != nil {
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request",
			}
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		}

		fmt.Println(newuser)
		//print ho rha hai

		resp := FindUser(newuser.Rollno, newuser.Password, w)

		json.NewEncoder(w).Encode(resp)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request",
		}
		json.NewEncoder(w).Encode(resp)

	}

}
