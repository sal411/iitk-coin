package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sal411/iitk-coin/database"
	"github.com/sal411/iitk-coin/models"
)

// function to get coins called from routes
func Coins(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getcoin" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == "GET" {
		var userCoin models.BankData
		err := json.NewDecoder(r.Body).Decode(&userCoin.Coin)
		if err != nil {
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, Cannot decode request",
			}
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		var tempCoin int
		tempCoin, err = database.GetCoinsFromRollno(userCoin.Rollno)

		userCoin.Coin = strconv.Itoa(tempCoin)
		if err == nil {
			var resp = map[string]interface{}{
				"status":  true,
				"message": "The user has : " + userCoin.Coin + " coins",
			}
			json.NewEncoder(w).Encode(resp)
			return

		} else {
			var resp = map[string]interface{}{
				"status":  false,
				"message": "Could not find user",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request, Supports only GET Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

}

// function to set coins for a roll number
func UpdateCoins(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/addcoin" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == "POST" {

		var newUserCoin models.BankData

		err := json.NewDecoder(r.Body).Decode(&newUserCoin)

		if err != nil {
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, Cannot decode request",
			}
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if newUserCoin.Rollno == "" {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, Please enter Rollno",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		_, err = strconv.Atoi(newUserCoin.Rollno)
		if err != nil {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, coin should be integer",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		err = database.WriteCoins(newUserCoin.Rollno, newUserCoin.Coin)
		if err != nil {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "User Not Found",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Coins added to  user ",
		}
		json.NewEncoder(w).Encode(resp)
		return

	} else {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request, Supports only POST Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

}

// handler to transfer coins
func TransferCoins(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/transfercoins" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == "POST" {

		var transferData models.TranfarCoin

		err := json.NewDecoder(r.Body).Decode(&transferData)
		if err != nil {
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, Cannot decode request",
			}
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		firstRollno := transferData.Account_1_Rollno
		secondRollno := transferData.Account_2_Rollno
		transferAmount := transferData.Amount

		if firstRollno == "" || secondRollno == "0" {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, Please enter Rollno",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		err = database.TransferCoin(firstRollno, secondRollno, transferAmount) // withdraw from first user and transfer to second
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Transaction Successful",
		}
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "This is an invalid Request, Supports only POST Request",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

}
