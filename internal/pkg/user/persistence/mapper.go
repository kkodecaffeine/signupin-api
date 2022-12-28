package persistence

import (
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"

	"github.com/kkodecaffeine/go-common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type entityMapper struct{}

func (e entityMapper) toDomainProps(ID primitive.ObjectID, model *user.User) *dto.GetUserResponse {

	id := utils.MapToStringID(ID)

	return &dto.GetUserResponse{
		Id:       id,
		Email:    model.Email,
		Name:     model.Name,
		NickName: model.NickName,
		Phone:    model.Phone,
	}
}

func (e entityMapper) toDomainProps2(ID primitive.ObjectID, model *user.User) *dto.GetUserWithTokenResponse {

	id := utils.MapToStringID(ID)

	return &dto.GetUserWithTokenResponse{
		Id:       id,
		Email:    model.Email,
		Name:     model.Name,
		NickName: model.NickName,
		Phone:    model.Phone,
	}
}
