package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"

	"gorm.io/gorm"
)

type SettingRepository interface {
	FindAll() ([]model.Pengaturan, error)
	Update(key string, value string) error
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db}
}

func (r *settingRepository) FindAll() ([]model.Pengaturan, error) {
	var settings []model.Pengaturan
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *settingRepository) Update(key string, value string) error {
	return r.db.Model(&model.Pengaturan{}).Where("key_setting = ?", key).Update("value_setting", value).Error
}
