package repository

import (
	"github.com/hiuncy/spp-payment-api/internal/model"

	"gorm.io/gorm"
)

type FindAllUsersParams struct {
	Limit  int
	Page   int
	RoleID uint
	Search string
}

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
	FindAll(params FindAllUsersParams) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll(params FindAllUsersParams) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

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
	err = query.Limit(params.Limit).Offset(offset).Preload("Role").Order("id desc").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}
