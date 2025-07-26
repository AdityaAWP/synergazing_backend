package migrations

import (
	"log"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Users{},
		&model.Profiles{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
