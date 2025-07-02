package model

import "time"

type Siswa struct {
	ID              uint       `gorm:"primaryKey"`
	UserID          uint       `gorm:"not null;unique"`
	NISN            string     `gorm:"type:varchar(20);not null;unique"`
	KelasID         uint       `gorm:"not null"`
	NamaLengkap     string     `gorm:"type:varchar(100);not null"`
	JenisKelamin    string     `gorm:"type:enum('L', 'P');not null"`
	TempatLahir     string     `gorm:"type:varchar(50)"`
	TanggalLahir    *time.Time `gorm:"type:date"`
	Alamat          string     `gorm:"type:text"`
	NamaOrangtua    string     `gorm:"type:varchar(100)"`
	TeleponOrangtua string     `gorm:"type:varchar(20)"`
	TahunMasuk      int        `gorm:"type:year"`
	Status          string     `gorm:"type:enum('aktif', 'pindah', 'lulus', 'keluar');default:'aktif'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	User            User  `gorm:"foreignKey:UserID"`
	Kelas           Kelas `gorm:"foreignKey:KelasID"`
}
