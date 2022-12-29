package user

import (
	"fmt"
	"signupin-api/internal/app/api/dto"
	"strings"
	"time"

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

type AuthNumber struct {
	mgm.DefaultModel `bson:",inline"`
	AuthNumber       string `json:"authnumber" bson:"authnumber"` // 인증번호
}

func newUser(req *dto.PostSignUpRequest, authnumber string) *User {
	result := compareAuthNumber(req.AuthNumber, authnumber)
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

func newAuthNumber() *AuthNumber {
	return &AuthNumber{
		AuthNumber: fmt.Sprint(time.Now().Nanosecond())[:6],
	}
}

func compareAuthNumber(request, authnumber string) bool {
	return strings.EqualFold(request, authnumber)
}
