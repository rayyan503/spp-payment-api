package dto

type CreateClassLevelInput struct {
	Tingkat     int
	NamaTingkat string
	BiayaSPP    float64
}

type UpdateClassLevelInput struct {
	Tingkat     int
	NamaTingkat string
	BiayaSPP    float64
	Status      string
}
