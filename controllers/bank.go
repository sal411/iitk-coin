package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sal411/iitk-coin/database"
	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"
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
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	rollno, _, _ := utils.ExtractTokenMetadata(tokenFromUser)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {

		coins, err := database.GetCoinsFromRollno(rollno)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, " -User not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Your coins are " + fmt.Sprintf("%f", coins),
		}
		json.NewEncoder(w).Encode(resp)
		return

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
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	_, Acctype, _ := utils.ExtractTokenMetadata(tokenFromUser)

	if Acctype == "member" {
		http.Error(w, "Unauthorized!! Only CTM and admins are allowed ", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {

		var coinsData models.BankData

		err := json.NewDecoder(r.Body).Decode(&coinsData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := coinsData.Rollno

		numberOfCoins := coinsData.Coins

		remarks := coinsData.Remarks

		if rollno == "" {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, Please enter a roll number",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		_, userAccType, _ := database.GetUserFromRollNo(rollno)
		if userAccType == "CTM" && Acctype == "CTM" {
			http.Error(w, "Unauthorized only admins are alowed ", http.StatusUnauthorized)
			return
		}
		if userAccType == "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		_, err = strconv.ParseFloat(numberOfCoins, 32)
		if err != nil {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request,Coins should be valid number",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		errorMessage, err := database.WriteCoins(rollno, numberOfCoins, remarks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, errorMessage)
			return
		}

		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": errorMessage + coinsData.Coins + " Coins added to user " + coinsData.Rollno,
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
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	userRollNo, _, _ := utils.ExtractTokenMetadata(tokenFromUser)

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {

		var transferData models.TranfarCoin

		err := json.NewDecoder(r.Body).Decode(&transferData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		transferTorollno := transferData.Roll_no
		transferAmount := transferData.Amount

		if transferTorollno == "" {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "Please enter a roll number",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		tax, err := database.TransferCoin(userRollNo, transferTorollno, transferAmount) // withdraw from first user and transfer to second
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Transaction of " + fmt.Sprintf("%.2f", transferAmount) + " Sucessfull !  Tax Decucted = " + fmt.Sprintf("%.2f", tax),
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
