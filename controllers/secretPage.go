package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sal411/iitk-coin/utils"
)

func ExtractMetadata(userToken string) (string, error) {
	token, err := utils.VerifyJWToken(userToken)
	if err != nil {
		return " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		rollNo, _ := claims["user_roll_no"].(string)

		return rollNo, err
	}
	return " ", err
}

func SecretPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Page not found",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if r.Method == "GET" {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// cookie is not set => return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)

				var resp = map[string]interface{}{
					"status":  false,
					"message": "Unauthorized Request",
				}
				json.NewEncoder(w).Encode(resp)

				return
			}

			//  if any other type of error => return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenFromUser := c.Value
		//fmt.Fprintf(w, tokenFromUser)
		userRollNo, _ := ExtractMetadata(tokenFromUser)
		fmt.Println(userRollNo, "Hello")
		var resp = map[string]interface{}{
			"status":  false,
			"message": "You are authorized and accessing the secret page, Roll no. : " + userRollNo,
		}
		json.NewEncoder(w).Encode(resp)

	} else {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{
			"status":  false,
			"message": "Only GET Methods are supported, please try again",
		}
		json.NewEncoder(w).Encode(resp)
	}

}
