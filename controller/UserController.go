package controller

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
	"synergazing.com/synergazing/service"
)

func ListAllUsers(c *fiber.Ctx) error {
	var users []model.Users
	db := service.GetAllUsersPaginated()

	paginationData, err := helper.Paginate(db, c, &users)
	if err != nil {
		return helper.Message500("Failed to retrieve users")
	}

	return helper.Message200(c, fiber.Map{
		"success":    true,
		"users":      users,
		"pagination": paginationData,
	}, "Successfully retrieved users")
}

func GetReadyUsers(c *fiber.Ctx) error {
	users, err := service.GetReadyUsers()
	if err != nil {
		return helper.Message500("Failed to retrieve ready users")
	}

	return helper.Message200(c, users, "Ready users retrieved successfully")
}
