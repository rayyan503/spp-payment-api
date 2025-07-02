package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"
)

type AuthHandler interface {
	Login(c *gin.Context)
	GetMe(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthHandler(authService service.AuthService, userService service.UserService) AuthHandler {
	return &authHandler{authService, userService}
}

func (h *authHandler) Login(c *gin.Context) {
	var req utils.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{"token": token})
}

func (h *authHandler) GetMe(c *gin.Context) {
	userID, ok := c.MustGet("userID").(uint)
	if !ok {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memproses ID pengguna dari token")
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	response := utils.UserResponse{
		ID:          user.ID,
		NamaLengkap: user.NamaLengkap,
		Email:       user.Email,
		Role:        user.Role.NamaRole,
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Profil pengguna berhasil diambil", response)
}
