package model

import "time"

type PeriodeSPP struct {
	ID             uint      `gorm:"primaryKey"`
	TahunAjaran    string    `gorm:"type:varchar(20);not null"`
	Bulan          int       `gorm:"not null"`
	NamaBulan      string    `gorm:"type:varchar(20);not null"`
	TanggalMulai   time.Time `gorm:"type:date;not null"`
	TanggalSelesai time.Time `gorm:"type:date;not null"`
	Status         string    `gorm:"type:enum('belum_aktif', 'aktif', 'selesai');default:'belum_aktif'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
