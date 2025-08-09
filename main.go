package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
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
		case "drop-column":
			log.Println("Dropping worker_type column from projects table...")
			err := migrations.DropWorkerTypeColumn(db)
			if err != nil {
				log.Fatalf("Failed to drop worker_type column: %v", err)
			}
			return
		}
	}

	log.Println("Running with auto migration...")
	migrations.AutoMigrate(db)

	app := fiber.New()
	app.Use(cors.New())
	app.Static("/storage", "./storage")

	if _, err := os.Stat("./Api.yml"); os.IsNotExist(err) {
		log.Println("Warning: Api.yml not found in root directory")
	}

	// Serve the OpenAPI YAML file BEFORE the swagger middleware
	app.Get("/api/docs/doc.yaml", func(c *fiber.Ctx) error {
		// Get absolute path
		absPath, err := filepath.Abs("./Api.yml")
		if err != nil {
			log.Printf("Error getting absolute path: %v", err)
			return c.Status(500).SendString("Error loading API spec")
		}

		// Check if file exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			log.Printf("Api.yml not found at: %s", absPath)
			return c.Status(404).SendString("API specification not found")
		}

		c.Set("Content-Type", "application/x-yaml")
		c.Set("Access-Control-Allow-Origin", "*")
		return c.SendFile("./Api.yml")
	})

	// Serve OpenAPI documentation
	app.Get("/api/docs/*", swagger.New(swagger.Config{
		URL:          "/api/docs/doc.yaml",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// Redirect /api/docs to /api/docs/
	app.Get("/api/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/")
	})
	routes.SetupAuthRoutes(app)
	routes.SetupProfileRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SkillRoutes(app)
	routes.SetupChatRoutes(app)
	routes.SetupTestRoutes(app)
	routes.SetupProjectRoutes(app)

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
