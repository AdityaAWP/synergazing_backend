package service

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

type NotificationService struct {
	DB *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{DB: db}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(userID uint, projectID *uint, notificationType, title, message string, data map[string]interface{}) (*model.Notification, error) {
	var dataJSON string
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal notification data: %v", err)
		}
		dataJSON = string(dataBytes)
	}

	notification := &model.Notification{
		UserID:    userID,
		ProjectID: projectID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		IsRead:    false,
		Data:      dataJSON,
	}

	if err := s.DB.Create(notification).Error; err != nil {
		return nil, fmt.Errorf("failed to create notification: %v", err)
	}

	return notification, nil
}

// GetUserNotifications retrieves notifications for a specific user
func (s *NotificationService) GetUserNotifications(userID uint, limit, offset int) ([]model.Notification, error) {
	var notifications []model.Notification

	query := s.DB.Where("user_id = ?", userID).
		Preload("Project").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %v", err)
	}

	return notifications, nil
}

// GetUnreadNotifications retrieves unread notifications for a user
func (s *NotificationService) GetUnreadNotifications(userID uint) ([]model.Notification, error) {
	var notifications []model.Notification

	if err := s.DB.Where("user_id = ? AND is_read = ?", userID, false).
		Preload("Project").
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %v", err)
	}

	return notifications, nil
}

// GetUnreadCount returns the count of unread notifications for a user
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	if err := s.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get unread count: %v", err)
	}
	return count, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(notificationID, userID uint) error {
	result := s.DB.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		return fmt.Errorf("failed to mark notification as read: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found or unauthorized")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	if err := s.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error; err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %v", err)
	}
	return nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(notificationID, userID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&model.Notification{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found or unauthorized")
	}
	return nil
}

// Project-specific notification methods

// NotifyDeadlineApproaching sends notifications to project members about approaching deadline
func (s *NotificationService) NotifyDeadlineApproaching(projectID uint) error {
	var project model.Project
	if err := s.DB.Preload("Members.User").First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	daysUntilDeadline := int(time.Until(project.RegistrationDeadline).Hours() / 24)

	title := "Project Deadline Approaching"
	message := fmt.Sprintf("The registration deadline for project '%s' is in %d days", project.Title, daysUntilDeadline)

	data := map[string]interface{}{
		"project_id":            project.ID,
		"project_title":         project.Title,
		"registration_deadline": project.RegistrationDeadline,
		"days_until_deadline":   daysUntilDeadline,
	}

	// Notify project creator
	if _, err := s.CreateNotification(project.CreatorID, &projectID, model.NotificationTypeDeadlineApproaching, title, message, data); err != nil {
		return fmt.Errorf("failed to notify project creator: %v", err)
	}

	// Notify all project members
	for _, member := range project.Members {
		if member.UserID != project.CreatorID { // Don't double-notify creator
			if _, err := s.CreateNotification(member.UserID, &projectID, model.NotificationTypeDeadlineApproaching, title, message, data); err != nil {
				return fmt.Errorf("failed to notify project member: %v", err)
			}
		}
	}

	return nil
}

// NotifyUserRegistered notifies project creator when someone registers for their project
func (s *NotificationService) NotifyUserRegistered(projectID, applicantUserID uint) error {
	var project model.Project
	var applicant model.Users

	if err := s.DB.First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	if err := s.DB.First(&applicant, applicantUserID).Error; err != nil {
		return fmt.Errorf("failed to find applicant: %v", err)
	}

	title := "New Project Application"
	message := fmt.Sprintf("%s has applied to join your project '%s'", applicant.Name, project.Title)

	data := map[string]interface{}{
		"project_id":      project.ID,
		"project_title":   project.Title,
		"applicant_id":    applicant.ID,
		"applicant_name":  applicant.Name,
		"applicant_email": applicant.Email,
	}

	_, err := s.CreateNotification(project.CreatorID, &projectID, model.NotificationTypeUserRegistered, title, message, data)
	return err
}

// NotifyUserAccepted notifies user when they're accepted into a project
func (s *NotificationService) NotifyUserAccepted(projectID, userID uint, roleTitle string) error {
	var project model.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	title := "Application Accepted!"
	message := fmt.Sprintf("Congratulations! You have been accepted into the project '%s'", project.Title)

	if roleTitle != "" {
		message += fmt.Sprintf(" as %s", roleTitle)
	}

	data := map[string]interface{}{
		"project_id":    project.ID,
		"project_title": project.Title,
		"role":          roleTitle,
		"status":        "accepted",
	}

	_, err := s.CreateNotification(userID, &projectID, model.NotificationTypeUserAccepted, title, message, data)
	return err
}

