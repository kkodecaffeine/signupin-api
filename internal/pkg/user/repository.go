package user

import (
	"signupin-api/internal/app/api/dto"
)

// Repository interface definition
type Repository interface {
	SaveOne(model *User) (string, error)

	// GET
	GetOne(email string, password ...string) (*dto.PostSignUpResponse, error)
	GetOneByID(ID string) (*dto.PostSignUpResponse, error)
}
