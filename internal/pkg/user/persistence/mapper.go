package persistence

import (
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"

	"github.com/kkodecaffeine/go-common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type entityMapper struct{}

func (e entityMapper) toDomainProps(ID primitive.ObjectID, model *user.User) *dto.PostSignUpResponse {

	id := utils.MapToStringID(ID)

	return &dto.PostSignUpResponse{
		Id:       id,
		Email:    model.Email,
		Name:     model.Name,
		NickName: model.NickName,
		Phone:    model.Phone,
	}
}
