package service

import (
	"errors"
	"time"

	"github.com/hiuncy/spp-payment-api/internal/dto"
	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"gorm.io/gorm"
)

type StudentService interface {
	CreateStudent(input dto.CreateStudentInput) (*model.Siswa, error)
	FindAllStudents(input dto.FindAllStudentsInput) ([]model.Siswa, int64, error)
	FindStudentByID(id uint) (*model.Siswa, error)
	UpdateStudent(id uint, input dto.UpdateStudentInput) (*model.Siswa, error)
	DeleteStudent(id uint) error
	GetStudentProfile(userID uint) (*model.Siswa, error)
}

type studentService struct {
	studentRepo repository.StudentRepository
	userRepo    repository.UserRepository
	db          *gorm.DB
}

func NewStudentService(studentRepo repository.StudentRepository, userRepo repository.UserRepository, db *gorm.DB) StudentService {
	return &studentService{studentRepo, userRepo, db}
}

func (s *studentService) CreateStudent(input dto.CreateStudentInput) (*model.Siswa, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userRepoTx := repository.NewUserRepository(tx)
	studentRepoTx := repository.NewStudentRepository(tx)

	_, err := userRepoTx.FindByEmail(input.Email)
	if err == nil {
		tx.Rollback()
		return nil, errors.New("email sudah terdaftar")
	}
	_, err = studentRepoTx.FindByNISN(input.NISN)
	if err == nil {
		tx.Rollback()
		return nil, errors.New("NISN sudah terdaftar")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newUser := &model.User{
		Email:       input.Email,
		Password:    hashedPassword,
		RoleID:      3,
		NamaLengkap: input.NamaLengkap,
		Status:      "aktif",
	}
	if err := userRepoTx.Create(newUser); err != nil {
		tx.Rollback()
		return nil, err
	}

	tglLahir, _ := time.Parse("2006-01-02", input.TanggalLahir)
	newStudent := &model.Siswa{
		UserID:          newUser.ID,
		NISN:            input.NISN,
		KelasID:         input.KelasID,
		NamaLengkap:     input.NamaLengkap,
		JenisKelamin:    input.JenisKelamin,
		TempatLahir:     input.TempatLahir,
		TanggalLahir:    &tglLahir,
		Alamat:          input.Alamat,
		NamaOrangtua:    input.NamaOrangTua,
		TeleponOrangtua: input.TeleponOrangTua,
		TahunMasuk:      input.TahunMasuk,
		Status:          "aktif",
	}
	if err := studentRepoTx.Create(newStudent); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.studentRepo.FindByID(newStudent.ID)
}

func (s *studentService) FindAllStudents(input dto.FindAllStudentsInput) ([]model.Siswa, int64, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	params := utils.FindAllStudentsParams{
		Page:    input.Page,
		Limit:   input.Limit,
		KelasID: input.KelasID,
		Search:  input.Search,
	}
	return s.studentRepo.FindAll(params)
}

func (s *studentService) FindStudentByID(id uint) (*model.Siswa, error) {
	return s.studentRepo.FindByID(id)
}

func (s *studentService) UpdateStudent(id uint, input dto.UpdateStudentInput) (*model.Siswa, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.NISN != student.NISN {
		existingStudent, err := s.studentRepo.FindByNISN(input.NISN)
		if err == nil && existingStudent.ID != student.ID {
			return nil, errors.New("NISN sudah terdaftar untuk siswa lain")
		}
	}

	if input.EmailUser != student.User.Email {
		existingUser, err := s.userRepo.FindByEmail(input.EmailUser)
		if err == nil && existingUser.ID != student.UserID {
			return nil, errors.New("email sudah terdaftar untuk pengguna lain")
		}
	}

	tglLahir, _ := time.Parse("2006-01-02", input.TanggalLahir)
	student.NISN = input.NISN
	student.KelasID = input.KelasID
	student.NamaLengkap = input.NamaLengkap
	student.JenisKelamin = input.JenisKelamin
	student.TempatLahir = input.TempatLahir
	student.TanggalLahir = &tglLahir
	student.Alamat = input.Alamat
	student.NamaOrangtua = input.NamaOrangTua
	student.TeleponOrangtua = input.TeleponOrangTua
	student.TahunMasuk = input.TahunMasuk
	student.Status = input.Status

	student.User.Email = input.EmailUser
	student.User.Status = input.StatusUser
	student.User.NamaLengkap = input.NamaLengkap

	err = s.studentRepo.Update(student)
	if err != nil {
		return nil, err
	}

	return s.studentRepo.FindByID(id)
}

func (s *studentService) DeleteStudent(id uint) error {
	return s.studentRepo.Delete(id)
}

func (s *studentService) GetStudentProfile(userID uint) (*model.Siswa, error) {
	return s.studentRepo.FindByUserID(userID)
}
