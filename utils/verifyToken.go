package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func VerifyJWToken(requestedToken string) (*jwt.Token, error) {

	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal("Error in loading .env file")
	}

	secretKey := os.Getenv("secretKey")

	tokenString := requestedToken
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
