package controller

import (
	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type ProfileController struct {
	ProfileService *service.ProfileService
}

func NewProfileController(s *service.ProfileService) *ProfileController {
	return &ProfileController{
		ProfileService: s,
	}
}

func (ctrl *ProfileController) GetUserProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	user, profile, err := ctrl.ProfileService.GetUserProfile(userId)
	if err != nil {
		return helper.Message404(err.Error())
	}

	responseData := fiber.Map{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"phone":           user.Phone,
		"profile_picture": helper.GetUrlFile(profile.ProfilePicture),
		"cv_file":         helper.GetUrlFile(profile.CVFile),
		"profile": fiber.Map{
			"about_me":      profile.AboutMe,
			"location":      profile.Location,
			"interests":     profile.Interests,
			"academic":      profile.Academic,
			"website_url":   profile.WebsiteURL,
			"github_url":    profile.GithubURL,
			"linkedin_url":  profile.LinkedInURL,
			"instagram_url": profile.InstagramURL,
			"portfolio_url": profile.PortfolioURL,
		},
	}

	return helper.Message200(c, responseData, "User profile retrieved successfully")
}

func (ctrl *ProfileController) UpdateProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	dto := new(service.UpdateProfileDTO)

	if name := c.FormValue("name"); name != "" {
		dto.Name = &name
	}
	if email := c.FormValue("email"); email != "" {
		dto.Email = &email
	}
	if phone := c.FormValue("phone"); phone != "" {
		dto.Phone = &phone
	}
	if password := c.FormValue("password"); password != "" {
		dto.Password = &password
	}
	if aboutMe := c.FormValue("about_me"); aboutMe != "" {
		dto.AboutMe = &aboutMe
	}
	if location := c.FormValue("location"); location != "" {
		dto.Location = &location
	}
	if interests := c.FormValue("interests"); interests != "" {
		dto.Interests = &interests
	}
	if academic := c.FormValue("academic"); academic != "" {
		dto.Academic = &academic
	}
	if websiteURL := c.FormValue("website_url"); websiteURL != "" {
		dto.WebsiteURL = &websiteURL
	}
	if githubURL := c.FormValue("github_url"); githubURL != "" {
		dto.GithubURL = &githubURL
	}
	if linkedInURL := c.FormValue("linkedin_url"); linkedInURL != "" {
		dto.LinkedInURL = &linkedInURL
	}
	if instagramURL := c.FormValue("instagram_url"); instagramURL != "" {
		dto.InstagramURL = &instagramURL
	}
	if portfolioURL := c.FormValue("portfolio_url"); portfolioURL != "" {
		dto.PortofolioURL = &portfolioURL
	}

	if profilePic, err := c.FormFile("profile_picture"); err == nil {
		dto.ProfilePicture = profilePic
	}
	if cv, err := c.FormFile("cv_file"); err == nil {
		dto.CVFile = cv
	}

	user, profile, err := ctrl.ProfileService.UpdateUserProfile(userId, dto)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"phone":           user.Phone,
		"profile_picture": helper.GetUrlFile(profile.ProfilePicture),
		"cv_file":         helper.GetUrlFile(profile.CVFile),
		"profile":         profile,
	}, "Profile updated successfully")
}
