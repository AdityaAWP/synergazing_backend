package service

import (
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
