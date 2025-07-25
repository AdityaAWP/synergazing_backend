package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://die4wp:root@127.0.0.1:5432/die4wp?sslmode=disable"
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
