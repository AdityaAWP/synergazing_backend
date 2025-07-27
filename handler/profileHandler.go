package handler

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
			"profile": fiber.Map{
				"id":         profile.ID,
				"created_at": profile.CreatedAt,
				"updated_at": profile.UpdatedAt,
			},
		},
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	// Get form values
	newName := c.FormValue("name")
	newEmail := c.FormValue("email")
	newPassword := c.FormValue("password")

	// Profile picture is optional
	file, err := c.FormFile("profile_picture")
	profilePictureProvided := err == nil

	db := config.GetDB()
	var user model.Users
	var profile model.Profiles

	// Find the user
	if err := db.First(&user, userId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Find or create profile
	if err := db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		profile = model.Profiles{
			UserID:         userId,
			ProfilePicture: "",
		}

		if err := db.Create(&profile).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create profile",
			})
		}
	}

	// Update user fields if provided
	userUpdates := make(map[string]interface{})

	if newName != "" {
		userUpdates["name"] = newName
	}

	if newEmail != "" {
		// Check if email already exists (excluding current user)
		var existingUser model.Users
		if err := db.Where("email = ? AND id != ?", newEmail, userId).First(&existingUser).Error; err == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		userUpdates["email"] = newEmail
	}

	if newPassword != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		userUpdates["password"] = string(hashedPassword)
	}

	// Update user if there are changes
	if len(userUpdates) > 0 {
		if err := db.Model(&user).Updates(userUpdates).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update user information",
			})
		}
	}

	// Handle profile picture update if provided
	var filePath string
	if profilePictureProvided {
		// Delete old profile picture if exists
		if profile.ProfilePicture != "" {
			helper.DeleteFile(profile.ProfilePicture)
		}

		// Upload new profile picture
		filePath, err = helper.UploadFile(file, "profile")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Update profile picture in database
		if err := db.Model(&profile).Update("profile_picture", filePath).Error; err != nil {
			helper.DeleteFile(filePath)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update profile picture",
			})
		}
	}

	// Get updated user data for response
	var updatedUser model.Users
	if err := db.First(&updatedUser, userId).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve updated user",
		})
	}
	updatedUser.Password = ""

	// Prepare response data
	responseData := fiber.Map{
		"id":    updatedUser.ID,
		"name":  updatedUser.Name,
		"email": updatedUser.Email,
	}

	// Add profile picture URL if updated or exists
	if profilePictureProvided {
		responseData["profile_picture"] = helper.GetUrlFile(filePath)
	} else if profile.ProfilePicture != "" {
		responseData["profile_picture"] = helper.GetUrlFile(profile.ProfilePicture)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"data":    responseData,
	})
}
