package dto

type CreateUserInput struct {
	NamaLengkap string
	Email       string
	Password    string
	RoleID      uint
}

type FindAllUsersInput struct {
	Page   int
	Limit  int
	RoleID uint
	Search string
}

type UpdateUserInput struct {
	NamaLengkap string
	Email       string
	RoleID      uint
}
