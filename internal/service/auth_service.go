package service

import (
	"errors"

	"github.com/hiuncy/spp-payment-api/internal/repository"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo     repository.UserRepository
	jwtSecretKey string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecretKey string) AuthService {
	return &authService{userRepo, jwtSecretKey}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("email atau password salah")
		}
		return "", err
	}

	err = utils.CheckPasswordHash(password, user.Password)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	token, err := utils.GenerateJWT(user.ID, user.Role.NamaRole, s.jwtSecretKey)
	if err != nil {
		return "", errors.New("gagal membuat token otentikasi")
	}

	return token, nil
}
