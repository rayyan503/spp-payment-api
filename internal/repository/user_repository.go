package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.Users, error)
	FindByID(id uint) (*model.Users, error)
	Create(user *model.Users) error
	FindAll(params utils.FindAllUsersParams) ([]model.Users, int64, error)
	Update(user *model.Users) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*model.Users, error) {
	var user model.Users
	err := r.db.Where("email = ?", email).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.Users, error) {
	var user model.Users
	err := r.db.Where("id = ?", id).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.Users) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll(params utils.FindAllUsersParams) ([]model.Users, int64, error) {
	var users []model.Users
	var total int64

	query := r.db.Model(&model.Users{})

	if params.RoleID != 0 {
		query = query.Where("role_id = ?", params.RoleID)
	}

	if params.Search != "" {
		query = query.Where("nama_lengkap LIKE ?", "%"+params.Search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err = query.Limit(params.Limit).Offset(offset).Preload("Role").Order("id asc").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *model.Users) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Users{}).Error
}
