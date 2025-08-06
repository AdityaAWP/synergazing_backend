package service

import (
	"gorm.io/gorm"
	"synergazing.com/synergazing/model"
)

type TagService struct {
	DB *gorm.DB
}

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{DB: db}
}

func (s *TagService) findOrCreate(tx *gorm.DB, names []string) ([]*model.Tag, error) {
	var tags []*model.Tag

	for _, name := range names {
		var tag model.Tag
		if err := tx.Where("LOWER(name) = LOWER(?)", name).First(&tag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				tag = model.Tag{Name: name}
				if err := tx.Create(&tag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}
