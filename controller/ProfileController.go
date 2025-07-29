package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
)

func GetUserProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	db := config.GetDB()
	var user model.Users
	var profile model.Profiles

	if err := db.First(&user, userId).Error; err != nil {
		return helper.Message400("User not found")
	}

	if err := db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		profile = model.Profiles{
			UserID:         userId,
			ProfilePicture: "",
		}
		if err := db.Create(&profile).Error; err != nil {
			return helper.Message500("Failed to create profile")
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
			"phone":           user.Phone,
			"profile": fiber.Map{
				"id":            profile.ID,
				"about_me":      profile.AboutMe,
				"location":      profile.Location,
				"interests":     profile.Interests,
				"academic":      profile.Academic,
				"website_url":   profile.WebsiteURL,
				"github_url":    profile.GithubURL,
				"linkedin_url":  profile.LinkedInURL,
				"instagram_url": profile.InstagramURL,
				"portfolio_url": profile.PortofolioURL,
				"created_at":    profile.CreatedAt,
				"updated_at":    profile.UpdatedAt,
			},
		},
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	newName := c.FormValue("name")
	newEmail := c.FormValue("email")
	newPassword := c.FormValue("password")

	newAboutMe := c.FormValue("about_me")
	newLocation := c.FormValue("location")
	newAcademic := c.FormValue("academic")
	newInterests := c.FormValue("interests")

	newWebsiteURL := c.FormValue("website_url")
	newGithubURL := c.FormValue("github_url")
	newLinkedInURL := c.FormValue("linkedin_url")
	newInstagramURL := c.FormValue("instagram_url")
	newPortofolioURL := c.FormValue("portfolio_url")

	file, err := c.FormFile("profile_picture")
	profilePictureProvided := err == nil

	db := config.GetDB()
	var user model.Users
	var profile model.Profiles

	if err := db.First(&user, userId).Error; err != nil {
		return helper.Message404("User not found")
	}

	if err := db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		profile = model.Profiles{
			UserID:         userId,
			ProfilePicture: "",
		}

		if err := db.Create(&profile).Error; err != nil {
			return helper.Message500("Failed to create profile")
		}
	}

	userUpdates := make(map[string]interface{})

	if newName != "" {
		userUpdates["name"] = newName
	}

	if newEmail != "" {
		var existingUser model.Users
		if err := db.Where("email = ? AND id != ?", newEmail, userId).First(&existingUser).Error; err == nil {
			return helper.Message400("Email already exist")
		}
		userUpdates["email"] = newEmail
	}

	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return helper.Message500("Failed to hash password")
		}
		userUpdates["password"] = string(hashedPassword)
	}

	if len(userUpdates) > 0 {
		if err := db.Model(&user).Updates(userUpdates).Error; err != nil {
			return helper.Message500("Failed to update user information")
		}
	}

	profileUpdates := make(map[string]interface{})

	if newAboutMe != "" {
		profileUpdates["about_me"] = newAboutMe
	}
	if newLocation != "" {
		profileUpdates["location"] = newLocation
	}
	if newAcademic != "" {
		profileUpdates["academic"] = newAcademic
	}
	if newInterests != "" {
		profileUpdates["interests"] = newInterests
	}
	if newWebsiteURL != "" {
		profileUpdates["website_url"] = newWebsiteURL
	}
	if newGithubURL != "" {
		profileUpdates["github_url"] = newGithubURL
	}
	if newLinkedInURL != "" {
		profileUpdates["linked_in_url"] = newLinkedInURL
	}
	if newInstagramURL != "" {
		profileUpdates["instagram_url"] = newInstagramURL
	}
	if newPortofolioURL != "" {
		profileUpdates["portofolio_url"] = newPortofolioURL
	}

	// Apply profile updates to the database
	fmt.Printf("DEBUG - Profile updates to apply: %+v\n", profileUpdates)
	if len(profileUpdates) > 0 {
		if err := db.Model(&profile).Updates(profileUpdates).Error; err != nil {
			fmt.Printf("DEBUG - Error updating profile: %v\n", err)
			return helper.Message500("Failed to update profile information")
		}
		fmt.Printf("DEBUG - Profile updates applied successfully\n")
	} else {
		fmt.Printf("DEBUG - No profile updates to apply\n")
	}

	var filePath string
	if profilePictureProvided {
		if profile.ProfilePicture != "" {
			helper.DeleteFile(profile.ProfilePicture)
		}

		filePath, err = helper.UploadFile(file, "profile")
		if err != nil {
			return helper.Message400(err.Error())
		}

		if err := db.Model(&profile).Update("profile_picture", filePath).Error; err != nil {
			helper.DeleteFile(filePath)
			return helper.Message500("Failed to update profile picture")
		}
	}

	var updatedUser model.Users
	if err := db.First(&updatedUser, userId).Error; err != nil {
		return helper.Message500("Failed to retrieve updated user")
	}
	updatedUser.Password = ""

	var updatedProfile model.Profiles
	if err := db.Where("user_id = ?", userId).First(&updatedProfile).Error; err != nil {
		return helper.Message500("Failed to retrieve updated profile")
	}

	responseData := fiber.Map{
		"id":    updatedUser.ID,
		"name":  updatedUser.Name,
		"email": updatedUser.Email,
		"phone": updatedUser.Phone,
		"profile": fiber.Map{
			"id":            updatedProfile.ID,
			"about_me":      updatedProfile.AboutMe,
			"location":      updatedProfile.Location,
			"interests":     updatedProfile.Interests,
			"academic":      updatedProfile.Academic,
			"website_url":   updatedProfile.WebsiteURL,
			"github_url":    updatedProfile.GithubURL,
			"linkedin_url":  updatedProfile.LinkedInURL,
			"instagram_url": updatedProfile.InstagramURL,
			"portfolio_url": updatedProfile.PortofolioURL,
			"created_at":    updatedProfile.CreatedAt,
			"updated_at":    updatedProfile.UpdatedAt,
		},
	}

	if profilePictureProvided {
		responseData["profile_picture"] = helper.GetUrlFile(filePath)
	} else if updatedProfile.ProfilePicture != "" {
		responseData["profile_picture"] = helper.GetUrlFile(updatedProfile.ProfilePicture)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"data":    responseData,
	})
}
