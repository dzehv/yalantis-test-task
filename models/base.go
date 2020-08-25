package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

// DB handler
var db *gorm.DB

func connect() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	// conn string
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	log.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Printf("WARNING: %v", err)
	}

	db = conn
	// migrate
	db.Debug().AutoMigrate(&Account{})
}

// return our db handler
func GetDB() *gorm.DB {
	connect()
	return db
}

func CloseDB() {
	db.Close()
}
