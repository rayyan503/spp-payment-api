package utils

type FindAllBillsParams struct {
	Limit            int
	Page             int
	PeriodeID        uint
	SiswaID          uint
	StatusPembayaran string
}

type FindAllStudentsParams struct {
	Limit   int
	Page    int
	KelasID uint
	Search  string
}

type FindAllUsersParams struct {
	Limit  int
	Page   int
	RoleID uint
	Search string
}
