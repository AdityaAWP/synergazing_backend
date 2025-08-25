package migrations

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	"otp":                  &model.OTP{},
	"otps":                 &model.OTP{},
}

func AutoMigrate(db *gorm.DB) {
	fmt.Println("Running Auto Migrate")

	// Create custom enums first
	if err := CreateCustomEnums(db); err != nil {
		log.Fatalf("Failed to create custom enums: %v", err)
	}

	err := db.AutoMigrate(
		&model.Users{}, &model.Role{}, &model.Permission{}, &model.Skill{}, &model.Tag{}, &model.Benefit{}, &model.Timeline{}, &model.OTP{},
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
		&model.Profiles{}, &model.SocialAuth{}, &model.UserSkill{}, &model.Project{}, &model.Chat{}, &model.OTP{},
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

	err = db.Exec("DO $$ BEGIN CREATE TYPE timeline_status AS ENUM ('not-started', 'in-progress', 'done'); EXCEPTION WHEN duplicate_object THEN null; END $$;").Error
	if err != nil {
		return err
	}

	err = db.Exec("DO $$ BEGIN CREATE TYPE otp_purpose AS ENUM ('registration', 'password_reset', 'email_change'); EXCEPTION WHEN duplicate_object THEN null; END $$;").Error
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

func MigrateOTP(db *gorm.DB) error {
	fmt.Println("Running OTP migration...")

	err := db.Exec("DO $$ BEGIN CREATE TYPE otp_purpose AS ENUM ('registration', 'password_reset', 'email_change'); EXCEPTION WHEN duplicate_object THEN null; END $$;").Error
	if err != nil {
		return fmt.Errorf("failed to create otp_purpose enum: %v", err)
	}

	err = db.AutoMigrate(&model.OTP{})
	if err != nil {
		return fmt.Errorf("failed to migrate OTP table: %v", err)
	}

	err = db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS is_email_verified BOOLEAN DEFAULT FALSE;").Error
	if err != nil {
		return fmt.Errorf("failed to add is_email_verified column: %v", err)
	}

	fmt.Println("OTP migration completed successfully")
	return nil
}

func CleanupExpiredOTPs(db *gorm.DB) {
	result := db.Where("expires_at < ?", time.Now()).Delete(&model.OTP{})
	if result.Error != nil {
		log.Printf("Error cleaning up expired OTPs: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired OTP records", result.RowsAffected)
	} else {
		log.Println("No expired OTP records found")
	}
}
