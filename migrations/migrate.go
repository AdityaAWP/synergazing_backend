package migrations

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

var allModels = []interface{}{
	&model.Users{},
	&model.Profiles{},
	&model.Role{},
	&model.Permission{},
	&model.SocialAuth{},
	&model.Skill{},
	&model.UserSkill{},
}

var modelMap = map[string]interface{}{
	"users":      &model.Users{},
	"profiles":   &model.Profiles{},
	"role":       &model.Role{},
	"permission": &model.Permission{},
	"socialauth": &model.SocialAuth{},
	"skill":      &model.Skill{},
	"userskill":  &model.UserSkill{},
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
func DropTableByName(db *gorm.DB, tableName string) {
	model, exist := modelMap[strings.ToLower(tableName)]
	if !exist {
		log.Fatalf("Model for table %s does not exist", tableName)
	}
	err := db.Migrator().DropTable(model)
	if err != nil {
		log.Fatalf("Failed to drop table %s: %v", tableName, err)
	}
	fmt.Printf("Table %s dropped successfully\n", tableName)
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
