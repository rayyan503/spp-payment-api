package dto

type FindAllBillsInput struct {
	Page             int
	Limit            int
	PeriodeID        uint
	SiswaID          uint
	StatusPembayaran string
}

type UpdateBillInput struct {
	JumlahTagihan    float64
	StatusPembayaran string
}
