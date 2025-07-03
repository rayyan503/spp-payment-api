package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hiuncy/spp-payment-api/internal/dto"
	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TreasurerHandler interface {
	CreateStudent(c *gin.Context)
	FindAllStudents(c *gin.Context)
	FindStudentByID(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudent(c *gin.Context)
	CreatePeriod(c *gin.Context)
	FindAllPeriods(c *gin.Context)
	FindPeriodByID(c *gin.Context)
	UpdatePeriod(c *gin.Context)
	DeletePeriod(c *gin.Context)
	GenerateBills(c *gin.Context)
	FindAllBills(c *gin.Context)
	FindBillByID(c *gin.Context)
	UpdateBill(c *gin.Context)
	DeleteBill(c *gin.Context)
	GetLaporanSiswa(c *gin.Context)
	GetLaporanKelas(c *gin.Context)
	GetLaporanKeseluruhan(c *gin.Context)
}

type treasurerHandler struct {
	studentService service.StudentService
	periodService  service.PeriodService
	billService    service.BillService
	reportService  service.ReportService
}

func NewTreasurerHandler(studentService service.StudentService, periodService service.PeriodService, billService service.BillService, reportService service.ReportService) TreasurerHandler {
	return &treasurerHandler{studentService, periodService, billService, reportService}
}

func (h *treasurerHandler) CreateStudent(c *gin.Context) {
	var req utils.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.CreateStudentInput{
		Email:           req.Email,
		Password:        req.Password,
		NISN:            req.NISN,
		KelasID:         req.KelasID,
		NamaLengkap:     req.NamaLengkap,
		JenisKelamin:    req.JenisKelamin,
		TempatLahir:     req.TempatLahir,
		TanggalLahir:    req.TanggalLahir,
		Alamat:          req.Alamat,
		NamaOrangTua:    req.NamaOrangTua,
		TeleponOrangTua: req.TeleponOrangTua,
		TahunMasuk:      req.TahunMasuk,
	}

	student, err := h.studentService.CreateStudent(input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Data siswa berhasil dibuat", student)
}

func (h *treasurerHandler) FindAllStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	kelasID, _ := strconv.Atoi(c.Query("kelas_id"))
	search := c.Query("search")

	input := dto.FindAllStudentsInput{
		Page:    page,
		Limit:   limit,
		KelasID: uint(kelasID),
		Search:  search,
	}

	students, total, err := h.studentService.FindAllStudents(input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data siswa")
		return
	}

	var responses []utils.StudentResponse
	for _, student := range students {
		responses = append(responses, utils.FormatStudentResponse(&student))
	}

	response := gin.H{
		"data": responses,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Data semua siswa berhasil diambil", response)
}

func (h *treasurerHandler) FindStudentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID siswa tidak valid")
		return
	}

	student, err := h.studentService.FindStudentByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Siswa tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail siswa")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Detail siswa berhasil diambil", utils.FormatStudentResponse(student))
}

func (h *treasurerHandler) UpdateStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID siswa tidak valid")
		return
	}

	var req utils.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.UpdateStudentInput{
		NISN:            req.NISN,
		KelasID:         req.KelasID,
		NamaLengkap:     req.NamaLengkap,
		JenisKelamin:    req.JenisKelamin,
		TempatLahir:     req.TempatLahir,
		TanggalLahir:    req.TanggalLahir,
		Alamat:          req.Alamat,
		NamaOrangTua:    req.NamaOrangTua,
		TeleponOrangTua: req.TeleponOrangTua,
		TahunMasuk:      req.TahunMasuk,
		Status:          req.Status,
		EmailUser:       req.EmailUser,
		StatusUser:      req.StatusUser,
	}

	student, err := h.studentService.UpdateStudent(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Siswa tidak ditemukan")
			return
		}
		if err.Error() == "NISN sudah terdaftar untuk siswa lain" || err.Error() == "email sudah terdaftar untuk pengguna lain" {
			utils.SendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data siswa")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data siswa berhasil diperbarui", utils.FormatStudentResponse(student))
}

func (h *treasurerHandler) DeleteStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID siswa tidak valid")
		return
	}

	err = h.studentService.DeleteStudent(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Siswa tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data siswa")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Data siswa berhasil dihapus", nil)
}

