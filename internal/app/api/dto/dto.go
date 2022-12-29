package dto

type PostSMSRequest struct {
	Phone string `json:"phone" binding:"required,customPhone"` // 전화번호
}

type PostSMSResponse struct {
	AuthNumber string `json:"authnumber"` // 인증번호
}

// 회원 가입
type PostSignUpRequest struct {
	AuthNumber string `json:"authnumber" binding:"required" validate:"len=6"` // 인증번호
	Email      string `json:"email" binding:"required,customEmail"`           // 이메일
	NickName   string `json:"nickname" binding:"required" validate:"min=2"`   // 닉네임
	Name       string `json:"name" binding:"required" validate:"min=2"`       // 이름
	Password   string `json:"password" binding:"required" validate:"min=8"`   // 비밀번호
	Phone      string `json:"phone" binding:"required,customPhone"`           // 전화번호
}

// 회원 로그인
type PostSignInRequest struct {
	Email    string `json:"email" binding:"customEmail"` // 이메일
	Password string `json:"password" binding:"required"` // 비밀번호
	Phone    string `json:"phone" binding:"customPhone"` // 전화번호
}

// 회원 조회
type GetUserResponse struct {
	Id       string `json:"id"`       // 아이디
	Email    string `json:"email"`    // 이메일
	NickName string `json:"nickname"` // 닉네임
	Name     string `json:"name"`     // 이름
	Phone    string `json:"phone"`    // 전화번호
}

type GetUserWithTokenResponse struct {
	AccessToken string `json:"accesstoken"` // 토큰
	Id          string `json:"id"`          // 아이디
	Email       string `json:"email"`       // 이메일
	NickName    string `json:"nickname"`    // 닉네임
	Name        string `json:"name"`        // 이름
	Phone       string `json:"phone"`       // 전화번호
}

type PostSignInResponse struct {
	Id string `json:"id"`
}

// 비밀번호 수정
type PutPasswordRequest struct {
	AuthNumber   string `json:"authnumber" binding:"required" validate:"len=6"`   // 인증번호
	Email        string `json:"email" binding:"required,customEmail"`             // 이메일
	Password     string `json:"password" binding:"required" validate:"min=8"`     // 비밀번호
	NewPassword  string `json:"newpassword" binding:"required" validate:"min=8"`  // 신규 비밀번호
	Confirmation string `json:"confirmation" binding:"required" validate:"min=8"` // 신규 비밀번호 확인
}
