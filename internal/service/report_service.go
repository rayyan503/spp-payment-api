package service

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
)

type ReportService interface {
	GetLaporanSiswa(tahunAjaran, nisn string) ([]model.LaporanSiswa, error)
	GetLaporanKelas(tahunAjaran, namaBulan string) ([]model.LaporanKelas, error)
	GetLaporanKeseluruhan(tahunAjaran string) ([]model.LaporanKeseluruhan, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo}
}

func (s *reportService) GetLaporanSiswa(tahunAjaran, nisn string) ([]model.LaporanSiswa, error) {
	return s.repo.GetLaporanSiswa(tahunAjaran, nisn)
}

func (s *reportService) GetLaporanKelas(tahunAjaran, namaBulan string) ([]model.LaporanKelas, error) {
	return s.repo.GetLaporanKelas(tahunAjaran, namaBulan)
}

func (s *reportService) GetLaporanKeseluruhan(tahunAjaran string) ([]model.LaporanKeseluruhan, error) {
	return s.repo.GetLaporanKeseluruhan(tahunAjaran)
}
