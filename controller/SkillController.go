package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type SkillController struct {
	SkillService *service.SkillService
}

func NewSkillController(skillService *service.SkillService) *SkillController {
	return &SkillController{
		SkillService: skillService,
	}
}
func (ctrl *SkillController) GetAllSkills(c *fiber.Ctx) error {
	skills, err := ctrl.SkillService.GetAllSkills()
	if err != nil {
		return helper.Message500("Failed to retrieve skills: " + err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"skills": skills,
	}, "All skills retrieved successfully")
}

func (ctrl *SkillController) GetUserSkills(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	user, err := ctrl.SkillService.GetUserSkills(userId)
	if err != nil {
		return helper.Message404(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"skills": user.UserSkills,
	}, "User skills retrieved successfully")
}

func (ctrl *SkillController) UpdateSkills(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	var skillNames []string
	var proficiencies []int

	skillsData := c.FormValue("skills")
	proficienciesData := c.FormValue("proficiencies")

	if skillsData != "" && proficienciesData != "" {
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
		form, err := c.MultipartForm()
		if err != nil {
			return helper.Message400("Failed to parse form data")
		}

		if skillValues, exists := form.Value["skill"]; exists {
			for _, skill := range skillValues {
				if skill != "" {
					skillNames = append(skillNames, strings.TrimSpace(skill))
				}
			}
		}

		proficiencyValues := []string{}
		if profValues, exists := form.Value["proficiency"]; exists {
			proficiencyValues = append(proficiencyValues, profValues...)
		}
		if profValues, exists := form.Value["proficiencies"]; exists {
			proficiencyValues = append(proficiencyValues, profValues...)
		}

		for _, profStr := range proficiencyValues {
			if prof, err := strconv.Atoi(strings.TrimSpace(profStr)); err == nil {
				if prof >= 0 && prof <= 100 {
					proficiencies = append(proficiencies, prof)
				}
			}
		}

		if len(skillNames) == 0 {
			return helper.Message400("At least one skill is required. Use either 'skills,proficiencies' format or individual 'skill,proficiency/proficiencies' pairs")
		}

		if len(skillNames) != len(proficiencies) {
			return helper.Message400("Each skill must have a corresponding proficiency value")
		}

		for i, prof := range proficiencies {
			if prof < 0 || prof > 100 {
				return helper.Message400("Proficiency must be between 0 and 100 for skill: " + skillNames[i])
			}
		}
	}

	if err := ctrl.SkillService.UpdateUserSkills(userId, skillNames, proficiencies); err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Skills updated successfully")
}

func (ctrl *SkillController) DeleteUserSkill(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)
	skillName := c.Params("skillName")

	if skillName == "" {
		return helper.Message400("Skill name is required")
	}

	if err := ctrl.SkillService.DeleteUserSkills(userId, skillName); err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "does not have this skill") {
			return helper.Message404(err.Error())
		}
		return helper.Message400(err.Error())
	}
	return helper.Message200(c, nil, "skill deleted succesfully")
}
