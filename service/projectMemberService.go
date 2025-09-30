package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

type ProjectMemberService struct {
	DB                  *gorm.DB
	NotificationService *NotificationService
}

func NewProjectMemberService(db *gorm.DB, notificationService *NotificationService) *ProjectMemberService {
	return &ProjectMemberService{
		DB:                  db,
		NotificationService: notificationService,
	}
}

// ApplicationData contains all the information for a project application
type ApplicationData struct {
	ProjectRoleID    uint   `json:"project_role_id"`
	WhyInterested    string `json:"why_interested"`
	SkillsExperience string `json:"skills_experience"`
	Contribution     string `json:"contribution"`
}

// ApplyToProject allows a user to apply for a project role with detailed information
func (s *ProjectMemberService) ApplyToProject(userID, projectID uint, applicationData ApplicationData) (*model.ProjectApplication, error) {
	// Check if project exists and is accepting applications
	var project model.Project
	if err := s.DB.First(&project, projectID).Error; err != nil {
		return nil, errors.New("project not found")
	}

	if project.Status != "published" {
		return nil, errors.New("project is not accepting applications")
	}

	if time.Now().After(project.RegistrationDeadline) {
		return nil, errors.New("registration deadline has passed")
	}

	// Check if role exists
	var role model.ProjectRole
	if err := s.DB.Where("id = ? AND project_id = ?", applicationData.ProjectRoleID, projectID).First(&role).Error; err != nil {
		return nil, errors.New("project role not found")
	}

	// Check if user already applied for this project
	var existingApplication model.ProjectApplication
	if err := s.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&existingApplication).Error; err == nil {
		return nil, errors.New("you have already applied to this project")
	}

	// Check if user is already a member
	var existingMember model.ProjectMember
	if err := s.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&existingMember).Error; err == nil {
		return nil, errors.New("you are already a member of this project")
	}

	// Check if user is the project creator
	if project.CreatorID == userID {
		return nil, errors.New("you cannot apply to your own project")
	}

	// Validate required fields
	if applicationData.WhyInterested == "" {
		return nil, errors.New("please explain why you're interested in this project")
	}
	if applicationData.SkillsExperience == "" {
		return nil, errors.New("please describe your relevant skills and experience")
	}
	if applicationData.Contribution == "" {
		return nil, errors.New("please describe what you can contribute to this project")
	}

	// Create application
	application := &model.ProjectApplication{
		ProjectID:        projectID,
		UserID:           userID,
		ProjectRoleID:    applicationData.ProjectRoleID,
		Status:           model.ApplicationStatusPending,
		WhyInterested:    applicationData.WhyInterested,
		SkillsExperience: applicationData.SkillsExperience,
		Contribution:     applicationData.Contribution,
		AppliedAt:        time.Now(),
	}

	if err := s.DB.Create(application).Error; err != nil {
		return nil, fmt.Errorf("failed to create application: %v", err)
	}

	// Load relations for response
	if err := s.DB.Preload("Project").Preload("User").Preload("ProjectRole").First(application, application.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load application details: %v", err)
	}

	// Send notification to project creator
	if err := s.NotificationService.NotifyUserRegistered(projectID, userID); err != nil {
		// Log error but don't fail the application
		fmt.Printf("Failed to send notification: %v\n", err)
	}

	return application, nil
}

// GetProjectApplications retrieves applications for a project (for project owners)
func (s *ProjectMemberService) GetProjectApplications(projectID, creatorID uint) ([]model.ProjectApplication, error) {
	// Verify the user is the project creator
	var project model.Project
	if err := s.DB.Where("id = ? AND creator_id = ?", projectID, creatorID).First(&project).Error; err != nil {
		return nil, errors.New("project not found or unauthorized")
	}

	var applications []model.ProjectApplication
	if err := s.DB.Where("project_id = ?", projectID).
		Preload("User").
		Preload("ProjectRole").
		Preload("Reviewer").
		Order("created_at DESC").
		Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("failed to get applications: %v", err)
	}

	return applications, nil
}

// GetUserApplications retrieves applications submitted by a user
func (s *ProjectMemberService) GetUserApplications(userID uint) ([]model.ProjectApplication, error) {
	var applications []model.ProjectApplication
	if err := s.DB.Where("user_id = ?", userID).
		Preload("Project").
		Preload("ProjectRole").
		Preload("Reviewer").
		Order("created_at DESC").
		Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("failed to get user applications: %v", err)
	}

	return applications, nil
}