func (h *treasurerHandler) CreatePeriod(c *gin.Context) {
	var req utils.PeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.CreatePeriodInput{
		TahunAjaran:    req.TahunAjaran,
		Bulan:          req.Bulan,
		NamaBulan:      req.NamaBulan,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
	}

	period, err := h.periodService.CreatePeriod(input)
	if err != nil {
		if err.Error() == "periode untuk tahun ajaran dan bulan tersebut sudah ada" {
			utils.SendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal membuat periode")
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Periode berhasil dibuat", period)
}

func (h *treasurerHandler) FindAllPeriods(c *gin.Context) {
	tahunAjaran := c.Query("tahun_ajaran")
	periods, err := h.periodService.FindAllPeriods(tahunAjaran)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data periode")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Data periode berhasil diambil", periods)
}

func (h *treasurerHandler) FindPeriodByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID periode tidak valid")
		return
	}

	period, err := h.periodService.FindPeriodByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Periode tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail periode")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Detail periode berhasil diambil", period)
}

func (h *treasurerHandler) UpdatePeriod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID periode tidak valid")
		return
	}

	var req utils.UpdatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.UpdatePeriodInput{
		TahunAjaran:    req.TahunAjaran,
		Bulan:          req.Bulan,
		NamaBulan:      req.NamaBulan,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
		Status:         req.Status,
	}

	period, err := h.periodService.UpdatePeriod(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Periode tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui periode")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Periode berhasil diperbarui", period)
}

func (h *treasurerHandler) DeletePeriod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID periode tidak valid")
		return
	}

	err = h.periodService.DeletePeriod(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Periode tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus periode")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Periode berhasil dihapus", nil)
}

func (h *treasurerHandler) GenerateBills(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID periode tidak valid")
		return
	}

	err = h.billService.GenerateBillsForPeriod(uint(id))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal men-generate tagihan: "+err.Error())
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Tagihan untuk periode terpilih berhasil di-generate", nil)
}

func (h *treasurerHandler) FindAllBills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	periodeID, _ := strconv.Atoi(c.Query("periode_id"))
	siswaID, _ := strconv.Atoi(c.Query("siswa_id"))
	status := c.Query("status_pembayaran")

	input := dto.FindAllBillsInput{
		Page:             page,
		Limit:            limit,
		PeriodeID:        uint(periodeID),
		SiswaID:          uint(siswaID),
		StatusPembayaran: status,
	}

	bills, total, err := h.billService.FindAllBills(input)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data tagihan")
		return
	}

	var responses []utils.BillResponse
	for _, bill := range bills {
		responses = append(responses, utils.FormatBillResponse(&bill))
	}

	response := gin.H{
		"data": responses,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Data tagihan berhasil diambil", response)
}

func (h *treasurerHandler) FindBillByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tagihan tidak valid")
		return
	}

	bill, err := h.billService.FindBillByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Tagihan tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail tagihan")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Detail tagihan berhasil diambil", utils.FormatBillResponse(bill))
}

func (h *treasurerHandler) UpdateBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tagihan tidak valid")
		return
	}

	var req utils.UpdateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := dto.UpdateBillInput{
		JumlahTagihan:    req.JumlahTagihan,
		StatusPembayaran: req.StatusPembayaran,
	}

	bill, err := h.billService.UpdateBill(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Tagihan tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui tagihan")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Tagihan berhasil diperbarui", utils.FormatBillResponse(bill))
}

func (h *treasurerHandler) DeleteBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tagihan tidak valid")
		return
	}

	err = h.billService.DeleteBill(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, "Tagihan tidak ditemukan")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus tagihan")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Tagihan berhasil dihapus", nil)
}

func (h *treasurerHandler) GetLaporanSiswa(c *gin.Context) {
	tahunAjaran := c.Query("tahun_ajaran")
	nisn := c.Query("nisn")
	result, err := h.reportService.GetLaporanSiswa(tahunAjaran, nisn)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil laporan per siswa")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Laporan per siswa berhasil diambil", result)
}

func (h *treasurerHandler) GetLaporanKelas(c *gin.Context) {
	tahunAjaran := c.Query("tahun_ajaran")
	namaBulan := c.Query("nama_bulan")
	result, err := h.reportService.GetLaporanKelas(tahunAjaran, namaBulan)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil laporan per kelas")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Laporan per kelas berhasil diambil", result)
}

func (h *treasurerHandler) GetLaporanKeseluruhan(c *gin.Context) {
	tahunAjaran := c.Query("tahun_ajaran")
	result, err := h.reportService.GetLaporanKeseluruhan(tahunAjaran)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil laporan keseluruhan")
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Laporan keseluruhan berhasil diambil", result)
}
