package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/hiuncy/spp-payment-api/internal/service"
	"github.com/hiuncy/spp-payment-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/hiuncy/spp-payment-api/internal/model"
	"gorm.io/gorm"
)

type PeriodRequest struct {
	TahunAjaran    string `json:"tahun_ajaran" binding:"required"`
	Bulan          int    `json:"bulan" binding:"required,gte=1,lte=12"`
	NamaBulan      string `json:"nama_bulan" binding:"required"`
	TanggalMulai   string `json:"tanggal_mulai" binding:"required"`
	TanggalSelesai string `json:"tanggal_selesai" binding:"required"`
}

type UpdatePeriodRequest struct {
	TahunAjaran    string `json:"tahun_ajaran" binding:"required"`
	Bulan          int    `json:"bulan" binding:"required,gte=1,lte=12"`
	NamaBulan      string `json:"nama_bulan" binding:"required"`
	TanggalMulai   string `json:"tanggal_mulai" binding:"required"`
	TanggalSelesai string `json:"tanggal_selesai" binding:"required"`
	Status         string `json:"status" binding:"required,oneof=belum_aktif aktif selesai"`
}

type CreateStudentRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	NISN            string `json:"nisn" binding:"required"`
	KelasID         uint   `json:"kelas_id" binding:"required"`
	NamaLengkap     string `json:"nama_lengkap" binding:"required"`
	JenisKelamin    string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TempatLahir     string `json:"tempat_lahir"`
	TanggalLahir    string `json:"tanggal_lahir"`
	Alamat          string `json:"alamat"`
	NamaOrangTua    string `json:"nama_orangtua"`
	TeleponOrangTua string `json:"telepon_orangtua"`
	TahunMasuk      int    `json:"tahun_masuk"`
}

type UpdateStudentRequest struct {
	NISN            string `json:"nisn" binding:"required"`
	KelasID         uint   `json:"kelas_id" binding:"required"`
	NamaLengkap     string `json:"nama_lengkap" binding:"required"`
	JenisKelamin    string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TempatLahir     string `json:"tempat_lahir"`
	TanggalLahir    string `json:"tanggal_lahir"`
	Alamat          string `json:"alamat"`
	NamaOrangTua    string `json:"nama_orangtua"`
	TeleponOrangTua string `json:"telepon_orangtua"`
	TahunMasuk      int    `json:"tahun_masuk"`
	Status          string `json:"status" binding:"required,oneof=aktif pindah lulus keluar"`
	EmailUser       string `json:"email" binding:"required,email"`
	StatusUser      string `json:"status_user" binding:"required,oneof=aktif nonaktif"`
}

type StudentResponse struct {
	ID              uint       `json:"id"`
	NISN            string     `json:"nisn"`
	NamaLengkap     string     `json:"nama_lengkap"`
	Email           string     `json:"email"`
	Status          string     `json:"status"`
	KelasID         uint       `json:"kelas_id"`
	NamaKelas       string     `json:"nama_kelas"`
	JenisKelamin    string     `json:"jenis_kelamin"`
	TempatLahir     string     `json:"tempat_lahir,omitempty"`
	TanggalLahir    *time.Time `json:"tanggal_lahir,omitempty"`
	Alamat          string     `json:"alamat,omitempty"`
	NamaOrangTua    string     `json:"nama_orangtua,omitempty"`
	TeleponOrangTua string     `json:"telepon_orangtua,omitempty"`
	TahunMasuk      int        `json:"tahun_masuk,omitempty"`
}

type UpdateBillRequest struct {
	JumlahTagihan    float64 `json:"jumlah_tagihan" binding:"required,gt=0"`
	StatusPembayaran string  `json:"status_pembayaran" binding:"required,oneof=belum_bayar pending lunas"`
}

type BillResponse struct {
	ID                uint      `json:"id"`
	SiswaID           uint      `json:"siswa_id"`
	NamaSiswa         string    `json:"nama_siswa"`
	PeriodeID         uint      `json:"periode_id"`
	NamaPeriode       string    `json:"nama_periode"`
	TahunAjaran       string    `json:"tahun_ajaran"`
	JumlahTagihan     float64   `json:"jumlah_tagihan"`
	StatusPembayaran  string    `json:"status_pembayaran"`
	TanggalJatuhTempo time.Time `json:"tanggal_jatuh_tempo"`
}

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

func formatStudentResponse(student *model.Siswa) StudentResponse {
	return StudentResponse{
		ID:              student.ID,
		NISN:            student.NISN,
		NamaLengkap:     student.NamaLengkap,
		Email:           student.User.Email,
		Status:          student.Status,
		KelasID:         student.KelasID,
		NamaKelas:       student.Kelas.NamaKelas,
		JenisKelamin:    student.JenisKelamin,
		TempatLahir:     student.TempatLahir,
		TanggalLahir:    student.TanggalLahir,
		Alamat:          student.Alamat,
		NamaOrangTua:    student.NamaOrangTua,
		TeleponOrangTua: student.TeleponOrangTua,
		TahunMasuk:      student.TahunMasuk,
	}
}

func formatBillResponse(bill *model.TagihanSPP) BillResponse {
	return BillResponse{
		ID:                bill.ID,
		SiswaID:           bill.SiswaID,
		NamaSiswa:         bill.Siswa.NamaLengkap,
		PeriodeID:         bill.PeriodeID,
		NamaPeriode:       bill.PeriodeSPP.NamaBulan,
		TahunAjaran:       bill.PeriodeSPP.TahunAjaran,
		JumlahTagihan:     bill.JumlahTagihan,
		StatusPembayaran:  bill.StatusPembayaran,
		TanggalJatuhTempo: bill.TanggalJatuhTempo,
	}
}

func (h *treasurerHandler) CreateStudent(c *gin.Context) {
	var req CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.CreateStudentInput{
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

	input := service.FindAllStudentsInput{
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

	var responses []StudentResponse
	for _, student := range students {
		responses = append(responses, formatStudentResponse(&student))
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

	utils.SendSuccessResponse(c, http.StatusOK, "Detail siswa berhasil diambil", formatStudentResponse(student))
}

func (h *treasurerHandler) UpdateStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID siswa tidak valid")
		return
	}

	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.UpdateStudentInput{
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

	utils.SendSuccessResponse(c, http.StatusOK, "Data siswa berhasil diperbarui", formatStudentResponse(student))
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
	var req PeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.CreatePeriodInput{
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

	var req UpdatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.UpdatePeriodInput{
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

	input := service.FindAllBillsInput{
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

	var responses []BillResponse
	for _, bill := range bills {
		responses = append(responses, formatBillResponse(&bill))
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
	utils.SendSuccessResponse(c, http.StatusOK, "Detail tagihan berhasil diambil", formatBillResponse(bill))
}

func (h *treasurerHandler) UpdateBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "ID tagihan tidak valid")
		return
	}

	var req UpdateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Input tidak valid: "+err.Error())
		return
	}

	input := service.UpdateBillInput{
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
	utils.SendSuccessResponse(c, http.StatusOK, "Tagihan berhasil diperbarui", formatBillResponse(bill))
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
