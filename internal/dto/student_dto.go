package dto

type CreateStudentInput struct {
	Email           string
	Password        string
	NISN            string
	KelasID         uint
	NamaLengkap     string
	JenisKelamin    string
	TempatLahir     string
	TanggalLahir    string
	Alamat          string
	NamaOrangTua    string
	TeleponOrangTua string
	TahunMasuk      int
}

type UpdateStudentInput struct {
	NISN            string
	KelasID         uint
	NamaLengkap     string
	JenisKelamin    string
	TempatLahir     string
	TanggalLahir    string
	Alamat          string
	NamaOrangTua    string
	TeleponOrangTua string
	TahunMasuk      int
	Status          string
	EmailUser       string
	StatusUser      string
}

type FindAllStudentsInput struct {
	Page    int
	Limit   int
	KelasID uint
	Search  string
}
