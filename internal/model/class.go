package model

import "time"

type TingkatKelas struct {
	ID          uint    `gorm:"primaryKey"`
	Tingkat     int     `gorm:"type:int;not null;unique"`
	NamaTingkat string  `gorm:"type:varchar(50);not null"`
	BiayaSPP    float64 `gorm:"type:decimal(12,2);not null"`
	Status      string  `gorm:"type:enum('aktif', 'nonaktif');default:'aktif'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Kelas struct {
	ID           uint   `gorm:"primaryKey"`
	TingkatID    uint   `gorm:"not null"`
	NamaKelas    string `gorm:"type:varchar(10);not null;unique"`
	WaliKelas    string `gorm:"type:varchar(100)"`
	Kapasitas    int    `gorm:"default:30"`
	Status       string `gorm:"type:enum('aktif', 'nonaktif');default:'aktif'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TingkatKelas TingkatKelas `gorm:"foreignKey:TingkatID"`
}
