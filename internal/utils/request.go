package utils

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

type ClassRequest struct {
	TingkatID uint   `json:"tingkat_id" binding:"required"`
	NamaKelas string `json:"nama_kelas" binding:"required"`
	WaliKelas string `json:"wali_kelas"`
	Kapasitas int    `json:"kapasitas" binding:"gte=0"`
}

type UpdateClassRequest struct {
	TingkatID uint   `json:"tingkat_id" binding:"required"`
	NamaKelas string `json:"nama_kelas" binding:"required"`
	WaliKelas string `json:"wali_kelas"`
	Kapasitas int    `json:"kapasitas" binding:"gte=0"`
	Status    string `json:"status" binding:"required,oneof=aktif nonaktif"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

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

type UpdateBillRequest struct {
	JumlahTagihan    float64 `json:"jumlah_tagihan" binding:"required,gt=0"`
	StatusPembayaran string  `json:"status_pembayaran" binding:"required,oneof=belum_bayar pending lunas"`
}
