package controller

import (
	"strconv"

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

func GetUserProfileByID(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user ID")
	}

	profile, err := service.GetUserProfileByID(uint(userID))
	if err != nil {
		if err.Error() == "record not found" {
			return helper.Message404("User not found or not ready for collaboration")
		}
		return helper.Message500("Failed to retrieve user profile")
	}

	return helper.Message200(c, profile, "User profile retrieved successfully")
}
