package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: File .env tidak ditemukan, menggunakan variabel environment bawaan sistem.")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	driver := os.Getenv("DB_USE")

	if driver == "" {
		driver = "postgres"
	}

	if host == "" || user == "" || name == "" || port == "" {
		log.Fatal("DB_HOST, DB_USER, DB_NAME, dan DB_PORT wajib diisi di dalam file .env")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)

	if err := goose.SetDialect(driver); err != nil {
		log.Fatalf("Gagal melakukan set dialect goose: %v", err)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal menghubungi database (Ping gagal). Pastikan MySQL menyala: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal(`Penggunaan: go run ./cmd/migrate <command> [args]
		Contoh: 
		- go run ./cmd/migrate up
		- go run ./cmd/migrate down
		- go run ./cmd/migrate create create_users_table sql`)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err := goose.Run(command, db, "migrations", args...); err != nil {
		log.Fatalf("Perintah goose '%s' gagal: %v", command, err)
	}

	fmt.Printf("Berhasil menjalankan perintah migrasi: %s\n", command)
}
