package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

type ProjectService struct {
	DB           *gorm.DB
	skillService *SkillService
	tagService   *TagService
}

type RoleDTO struct {
	Name           string   `json:"name"`
	SlotsAvailable int      `json:"slots_available"`
	Description    string   `json:"description"`
	SkillNames     []string `json:"skill_names"`
}

type MemberDTO struct {
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
}

func NewProjectService(db *gorm.DB, skillService *SkillService, tagService *TagService) *ProjectService {
	return &ProjectService{
		DB:           db,
		skillService: skillService,
		tagService:   tagService,
	}
}

func (s *ProjectService) getProjectForUpdate(tx *gorm.DB, projectID, userID uint, requiredStage int) (*model.Project, error) {
	var project model.Project

	if err := tx.First(&project, projectID).Error; err != nil {
		return nil, errors.New("project not found")
	}
	if project.CreatorID != userID {
		return nil, errors.New("you are not authorized for this project")
	}
	if project.CompletionStage < requiredStage {
		return nil, fmt.Errorf("you must complete the previous stage", requiredStage)
	}
	return &project, nil
}

// loadProjectWithRelationships loads a project with all its relationships
func (s *ProjectService) loadProjectWithRelationships(projectID uint) (*model.Project, error) {
	var project model.Project
	if err := s.DB.Preload("Creator").
		Preload("RequiredSkills.Skill").
		Preload("Conditions").
		Preload("Roles.RequiredSkills.Skill").
		Preload("Members.User").
		Preload("Members.ProjectRole").
		Preload("Tags").
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
	project.WorkerType = details.WorkerType
	project.Budget = details.Budget
	project.RegistrationDeadline = details.RegistrationDeadline
	project.CompletionStage = 2

	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return project, tx.Commit().Error
}

func (s *ProjectService) UpdateStage3(projectID, userID uint, timeCommitment string, skillNames []string, conditionDescriptions []string) (*model.Project, error) {
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

	// Delete existing required skills for this project
	if err := tx.Where("project_id = ?", projectID).Delete(&model.ProjectRequiredSkill{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new required skills
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

	// Commit the transaction first
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Now preload the relationships and return the updated project
	return s.loadProjectWithRelationships(project.ID)
}

func (s *ProjectService) UpdateStage4(projectID, userID uint, roles []RoleDTO, members []MemberDTO) (*model.Project, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Where("project_id = ?", projectID).Delete(&model.ProjectMember{})
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
		if err := tx.Where("email = ?", memberData.Email).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("user to invite not found: " + memberData.Email)
		}
		roleID, ok := roleMap[memberData.RoleName]
		if !ok {
			tx.Rollback()
			return nil, errors.New("role specified for member does not exist: " + memberData.RoleName)
		}
		member := model.ProjectMember{
			ProjectID:     project.ID,
			UserID:        user.ID,
			ProjectRoleID: roleID,
			Status:        "invited",
		}
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	project.CompletionStage = 4
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit the transaction first
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Now preload the relationships and return the updated project
	return s.loadProjectWithRelationships(project.ID)
}

func (s *ProjectService) UpdateStage5(projectID, userID uint, benefits, timeline string, tagNames []string) (*model.Project, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 4)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if benefits == "" {
		tx.Rollback()
		return nil, errors.New("benefits field is required")
	}

	project.Benefits = benefits
	project.Timeline = timeline
	project.CompletionStage = 5
	project.Status = "published"

	if len(tagNames) > 0 {
		tags, err := s.tagService.findOrCreate(tx, tagNames)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(project).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit the transaction first
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Now preload the relationships and return the updated project
	return s.loadProjectWithRelationships(project.ID)
}

func (s *ProjectService) CreateRolesOnly(projectID, userID uint, roles []RoleDTO) (*model.Project, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Clear existing roles
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

	// Commit the transaction first
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Now preload the relationships and return the updated project
	return s.loadProjectWithRelationships(project.ID)
}

func (s *ProjectService) AddMembersOnly(projectID, userID uint, members []MemberDTO) (*model.Project, error) {
	tx := s.DB.Begin()
	project, err := s.getProjectForUpdate(tx, projectID, userID, 3)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Get existing roles for this project
	var existingRoles []model.ProjectRole
	if err := tx.Where("project_id = ?", projectID).Find(&existingRoles).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("failed to fetch project roles")
	}

	// Build role map
	roleMap := make(map[string]uint)
	for _, role := range existingRoles {
		roleMap[role.Name] = role.ID
	}

	// Clear existing members
	tx.Where("project_id = ?", projectID).Delete(&model.ProjectMember{})

	for _, memberData := range members {
		var user model.Users
		if err := tx.Where("email = ?", memberData.Email).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("user to invite not found: " + memberData.Email)
		}

		roleID, ok := roleMap[memberData.RoleName]
		if !ok {
			tx.Rollback()
			return nil, errors.New("role specified for member does not exist: " + memberData.RoleName)
		}

		member := model.ProjectMember{
			ProjectID:     project.ID,
			UserID:        user.ID,
			ProjectRoleID: roleID,
			Status:        "invited",
		}
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	project.CompletionStage = 4
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit the transaction first
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Now preload the relationships and return the updated project
	return s.loadProjectWithRelationships(project.ID)
}
