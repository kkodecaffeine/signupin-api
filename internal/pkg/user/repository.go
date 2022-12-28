package user

import (
	"signupin-api/internal/app/api/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository interface definition
type Repository interface {
	SaveOne(model *User) (string, error)

	// GET
	GetOne(email string, password ...string) (*dto.GetUserWithTokenResponse, error)
	GetOneByID(ID string) (*dto.GetUserResponse, error)

	// UPDATE
	UpdatePassword(ID primitive.ObjectID, newpassword string) (*dto.GetUserResponse, error)
}
