package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/utils"
	"gorm.io/gorm"
)

type StudentRepository interface {
	Create(student *model.Siswa) error
	FindAll(params utils.FindAllStudentsParams) ([]model.Siswa, int64, error)
	FindByID(id uint) (*model.Siswa, error)
	FindByNISN(nisn string) (*model.Siswa, error)
	Update(student *model.Siswa) error
	Delete(id uint) error
	FindByUserID(userID uint) (*model.Siswa, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db}
}

func (r *studentRepository) Create(student *model.Siswa) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) FindAll(params utils.FindAllStudentsParams) ([]model.Siswa, int64, error) {
	var students []model.Siswa
	var total int64

	query := r.db

	if params.KelasID != 0 {
		query = query.Where("kelas_id = ?", params.KelasID)
	}
	if params.Search != "" {
		query = query.Where("nama_lengkap LIKE ? OR nisn LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err = query.Limit(params.Limit).Offset(offset).
		Preload("User").
		Preload("Kelas.TingkatKelas").
		Order("id desc").
		Find(&students).Error
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (r *studentRepository) FindByID(id uint) (*model.Siswa, error) {
	var student model.Siswa
	err := r.db.Preload("User").Preload("Kelas.TingkatKelas").Where("id = ?", id).First(&student).Error
	return &student, err
}

func (r *studentRepository) FindByNISN(nisn string) (*model.Siswa, error) {
	var student model.Siswa
	err := r.db.Where("nisn = ?", nisn).First(&student).Error
	return &student, err
}

func (r *studentRepository) Update(student *model.Siswa) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id uint) error {
	student, err := r.FindByID(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&model.Users{}, student.UserID).Error
}

func (r *studentRepository) FindByUserID(userID uint) (*model.Siswa, error) {
	var student model.Siswa
	err := r.db.Preload("User").Preload("Kelas.TingkatKelas").Where("user_id = ?", userID).First(&student).Error
	return &student, err
}
