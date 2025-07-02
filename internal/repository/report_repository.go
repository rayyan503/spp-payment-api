package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetLaporanSiswa(tahunAjaran, nisn string) ([]model.LaporanSiswa, error)
	GetLaporanKelas(tahunAjaran, namaBulan string) ([]model.LaporanKelas, error)
	GetLaporanKeseluruhan(tahunAjaran string) ([]model.LaporanKeseluruhan, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetLaporanSiswa(tahunAjaran, nisn string) ([]model.LaporanSiswa, error) {
	var results []model.LaporanSiswa
	query := r.db.Table("v_laporan_siswa")
	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}
	if nisn != "" {
		query = query.Where("nisn = ?", nisn)
	}
	err := query.Find(&results).Error
	return results, err
}

func (r *reportRepository) GetLaporanKelas(tahunAjaran, namaBulan string) ([]model.LaporanKelas, error) {
	var results []model.LaporanKelas
	query := r.db.Table("v_laporan_kelas")
	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}
	if namaBulan != "" {
		query = query.Where("nama_bulan = ?", namaBulan)
	}
	err := query.Find(&results).Error
	return results, err
}

func (r *reportRepository) GetLaporanKeseluruhan(tahunAjaran string) ([]model.LaporanKeseluruhan, error) {
	var results []model.LaporanKeseluruhan
	query := r.db.Table("v_laporan_keseluruhan")
	if tahunAjaran != "" {
		query = query.Where("tahun_ajaran = ?", tahunAjaran)
	}
	err := query.Find(&results).Error
	return results, err
}