// GetUserInvitations retrieves project invitations received by a user
func (s *ProjectMemberService) GetUserInvitations(userID uint) ([]model.ProjectMember, error) {
	var invitations []model.ProjectMember
	if err := s.DB.Where("user_id = ? AND status = ?", userID, "invited").
		Preload("Project").
		Preload("ProjectRole").
		Order("created_at DESC").
		Find(&invitations).Error; err != nil {
		return nil, fmt.Errorf("failed to get user invitations: %v", err)
	}

	return invitations, nil
}

// ReviewApplicationData contains the review decision and notes
type ReviewApplicationData struct {
	Action      string `json:"action"`
	ReviewNotes string `json:"review_notes"`
}

// ReviewApplication allows project creator to accept or reject an application
func (s *ProjectMemberService) ReviewApplication(applicationID, reviewerID uint, reviewData ReviewApplicationData) error {
	if reviewData.Action != "accept" && reviewData.Action != "reject" {
		return errors.New("invalid action. Must be 'accept' or 'reject'")
	}

	var application model.ProjectApplication
	if err := s.DB.Preload("Project").Preload("ProjectRole").First(&application, applicationID).Error; err != nil {
		return errors.New("application not found")
	}

	// Verify the reviewer is the project creator
	if application.Project.CreatorID != reviewerID {
		return errors.New("unauthorized to review this application")
	}

	if application.Status != model.ApplicationStatusPending {
		return errors.New("application has already been reviewed")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	newStatus := model.ApplicationStatusRejected
	if reviewData.Action == "accept" {
		newStatus = model.ApplicationStatusAccepted
	}

	// Update application status
	if err := tx.Model(&application).Updates(map[string]interface{}{
		"status":       newStatus,
		"reviewed_at":  &now,
		"reviewed_by":  reviewerID,
		"review_notes": reviewData.ReviewNotes,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update application: %v", err)
	}

	if reviewData.Action == "accept" {
		// Check if there's still room in the role
		var currentMembers int64
		if err := tx.Model(&model.ProjectMember{}).
			Where("project_id = ? AND project_role_id = ?", application.ProjectID, application.ProjectRoleID).
			Count(&currentMembers).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to count current members: %v", err)
		}

		if int(currentMembers) >= application.ProjectRole.SlotsAvailable {
			tx.Rollback()
			return errors.New("no more slots available for this role")
		}

		// Create project member
		member := &model.ProjectMember{
			ProjectID:     application.ProjectID,
			UserID:        application.UserID,
			ProjectRoleID: application.ProjectRoleID,
			Status:        "accepted",
		}

		if err := tx.Create(member).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create project member: %v", err)
		}

		// Send acceptance notification
		if err := s.NotificationService.NotifyUserAccepted(application.ProjectID, application.UserID, application.ProjectRole.Name); err != nil {
			fmt.Printf("Failed to send acceptance notification: %v\n", err)
		}
	} else {
		// Send rejection notification
		if err := s.NotificationService.NotifyUserRejected(application.ProjectID, application.UserID); err != nil {
			fmt.Printf("Failed to send rejection notification: %v\n", err)
		}
	}

	return tx.Commit().Error
}

// WithdrawApplication allows a user to withdraw their application
func (s *ProjectMemberService) WithdrawApplication(applicationID, userID uint) error {
	var application model.ProjectApplication
	if err := s.DB.Where("id = ? AND user_id = ?", applicationID, userID).First(&application).Error; err != nil {
		return errors.New("application not found or unauthorized")
	}

	if application.Status != model.ApplicationStatusPending {
		return errors.New("can only withdraw pending applications")
	}

	if err := s.DB.Model(&application).Update("status", model.ApplicationStatusWithdrawn).Error; err != nil {
		return fmt.Errorf("failed to withdraw application: %v", err)
	}

	return nil
}

// RemoveMember allows project creator to remove a member from the project
func (s *ProjectMemberService) RemoveMember(projectID, memberUserID, creatorID uint) error {
	// Verify the user is the project creator
	var project model.Project
	if err := s.DB.Where("id = ? AND creator_id = ?", projectID, creatorID).First(&project).Error; err != nil {
		return errors.New("project not found or unauthorized")
	}

	// Cannot remove the creator
	if memberUserID == creatorID {
		return errors.New("cannot remove project creator")
	}

	// Find and remove the member
	var member model.ProjectMember
	if err := s.DB.Where("project_id = ? AND user_id = ?", projectID, memberUserID).First(&member).Error; err != nil {
		return errors.New("member not found")
	}

	if err := s.DB.Delete(&member).Error; err != nil {
		return fmt.Errorf("failed to remove member: %v", err)
	}

	return nil
}

