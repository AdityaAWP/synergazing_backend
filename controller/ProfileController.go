package controller

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type PublicProfileResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	CVFile         string `json:"cv_file"`
	AboutMe        string `json:"about_me"`
	Location       string `json:"location"`
	Interests      string `json:"interests"`
	Academic       string `json:"academic"`
	WebsiteURL     string `json:"website_url"`
	GithubURL      string `json:"github_url"`
	LinkedInURL    string `json:"linkedin_url"`
	InstagramURL   string `json:"instagram_url"`
	PortofolioURL  string `json:"portfolio_url"`
}

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

func (ctrl *ProfileController) GetPublicUserProfile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user ID")
	}

	user, profile, err := ctrl.ProfileService.GetPublicUserProfile(uint(userId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.Message404("User not found")
		}
		return helper.Message500("Could not retrieve user profile")
	}

	publicResponse := PublicProfileResponse{
		ID:             user.ID,
		Name:           user.Name,
		ProfilePicture: helper.GetUrlFile(profile.ProfilePicture),
		CVFile:         helper.GetUrlFile(profile.CVFile),
		AboutMe:        profile.AboutMe,
		Location:       profile.Location,
		Interests:      profile.Interests,
		Academic:       profile.Academic,
		WebsiteURL:     profile.WebsiteURL,
		GithubURL:      profile.GithubURL,
		LinkedInURL:    profile.LinkedInURL,
		InstagramURL:   profile.InstagramURL,
		PortofolioURL:  profile.PortfolioURL,
	}

	return helper.Message200(c, publicResponse, "Profile retrieved successfully")
}

func (ctrl *ProfileController) GetCVFile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Message400("Invalid user ID")
	}

	filePath, err := ctrl.ProfileService.GetCVFilePath(uint(userId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return helper.Message404("Profile not found for this user")
		}
		return helper.Message500("Could not retrieve profile")
	}

	if filePath == "" {
		return helper.Message404("User has not uploaded a CV")
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return helper.Message404("CV file not found on server")
	}

	action := c.Query("action")

	if action == "download" {
		return c.Download(filePath)
	}

	return c.SendFile(filePath, false)
}

func (ctrl *ProfileController) DeleteProfilePicture(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	if err := ctrl.ProfileService.DeleteProfilePicture(userId); err != nil {
		if err.Error() == "no profile picture to delete" {
			return helper.Message404(err.Error())
		}
		return helper.Message500(err.Error())
	}

	return helper.Message200(c, nil, "Profile picture deleted successfully")
}

func (ctrl *ProfileController) DeleteCVFile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	if err := ctrl.ProfileService.DeleteCVFile(userId); err != nil {
		if err.Error() == "no CV file to delete" {
			return helper.Message404(err.Error())
		}
		return helper.Message500(err.Error())
	}

	return helper.Message200(c, nil, "CV file deleted successfully")
}
