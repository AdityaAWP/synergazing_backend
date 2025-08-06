package model

import "time"

type Project struct {
	ID        uint  `json:"id" gorm:"primaryKey"`
	CreatorID uint  `json:"creator_id" gorm:"not null"`
	Creator   Users `json:"creator" gorm:"foreignKey:CreatorID"`

	Title       string `json:"title" gorm:"not null"`
	ProjectType string `json:"project_type" gorm:"not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	PictureURL  string `json:"picture_url" gorm:"type:text"`

	Duration             string    `json:"duration"`
	TotalTeam            int       `json:"total_team"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	Location             string    `json:"location"`
	WorkerType           string    `json:"worker_type"`
	Budget               float64   `json:"budget"`
	RegistrationDeadline time.Time `json:"registration_deadline"`

	TimeCommitment string `json:"time_commitment"`

	Benefits string `json:"benefits" gorm:"type:text;not null"`
	Timeline string `json:"timeline" gorm:"type:text"`

	RequiredSkills []*Skill            `json:"required_skills" gorm:"many2many:project_required_skills;"`
	Conditions     []*ProjectCondition `json:"conditions" gorm:"foreignKey:ProjectID"`
	Roles          []*ProjectRole      `json:"roles" gorm:"foreignKey:ProjectID"`
	Members        []*ProjectMember    `json:"members" gorm:"foreignKey:ProjectID"`
	Tags           []*Tag              `json:"tags" gorm:"many2many:project_tags;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Project) TableName() string {
	return "projects"
}
