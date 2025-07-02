package model

import "time"

type LaporanSiswa struct {
	NISN              string     `gorm:"column:nisn" json:"nisn"`
	NamaLengkap       string     `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	NamaKelas         string     `gorm:"column:nama_kelas" json:"nama_kelas"`
	NamaTingkat       string     `gorm:"column:nama_tingkat" json:"nama_tingkat"`
	TahunAjaran       string     `gorm:"column:tahun_ajaran" json:"tahun_ajaran"`
	NamaBulan         string     `gorm:"column:nama_bulan" json:"nama_bulan"`
	JumlahTagihan     float64    `gorm:"column:jumlah_tagihan" json:"jumlah_tagihan"`
	StatusPembayaran  string     `gorm:"column:status_pembayaran" json:"status_pembayaran"`
	TanggalJatuhTempo time.Time  `gorm:"column:tanggal_jatuh_tempo" json:"tanggal_jatuh_tempo"`
	TanggalSettlement *time.Time `gorm:"column:tanggal_settlement" json:"tanggal_settlement"`
	MetodePembayaran  *string    `gorm:"column:metode_pembayaran" json:"metode_pembayaran"`
}

type LaporanKelas struct {
	NamaKelas       string  `gorm:"column:nama_kelas" json:"nama_kelas"`
	NamaTingkat     string  `gorm:"column:nama_tingkat" json:"nama_tingkat"`
	TahunAjaran     string  `gorm:"column:tahun_ajaran" json:"tahun_ajaran"`
	NamaBulan       string  `gorm:"column:nama_bulan" json:"nama_bulan"`
	TotalSiswa      int     `gorm:"column:total_siswa" json:"total_siswa"`
	SiswaLunas      int     `gorm:"column:siswa_lunas" json:"siswa_lunas"`
	SiswaBelumBayar int     `gorm:"column:siswa_belum_bayar" json:"siswa_belum_bayar"`
	SiswaPending    int     `gorm:"column:siswa_pending" json:"siswa_pending"`
	TotalTagihan    float64 `gorm:"column:total_tagihan" json:"total_tagihan"`
	TotalTerbayar   float64 `gorm:"column:total_terbayar" json:"total_terbayar"`
}

type LaporanKeseluruhan struct {
	TahunAjaran          string  `gorm:"column:tahun_ajaran" json:"tahun_ajaran"`
	NamaBulan            string  `gorm:"column:nama_bulan" json:"nama_bulan"`
	TotalTagihan         int     `gorm:"column:total_tagihan" json:"total_tagihan"`
	TotalLunas           int     `gorm:"column:total_lunas" json:"total_lunas"`
	TotalBelumBayar      int     `gorm:"column:total_belum_bayar" json:"total_belum_bayar"`
	TotalPending         int     `gorm:"column:total_pending" json:"total_pending"`
	TotalNominalTagihan  float64 `gorm:"column:total_nominal_tagihan" json:"total_nominal_tagihan"`
	TotalNominalTerbayar float64 `gorm:"column:total_nominal_terbayar" json:"total_nominal_terbayar"`
	PersentasePembayaran float64 `gorm:"column:persentase_pembayaran" json:"persentase_pembayaran"`
}
