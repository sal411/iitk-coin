package utils

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/sal411/iitk-coin/models"
)

func GenerateToken(userRollNo string) (time.Time, string, error) {
	var err error
	//Creating Access Token

	errenv := godotenv.Load()
	if errenv != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("secretKey")
	expireTime := time.Now().Add(time.Minute * 10)
	expiresAt := expireTime.Unix()

	tk := &models.Token{
		Rollno:     userRollNo,
		Authorized: true,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(secretKey))
	PrintError(err)
	if err != nil {
		return time.Now(), "", err
	}

	return expireTime, tokenString, nil

}
