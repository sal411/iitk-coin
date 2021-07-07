package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sal411/iitk-coin/database"
	"github.com/sal411/iitk-coin/models"
	"github.com/sal411/iitk-coin/utils"
)

func AddItems(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/additems" {
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
		var itemData models.ItemsData

		err := json.NewDecoder(r.Body).Decode(&itemData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item_id := itemData.Item_id

		cost := itemData.Cost
		number := itemData.Number

		w.Header().Set("Content-Type", "application/json")

		message, err := database.WriteItems(item_id, cost, number)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		var resp = map[string]interface{}{
			"status":  false,
			"message": message + " updated in database ",
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
