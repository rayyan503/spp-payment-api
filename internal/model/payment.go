package model

import "time"

type Pembayaran struct {
	ID                uint    `gorm:"primaryKey"`
	TagihanID         uint    `gorm:"not null"`
	SiswaID           uint    `gorm:"not null"`
	OrderID           string  `gorm:"type:varchar(100);not null;unique"`
	TransactionID     *string `gorm:"type:varchar(100)"`
	JumlahBayar       float64 `gorm:"type:decimal(12,2);not null"`
	MetodePembayaran  *string `gorm:"type:varchar(50)"`
	StatusPembayaran  string  `gorm:"type:enum('pending', 'settlement', 'cancel', 'expire', 'failure');default:'pending'"`
	TanggalPembayaran *time.Time
	TanggalSettlement *time.Time
	MidtransResponse  *string `gorm:"type:json"`
	Keterangan        *string `gorm:"type:text"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	TagihanSPP        TagihanSPP `gorm:"foreignKey:TagihanID"`
	Siswa             Siswa      `gorm:"foreignKey:SiswaID"`
}
