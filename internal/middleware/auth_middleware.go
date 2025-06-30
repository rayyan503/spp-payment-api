package middleware

import (
	"net/http"
	"strings"

	"github.com/hiuncy/spp-payment-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Header otorisasi tidak ditemukan")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Format token salah, harus 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenString, secretKey)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau telah kedaluwarsa")
			c.Abort()
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.SendErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki akses ke sumber daya ini")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
