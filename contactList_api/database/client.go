package database

import (
	"log"
	"rest-go-demo/entity"

	"github.com/jinzhu/gorm"
	//"github.com/go-sql-driver/mysql"
)

//Connector variable used for CRUD operation's
var Connector *gorm.DB

//Connect creates MySQL connection
func Connect(connectionString string) (*gorm.DB, error) {
	var err error
	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	log.Println("Connection was successful!!")
	return Connector, nil
}

//Migrate create/updates database table
func Migrate(table *entity.Contact) {
	Connector.AutoMigrate(&table)
	log.Println("Table migrated")
}
