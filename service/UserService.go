package service

import (
	"gorm.io/gorm"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
)

// ReadyUserResponse represents the response for ready users
type ReadyUserResponse struct {
	ID             uint               `json:"id"`
	Name           string             `json:"name"`
	ProfilePicture string             `json:"profile_picture"`
	AboutMe        string             `json:"about_me"`
	Location       string             `json:"location"`
	Interests      string             `json:"interests"`
	Academic       string             `json:"academic"`
	Skills         []*model.UserSkill `json:"skills"`
}

// UserProfileResponse represents the full user profile response
type UserProfileResponse struct {
	ID             uint               `json:"id"`
	Name           string             `json:"name"`
	ProfilePicture string             `json:"profile_picture"`
	CVFile         string             `json:"cv_file"`
	AboutMe        string             `json:"about_me"`
	Location       string             `json:"location"`
	Interests      string             `json:"interests"`
	Academic       string             `json:"academic"`
	WebsiteURL     string             `json:"website_url"`
	GithubURL      string             `json:"github_url"`
	LinkedInURL    string             `json:"linkedin_url"`
	InstagramURL   string             `json:"instagram_url"`
	PortofolioURL  string             `json:"portofolio_url"`
	Skills         []*model.UserSkill `json:"skills"`
}

func GetAllUser() ([]model.Users, error) {
	var user []model.Users
	result := config.DB.Find(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func GetAllUsersPaginated() *gorm.DB {
	return config.DB.Model(&model.Users{})
}

func GetReadyUsers() ([]ReadyUserResponse, error) {
	var users []model.Users
	var profiles []model.Profiles

	// Get users with ready status
	result := config.DB.Preload("UserSkills.Skill").Where("status_collaboration = ?", "ready").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// Get user IDs
	var userIDs []uint
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	// Get profiles for these users
	profileResult := config.DB.Where("user_id IN ?", userIDs).Find(&profiles)
	if profileResult.Error != nil {
		return nil, profileResult.Error
	}

	// Create a map of profiles by user_id for quick lookup
	profileMap := make(map[uint]model.Profiles)
	for _, profile := range profiles {
		profileMap[profile.UserID] = profile
	}

	// Transform to response format
	var response []ReadyUserResponse
	for _, user := range users {
		profile, exists := profileMap[user.ID]
		readyUser := ReadyUserResponse{
			ID:             user.ID,
			Name:           user.Name,
			ProfilePicture: "",
			AboutMe:        "",
			Location:       "",
			Interests:      "",
			Academic:       "",
			Skills:         user.UserSkills,
		}

		if exists {
			readyUser.ProfilePicture = helper.GetUrlFile(profile.ProfilePicture)
			readyUser.AboutMe = profile.AboutMe
			readyUser.Location = profile.Location
			readyUser.Interests = profile.Interests
			readyUser.Academic = profile.Academic
		}

		response = append(response, readyUser)
	}

	return response, nil
}

func GetUserProfileByID(userID uint) (*UserProfileResponse, error) {
	var user model.Users
	var profile model.Profiles

	// Get user with ready status
	userResult := config.DB.Preload("UserSkills.Skill").Where("id = ? AND status_collaboration = ?", userID, "ready").First(&user)
	if userResult.Error != nil {
		return nil, userResult.Error
	}

	// Get profile for this user
	profileResult := config.DB.Where("user_id = ?", userID).First(&profile)
	if profileResult.Error != nil {
		return nil, profileResult.Error
	}

	// Transform to response format
	response := &UserProfileResponse{
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
		Skills:         user.UserSkills,
	}

	return response, nil
}

func GetReadyUsersPaginated() *gorm.DB {
	return config.DB.Preload("UserSkills.Skill").Where("status_collaboration = ?", "ready")
}

func GetReadyUsersPaginatedWithTransform(page, perPage int) ([]ReadyUserResponse, *helper.PaginationData, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	var users []model.Users
	var profiles []model.Profiles
	var totalRecords int64

	// Get total count
	countResult := config.DB.Model(&model.Users{}).Where("status_collaboration = ?", "ready").Count(&totalRecords)
	if countResult.Error != nil {
		return nil, nil, countResult.Error
	}

	// Calculate pagination
	totalPages := int((totalRecords + int64(perPage) - 1) / int64(perPage))
	offset := (page - 1) * perPage

	// Get paginated users with ready status
	result := config.DB.Preload("UserSkills.Skill").Where("status_collaboration = ?", "ready").
		Limit(perPage).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	// Get user IDs
	var userIDs []uint
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	// Get profiles for these users
	if len(userIDs) > 0 {
		profileResult := config.DB.Where("user_id IN ?", userIDs).Find(&profiles)
		if profileResult.Error != nil {
			return nil, nil, profileResult.Error
		}
	}

	// Create a map of profiles by user_id for quick lookup
	profileMap := make(map[uint]model.Profiles)
	for _, profile := range profiles {
		profileMap[profile.UserID] = profile
	}

	// Transform to response format
	var response []ReadyUserResponse
	for _, user := range users {
		profile, exists := profileMap[user.ID]
		readyUser := ReadyUserResponse{
			ID:             user.ID,
			Name:           user.Name,
			ProfilePicture: "",
			AboutMe:        "",
			Location:       "",
			Interests:      "",
			Academic:       "",
			Skills:         user.UserSkills,
		}

		if exists {
			readyUser.ProfilePicture = helper.GetUrlFile(profile.ProfilePicture)
			readyUser.AboutMe = profile.AboutMe
			readyUser.Location = profile.Location
			readyUser.Interests = profile.Interests
			readyUser.Academic = profile.Academic
		}

		response = append(response, readyUser)
	}

	// Create pagination data
	paginationData := &helper.PaginationData{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PerPage:      perPage,
		NextPage:     nil,
		PrevPage:     nil,
	}

	if page < totalPages {
		next := page + 1
		paginationData.NextPage = &next
	}

	if page > 1 {
		prev := page - 1
		paginationData.PrevPage = &prev
	}

	return response, paginationData, nil
}
