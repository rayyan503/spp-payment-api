// File: internal/service/class_service.go (File Baru)

package service

import (
	"errors"

	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"

	"gorm.io/gorm"
)

type CreateClassLevelInput struct {
	Tingkat     int
	NamaTingkat string
	BiayaSPP    float64
}

type UpdateClassLevelInput struct {
	Tingkat     int
	NamaTingkat string
	BiayaSPP    float64
	Status      string
}

type ClassLevelService interface {
	CreateClassLevel(input CreateClassLevelInput) (*model.TingkatKelas, error)
	FindAllClassLevels() ([]model.TingkatKelas, error)
	FindClassLevelByID(id uint) (*model.TingkatKelas, error)
	UpdateClassLevel(id uint, input UpdateClassLevelInput) (*model.TingkatKelas, error)
	DeleteClassLevel(id uint) error
}

type classLevelService struct {
	repo repository.ClassLevelRepository
}

func NewClassLevelService(repo repository.ClassLevelRepository) ClassLevelService {
	return &classLevelService{repo}
}

func (s *classLevelService) CreateClassLevel(input CreateClassLevelInput) (*model.TingkatKelas, error) {
	_, err := s.repo.FindByTingkat(input.Tingkat)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return nil, errors.New("tingkat kelas sudah ada")
		}
		return nil, err
	}

	classLevel := &model.TingkatKelas{
		Tingkat:     input.Tingkat,
		NamaTingkat: input.NamaTingkat,
		BiayaSPP:    input.BiayaSPP,
		Status:      "aktif",
	}

	err = s.repo.Create(classLevel)
	if err != nil {
		return nil, err
	}

	return classLevel, nil
}

func (s *classLevelService) FindAllClassLevels() ([]model.TingkatKelas, error) {
	classLevels, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return classLevels, nil
}

func (s *classLevelService) FindClassLevelByID(id uint) (*model.TingkatKelas, error) {
	classLevel, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return classLevel, nil
}

func (s *classLevelService) UpdateClassLevel(id uint, input UpdateClassLevelInput) (*model.TingkatKelas, error) {
	classLevel, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Tingkat != classLevel.Tingkat {
		existing, err := s.repo.FindByTingkat(input.Tingkat)
		if err == nil && existing.ID != classLevel.ID {
			return nil, errors.New("tingkat kelas sudah ada")
		}
	}

	classLevel.Tingkat = input.Tingkat
	classLevel.NamaTingkat = input.NamaTingkat
	classLevel.BiayaSPP = input.BiayaSPP
	classLevel.Status = input.Status

	err = s.repo.Update(classLevel)
	if err != nil {
		return nil, err
	}

	return classLevel, nil
}

func (s *classLevelService) DeleteClassLevel(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
