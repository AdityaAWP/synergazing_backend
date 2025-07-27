package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

var allModels = []interface{}{
	&model.Users{},
	&model.Profiles{},
	&model.Role{},
	&model.Permission{},
}

func AutoMigrate(db *gorm.DB) {
	fmt.Println("Running Auto Migrate")
	err := db.AutoMigrate(
		allModels...,
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Success run Auto-migrate")
}

func MigrateFresh(db *gorm.DB) {
	fmt.Println("Running Migrate Fresh")
	err := db.Migrator().DropTable(allModels...)

	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	fmt.Println("Success drop all table")
	err = db.AutoMigrate(
		allModels...,
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
	fmt.Println("Success run Auto-migrate")
}
