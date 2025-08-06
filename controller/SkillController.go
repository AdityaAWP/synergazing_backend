package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type SkillController struct {
	ProfileService *service.ProfileService
}

func NewSkillController(profileService *service.ProfileService) *SkillController {
	return &SkillController{
		ProfileService: profileService,
	}
}

func (ctrl *SkillController) UpdateSkills(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	var skillNames []string
	var proficiencies []int

	// Debug: Print all form data received
	fmt.Println("=== DEBUG: Form data received ===")
	args := c.Request().PostArgs()
	args.VisitAll(func(key, value []byte) {
		fmt.Printf("Key: '%s', Value: '%s'\n", string(key), string(value))
	})
	fmt.Println("=== END DEBUG ===")

	// Check if using comma-separated format (legacy support)
	skillsData := c.FormValue("skills")
	proficienciesData := c.FormValue("proficiencies")

	if skillsData != "" && proficienciesData != "" {
		// Handle comma-separated format: skills="react,java,js" proficiencies="90,85,80"
		skillNames = strings.Split(skillsData, ",")
		for i, skill := range skillNames {
			skillNames[i] = strings.TrimSpace(skill)
		}

		proficiencyStrs := strings.Split(proficienciesData, ",")
		if len(skillNames) != len(proficiencyStrs) {
			return helper.Message400("Skills and proficiencies count must match")
		}

		proficiencies = make([]int, len(proficiencyStrs))
		for i, profStr := range proficiencyStrs {
			profStr = strings.TrimSpace(profStr)
			if prof, err := strconv.Atoi(profStr); err == nil {
				if prof < 0 || prof > 100 {
					return helper.Message400("Proficiency must be between 0 and 100")
				}
				proficiencies[i] = prof
			} else {
				return helper.Message400("Invalid proficiency value: " + profStr)
			}
		}
	} else {
		// Handle individual field format: skill=react proficiency=90 skill=java proficiency=85
		// Get all form values (handles multiple values with same key)
		form, err := c.MultipartForm()
		if err != nil {
			return helper.Message400("Failed to parse form data")
		}

		// Extract skills
		if skillValues, exists := form.Value["skill"]; exists {
			for _, skill := range skillValues {
				if skill != "" {
					skillNames = append(skillNames, strings.TrimSpace(skill))
				}
			}
		}

		// Extract proficiencies (handle both 'proficiency' and 'proficiencies')
		proficiencyValues := []string{}
		if profValues, exists := form.Value["proficiency"]; exists {
			proficiencyValues = append(proficiencyValues, profValues...)
		}
		if profValues, exists := form.Value["proficiencies"]; exists {
			proficiencyValues = append(proficiencyValues, profValues...)
		}

		// Convert proficiencies to integers
		for _, profStr := range proficiencyValues {
			if prof, err := strconv.Atoi(strings.TrimSpace(profStr)); err == nil {
				if prof >= 0 && prof <= 100 {
					proficiencies = append(proficiencies, prof)
				}
			}
		}

		// Validate input for individual field format
		if len(skillNames) == 0 {
			return helper.Message400("At least one skill is required. Use either 'skills,proficiencies' format or individual 'skill,proficiency/proficiencies' pairs")
		}

		if len(skillNames) != len(proficiencies) {
			return helper.Message400("Each skill must have a corresponding proficiency value")
		}

		// Validate proficiencies
		for i, prof := range proficiencies {
			if prof < 0 || prof > 100 {
				return helper.Message400("Proficiency must be between 0 and 100 for skill: " + skillNames[i])
			}
		}
	}

	if err := ctrl.ProfileService.UpdateUserSkills(userId, skillNames, proficiencies); err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Skills updated successfully")
}

func (ctrl *SkillController) GetUserSkills(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	user, _, err := ctrl.ProfileService.GetUserProfile(userId)
	if err != nil {
		return helper.Message404(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"skills": user.UserSkills,
	}, "User skills retrieved successfully")
}

// func (ctrl *SkillController) DeleteSkill(c *fiber.Ctx) error {
// 	userId := c.Locals("user_id").(uint)
// 	skillIdStr := c.Params("skillId")

// 	skillId, err := strconv.ParseUint(skillIdStr, 10, 32)
// 	if err != nil {
// 		return helper.Message400("Invalid skill ID")
// 	}

// 	// You can implement a DeleteUserSkill method in ProfileService if needed
// 	// For now, return a message indicating the functionality needs to be implemented
// 	return helper.Message400("Delete individual skill functionality needs to be implemented in ProfileService")
// }

func (ctrl *SkillController) GetAllSkills(c *fiber.Ctx) error {
	skills, err := ctrl.ProfileService.SkillService.GetAllSkills()
	if err != nil {
		return helper.Message500("Failed to retrieve skills: " + err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"skills": skills,
	}, "All skills retrieved successfully")
}
