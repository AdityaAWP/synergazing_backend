package migrations

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

var modelMap = map[string]interface{}{
	"users":                &model.Users{},
	"profiles":             &model.Profiles{},
	"role":                 &model.Role{},
	"permission":           &model.Permission{},
	"socialauth":           &model.SocialAuth{},
	"skill":                &model.Skill{},
	"userskill":            &model.UserSkill{},
	"project":              &model.Project{},
	"projectcondition":     &model.ProjectCondition{},
	"tag":                  &model.Tag{},
	"benefit":              &model.Benefit{},
	"timeline":             &model.Timeline{},
	"projecttag":           &model.ProjectTag{},
	"projectbenefit":       &model.ProjectBenefit{},
	"projecttimeline":      &model.ProjectTimeline{},
	"projectrequiredskill": &model.ProjectRequiredSkill{},
	"projectrole":          &model.ProjectRole{},
	"projectroleskill":     &model.ProjectRoleSkill{},
	"projectmember":        &model.ProjectMember{},
	"projectmemberskill":   &model.ProjectMemberSkill{},
	"chat":                 &model.Chat{},
	"chats":                &model.Chat{},
	"message":              &model.Message{},
	"messages":             &model.Message{},
}

func AutoMigrate(db *gorm.DB) {
	fmt.Println("Running Auto Migrate")

	err := db.AutoMigrate(
		&model.Users{}, &model.Role{}, &model.Permission{}, &model.Skill{}, &model.Tag{}, &model.Benefit{}, &model.Timeline{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate primary tables: %v", err)
	}

	err = db.AutoMigrate(
		&model.Profiles{}, &model.SocialAuth{}, &model.UserSkill{}, &model.Project{}, &model.Chat{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate dependent tables: %v", err)
	}

	err = db.AutoMigrate(
		&model.ProjectCondition{}, &model.ProjectRequiredSkill{}, &model.ProjectTag{}, &model.ProjectBenefit{}, &model.ProjectTimeline{}, &model.ProjectRole{}, &model.ProjectRoleSkill{}, &model.ProjectMember{}, &model.ProjectMemberSkill{}, &model.Message{},
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

func MigrateFresh(db *gorm.DB) {
	fmt.Println("Running Migrate Fresh")

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	modelsToDrop := []interface{}{
		&model.ProjectMemberSkill{}, &model.ProjectMember{}, &model.ProjectRoleSkill{}, &model.ProjectCondition{}, &model.ProjectRequiredSkill{}, &model.ProjectTag{}, &model.ProjectBenefit{}, &model.ProjectTimeline{}, &model.ProjectRole{}, &model.Message{},
	}
	if err := tx.Migrator().DropTable(modelsToDrop...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop junction/dependent tables: %v", err)
	}

	modelsToDrop = []interface{}{
		&model.Profiles{}, &model.SocialAuth{}, &model.UserSkill{}, &model.Project{}, &model.Chat{},
	}
	if err := tx.Migrator().DropTable(modelsToDrop...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop main dependent tables: %v", err)
	}

	modelsToDrop = []interface{}{
		&model.Users{}, &model.Role{}, &model.Permission{}, &model.Skill{}, &model.Tag{}, &model.Benefit{}, &model.Timeline{},
	}
	if err := tx.Migrator().DropTable(modelsToDrop...); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop primary tables: %v", err)
	}

	fmt.Println("Success dropping all tables")

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

func DropWorkerTypeColumn(db *gorm.DB) error {
	fmt.Println("Dropping worker_type column from projects table...")
	err := db.Exec("ALTER TABLE projects DROP COLUMN IF EXISTS worker_type;").Error
	if err != nil {
		return fmt.Errorf("failed to drop worker_type column: %v", err)
	}
	fmt.Println("Successfully dropped worker_type column")
	return nil
}
