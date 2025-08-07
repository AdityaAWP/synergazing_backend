package model

import (
	"encoding/json"
	"time"

	"synergazing.com/synergazing/helper"
)

type Project struct {
	ID        uint  `json:"id" gorm:"primaryKey"`
	CreatorID uint  `json:"creator_id" gorm:"not null"`
	Creator   Users `json:"creator" gorm:"foreignKey:CreatorID"`

	Status          string `json:"status" gorm:"not null;default:'draft'"`
	CompletionStage int    `json:"completion_stage" gorm:"not null;default:0"`

	Title       string `json:"title" gorm:"not null"`
	ProjectType string `json:"project_type" gorm:"not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	PictureURL  string `json:"picture_url" gorm:"type:text"`

	Duration             string    `json:"duration"`
	TotalTeam            int       `json:"total_team"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Location             string    `json:"location"`
	Budget               float64   `json:"budget"`
	RegistrationDeadline time.Time `json:"registration_deadline"`

	TimeCommitment string `json:"time_commitment"`

	Benefits []*ProjectBenefit  `json:"benefits" gorm:"foreignKey:ProjectID"`
	Timeline []*ProjectTimeline `json:"timeline" gorm:"foreignKey:ProjectID"`

	RequiredSkills []*ProjectRequiredSkill `json:"required_skills" gorm:"foreignKey:ProjectID"`
	Conditions     []*ProjectCondition     `json:"conditions" gorm:"foreignKey:ProjectID"`
	Roles          []*ProjectRole          `json:"roles" gorm:"foreignKey:ProjectID"`
	Members        []*ProjectMember        `json:"members" gorm:"foreignKey:ProjectID"`
	Tags           []*ProjectTag           `json:"tags" gorm:"foreignKey:ProjectID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Project) TableName() string {
	return "projects"
}

func (p Project) MarshalJSON() ([]byte, error) {
	type Alias Project
	return json.Marshal(&struct {
		PictureURL string `json:"picture_url"`
		*Alias
	}{
		PictureURL: helper.GetUrlFile(p.PictureURL),
		Alias:      (*Alias)(&p),
	})
}
