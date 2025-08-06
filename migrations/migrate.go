package migrations

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

// List of all models for reference in DropTableByName
var modelMap = map[string]interface{}{
	"users":            &model.Users{},
	"profiles":         &model.Profiles{},
	"role":             &model.Role{},
	"permission":       &model.Permission{},
	"socialauth":       &model.SocialAuth{},
	"skill":            &model.Skill{},
	"userskill":        &model.UserSkill{},
	"project":          &model.Project{},
	"projectcondition": &model.ProjectCondition{},
	"tag":              &model.Tag{},
	"projectrole":      &model.ProjectRole{},
	"projectmember":    &model.ProjectMember{},
	"chat":             &model.Chat{},
	"chats":            &model.Chat{},
	"message":          &model.Message{},
	"messages":         &model.Message{},
}

// AutoMigrate contains the correctly ordered logic for table creation.
func AutoMigrate(db *gorm.DB) {
	fmt.Println("Running Auto Migrate")

	// 1. Migrate primary tables first (no dependencies on other models)
	err := db.AutoMigrate(
		&model.Users{}, &model.Role{}, &model.Permission{}, &model.Skill{}, &model.Tag{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate primary tables: %v", err)
	}

	// 2. Migrate tables that have foreign keys pointing to the primary tables
	err = db.AutoMigrate(
		&model.Profiles{}, &model.SocialAuth{}, &model.UserSkill{}, &model.Project{}, &model.Chat{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate dependent tables: %v", err)
	}

	// 3. Migrate the final layer of tables that depend on the second layer
	err = db.AutoMigrate(
		&model.ProjectCondition{}, &model.ProjectRole{}, &model.ProjectMember{}, &model.Message{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate final tables: %v", err)
	}

	fmt.Println("Success run Auto-migrate")
}

func DropTableByName(db *gorm.DB, tableName string) {
	model, exist := modelMap[strings.ToLower(tableName)]
	if !exist {
		log.Fatalf("Model for table %s does not exist", tableName)
	}
	if err := db.Migrator().DropTable(model); err != nil {
		log.Fatalf("Failed to drop table %s: %v", tableName, err)
	}
	fmt.Printf("Table %s dropped successfully\n", tableName)
}

// MigrateFresh will now correctly drop and create tables in the proper order.
func MigrateFresh(db *gorm.DB) {
	fmt.Println("Running Migrate Fresh")

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// --- THIS IS THE CRITICAL CHANGE ---
	// Drop tables in the REVERSE order of creation to respect foreign key constraints.

	// 1. Drop the most dependent tables first
	modelsToDrop := []interface{}{
		&model.ProjectMember{}, &model.ProjectCondition{}, &model.ProjectRole{}, &model.Message{},
	}
	if err := tx.Migrator().DropTable(modelsToDrop...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop junction/dependent tables: %v", err)
	}

	// 2. Drop the next layer of tables
	modelsToDrop = []interface{}{
		&model.Profiles{}, &model.SocialAuth{}, &model.UserSkill{}, &model.Project{}, &model.Chat{},
	}
	if err := tx.Migrator().DropTable(modelsToDrop...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop main dependent tables: %v", err)
	}

	// 3. Drop the primary tables last
	modelsToDrop = []interface{}{
		&model.Users{}, &model.Role{}, &model.Permission{}, &model.Skill{}, &model.Tag{},
	}
	if err := tx.Migrator().DropTable(modelsToDrop...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop primary tables: %v", err)
	}

	fmt.Println("Success dropping all tables")
	// --- END OF CRITICAL CHANGE ---

	// Now, call the ordered creation function
	AutoMigrate(tx)

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Migrate Fresh completed successfully")
}
func CreateCustomEnums(db *gorm.DB) error {
	err := db.Exec("DO $$ BEGIN CREATE TYPE collaboration_status AS ENUM ('not ready', 'ready'); EXCEPTION WHEN duplicate_object THEN null; END $$;").Error
	if err != nil {
		return err
	}
	return nil
}
