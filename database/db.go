package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const DSN = "doadmin:AVNS_1SUWMwvNP7bhHjPYpC1@tcp(db-mysql-sgp1-11646-do-user-6185766-0.m.db.ondigitalocean.com:25060)/say_right?tls=do&parseTime=true&charset=utf8mb4"

func Connect() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connection established successfully (MySQL)")
}
