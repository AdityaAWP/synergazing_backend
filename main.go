package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Build connection string from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal open koneksi:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	log.Println("Database terkoneksi!")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	log.Fatal(app.Listen(":3001"))
}
