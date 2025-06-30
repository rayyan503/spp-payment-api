package model

import "time"

type Pengaturan struct {
	ID           uint    `gorm:"primaryKey"`
	KeySetting   string  `gorm:"type:varchar(100);not null;unique"`
	ValueSetting string  `gorm:"type:text;not null"`
	Deskripsi    *string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
