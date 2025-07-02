package dto

type CreatePeriodInput struct {
	TahunAjaran    string
	Bulan          int
	NamaBulan      string
	TanggalMulai   string
	TanggalSelesai string
}

type UpdatePeriodInput struct {
	TahunAjaran    string
	Bulan          int
	NamaBulan      string
	TanggalMulai   string
	TanggalSelesai string
	Status         string
}
