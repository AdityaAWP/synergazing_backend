package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/model"
)

type ProjectService struct {
	DB              *gorm.DB
	skillService    *SkillService
	tagService      *TagService
	benefitService  *BenefitService
	timelineService *TimelineService
}

type RoleDTO struct {
	Name           string   `json:"name"`
	SlotsAvailable int      `json:"slots_available"`
	Description    string   `json:"description"`
	SkillNames     []string `json:"skill_names"`
}

type MemberDTO struct {
	Name            string   `json:"name"`
	RoleName        string   `json:"role_name"`
	RoleDescription string   `json:"role_description"`
	SkillNames      []string `json:"skill_names"`
}

type MemberResponse struct {
	Name            string   `json:"name"`
	RoleDescription string   `json:"role_description"`
	RoleName        string   `json:"role_name"`
	SkillNames      []string `json:"skill_names"`
}

type ProjectResponse struct {
	*model.Project
	Members []MemberResponse `json:"members"`
}

type ProjectResponseForMarshal struct {
	ID                   uint                          `json:"id"`
	Title                string                        `json:"title"`
	Description          string                        `json:"description"`
	CompletionStage      int                           `json:"completion_stage"`
	CreatorID            uint                          `json:"creator_id"`
	Creator              model.Users                   `json:"creator"`
	Status               string                        `json:"status"`
	ProjectType          string                        `json:"project_type"`
	PictureURL           string                        `json:"picture_url"`
	Duration             string                        `json:"duration"`
	TotalTeam            int                           `json:"total_team"`
	FilledTeam           int                           `json:"filled_team"`
	RemainingTeam        int                           `json:"remaining_team"`
	StartDate            string                        `json:"start_date"`
	EndDate              string                        `json:"end_date"`
	Location             string                        `json:"location"`
	Budget               string                        `json:"budget"`
	RegistrationDeadline string                        `json:"registration_deadline"`
	TimeCommitment       string                        `json:"time_commitment"`
	Benefits             []*model.ProjectBenefit       `json:"benefits"`
	Timeline             []*model.ProjectTimeline      `json:"timeline"`
	RequiredSkills       []*model.ProjectRequiredSkill `json:"required_skills"`
	Conditions           []*model.ProjectCondition     `json:"conditions"`
	Tags                 []*model.ProjectTag           `json:"tags"`
	Members              []MemberResponse              `json:"members"`
	Roles                []*model.ProjectRole          `json:"roles"`
	CreatedAt            string                        `json:"created_at"`
	UpdatedAt            string                        `json:"updated_at"`
}

func NewProjectService(db *gorm.DB, skillService *SkillService, tagService *TagService, benefitService *BenefitService, timelineService *TimelineService) *ProjectService {
	return &ProjectService{
		DB:              db,
		skillService:    skillService,
		tagService:      tagService,
		benefitService:  benefitService,
		timelineService: timelineService,
	}
}

func (s *ProjectService) getProjectForUpdate(tx *gorm.DB, projectID, userID uint, requiredStage int) (model.Project, error) {
	var project model.Project
	if err := tx.First(&project, projectID).Error; err != nil {
		return project, errors.New("project not found")
	}
	if project.CreatorID != userID {
		return project, errors.New("you are not authorized to edit this project")
	}
	if project.CompletionStage < requiredStage {
		return project, fmt.Errorf("you must complete the previous stage %d", requiredStage)
	}
	return project, nil
}

