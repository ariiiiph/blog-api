package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath, dbName string) {
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	dbFile := filepath.Join(dbPath, dbName)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping to databse: %v", err)
	}
	DB = db

	log.Println("DB initialised")
}

func CloseDB() {
	if DB == nil {
		return
	}

	err := DB.Close()
	if err != nil {
		fmt.Println("Error closing DB:", err)
	} else {
		fmt.Println("DB Closed")
	}
}
