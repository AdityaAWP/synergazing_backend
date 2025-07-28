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
	&model.SocialAuth{},
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

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	if err := tx.Migrator().DropTable(allModels...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop tables: %v", err)
	}
	fmt.Println("Success dropping all tables")

	if err := tx.AutoMigrate(allModels...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
	fmt.Println("Success running Auto-migrate")

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Migrate Fresh completed successfully")
}
