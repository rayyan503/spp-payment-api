package service

import (
	"errors"
	"time"

	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
)

type CreatePeriodInput struct {
	TahunAjaran    string
	Bulan          int
	NamaBulan      string
	TanggalMulai   string
	TanggalSelesai string
}

type UpdatePeriodInput struct {
	TahunAjaran    string
	Bulan          int
	NamaBulan      string
	TanggalMulai   string
	TanggalSelesai string
	Status         string
}

type PeriodService interface {
	CreatePeriod(input CreatePeriodInput) (*model.PeriodeSPP, error)
	FindAllPeriods(tahunAjaran string) ([]model.PeriodeSPP, error)
	FindPeriodByID(id uint) (*model.PeriodeSPP, error)
	UpdatePeriod(id uint, input UpdatePeriodInput) (*model.PeriodeSPP, error)
	DeletePeriod(id uint) error
}

type periodService struct {
	repo repository.PeriodRepository
}

func NewPeriodService(repo repository.PeriodRepository) PeriodService {
	return &periodService{repo}
}

func (s *periodService) CreatePeriod(input CreatePeriodInput) (*model.PeriodeSPP, error) {
	_, err := s.repo.FindByTahunAjaranAndBulan(input.TahunAjaran, input.Bulan)
	if err == nil {
		return nil, errors.New("periode untuk tahun ajaran dan bulan tersebut sudah ada")
	}

	tglMulai, _ := time.Parse("2006-01-02", input.TanggalMulai)
	tglSelesai, _ := time.Parse("2006-01-02", input.TanggalSelesai)

	newPeriod := &model.PeriodeSPP{
		TahunAjaran:    input.TahunAjaran,
		Bulan:          input.Bulan,
		NamaBulan:      input.NamaBulan,
		TanggalMulai:   tglMulai,
		TanggalSelesai: tglSelesai,
		Status:         "belum_aktif",
	}

	if err := s.repo.Create(newPeriod); err != nil {
		return nil, err
	}
	return newPeriod, nil
}

func (s *periodService) FindAllPeriods(tahunAjaran string) ([]model.PeriodeSPP, error) {
	return s.repo.FindAll(tahunAjaran)
}

func (s *periodService) FindPeriodByID(id uint) (*model.PeriodeSPP, error) {
	return s.repo.FindByID(id)
}

func (s *periodService) UpdatePeriod(id uint, input UpdatePeriodInput) (*model.PeriodeSPP, error) {
	period, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	tglMulai, _ := time.Parse("2006-01-02", input.TanggalMulai)
	tglSelesai, _ := time.Parse("2006-01-02", input.TanggalSelesai)

	period.TahunAjaran = input.TahunAjaran
	period.Bulan = input.Bulan
	period.NamaBulan = input.NamaBulan
	period.TanggalMulai = tglMulai
	period.TanggalSelesai = tglSelesai
	period.Status = input.Status

	if err := s.repo.Update(period); err != nil {
		return nil, err
	}
	return period, nil
}

func (s *periodService) DeletePeriod(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
