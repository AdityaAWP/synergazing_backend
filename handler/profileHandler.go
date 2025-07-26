package handler

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
)

func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	db := config.GetDB()
	var user model.Users
	var profile model.Profiles

	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		profile = model.Profiles{
			UserID:         userID,
			ProfilePicture: "",
		}

		if err := db.Create(&profile).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create profile",
			})
		}
	}
	user.Password = ""
	profilePictureURL := helper.GetUrlFile(profile.ProfilePicture)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User profile retrieved successfully",
		"data": fiber.Map{
			"id":              user.ID,
			"name":            user.Name,
			"email":           user.Email,
			"profile_picture": profilePictureURL,
			// "created_at":      user.CreatedAt,
			// "updated_at":      user.UpdatedAt,
			// "profile": fiber.Map{
			// 	"id":         profile.ID,
			// 	"created_at": profile.CreatedAt,
			// 	"updated_at": profile.UpdatedAt,
			// },
		},
	})
}
