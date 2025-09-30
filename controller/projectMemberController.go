package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type ProjectMemberController struct {
	projectMemberService *service.ProjectMemberService
}

func NewProjectMemberController(pms *service.ProjectMemberService) *ProjectMemberController {
	return &ProjectMemberController{projectMemberService: pms}
}

// ApplyToProject allows a user to apply for a project role
func (ctrl *ProjectMemberController) ApplyToProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, err := strconv.ParseUint(c.Params("project_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid project ID")
	}

	// Parse form data
	projectRoleIDStr := c.FormValue("project_role_id")
	whyInterested := c.FormValue("why_interested")
	skillsExperience := c.FormValue("skills_experience")
	contribution := c.FormValue("contribution")

	// Validate required fields
	if projectRoleIDStr == "" {
		return helper.Message400("Project role ID is required")
	}
	if whyInterested == "" {
		return helper.Message400("Please explain why you're interested in this project")
	}
	if skillsExperience == "" {
		return helper.Message400("Please describe your relevant skills and experience")
	}
	if contribution == "" {
		return helper.Message400("Please describe what you can contribute to this project")
	}

	projectRoleID, err := strconv.ParseUint(projectRoleIDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid project role ID")
	}

	applicationData := service.ApplicationData{
		ProjectRoleID:    uint(projectRoleID),
		WhyInterested:    whyInterested,
		SkillsExperience: skillsExperience,
		Contribution:     contribution,
	}

	application, err := ctrl.projectMemberService.ApplyToProject(userID, uint(projectID), applicationData)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message201(c, application, "Application submitted successfully")
}

// GetProjectApplications retrieves applications for a project (for project owners)
func (ctrl *ProjectMemberController) GetProjectApplications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, err := strconv.ParseUint(c.Params("project_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid project ID")
	}

	applications, err := ctrl.projectMemberService.GetProjectApplications(uint(projectID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, applications, "Project applications retrieved successfully")
}

// GetUserApplications retrieves applications submitted by a user
func (ctrl *ProjectMemberController) GetUserApplications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	applications, err := ctrl.projectMemberService.GetUserApplications(userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, applications, "User applications retrieved successfully")
}

// GetUserInvitations retrieves project invitations received by a user
func (ctrl *ProjectMemberController) GetUserInvitations(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	invitations, err := ctrl.projectMemberService.GetUserInvitations(userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, invitations, "User invitations retrieved successfully")
}

// ReviewApplication allows project creator to accept or reject an application
func (ctrl *ProjectMemberController) ReviewApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	applicationID, err := strconv.ParseUint(c.Params("application_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid application ID")
	}

	// Parse form data
	action := c.FormValue("action")
	reviewNotes := c.FormValue("review_notes")

	if action != "accept" && action != "reject" {
		return helper.Message400("Action must be 'accept' or 'reject'")
	}

	reviewData := service.ReviewApplicationData{
		Action:      action,
		ReviewNotes: reviewNotes,
	}

	err = ctrl.projectMemberService.ReviewApplication(uint(applicationID), userID, reviewData)
	if err != nil {
		return helper.Message400(err.Error())
	}

	message := "Application accepted successfully"
	if action == "reject" {
		message = "Application rejected successfully"
	}

	return helper.Message200(c, nil, message)
}

// WithdrawApplication allows a user to withdraw their application
func (ctrl *ProjectMemberController) WithdrawApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	applicationID, err := strconv.ParseUint(c.Params("application_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid application ID")
	}

	err = ctrl.projectMemberService.WithdrawApplication(uint(applicationID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Application withdrawn successfully")
}

// RemoveMember allows project creator to remove a member from the project
func (ctrl *ProjectMemberController) RemoveMember(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, err := strconv.ParseUint(c.Params("project_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid project ID")
	}

	memberUserID, err := strconv.ParseUint(c.Params("user_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid member user ID")
	}

	err = ctrl.projectMemberService.RemoveMember(uint(projectID), uint(memberUserID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Member removed successfully")
}

// InviteMember allows project creator to invite a user to join the project
func (ctrl *ProjectMemberController) InviteMember(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, err := strconv.ParseUint(c.Params("project_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid project ID")
	}

	// Parse form data
	inviteUserIDStr := c.FormValue("user_id")
	projectRoleIDStr := c.FormValue("project_role_id")

	if inviteUserIDStr == "" {
		return helper.Message400("User ID is required")
	}

	if projectRoleIDStr == "" {
		return helper.Message400("Project role ID is required")
	}

	inviteUserID, err := strconv.ParseUint(inviteUserIDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user ID")
	}

	projectRoleID, err := strconv.ParseUint(projectRoleIDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid project role ID")
	}

	err = ctrl.projectMemberService.InviteMember(uint(projectID), uint(inviteUserID), uint(projectRoleID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Invitation sent successfully")
}

// RespondToInvitation allows a user to accept or decline an invitation
func (ctrl *ProjectMemberController) RespondToInvitation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, err := strconv.ParseUint(c.Params("project_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid project ID")
	}

	// Parse form data
	response := c.FormValue("response")

	if response != "accept" && response != "decline" {
		return helper.Message400("Response must be 'accept' or 'decline'")
	}

	err = ctrl.projectMemberService.RespondToInvitation(uint(projectID), userID, response)
	if err != nil {
		return helper.Message400(err.Error())
	}

	message := "Invitation accepted successfully"
	if response == "decline" {
		message = "Invitation declined successfully"
	}

	return helper.Message200(c, nil, message)
}

// GetProjectMembers retrieves all members of a project
func (ctrl *ProjectMemberController) GetProjectMembers(c *fiber.Ctx) error {
	projectID, err := strconv.ParseUint(c.Params("project_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid project ID")
	}

	members, err := ctrl.projectMemberService.GetProjectMembers(uint(projectID))
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, members, "Project members retrieved successfully")
}

// GetApplicationDetails retrieves detailed information about a specific application
func (ctrl *ProjectMemberController) GetApplicationDetails(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	applicationID, err := strconv.ParseUint(c.Params("application_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid application ID")
	}

	application, err := ctrl.projectMemberService.GetApplicationDetails(uint(applicationID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, application, "Application details retrieved successfully")
}

// GetApplicationSummary retrieves a summary of an application
func (ctrl *ProjectMemberController) GetApplicationSummary(c *fiber.Ctx) error {
	applicationID, err := strconv.ParseUint(c.Params("application_id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid application ID")
	}

	summary, err := ctrl.projectMemberService.GetApplicationSummary(uint(applicationID))
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, summary, "Application summary retrieved successfully")
}
