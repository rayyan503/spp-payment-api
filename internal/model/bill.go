package model

import "time"

type TagihanSPP struct {
	ID                uint      `gorm:"primaryKey"`
	SiswaID           uint      `gorm:"not null"`
	PeriodeID         uint      `gorm:"not null"`
	JumlahTagihan     float64   `gorm:"type:decimal(12,2);not null"`
	StatusPembayaran  string    `gorm:"type:enum('belum_bayar', 'pending', 'lunas');default:'belum_bayar'"`
	TanggalJatuhTempo time.Time `gorm:"type:date;not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Siswa             Siswa      `gorm:"foreignKey:SiswaID"`
	PeriodeSPP        PeriodeSPP `gorm:"foreignKey:PeriodeID"`
}
