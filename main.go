package main

import (
	"log"
	"net/url"
	"os"

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

	app.Static("/storage", "./storage")

	routes.SetupAuthRoutes(app)
	// routes.SetupProtectedRoutes(app)
	routes.SetupProfileRoutes(app)

	// routes.SetupProfileRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World - GORM Connected!")
	})

	appURL := os.Getenv("APP_URL")

	parsedURL, err := url.Parse(appURL)
	if err != nil {
		log.Fatalf("Failed to parse APP_URL: %v", err)
	}

	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	log.Fatal(app.Listen(":" + port))
}
