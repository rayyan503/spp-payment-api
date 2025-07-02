package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type FindAllBillsParams struct {
	Limit            int
	Page             int
	PeriodeID        uint
	SiswaID          uint
	StatusPembayaran string
}

type BillRepository interface {
	GenerateBills(periodID uint) error
	FindAll(params FindAllBillsParams) ([]model.TagihanSPP, int64, error)
	FindByID(id uint) (*model.TagihanSPP, error)
	Update(bill *model.TagihanSPP) error
	Delete(id uint) error
}

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db}
}

func (r *billRepository) GenerateBills(periodID uint) error {
	return r.db.Exec("CALL GenerateTagihanSPP(?)", periodID).Error
}

func (r *billRepository) FindAll(params FindAllBillsParams) ([]model.TagihanSPP, int64, error) {
	var bills []model.TagihanSPP
	var total int64

	query := r.db.Model(&model.TagihanSPP{})

	if params.PeriodeID != 0 {
		query = query.Where("periode_id = ?", params.PeriodeID)
	}
	if params.SiswaID != 0 {
		query = query.Where("siswa_id = ?", params.SiswaID)
	}
	if params.StatusPembayaran != "" {
		query = query.Where("status_pembayaran = ?", params.StatusPembayaran)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := query.Limit(params.Limit).Offset(offset).
		Preload("Siswa").
		Preload("PeriodeSPP").
		Order("id desc").
		Find(&bills).Error

	return bills, total, err
}

func (r *billRepository) FindByID(id uint) (*model.TagihanSPP, error) {
	var bill model.TagihanSPP
	err := r.db.Preload("Siswa").Preload("PeriodeSPP").Where("id = ?", id).First(&bill).Error
	return &bill, err
}

func (r *billRepository) Update(bill *model.TagihanSPP) error {
	return r.db.Save(bill).Error
}

func (r *billRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.TagihanSPP{}).Error
}
