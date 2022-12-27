package dto

type PostSMSRequest struct {
	Phone string `json:"phone" binding:"required,customPhone"`
}

type PostSMSResponse struct {
	AuthNumber string `json:"authnumber"`
}

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
