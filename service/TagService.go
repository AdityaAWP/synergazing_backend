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

type BenefitService struct {
	DB *gorm.DB
}

func NewBenefitService(db *gorm.DB) *BenefitService {
	return &BenefitService{DB: db}
}

func (s *BenefitService) findOrCreate(tx *gorm.DB, names []string) ([]*model.Benefit, error) {
	var benefits []*model.Benefit

	for _, name := range names {
		var benefit model.Benefit
		if err := tx.Where("LOWER(name) = LOWER(?)", name).First(&benefit).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				benefit = model.Benefit{Name: name}
				if err := tx.Create(&benefit).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		benefits = append(benefits, &benefit)
	}
	return benefits, nil
}

type TimelineService struct {
	DB *gorm.DB
}

func NewTimelineService(db *gorm.DB) *TimelineService {
	return &TimelineService{DB: db}
}

func (s *TimelineService) findOrCreate(tx *gorm.DB, names []string) ([]*model.Timeline, error) {
	var timelines []*model.Timeline

	for _, name := range names {
		var timeline model.Timeline
		if err := tx.Where("LOWER(name) = LOWER(?)", name).First(&timeline).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				timeline = model.Timeline{Name: name}
				if err := tx.Create(&timeline).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		timelines = append(timelines, &timeline)
	}
	return timelines, nil
}
