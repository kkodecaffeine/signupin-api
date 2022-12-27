package user

import (
	"os"
	"signupin-api/internal/app/api/dto"
	"strings"

	"github.com/kamva/mgm/v3"
)

// User is
type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`       // 이메일
	Name             string `json:"name" bson:"name"`         // 이름
	NickName         string `json:"nickname" bson:"nickname"` // 닉네임
	Password         string `json:"password" bson:"password"` // 비밀번호
	Phone            string `json:"phone" bson:"phone"`       // 전화번혼
}

func newUser(req *dto.PostSignUpRequest) *User {
	result := compareAuthNumber(req.AuthNumber)
	if !result {
		return nil
	}

	return &User{
		Email:    req.Email,
		Name:     req.Name,
		NickName: req.NickName,
		Password: req.Password,
		Phone:    req.Phone,
	}
}

func compareAuthNumber(authnumber string) bool {
	return strings.EqualFold(authnumber, os.Getenv("AUTH_NUMBER"))
}
