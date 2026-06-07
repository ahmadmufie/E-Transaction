package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

var DB *sql.DB

func ConnectDB() {
	// 1. Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Ambil data dari .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 3. Rangkai string koneksi
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// 4. Buka koneksi ke Database
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Gagal membuka database: ", err)
	}

	// 5. Tes ping ke Database untuk memastikan benar-benar terkoneksi
	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak merespons: ", err)
	}

	fmt.Println("Berhasil terhubung ke Database PostgreSQL!")
}
