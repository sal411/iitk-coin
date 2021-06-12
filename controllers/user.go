package controllers

import (
	"gorm.io/gorm"
)

type UserData struct {
	gorm.Model

	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

type Users struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *Users {

	db.AutoMigrate(&UserData{})

	return &Users{
		DB: db,
	}
}

func (user *Users) Add(userdata UserData) {

	user.DB.Create(&UserData{Name: userdata.Name, Rollno: userdata.Rollno})

}
