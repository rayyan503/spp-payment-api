package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHistoryResponse struct {
	OrderID           string     `json:"order_id"`
	NamaPeriode       string     `json:"nama_periode"`
	TahunAjaran       string     `json:"tahun_ajaran"`
	JumlahBayar       float64    `json:"jumlah_bayar"`
	StatusPembayaran  string     `json:"status_pembayaran"`
	MetodePembayaran  *string    `json:"metode_pembayaran,omitempty"`
	TanggalPembayaran *time.Time `json:"tanggal_pembayaran,omitempty"`
}

type StudentHandler interface {
	GetProfile(c *gin.Context)
	FindMyBills(c *gin.Context)
	InitiatePayment(c *gin.Context)
	GetPaymentHistory(c *gin.Context)
}

type studentHandler struct {
	studentService service.StudentService
	billService    service.BillService
	paymentService service.PaymentService
}

func NewStudentHandler(studentService service.StudentService, billService service.BillService, paymentService service.PaymentService) StudentHandler {
	return &studentHandler{studentService, billService, paymentService}
}

func (h *studentHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	student, err := h.studentService.GetStudentProfile(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Profil siswa untuk pengguna ini tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil profil siswa")
		return
	}

	response := formatStudentResponse(student)
	utils.SendSuccessResponse(c, http.StatusOK, "Profil siswa berhasil diambil", response)
}

func (h *studentHandler) FindMyBills(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	student, err := h.studentService.GetStudentProfile(userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Profil siswa untuk pengguna ini tidak ditemukan")
		return
	}

	status := c.Query("status")

	input := service.FindAllBillsInput{
		SiswaID:          student.ID,
		StatusPembayaran: status,
		Limit:            100,
		Page:             1,
	}

	bills, _, err := h.billService.FindAllBills(input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data tagihan")
		return
	}

	var responses []BillResponse
	for _, bill := range bills {
		responses = append(responses, FormatBillResponse(&bill))
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data tagihan berhasil diambil", responses)
}

func (h *studentHandler) InitiatePayment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	billID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tagihan tidak valid")
		return
	}

	snapToken, err := h.paymentService.InitiatePayment(uint(billID), userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Token pembayaran berhasil dibuat", gin.H{"snap_token": snapToken})
}

func (h *studentHandler) GetPaymentHistory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	payments, err := h.paymentService.GetPaymentHistory(userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Format response
	var responses []PaymentHistoryResponse
	for _, p := range payments {
		responses = append(responses, PaymentHistoryResponse{
			OrderID:           p.OrderID,
			NamaPeriode:       p.TagihanSPP.PeriodeSPP.NamaBulan,
			TahunAjaran:       p.TagihanSPP.PeriodeSPP.TahunAjaran,
			JumlahBayar:       p.JumlahBayar,
			StatusPembayaran:  p.StatusPembayaran,
			MetodePembayaran:  p.MetodePembayaran,
			TanggalPembayaran: p.TanggalPembayaran,
		})
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Riwayat pembayaran berhasil diambil", responses)
}
