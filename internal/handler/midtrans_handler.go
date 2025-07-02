// File: internal/handler/midtrans_handler.go (File Baru)

package handler

import (
	"net/http"

	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type MidtransHandler interface {
	HandleNotification(c *gin.Context)
}

type midtransHandler struct {
	paymentService service.PaymentService
}

func NewMidtransHandler(paymentService service.PaymentService) MidtransHandler {
	return &midtransHandler{paymentService}
}

func (h *midtransHandler) HandleNotification(c *gin.Context) {
	var notificationPayload map[string]interface{}

	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Payload notifikasi tidak valid")
		return
	}

	err := h.paymentService.ProcessMidtransNotification(notificationPayload)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
