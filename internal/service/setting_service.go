package service

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"gorm.io/gorm"
)

type SettingService interface {
	FindAllSettings() ([]model.Pengaturan, error)
	UpdateSettings(input map[string]string) error
}

type settingService struct {
	repo repository.SettingRepository
	db   *gorm.DB
}

func NewSettingService(repo repository.SettingRepository, db *gorm.DB) SettingService {
	return &settingService{repo: repo, db: db}
}

func (s *settingService) FindAllSettings() ([]model.Pengaturan, error) {
	return s.repo.FindAll()
}

func (s *settingService) UpdateSettings(input map[string]string) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for key, value := range input {
		repoTx := repository.NewSettingRepository(tx)
		err := repoTx.Update(key, value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
