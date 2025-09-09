package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	
	_ "github.com/lib/pq"
	"github.com/go-vgo/robotgo"
)

func main() {
	// Membaca environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "tracking_db")
	
	// Membuat connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
		dbUser, dbPassword, dbHost, dbPort, dbName)
	
	// Inisialisasi koneksi database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	
	// Membuat tabel jika belum ada
	createTables(db)
	
	// Jalankan API di goroutine terpisah
	go setupAPI(db)
	
	// Memulai tracking
	go trackMouse(db)
	trackKeyboard(db)
}

// Helper function untuk mendapatkan environment variable dengan nilai default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}