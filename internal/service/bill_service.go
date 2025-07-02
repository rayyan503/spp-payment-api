package service

import (
	"github.com/hiuncy/spp-payment-api/internal/dto"
	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"github.com/hiuncy/spp-payment-api/internal/utils"
)

type BillService interface {
	GenerateBillsForPeriod(periodID uint) error
	FindAllBills(input dto.FindAllBillsInput) ([]model.TagihanSPP, int64, error)
	FindBillByID(id uint) (*model.TagihanSPP, error)
	UpdateBill(id uint, input dto.UpdateBillInput) (*model.TagihanSPP, error)
	DeleteBill(id uint) error
}

type billService struct {
	repo repository.BillRepository
}

func NewBillService(repo repository.BillRepository) BillService {
	return &billService{repo}
}

func (s *billService) GenerateBillsForPeriod(periodID uint) error {
	return s.repo.GenerateBills(periodID)
}

func (s *billService) FindAllBills(input dto.FindAllBillsInput) ([]model.TagihanSPP, int64, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	params := utils.FindAllBillsParams{
		Page:             input.Page,
		Limit:            input.Limit,
		PeriodeID:        input.PeriodeID,
		SiswaID:          input.SiswaID,
		StatusPembayaran: input.StatusPembayaran,
	}
	return s.repo.FindAll(params)
}

func (s *billService) FindBillByID(id uint) (*model.TagihanSPP, error) {
	return s.repo.FindByID(id)
}

func (s *billService) UpdateBill(id uint, input dto.UpdateBillInput) (*model.TagihanSPP, error) {
	bill, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	bill.JumlahTagihan = input.JumlahTagihan
	bill.StatusPembayaran = input.StatusPembayaran
	if err := s.repo.Update(bill); err != nil {
		return nil, err
	}
	return bill, nil
}

func (s *billService) DeleteBill(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
