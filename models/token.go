package models

import "github.com/dgrijalva/jwt-go"

type Token struct {
	Rollno     string
	Authorized bool
	*jwt.StandardClaims
}
