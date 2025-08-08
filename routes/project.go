package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SetupProjectRoutes(app *fiber.App) {
	db := config.GetDB()
	skillService := service.NewSkillService(db)
	tagService := service.NewTagService(db)
	benefitService := service.NewBenefitService(db)
	timelineService := service.NewTimelineService(db)
	ProjectService := service.NewProjectService(db, skillService, tagService, benefitService, timelineService)
	projectController := controller.NewProjectController(ProjectService)

	publicProjects := app.Group("/api/projects")
	publicProjects.Get("/all", projectController.GetAllProjects)
	publicProjects.Get("/public/:id", projectController.GetProjectByID)

	project := app.Group("/api/projects", middleware.AuthMiddleware())

	project.Post("/stage1", projectController.CreateStage1)
	project.Put("/:id/stage2", projectController.UpdateStage2)
	project.Put("/:id/stage3", projectController.UpdateStage3)
	project.Put("/:id/stage4", projectController.UpdateStage4)
	project.Put("/:id/stage5", projectController.UpdateStage5)

	project.Get("/", projectController.GetUserProjects)
	project.Get("/created", projectController.GetMyCreatedProjects)
	project.Get("/member", projectController.GetMyMemberProjects)
	project.Get("/:id", projectController.GetUserProject)
}
