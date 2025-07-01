package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(class *model.Kelas) error
	FindAll() ([]model.Kelas, error)
	FindByID(id uint) (*model.Kelas, error)
	FindByName(name string) (*model.Kelas, error)
	Update(class *model.Kelas) error
	Delete(id uint) error
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db}
}

func (r *classRepository) Create(class *model.Kelas) error {
	return r.db.Create(class).Error
}

func (r *classRepository) FindAll() ([]model.Kelas, error) {
	var classes []model.Kelas
	err := r.db.Preload("TingkatKelas").Order("nama_kelas asc").Find(&classes).Error
	return classes, err
}

func (r *classRepository) FindByID(id uint) (*model.Kelas, error) {
	var class model.Kelas
	err := r.db.Preload("TingkatKelas").Where("id = ?", id).First(&class).Error
	return &class, err
}

func (r *classRepository) FindByName(name string) (*model.Kelas, error) {
	var class model.Kelas
	err := r.db.Where("nama_kelas = ?", name).First(&class).Error
	return &class, err
}

func (r *classRepository) Update(class *model.Kelas) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Kelas{}).Error
}
