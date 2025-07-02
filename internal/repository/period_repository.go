package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type PeriodRepository interface {
	Create(period *model.PeriodeSPP) error
	FindAll(tahunAjaran string) ([]model.PeriodeSPP, error)
	FindByID(id uint) (*model.PeriodeSPP, error)
	FindByTahunAjaranAndBulan(tahunAjaran string, bulan int) (*model.PeriodeSPP, error)
	Update(period *model.PeriodeSPP) error
	Delete(id uint) error
}

type periodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) PeriodRepository {
	return &periodRepository{db}
}

func (r *periodRepository) Create(period *model.PeriodeSPP) error {
	return r.db.Create(period).Error
}

func (r *periodRepository) FindAll(tahunAjaran string) ([]model.PeriodeSPP, error) {
	var periods []model.PeriodeSPP
	query := r.db

	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}

	err := query.Order("tahun_ajaran desc, bulan asc").Find(&periods).Error
	return periods, err
}

func (r *periodRepository) FindByID(id uint) (*model.PeriodeSPP, error) {
	var period model.PeriodeSPP
	err := r.db.Where("id = ?", id).First(&period).Error
	return &period, err
}

func (r *periodRepository) FindByTahunAjaranAndBulan(tahunAjaran string, bulan int) (*model.PeriodeSPP, error) {
	var period model.PeriodeSPP
	err := r.db.Where("tahun_ajaran = ? AND bulan = ?", tahunAjaran, bulan).First(&period).Error
	return &period, err
}

func (r *periodRepository) Update(period *model.PeriodeSPP) error {
	return r.db.Save(period).Error
}

func (r *periodRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.PeriodeSPP{}).Error
}
