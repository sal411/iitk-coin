package models

import (
	"gorm.io/gorm"
)

type UserData struct {
	gorm.Model

	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}
