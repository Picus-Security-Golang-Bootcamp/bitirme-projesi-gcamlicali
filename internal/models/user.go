package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
)

type User struct {
	gorm.Model
	Mail      *string `gorm:"unique" json:"email"`
	Password  *string `json:"password"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Mobile    string
	IsAdmin   bool
}

func (User) TableName() string {
	//default table name
	return "user"
}

func GetAdmin() User {
	admin := User{}

	jsonFile, err := os.Open("./Docs/admin.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal("error while opening json: ", err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//Parse to Authors struct
	err = json.Unmarshal(byteValue, &admin)

	//fmt.Println(byteValue)
	if err != nil {
		log.Fatal("Error while unmarshal json: ", err)
	}
	admin.IsAdmin = true
	return admin
}
