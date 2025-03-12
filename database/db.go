package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// Create table if not exists
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS people (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			phone VARCHAR(20) NOT NULL,
			stack VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Println("Successfully connected to database")
} 