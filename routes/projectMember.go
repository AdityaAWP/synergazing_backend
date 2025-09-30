package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SetupProjectMemberRoutes(app *fiber.App) {
	db := config.GetDB()
	notificationService := service.NewNotificationService(db)
	projectMemberService := service.NewProjectMemberService(db, notificationService)
	projectMemberController := controller.NewProjectMemberController(projectMemberService)

	// Protected routes - authentication required
	api := app.Group("/api/projects", middleware.AuthMiddleware())

	// Application management
	api.Post("/:project_id/apply", projectMemberController.ApplyToProject)
	api.Get("/:project_id/applications", projectMemberController.GetProjectApplications)
	api.Get("/applications/:application_id", projectMemberController.GetApplicationDetails)
	api.Get("/applications/:application_id/summary", projectMemberController.GetApplicationSummary)
	api.Put("/applications/:application_id/review", projectMemberController.ReviewApplication)
	api.Put("/applications/:application_id/withdraw", projectMemberController.WithdrawApplication)

	// Member management
	api.Get("/:project_id/members", projectMemberController.GetProjectMembers)
	api.Post("/:project_id/invite", projectMemberController.InviteMember)
	api.Put("/:project_id/invitation/respond", projectMemberController.RespondToInvitation)
	api.Delete("/:project_id/members/:user_id", projectMemberController.RemoveMember)

	// User's own applications and invitations
	userApi := app.Group("/api/user", middleware.AuthMiddleware())
	userApi.Get("/applications", projectMemberController.GetUserApplications)
	userApi.Get("/project-invitations", projectMemberController.GetUserInvitations)
}