// InviteMember allows project creator to invite a user to join the project
func (s *ProjectMemberService) InviteMember(projectID, userID, roleID, creatorID uint) error {
	// Verify the user is the project creator
	var project model.Project
	if err := s.DB.Where("id = ? AND creator_id = ?", projectID, creatorID).First(&project).Error; err != nil {
		return errors.New("project not found or unauthorized")
	}

	// Check if role exists
	var role model.ProjectRole
	if err := s.DB.Where("id = ? AND project_id = ?", roleID, projectID).First(&role).Error; err != nil {
		return errors.New("project role not found")
	}

	// Check if user exists
	var user model.Users
	if err := s.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Check if user is already a member
	var existingMember model.ProjectMember
	if err := s.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&existingMember).Error; err == nil {
		return errors.New("user is already a member of this project")
	}

	// Check if user already has a pending application
	var existingApplication model.ProjectApplication
	if err := s.DB.Where("user_id = ? AND project_id = ? AND status = ?", userID, projectID, model.ApplicationStatusPending).First(&existingApplication).Error; err == nil {
		return errors.New("user already has a pending application for this project")
	}

	// Check if there's room in the role
	var currentMembers int64
	if err := s.DB.Model(&model.ProjectMember{}).
		Where("project_id = ? AND project_role_id = ?", projectID, roleID).
		Count(&currentMembers).Error; err != nil {
		return fmt.Errorf("failed to count current members: %v", err)
	}

	if int(currentMembers) >= role.SlotsAvailable {
		return errors.New("no more slots available for this role")
	}

	// Create project member with invited status
	member := &model.ProjectMember{
		ProjectID:     projectID,
		UserID:        userID,
		ProjectRoleID: roleID,
		Status:        "invited",
	}

	if err := s.DB.Create(member).Error; err != nil {
		return fmt.Errorf("failed to create invitation: %v", err)
	}

	// Send invitation notification
	if err := s.NotificationService.NotifyInvitationReceived(projectID, userID, role.Name); err != nil {
		fmt.Printf("Failed to send invitation notification: %v", err)
	}

	return nil
}

// RespondToInvitation allows a user to accept or decline an invitation
func (s *ProjectMemberService) RespondToInvitation(projectID, userID uint, response string) error {
	if response != "accept" && response != "decline" {
		return errors.New("invalid response. Must be 'accept' or 'decline'")
	}

	var member model.ProjectMember
	if err := s.DB.Preload("ProjectRole").Where("project_id = ? AND user_id = ? AND status = ?", projectID, userID, "invited").First(&member).Error; err != nil {
		return errors.New("invitation not found")
	}

	newStatus := "declined"
	if response == "accept" {
		newStatus = "accepted"

		// Send role assignment notification
		if err := s.NotificationService.NotifyRoleAssigned(projectID, userID, member.ProjectRole.Name); err != nil {
			fmt.Printf("Failed to send role assignment notification: %v", err)
		}
	}

	if err := s.DB.Model(&member).Update("status", newStatus).Error; err != nil {
		return fmt.Errorf("failed to update invitation status: %v", err)
	}

	return nil
}

// GetProjectMembers retrieves all members of a project
func (s *ProjectMemberService) GetProjectMembers(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	if err := s.DB.Where("project_id = ?", projectID).
		Preload("User").
		Preload("ProjectRole").
		Find(&members).Error; err != nil {
		return nil, fmt.Errorf("failed to get project members: %v", err)
	}

	return members, nil
}

// GetApplicationSummary gets a brief summary of an application for listings
func (s *ProjectMemberService) GetApplicationSummary(applicationID uint) (map[string]interface{}, error) {
	var application model.ProjectApplication
	if err := s.DB.Preload("User").Preload("ProjectRole").Preload("Project").First(&application, applicationID).Error; err != nil {
		return nil, fmt.Errorf("application not found: %v", err)
	}

	summary := map[string]interface{}{
		"id":                        application.ID,
		"status":                    application.Status,
		"applied_at":                application.AppliedAt,
		"applicant_name":            application.User.Name,
		"applicant_email":           application.User.Email,
		"role_name":                 application.ProjectRole.Name,
		"project_title":             application.Project.Title,
		"skills_experience_summary": truncateText(application.SkillsExperience, 150),
	}

	return summary, nil
}

// GetApplicationDetails gets full details of an application for review
func (s *ProjectMemberService) GetApplicationDetails(applicationID, requesterID uint) (*model.ProjectApplication, error) {
	var application model.ProjectApplication
	if err := s.DB.Preload("User").Preload("ProjectRole").Preload("Project").Preload("Reviewer").First(&application, applicationID).Error; err != nil {
		return nil, fmt.Errorf("application not found: %v", err)
	}

	// Check if requester has permission to view (project owner or applicant)
	if application.Project.CreatorID != requesterID && application.UserID != requesterID {
		return nil, errors.New("unauthorized to view this application")
	}

	return &application, nil
}

// Helper function to truncate text for summaries
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}