// NotifyUserRejected notifies user when their application is rejected
func (s *NotificationService) NotifyUserRejected(projectID, userID uint) error {
	var project model.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	title := "Application Update"
	message := fmt.Sprintf("Thank you for your interest in '%s'. Unfortunately, your application was not selected this time.", project.Title)

	data := map[string]interface{}{
		"project_id":    project.ID,
		"project_title": project.Title,
		"status":        "rejected",
	}

	_, err := s.CreateNotification(userID, &projectID, model.NotificationTypeUserRejected, title, message, data)
	return err
}

// NotifyRoleAssigned notifies user when they're assigned a specific role
func (s *NotificationService) NotifyRoleAssigned(projectID, userID uint, roleTitle string) error {
	var project model.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	title := "Role Assigned"
	message := fmt.Sprintf("You have been assigned the role '%s' in project '%s'", roleTitle, project.Title)

	data := map[string]interface{}{
		"project_id":    project.ID,
		"project_title": project.Title,
		"role":          roleTitle,
	}

	_, err := s.CreateNotification(userID, &projectID, model.NotificationTypeRoleAssigned, title, message, data)
	return err
}

// NotifyProjectStatusChange notifies relevant users about project status changes
func (s *NotificationService) NotifyProjectStatusChange(projectID uint, newStatus string) error {
	var project model.Project
	if err := s.DB.Preload("Members.User").First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	title := "Project Status Update"
	message := fmt.Sprintf("Project '%s' status has been updated to: %s", project.Title, newStatus)

	data := map[string]interface{}{
		"project_id":    project.ID,
		"project_title": project.Title,
		"new_status":    newStatus,
		"old_status":    project.Status,
	}

	// Notify all project members
	for _, member := range project.Members {
		if _, err := s.CreateNotification(member.UserID, &projectID, model.NotificationTypeProjectStatusChange, title, message, data); err != nil {
			return fmt.Errorf("failed to notify project member: %v", err)
		}
	}

	return nil
}

// NotifyInvitationReceived notifies user when they receive a project invitation
func (s *NotificationService) NotifyInvitationReceived(projectID, userID uint, roleTitle string) error {
	var project model.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return fmt.Errorf("failed to find project: %v", err)
	}

	title := "Project Invitation"
	message := fmt.Sprintf("You have been invited to join project '%s'", project.Title)

	if roleTitle != "" {
		message += fmt.Sprintf(" as %s", roleTitle)
	}

	data := map[string]interface{}{
		"project_id":    project.ID,
		"project_title": project.Title,
		"role":          roleTitle,
		"invitation":    true,
	}

	_, err := s.CreateNotification(userID, &projectID, model.NotificationTypeInvitationReceived, title, message, data)
	return err
}

// CheckAndNotifyApproachingDeadlines checks for projects with approaching deadlines
func (s *NotificationService) CheckAndNotifyApproachingDeadlines() error {
	// Check for deadlines in 1, 3, and 7 days
	deadlineDays := []int{1, 3, 7}

	for _, days := range deadlineDays {
		targetDate := time.Now().AddDate(0, 0, days)
		startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var projects []model.Project
		if err := s.DB.Where("registration_deadline >= ? AND registration_deadline < ? AND status != ?",
			startOfDay, endOfDay, "completed").Find(&projects).Error; err != nil {
			return fmt.Errorf("failed to find projects with approaching deadlines: %v", err)
		}

		for _, project := range projects {
			// Check if we already sent a notification for this deadline
			var existingNotification model.Notification
			err := s.DB.Where("project_id = ? AND type = ? AND created_at >= ?",
				project.ID, model.NotificationTypeDeadlineApproaching, startOfDay).First(&existingNotification).Error

			if err == gorm.ErrRecordNotFound {
				// No notification sent yet, send one
				if err := s.NotifyDeadlineApproaching(project.ID); err != nil {
					return fmt.Errorf("failed to notify approaching deadline for project %d: %v", project.ID, err)
				}
			}
		}
	}

	return nil
}
