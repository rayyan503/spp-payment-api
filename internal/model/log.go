package model

import "time"

type LogAktivitas struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    *uint   `gorm:"null"`
	Aktivitas string  `gorm:"type:varchar(255);not null"`
	Detail    *string `gorm:"type:text"`
	IPAddress *string `gorm:"type:varchar(45)"`
	UserAgent *string `gorm:"type:text"`
	CreatedAt time.Time
	User      Users `gorm:"foreignKey:UserID"`
}
