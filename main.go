package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/migrations"
	"synergazing.com/synergazing/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load ENV")
	}

	config.ConnectEnvDBConfig()
	migrations.AutoMigrate(config.GetDB())

	app := fiber.New()

	routes.SetupAuthRoutes(app)
	routes.SetupProtectedRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World - GORM Connected!")
	})

	log.Fatal(app.Listen(":3002"))
}
