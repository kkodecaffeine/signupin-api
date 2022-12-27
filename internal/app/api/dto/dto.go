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

type PostSignUpResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

// 회원 로그인
type PostSignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PostSignInResponse struct {
	Id string `json:"id"`
}
