package dto

type CreateClassInput struct {
	TingkatID uint
	NamaKelas string
	WaliKelas string
	Kapasitas int
}

type UpdateClassInput struct {
	TingkatID uint
	NamaKelas string
	WaliKelas string
	Kapasitas int
	Status    string
}
