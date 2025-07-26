package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/migrations"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load ENV")
	}

	config.ConnectEnvDBConfig()
	migrations.AutoMigrate(config.GetDB())

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World - GORM Connected!")
	})

	log.Fatal(app.Listen(":5000"))
}
