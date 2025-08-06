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
	db := config.GetDB()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "fresh":
			log.Println("Running with fresh migration...")
			migrations.MigrateFresh(db)
			return
		case "drop":
			if len(os.Args) < 3 {
				log.Fatal("Please provide the table name to drop: e.g., `go run main.go drop role`")
			}
			tableName := os.Args[2]
			log.Printf("Dropping table: %s", tableName)
			migrations.DropTableByName(db, tableName)
			return
		}
	}

	log.Println("Running with auto migration...")
	migrations.AutoMigrate(db)

	app := fiber.New()

	app.Static("/storage", "./storage")

	routes.SetupAuthRoutes(app)
	routes.SetupProfileRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SkillRoutes(app)
	routes.SetupChatRoutes(app)
	routes.SetupTestRoutes(app)

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
