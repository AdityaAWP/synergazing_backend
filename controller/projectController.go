package controller

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
	"synergazing.com/synergazing/service"
)

type ProjectController struct {
	projectService *service.ProjectService
}

func NewProjectController(ps *service.ProjectService) *ProjectController {
	return &ProjectController{projectService: ps}
}

func (ctrl *ProjectController) CreateStage1(c *fiber.Ctx) error {
	creatorID := c.Locals("user_id").(uint)
	title := c.FormValue("title")
	projectType := c.FormValue("project_type")
	description := c.FormValue("description")

	file, _ := c.FormFile("picture")
	var pictureURL string
	if file != nil {
		filePath, uploadErr := helper.UploadFile(file, "post")
		if uploadErr != nil {
			return helper.Message400(uploadErr.Error())
		}
		pictureURL = filePath
	}

	project, err := ctrl.projectService.CreateProjectStage1(creatorID, title, projectType, description, pictureURL)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message201(c, project, "Project draft created. Proceed to stage 2.")
}

func (ctrl *ProjectController) UpdateStage2(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var details model.Project
	details.Duration = c.FormValue("duration")
	details.TotalTeam, _ = strconv.Atoi(c.FormValue("total_team"))
	details.StartDate, _ = time.Parse(time.RFC3339, c.FormValue("start_date"))
	details.EndDate, _ = time.Parse(time.RFC3339, c.FormValue("end_date"))
	details.Location = c.FormValue("location")
	details.WorkerType = c.FormValue("worker_type")
	details.Budget, _ = strconv.ParseFloat(c.FormValue("budget"), 64)
	details.RegistrationDeadline, _ = time.Parse(time.RFC3339, c.FormValue("registration_deadline"))

	project, err := ctrl.projectService.CreateProjectStage2(uint(projectID), userID, details)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, project, "Stage 2 completed. Proceed to stage 3.")
}

func (ctrl *ProjectController) UpdateStage3(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	timeCommitment := c.FormValue("time_commitment")
	skillNames, _ := helper.ParseStringSlice(c.FormValue("required_skills"))
	conditions, _ := helper.ParseStringSlice(c.FormValue("conditions"))

	project, err := ctrl.projectService.UpdateStage3(uint(projectID), userID, timeCommitment, skillNames, conditions)
	if err != nil {
		return helper.Message400(err.Error())
	}
	return helper.Message200(c, project, "Stage 3 completed. Proceed to stage 4.")
}

func (ctrl *ProjectController) UpdateStage4(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var requestData struct {
		Roles   []service.RoleDTO   `json:"roles"`
		Members []service.MemberDTO `json:"members"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return helper.Message400("Invalid JSON format: " + err.Error())
	}

	project, err := ctrl.projectService.UpdateStage4(uint(projectID), userID, requestData.Roles, requestData.Members)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, project, "Stage 4 completed. Proceed to finalization.")
}

func (ctrl *ProjectController) UpdateStage5(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	projectID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	benefits := c.FormValue("benefits")
	timeline := c.FormValue("timeline")

	tagsRaw := c.FormValue("tags")
	var tagNames []string

	if tagsRaw != "" {
		if jsonTags, err := helper.ParseStringSlice(tagsRaw); err == nil && len(jsonTags) > 0 {
			tagNames = jsonTags
		} else {
			tagNames = strings.Split(strings.TrimSpace(tagsRaw), ",")
			for i, tag := range tagNames {
				tagNames[i] = strings.TrimSpace(tag)
			}
			var cleanTags []string
			for _, tag := range tagNames {
				if tag != "" {
					cleanTags = append(cleanTags, tag)
				}
			}
			tagNames = cleanTags
		}
	}

	project, err := ctrl.projectService.UpdateStage5(uint(projectID), userID, benefits, timeline, tagNames)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, project, "Project successfully published!")
}
