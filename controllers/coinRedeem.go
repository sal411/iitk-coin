package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sal411/iitk-coin/database"
	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"
)

func RedeemCoins(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/redeem" {
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

	if r.Method == "POST" {
		var redeemData models.RedeemCoinsData

		err := json.NewDecoder(r.Body).Decode(&redeemData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		item_id := redeemData.Item_id

		if rollno == "" {
			w.WriteHeader(401)
			var resp = map[string]interface{}{
				"status":  false,
				"message": "This is an invalid Request, enter a roll number",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		coins, err := database.RedeemCoinsDb(rollno, item_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Sucessfully redeemed item " + fmt.Sprintf("%d", item_id) + " .Coins remaining are " + fmt.Sprintf("%.2f", coins),
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
