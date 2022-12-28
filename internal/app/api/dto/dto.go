package dto

type PostSMSRequest struct {
	Phone string `json:"phone" binding:"required,customPhone"`
}

type PostSMSResponse struct {
	AuthNumber string `json:"authnumber"`
}

// 회원 가입
type PostSignUpRequest struct {
	AuthNumber string `json:"authnumber" binding:"required" validate:"len=6"`
	Email      string `json:"email" binding:"required"`
	NickName   string `json:"nickname" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Phone      string `json:"phone" binding:"required,customPhone"`
}

// 회원 로그인
type PostSignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 회원 조회
type GetUserResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type GetUserWithTokenResponse struct {
	AccessToken string `json:"accesstoken"`
	Id          string `json:"id"`
	Email       string `json:"email"`
	NickName    string `json:"nickname"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
}

type PostSignInResponse struct {
	Id string `json:"id"`
}

// 비밀번호 수정
type PutPasswordRequest struct {
	AuthNumber   string `json:"authnumber" binding:"required" validate:"len=6"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	NewPassword  string `json:"newpassword" binding:"required"`
	Confirmation string `json:"confirmation" binding:"required"`
}
