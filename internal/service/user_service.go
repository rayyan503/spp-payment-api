package service

import (
	"errors"

	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/repository"
	"github.com/hiuncy/spp-payment-api/internal/utils"
	"gorm.io/gorm"
)

type CreateUserInput struct {
	NamaLengkap string
	Email       string
	Password    string
	RoleID      uint
}

type FindAllUsersInput struct {
	Page   int
	Limit  int
	RoleID uint
	Search string
}

type UpdateUserInput struct {
	NamaLengkap string
	Email       string
	RoleID      uint
}

type UserService interface {
	GetUserProfile(userID uint) (*model.User, error)
	CreateUser(input CreateUserInput) (*model.User, error)
	FindAllUsers(input FindAllUsersInput) ([]model.User, int64, error)
	UpdateUser(id uint, input UpdateUserInput) (*model.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) GetUserProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) CreateUser(input CreateUserInput) (*model.User, error) {
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		if err == nil {
			return nil, errors.New("email sudah terdaftar")
		}
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		NamaLengkap: input.NamaLengkap,
		Email:       input.Email,
		Password:    hashedPassword,
		RoleID:      input.RoleID,
		Status:      "aktif",
	}

	err = s.userRepo.Create(&newUser)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.FindByID(newUser.ID)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) FindAllUsers(input FindAllUsersInput) ([]model.User, int64, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	params := repository.FindAllUsersParams{
		Page:   input.Page,
		Limit:  input.Limit,
		RoleID: input.RoleID,
		Search: input.Search,
	}

	users, total, err := s.userRepo.FindAll(params)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *userService) UpdateUser(id uint, input UpdateUserInput) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(input.Email)
		if err == nil && existingUser.ID != user.ID {
			return nil, errors.New("email sudah terdaftar untuk pengguna lain")
		}
	}

	user.NamaLengkap = input.NamaLengkap
	user.Email = input.Email
	user.RoleID = input.RoleID

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
