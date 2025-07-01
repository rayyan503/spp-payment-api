package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	RoleID      uint   `json:"role_id" binding:"required"`
}

type UpdateUserRequest struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	RoleID      uint   `json:"role_id" binding:"required"`
}

type CreateClassLevelRequest struct {
	Tingkat     int     `json:"tingkat" binding:"required,gte=1,lte=6"`
	NamaTingkat string  `json:"nama_tingkat" binding:"required"`
	BiayaSPP    float64 `json:"biaya_spp" binding:"required,gt=0"`
}

type UpdateClassLevelRequest struct {
	Tingkat     int     `json:"tingkat" binding:"required,gte=1,lte=6"`
	NamaTingkat string  `json:"nama_tingkat" binding:"required"`
	BiayaSPP    float64 `json:"biaya_spp" binding:"required,gt=0"`
	Status      string  `json:"status" binding:"required,oneof=aktif nonaktif"`
}

type AdminHandler interface {
	CreateUser(c *gin.Context)
	FindAllUsers(c *gin.Context)
	FindUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	CreateClassLevel(c *gin.Context)
	FindAllClassLevels(c *gin.Context)
	FindClassLevelByID(c *gin.Context)
	UpdateClassLevel(c *gin.Context)
	DeleteClassLevel(c *gin.Context)
}

type adminHandler struct {
	userService       service.UserService
	classLevelService service.ClassLevelService
}

func NewAdminHandler(userService service.UserService, classLevelService service.ClassLevelService) AdminHandler {
	return &adminHandler{userService, classLevelService}
}

func (h *adminHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.CreateUserInput{
		NamaLengkap: req.NamaLengkap,
		Email:       req.Email,
		Password:    req.Password,
		RoleID:      req.RoleID,
	}

	createdUser, err := h.userService.CreateUser(input)
	if err != nil {
		if err.Error() == "email sudah terdaftar" {
			utils.SendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuat pengguna: "+err.Error())
		return
	}

	response := UserResponse{
		ID:          createdUser.ID,
		NamaLengkap: createdUser.NamaLengkap,
		Email:       createdUser.Email,
		Role:        createdUser.Role.NamaRole,
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Pengguna berhasil dibuat", response)
}

func (h *adminHandler) FindAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	roleID, _ := strconv.Atoi(c.Query("role_id"))
	search := c.Query("search")

	input := service.FindAllUsersInput{
		Page:   page,
		Limit:  limit,
		RoleID: uint(roleID),
		Search: search,
	}

	users, total, err := h.userService.FindAllUsers(input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pengguna")
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:          user.ID,
			NamaLengkap: user.NamaLengkap,
			Email:       user.Email,
			Role:        user.Role.NamaRole,
		})
	}

	response := gin.H{
		"data": userResponses,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data pengguna berhasil diambil", response)
}

func (h *adminHandler) FindUserByID(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	user, err := h.userService.GetUserProfile(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Pengguna dengan ID tersebut tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pengguna")
		return
	}

	response := UserResponse{
		ID:          user.ID,
		NamaLengkap: user.NamaLengkap,
		Email:       user.Email,
		Role:        user.Role.NamaRole,
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Detail pengguna berhasil diambil", response)
}

func (h *adminHandler) UpdateUser(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.UpdateUserInput{
		NamaLengkap: req.NamaLengkap,
		Email:       req.Email,
		RoleID:      req.RoleID,
	}

	updatedUser, err := h.userService.UpdateUser(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Pengguna dengan ID tersebut tidak ditemukan")
			return
		}
		if err.Error() == "email sudah terdaftar untuk pengguna lain" {
			utils.SendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data pengguna")
		return
	}

	response := UserResponse{
		ID:          updatedUser.ID,
		NamaLengkap: updatedUser.NamaLengkap,
		Email:       updatedUser.Email,
		Role:        updatedUser.Role.NamaRole,
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data pengguna berhasil diperbarui", response)
}

func (h *adminHandler) DeleteUser(c *gin.Context) {
	idString := c.Param("id")
	idToDelete, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	adminID := c.MustGet("userID").(uint)

	if adminID == uint(idToDelete) {
		utils.SendErrorResponse(c, http.StatusForbidden, "Anda tidak dapat menghapus akun Anda sendiri")
		return
	}

	err = h.userService.DeleteUser(uint(idToDelete))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Pengguna dengan ID tersebut tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus pengguna")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Pengguna berhasil dihapus", nil)
}

func (h *adminHandler) CreateClassLevel(c *gin.Context) {
	var req CreateClassLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.CreateClassLevelInput{
		Tingkat:     req.Tingkat,
		NamaTingkat: req.NamaTingkat,
		BiayaSPP:    req.BiayaSPP,
	}

	classLevel, err := h.classLevelService.CreateClassLevel(input)
	if err != nil {
		if err.Error() == "tingkat kelas sudah ada" {
			utils.SendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuat tingkat kelas")
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Tingkat kelas berhasil dibuat", classLevel)
}

func (h *adminHandler) FindAllClassLevels(c *gin.Context) {
	classLevels, err := h.classLevelService.FindAllClassLevels()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data tingkat kelas")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data tingkat kelas berhasil diambil", classLevels)
}

func (h *adminHandler) FindClassLevelByID(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tingkat kelas tidak valid")
		return
	}

	classLevel, err := h.classLevelService.FindClassLevelByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Tingkat kelas dengan ID tersebut tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data tingkat kelas")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Detail tingkat kelas berhasil diambil", classLevel)
}

func (h *adminHandler) UpdateClassLevel(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tingkat kelas tidak valid")
		return
	}

	var req UpdateClassLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.UpdateClassLevelInput{
		Tingkat:     req.Tingkat,
		NamaTingkat: req.NamaTingkat,
		BiayaSPP:    req.BiayaSPP,
		Status:      req.Status,
	}

	updated, err := h.classLevelService.UpdateClassLevel(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Tingkat kelas dengan ID tersebut tidak ditemukan")
			return
		}
		if err.Error() == "tingkat kelas sudah ada" {
			utils.SendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui tingkat kelas")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data tingkat kelas berhasil diperbarui", updated)
}

func (h *adminHandler) DeleteClassLevel(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tingkat kelas tidak valid")
		return
	}

	err = h.classLevelService.DeleteClassLevel(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Tingkat kelas dengan ID tersebut tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus tingkat kelas")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data tingkat kelas berhasil dihapus", nil)
}
