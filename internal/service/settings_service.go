package service

import (
	"yoosoPosV3/internal/config"
	"yoosoPosV3/internal/model"
)

func GetSetting(key string) string {
	var setting model.Setting
	err := config.DB.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return ""
	}
	return setting.Value
}

func SetSetting(key, value string) error {
	var setting model.Setting
	result := config.DB.Where("key = ?", key).First(&setting)
	if result.Error != nil {
		return config.DB.Create(&model.Setting{Key: key, Value: value}).Error
	}

	if err := config.DB.Model(&setting).Update("value", value).Error; err != nil {
		return err
	}

	if key == "store_name" || key == "secondary_password" {
		go CheckVersion()
	}
	return nil
}

func GetStoreName() (string, error) {
	return GetSetting("store_name"), nil
}

func SetStoreName(name string) error {
	return SetSetting("store_name", name)
}
