package service

import (
	"gorm.io/gorm"
	"synergazing.com/synergazing/config"
	"synergazing.com/synergazing/model"
)

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
