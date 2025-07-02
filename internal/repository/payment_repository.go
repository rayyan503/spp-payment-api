package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Pembayaran) error
	Delete(id uint) error
	FindAllBySiswaID(siswaID uint) ([]model.Pembayaran, error)
	FindByOrderID(orderID string) (*model.Pembayaran, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(payment *model.Pembayaran) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Pembayaran{}).Error
}

func (r *paymentRepository) FindAllBySiswaID(siswaID uint) ([]model.Pembayaran, error) {
	var payments []model.Pembayaran
	err := r.db.Where("siswa_id = ?", siswaID).
		Preload("TagihanSPP.PeriodeSPP").
		Order("id desc").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) FindByOrderID(orderID string) (*model.Pembayaran, error) {
	var payment model.Pembayaran
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}