func (s *ProjectService) transformProjectToResponse(project *model.Project) interface{} {
	memberResponses := make([]MemberResponse, len(project.Members))
	for i, member := range project.Members {
		skillNames := make([]string, len(member.MemberSkills))
		for j, memberSkill := range member.MemberSkills {
			skillNames[j] = memberSkill.Skill.Name
		}

		memberResponses[i] = MemberResponse{
			Name:            member.User.Name,
			RoleDescription: member.RoleDescription,
			RoleName:        member.ProjectRole.Name,
			SkillNames:      skillNames,
		}
	}

	// Calculate team capacity
	filledTeam, _, remainingTeam := s.calculateTeamCapacity(project)

	// Format dates to RFC3339 strings
	var startDateStr, endDateStr, registrationDeadlineStr, createdAtStr, updatedAtStr string
	if !project.StartDate.IsZero() {
		startDateStr = project.StartDate.Format("2006-01-02T15:04:05Z07:00")
	}
	if !project.EndDate.IsZero() {
		endDateStr = project.EndDate.Format("2006-01-02T15:04:05Z07:00")
	}
	if !project.RegistrationDeadline.IsZero() {
		registrationDeadlineStr = project.RegistrationDeadline.Format("2006-01-02T15:04:05Z07:00")
	}
	if !project.CreatedAt.IsZero() {
		createdAtStr = project.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if !project.UpdatedAt.IsZero() {
		updatedAtStr = project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	response := ProjectResponseForMarshal{
		ID:                   project.ID,
		Title:                project.Title,
		Description:          project.Description,
		CompletionStage:      project.CompletionStage,
		CreatorID:            project.CreatorID,
		Creator:              project.Creator,
		Status:               project.Status,
		ProjectType:          project.ProjectType,
		PictureURL:           helper.GetUrlFile(project.PictureURL),
		Duration:             project.Duration,
		TotalTeam:            project.TotalTeam,
		FilledTeam:           filledTeam,
		RemainingTeam:        remainingTeam,
		StartDate:            startDateStr,
		EndDate:              endDateStr,
		Location:             project.Location,
		Budget:               project.Budget,
		RegistrationDeadline: registrationDeadlineStr,
		TimeCommitment:       project.TimeCommitment,
		Benefits:             project.Benefits,
		Timeline:             project.Timeline,
		RequiredSkills:       project.RequiredSkills,
		Conditions:           project.Conditions,
		Tags:                 project.Tags,
		Members:              memberResponses,
		Roles:                project.Roles,
		CreatedAt:            createdAtStr,
		UpdatedAt:            updatedAtStr,
	}

	return response
}

func (s *ProjectService) loadProjectWithRelationships(projectID uint) (*model.Project, error) {
	var project model.Project
	if err := s.DB.Preload("Creator").
		Preload("RequiredSkills.Skill").
		Preload("Conditions").
		Preload("Roles.RequiredSkills.Skill").
		Preload("Members.User").
		Preload("Members.ProjectRole.RequiredSkills.Skill").
		Preload("Members.MemberSkills.Skill").
		Preload("Tags.Tag").
		Preload("Benefits.Benefit").
		Preload("Timeline.Timeline").
		First(&project, projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) CreateProjectStage1(userID uint, title, projectType, description, pictureURL string) (*model.Project, error) {
	if title == "" {
		return nil, errors.New("title are required")
	}
	if description == "" {
		return nil, errors.New("description are required")
	}
	if projectType == "" {
		return nil, errors.New("project type are required")
	}

	project := model.Project{
		CreatorID:       userID,
		Title:           title,
		ProjectType:     projectType,
		Description:     description,
		PictureURL:      pictureURL,
		Status:          "draft",
		CompletionStage: 1,
	}
	if err := s.DB.Create(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) CreateProjectStage2(ProjectID, userID uint, details model.Project) (*model.Project, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, ProjectID, userID, 1)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	project.Duration = details.Duration
	project.TotalTeam = details.TotalTeam
	project.StartDate = details.StartDate
	project.EndDate = details.EndDate
	project.Location = details.Location
	project.Budget = details.Budget
	project.RegistrationDeadline = details.RegistrationDeadline
	project.CompletionStage = 2

	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return &project, tx.Commit().Error
}

func (s *ProjectService) UpdateStage3(projectID, userID uint, timeCommitment string, skillNames []string, conditionDescriptions []string) (interface{}, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 2)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if timeCommitment == "" {
		tx.Rollback()
		return nil, errors.New("Time commitment are required for stage 3")
	}
	if len(skillNames) == 0 {
		tx.Rollback()
		return nil, errors.New("skill are required for stage 3")
	}

	project.TimeCommitment = timeCommitment

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectRequiredSkill{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, skillName := range skillNames {
		skill, err := s.skillService.FindOrCreateWithTx(tx, skillName)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		projectSkill := model.ProjectRequiredSkill{
			ProjectID: projectID,
			SkillID:   skill.ID,
		}
		if err := tx.Create(&projectSkill).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Where("project_id = ?", projectID).Delete(&model.ProjectCondition{})
	for _, desc := range conditionDescriptions {
		condition := model.ProjectCondition{ProjectID: project.ID, Description: desc}
		if err := tx.Create(&condition).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	project.CompletionStage = 3
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	projectResult, err := s.loadProjectWithRelationships(project.ID)
	if err != nil {
		return nil, err
	}

	return s.transformProjectToResponse(projectResult), nil
}

func (s *ProjectService) UpdateStage4(projectID, userID uint, roles []RoleDTO, members []MemberDTO) (interface{}, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Validate team capacity
	totalMembers := len(members)
	totalRoleSlots := 0
	for _, role := range roles {
		totalRoleSlots += role.SlotsAvailable
	}

	if totalMembers+totalRoleSlots > project.TotalTeam {
		tx.Rollback()
		return nil, fmt.Errorf("Cannot add %d members and %d role slots. Total capacity is %d, but you're trying to allocate %d positions. Remaining capacity: %d",
			totalMembers, totalRoleSlots, project.TotalTeam, totalMembers+totalRoleSlots, project.TotalTeam-(totalMembers+totalRoleSlots))
	}

	if totalMembers+totalRoleSlots == 0 {
		tx.Rollback()
		return nil, errors.New("You must add at least one member or create at least one role for team recruitment")
	}

	var existingMembers []model.ProjectMember
	if err := tx.Where("project_id = ?", projectID).Find(&existingMembers).Error; err == nil {
		for _, member := range existingMembers {
			tx.Where("project_member_id = ?", member.ID).Delete(&model.ProjectMemberSkill{})
		}
	}
	tx.Where("project_id = ?", projectID).Delete(&model.ProjectMember{})

	var existingRoles []model.ProjectRole
	if err := tx.Where("project_id = ?", projectID).Find(&existingRoles).Error; err == nil {
		for _, role := range existingRoles {
			tx.Where("project_role_id = ?", role.ID).Delete(&model.ProjectRoleSkill{})
		}
	}

	tx.Where("project_id = ?", projectID).Delete(&model.ProjectRole{})

	roleMap := make(map[string]uint)

	for _, roleData := range roles {
		role := model.ProjectRole{
			ProjectID:      project.ID,
			Name:           roleData.Name,
			SlotsAvailable: roleData.SlotsAvailable,
			Description:    roleData.Description,
		}
		if err := tx.Create(&role).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		roleMap[role.Name] = role.ID

		if len(roleData.SkillNames) > 0 {
			for _, skillName := range roleData.SkillNames {
				skill, err := s.skillService.FindOrCreateWithTx(tx, skillName)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				roleSkill := model.ProjectRoleSkill{
					ProjectRoleID: role.ID,
					SkillID:       skill.ID,
				}
				if err := tx.Create(&roleSkill).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	for _, memberData := range members {
		var user model.Users
		if err := tx.Where("name = ?", memberData.Name).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("user to invite not found: " + memberData.Name)
		}
		roleID, ok := roleMap[memberData.RoleName]
		if !ok {
			tx.Rollback()
			return nil, errors.New("role specified for member does not exist: " + memberData.RoleName)
		}
		member := model.ProjectMember{
			ProjectID:       project.ID,
			UserID:          user.ID,
			ProjectRoleID:   roleID,
			Status:          "invited",
			RoleDescription: memberData.RoleDescription,
		}
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(memberData.SkillNames) > 0 {
			for _, skillName := range memberData.SkillNames {
				skill, err := s.skillService.FindOrCreateWithTx(tx, skillName)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				memberSkill := model.ProjectMemberSkill{
					ProjectMemberID: member.ID,
					SkillID:         skill.ID,
				}
				if err := tx.Create(&memberSkill).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	project.CompletionStage = 4
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	projectResult, err := s.loadProjectWithRelationships(project.ID)
	if err != nil {
		return nil, err
	}

	return s.transformProjectToResponse(projectResult), nil
}

// calculateTeamCapacity calculates team statistics for a project
func (s *ProjectService) calculateTeamCapacity(project *model.Project) (filledTeam, totalRoleSlots, remainingTeam int) {
	filledTeam = len(project.Members)
	totalRoleSlots = 0
	for _, role := range project.Roles {
		totalRoleSlots += role.SlotsAvailable
	}
	remainingTeam = project.TotalTeam - filledTeam - totalRoleSlots
	return
}

// GetProjectTeamCapacity returns team capacity information for a project
func (s *ProjectService) GetProjectTeamCapacity(projectID, userID uint) (map[string]interface{}, error) {
	_, err := s.GetUserProject(userID, projectID)
	if err != nil {
		return nil, err
	}

	// Load the full project with relationships
	var fullProject model.Project
	if err := s.DB.Preload("Members").Preload("Roles").Where("id = ?", projectID).First(&fullProject).Error; err != nil {
		return nil, errors.New("project not found")
	}

	filledTeam, totalRoleSlots, remainingTeam := s.calculateTeamCapacity(&fullProject)

	return map[string]interface{}{
		"total_team":       fullProject.TotalTeam,
		"filled_team":      filledTeam,
		"total_role_slots": totalRoleSlots,
		"remaining_team":   remainingTeam,
		"members":          len(fullProject.Members),
		"roles":            len(fullProject.Roles),
	}, nil
}

func (s *ProjectService) UpdateStage5(projectID, userID uint, benefitNames, timelineNames, tagNames []string) (interface{}, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 4)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(benefitNames) == 0 {
		tx.Rollback()
		return nil, errors.New("at least one benefit is required")
	}

	project.CompletionStage = 5
	project.Status = "published"

	if len(benefitNames) > 0 {
		benefits, err := s.benefitService.findOrCreate(tx, benefitNames)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Where("project_id = ?", project.ID).Delete(&model.ProjectBenefit{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, benefit := range benefits {
			projectBenefit := &model.ProjectBenefit{
				ProjectID: project.ID,
				BenefitID: benefit.ID,
			}
			if err := tx.Create(projectBenefit).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if len(timelineNames) > 0 {
		timelines, err := s.timelineService.findOrCreate(tx, timelineNames)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Where("project_id = ?", project.ID).Delete(&model.ProjectTimeline{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, timeline := range timelines {
			projectTimeline := &model.ProjectTimeline{
				ProjectID:  project.ID,
				TimelineID: timeline.ID,
			}
			if err := tx.Create(projectTimeline).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if len(tagNames) > 0 {
		tags, err := s.tagService.findOrCreate(tx, tagNames)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Where("project_id = ?", project.ID).Delete(&model.ProjectTag{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, tag := range tags {
			projectTag := &model.ProjectTag{
				ProjectID: project.ID,
				TagID:     tag.ID,
			}
			if err := tx.Create(projectTag).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	projectResult, err := s.loadProjectWithRelationships(project.ID)
	if err != nil {
		return nil, err
	}

	return s.transformProjectToResponse(projectResult), nil
}

func (s *ProjectService) CreateRolesOnly(projectID, userID uint, roles []RoleDTO) (interface{}, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var existingRoles []model.ProjectRole
	if err := tx.Where("project_id = ?", projectID).Find(&existingRoles).Error; err == nil {
		for _, role := range existingRoles {
			tx.Where("project_role_id = ?", role.ID).Delete(&model.ProjectRoleSkill{})
		}
	}
	tx.Where("project_id = ?", projectID).Delete(&model.ProjectRole{})

	for _, roleData := range roles {
		role := model.ProjectRole{
			ProjectID:      project.ID,
			Name:           roleData.Name,
			SlotsAvailable: roleData.SlotsAvailable,
			Description:    roleData.Description,
		}
		if err := tx.Create(&role).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(roleData.SkillNames) > 0 {
			for _, skillName := range roleData.SkillNames {
				skill, err := s.skillService.FindOrCreateWithTx(tx, skillName)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				roleSkill := model.ProjectRoleSkill{
					ProjectRoleID: role.ID,
					SkillID:       skill.ID,
				}
				if err := tx.Create(&roleSkill).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	projectResult, err := s.loadProjectWithRelationships(project.ID)
	if err != nil {
		return nil, err
	}

	return s.transformProjectToResponse(projectResult), nil
}

func (s *ProjectService) AddMembersOnly(projectID, userID uint, members []MemberDTO) (interface{}, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var existingRoles []model.ProjectRole
	if err := tx.Where("project_id = ?", projectID).Find(&existingRoles).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch project roles")
	}

	roleMap := make(map[string]uint)
	for _, role := range existingRoles {
		roleMap[role.Name] = role.ID
	}

	var existingMembers []model.ProjectMember
	if err := tx.Where("project_id = ?", projectID).Find(&existingMembers).Error; err == nil {
		for _, member := range existingMembers {
			tx.Where("project_member_id = ?", member.ID).Delete(&model.ProjectMemberSkill{})
		}
	}
	tx.Where("project_id = ?", projectID).Delete(&model.ProjectMember{})

	for _, memberData := range members {
		var user model.Users
		if err := tx.Where("name = ?", memberData.Name).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("user to invite not found: " + memberData.Name)
		}

		roleID, ok := roleMap[memberData.RoleName]
		if !ok {
			tx.Rollback()
			return nil, errors.New("role specified for member does not exist: " + memberData.RoleName)
		}

		member := model.ProjectMember{
			ProjectID:       project.ID,
			UserID:          user.ID,
			ProjectRoleID:   roleID,
			Status:          "invited",
			RoleDescription: memberData.RoleDescription,
		}
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(memberData.SkillNames) > 0 {
			for _, skillName := range memberData.SkillNames {
				skill, err := s.skillService.FindOrCreateWithTx(tx, skillName)
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				memberSkill := model.ProjectMemberSkill{
					ProjectMemberID: member.ID,
					SkillID:         skill.ID,
				}
				if err := tx.Create(&memberSkill).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	project.CompletionStage = 4
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	projectResult, err := s.loadProjectWithRelationships(project.ID)
	if err != nil {
		return nil, err
	}

	return s.transformProjectToResponse(projectResult), nil
}

func (s *ProjectService) GetUserProjects(userID uint) ([]interface{}, error) {
	var projects []model.Project

	err := s.DB.Preload("Creator").
		Preload("RequiredSkills.Skill").
		Preload("Conditions").
		Preload("Roles.RequiredSkills.Skill").
		Preload("Members.User").
		Preload("Members.ProjectRole.RequiredSkills.Skill").
		Preload("Members.MemberSkills.Skill").
		Preload("Tags.Tag").
		Preload("Benefits.Benefit").
		Preload("Timeline.Timeline").
		Where("creator_id = ? OR id IN (SELECT project_id FROM project_members WHERE user_id = ?)", userID, userID).
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user projects: %w", err)
	}

	var responses []interface{}
	for _, project := range projects {
		response := s.transformProjectToResponse(&project)
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *ProjectService) GetMyCreatedProjects(userID uint) ([]interface{}, error) {
	var projects []model.Project

	err := s.DB.Preload("Creator").
		Preload("RequiredSkills.Skill").
		Preload("Conditions").
		Preload("Roles.RequiredSkills.Skill").
		Preload("Members.User").
		Preload("Members.ProjectRole.RequiredSkills.Skill").
		Preload("Members.MemberSkills.Skill").
		Preload("Tags.Tag").
		Preload("Benefits.Benefit").
		Preload("Timeline.Timeline").
		Where("creator_id = ?", userID).
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created projects: %w", err)
	}

	var responses []interface{}
	for _, project := range projects {
		response := s.transformProjectToResponse(&project)
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *ProjectService) GetMyMemberProjects(userID uint) ([]interface{}, error) {
	var projects []model.Project

	err := s.DB.Preload("Creator").
		Preload("RequiredSkills.Skill").
		Preload("Conditions").
		Preload("Roles.RequiredSkills.Skill").
		Preload("Members.User").
		Preload("Members.ProjectRole.RequiredSkills.Skill").
		Preload("Members.MemberSkills.Skill").
		Preload("Tags.Tag").
		Preload("Benefits.Benefit").
		Preload("Timeline.Timeline").
		Where("id IN (SELECT project_id FROM project_members WHERE user_id = ?)", userID).
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve member projects: %w", err)
	}

	var responses []interface{}
	for _, project := range projects {
		response := s.transformProjectToResponse(&project)
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *ProjectService) GetUserProject(userID, projectID uint) (interface{}, error) {
	var project model.Project

	err := s.DB.Where("id = ?", projectID).First(&project).Error
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	var hasAccess bool = false

	if project.CreatorID == userID {
		hasAccess = true
	} else {
		var memberCount int64
		s.DB.Model(&model.ProjectMember{}).
			Where("project_id = ? AND user_id = ?", projectID, userID).
			Count(&memberCount)
		if memberCount > 0 {
			hasAccess = true
		}
	}

	if !hasAccess {
		return nil, fmt.Errorf("project not found or access denied")
	}

	projectResult, err := s.loadProjectWithRelationships(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	return s.transformProjectToResponse(projectResult), nil
}

func (s *ProjectService) GetAllProjects() ([]interface{}, error) {
	var projects []model.Project

	err := s.DB.Preload("Creator").
		Preload("RequiredSkills.Skill").
		Preload("Conditions").
		Preload("Roles.RequiredSkills.Skill").
		Preload("Members.User").
		Preload("Members.ProjectRole.RequiredSkills.Skill").
		Preload("Members.MemberSkills.Skill").
		Preload("Tags.Tag").
		Preload("Benefits.Benefit").
		Preload("Timeline.Timeline").
		Where("status != ?", "draft").
		Find(&projects).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all projects: %w", err)
	}

	var responses []interface{}
	for _, project := range projects {
		response := s.transformProjectToResponse(&project)
		responses = append(responses, response)
	}

	return responses, nil
}

// GetProjectByID retrieves a single project by ID without authentication (public access)
func (s *ProjectService) GetProjectByID(projectID uint) (interface{}, error) {
	var project model.Project

	err := s.DB.Where("id = ? AND status != ?", projectID, "draft").First(&project).Error
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	projectResult, err := s.loadProjectWithRelationships(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	return s.transformProjectToResponse(projectResult), nil
}

// DeleteProject deletes a project by ID, only if the user is the creator
func (s *ProjectService) DeleteProject(projectID, userID uint) error {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First check if project exists and user is authorized
	var project model.Project
	if err := tx.First(&project, projectID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return fmt.Errorf("failed to find project: %w", err)
	}

	// Check if user is the creator
	if project.CreatorID != userID {
		tx.Rollback()
		return errors.New("you are not authorized to delete this project")
	}

	// Delete related records first (due to foreign key constraints)
	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectBenefit{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project benefits: %w", err)
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectTimeline{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project timeline: %w", err)
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectRequiredSkill{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project required skills: %w", err)
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectCondition{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project conditions: %w", err)
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectRole{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project roles: %w", err)
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectMember{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project members: %w", err)
	}

	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectTag{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project tags: %w", err)
	}

	// Finally delete the project itself
	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ProfileService) UpdateCollaborationStatus(userId uint, status string) (*model.Users, error) {
	var user model.Users
	if err := s.DB.First(&user, userId).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	user.StatusCollaboration = status
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update collaboration status: %w", err)
	}

	user.Password = ""
	return &user, nil
}
