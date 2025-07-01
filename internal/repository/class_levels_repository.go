package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type ClassLevelRepository interface {
	Create(classLevel *model.TingkatKelas) error
	FindAll() ([]model.TingkatKelas, error)
	FindByID(id uint) (*model.TingkatKelas, error)
	FindByTingkat(tingkat int) (*model.TingkatKelas, error)
	Update(classLevel *model.TingkatKelas) error
	Delete(id uint) error
}

type classLevelRepository struct {
	db *gorm.DB
}

func NewClassLevelRepository(db *gorm.DB) ClassLevelRepository {
	return &classLevelRepository{db}
}

func (r *classLevelRepository) Create(classLevel *model.TingkatKelas) error {
	return r.db.Create(classLevel).Error
}

func (r *classLevelRepository) FindByTingkat(tingkat int) (*model.TingkatKelas, error) {
	var classLevel model.TingkatKelas
	err := r.db.Where("tingkat = ?", tingkat).First(&classLevel).Error
	return &classLevel, err
}

func (r *classLevelRepository) FindAll() ([]model.TingkatKelas, error) {
	var classLevels []model.TingkatKelas
	err := r.db.Order("tingkat asc").Find(&classLevels).Error
	return classLevels, err
}

func (r *classLevelRepository) FindByID(id uint) (*model.TingkatKelas, error) {
	var classLevel model.TingkatKelas
	err := r.db.Where("id = ?", id).First(&classLevel).Error
	return &classLevel, err
}

func (r *classLevelRepository) Update(classLevel *model.TingkatKelas) error {
	return r.db.Save(classLevel).Error
}

func (r *classLevelRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.TingkatKelas{}).Error
}
