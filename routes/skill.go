package routes

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/controller"
	"synergazing.com/synergazing/middleware"
	"synergazing.com/synergazing/service"
)

func SkillRoutes(app *fiber.App) {
	profileService := service.NewProfileService()

	skillController := controller.NewSkillController(profileService)

	skillGroup := app.Group("/api/skills")

	skillGroup.Get("/all", skillController.GetAllSkills)

	skillGroup.Use(middleware.AuthMiddleware())

	skillGroup.Post("/", skillController.UpdateSkills)

	skillGroup.Get("/", skillController.GetUserSkills)

	// skillGroup.Delete("/:skillId", skillController.DeleteSkill)
}
