package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hiuncy/spp-payment-api/internal/dto"
	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

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
	CreateClass(c *gin.Context)
	FindAllClasses(c *gin.Context)
	FindClassByID(c *gin.Context)
	UpdateClass(c *gin.Context)
	DeleteClass(c *gin.Context)
	FindAllSettings(c *gin.Context)
	UpdateSettings(c *gin.Context)
}

type adminHandler struct {
	userService       service.UserService
	classLevelService service.ClassLevelService
	classService      service.ClassService
	settingService    service.SettingService
}

func NewAdminHandler(userService service.UserService, classLevelService service.ClassLevelService, classService service.ClassService, settingService service.SettingService) AdminHandler {
	return &adminHandler{userService, classLevelService, classService, settingService}
}

func (h *adminHandler) CreateUser(c *gin.Context) {
	var req utils.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.CreateUserInput{
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

	response := utils.UserResponse{
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

	input := dto.FindAllUsersInput{
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

	var userResponses []utils.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, utils.UserResponse{
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

	response := utils.UserResponse{
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

	var req utils.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.UpdateUserInput{
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

	response := utils.UserResponse{
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
	var req utils.CreateClassLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.CreateClassLevelInput{
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

	var req utils.UpdateClassLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.UpdateClassLevelInput{
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

func (h *adminHandler) CreateClass(c *gin.Context) {
	var req utils.ClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}
	input := dto.CreateClassInput{
		TingkatID: req.TingkatID,
		NamaKelas: req.NamaKelas,
		WaliKelas: req.WaliKelas,
		Kapasitas: req.Kapasitas,
	}
	newClass, err := h.classService.CreateClass(input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccessResponse(c, http.StatusCreated, "Kelas berhasil dibuat", utils.FormatClassResponse(newClass))
}

func (h *adminHandler) FindAllClasses(c *gin.Context) {
	classes, err := h.classService.FindAllClasses()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data kelas")
		return
	}
	var responses []utils.ClassResponse
	for _, class := range classes {
		responses = append(responses, utils.FormatClassResponse(&class))
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Data kelas berhasil diambil", responses)
}

func (h *adminHandler) FindClassByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	class, err := h.classService.FindClassByID(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Kelas tidak ditemukan")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Detail kelas berhasil diambil", utils.FormatClassResponse(class))
}

func (h *adminHandler) UpdateClass(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req utils.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}
	input := dto.UpdateClassInput{
		TingkatID: req.TingkatID,
		NamaKelas: req.NamaKelas,
		WaliKelas: req.WaliKelas,
		Kapasitas: req.Kapasitas,
		Status:    req.Status,
	}
	updatedClass, err := h.classService.UpdateClass(uint(id), input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Kelas berhasil diperbarui", utils.FormatClassResponse(updatedClass))
}

func (h *adminHandler) DeleteClass(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := h.classService.DeleteClass(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Kelas tidak ditemukan")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Kelas berhasil dihapus", nil)
}

func (h *adminHandler) FindAllSettings(c *gin.Context) {
	settings, err := h.settingService.FindAllSettings()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pengaturan")
		return
	}

	response := make(map[string]string)
	for _, setting := range settings {
		response[setting.KeySetting] = setting.ValueSetting
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data pengaturan berhasil diambil", response)
}

func (h *adminHandler) UpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	err := h.settingService.UpdateSettings(req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui pengaturan")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Pengaturan berhasil diperbarui", nil)
}
