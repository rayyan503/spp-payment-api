package service

import (
	"errors"

	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"gorm.io/gorm"
)

type CreateClassInput struct {
	TingkatID uint
	NamaKelas string
	WaliKelas string
	Kapasitas int
}

type UpdateClassInput struct {
	TingkatID uint
	NamaKelas string
	WaliKelas string
	Kapasitas int
	Status    string
}

type ClassService interface {
	CreateClass(input CreateClassInput) (*model.Kelas, error)
	FindAllClasses() ([]model.Kelas, error)
	FindClassByID(id uint) (*model.Kelas, error)
	UpdateClass(id uint, input UpdateClassInput) (*model.Kelas, error)
	DeleteClass(id uint) error
}

type classService struct {
	repo repository.ClassRepository
}

func NewClassService(repo repository.ClassRepository) ClassService {
	return &classService{repo}
}

func (s *classService) CreateClass(input CreateClassInput) (*model.Kelas, error) {
	_, err := s.repo.FindByName(input.NamaKelas)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return nil, errors.New("nama kelas sudah ada")
		}
		return nil, err
	}

	newClass := &model.Kelas{
		TingkatID: input.TingkatID,
		NamaKelas: input.NamaKelas,
		WaliKelas: input.WaliKelas,
		Kapasitas: input.Kapasitas,
		Status:    "aktif",
	}

	err = s.repo.Create(newClass)
	if err != nil {
		return nil, err
	}

	// Ambil kembali data lengkap dengan preload
	return s.repo.FindByID(newClass.ID)
}

func (s *classService) FindAllClasses() ([]model.Kelas, error) {
	return s.repo.FindAll()
}

func (s *classService) FindClassByID(id uint) (*model.Kelas, error) {
	return s.repo.FindByID(id)
}

func (s *classService) UpdateClass(id uint, input UpdateClassInput) (*model.Kelas, error) {
	class, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.NamaKelas != class.NamaKelas {
		existing, err := s.repo.FindByName(input.NamaKelas)
		if err == nil && existing.ID != class.ID {
			return nil, errors.New("nama kelas sudah ada")
		}
	}

	class.TingkatID = input.TingkatID
	class.NamaKelas = input.NamaKelas
	class.WaliKelas = input.WaliKelas
	class.Kapasitas = input.Kapasitas
	class.Status = input.Status

	err = s.repo.Update(class)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *classService) DeleteClass(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
