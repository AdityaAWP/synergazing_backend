package controller

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

func ListAllUsers(c *fiber.Ctx) error {
	users, err := service.GetAllUser()
	if err != nil {
		return helper.Message500("Failed to retrive users")
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  true,
		"message": "Succesfull to retrieve users",
		"users":   users,
	})
}
