package model

import (
	"time"
)

type Users struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"type:varchar(100);unique;not null"`
	Password    string `gorm:"type:varchar(255);not null"`
	RoleID      uint   `gorm:"not null"`
	NamaLengkap string `gorm:"type:varchar(100);not null"`
	Status      string `gorm:"type:enum('aktif','nonaktif');default:'aktif'"`
	LastLogin   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Role        Role `gorm:"foreignKey:RoleID"`
}

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	NamaRole  string `gorm:"type:varchar(50);unique;not null"`
	Deskripsi string `gorm:"type:text"`
	CreatedAt time.Time
}
