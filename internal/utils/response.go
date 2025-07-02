package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hiuncy/spp-payment-api/internal/model"
)

type ClassResponse struct {
	ID          uint    `json:"id"`
	NamaKelas   string  `json:"nama_kelas"`
	WaliKelas   string  `json:"wali_kelas"`
	Kapasitas   int     `json:"kapasitas"`
	Status      string  `json:"status"`
	TingkatID   uint    `json:"tingkat_id"`
	NamaTingkat string  `json:"nama_tingkat"`
	BiayaSPP    float64 `json:"biaya_spp"`
}

type UserResponse struct {
	ID          uint   `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

type PaymentHistoryResponse struct {
	OrderID           string     `json:"order_id"`
	NamaPeriode       string     `json:"nama_periode"`
	TahunAjaran       string     `json:"tahun_ajaran"`
	JumlahBayar       float64    `json:"jumlah_bayar"`
	StatusPembayaran  string     `json:"status_pembayaran"`
	MetodePembayaran  *string    `json:"metode_pembayaran,omitempty"`
	TanggalPembayaran *time.Time `json:"tanggal_pembayaran,omitempty"`
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

func FormatClassResponse(class *model.Kelas) ClassResponse {
	return ClassResponse{
		ID:          class.ID,
		NamaKelas:   class.NamaKelas,
		WaliKelas:   class.WaliKelas,
		Kapasitas:   class.Kapasitas,
		Status:      class.Status,
		TingkatID:   class.TingkatID,
		NamaTingkat: class.TingkatKelas.NamaTingkat,
		BiayaSPP:    class.TingkatKelas.BiayaSPP,
	}
}

func FormatStudentResponse(student *model.Siswa) StudentResponse {
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
		NamaOrangTua:    student.NamaOrangtua,
		TeleponOrangTua: student.TeleponOrangtua,
		TahunMasuk:      student.TahunMasuk,
	}
}

func FormatBillResponse(bill *model.TagihanSPP) BillResponse {
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

func SendSuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func SendErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
