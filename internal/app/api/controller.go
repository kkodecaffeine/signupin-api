package api

import (
	"fmt"
	"net/http"
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
	"github.com/kkodecaffeine/go-common/errorcode"
	"github.com/kkodecaffeine/go-common/middleware/token"

	"github.com/kkodecaffeine/go-common/rest"
)

type Controller struct {
	v       *validator.Validate
	usecase user.Usecase
}

// NewController returns new controller instance
func NewController(e *gin.Engine, v *validator.Validate, uc user.Usecase) Controller {
	ctrl := Controller{v, uc}

	v1 := e.Group("/v1")
	v1.POST("/auth/sms", ctrl.SendSMS)
	v1.POST("/auth/sign-up", ctrl.SignUp)
	v1.POST("/auth/sign-in", ctrl.SignIn)

	authorized := v1.Group("/")
	authorized.Use(token.JwtAuthMiddleware()).GET("/users/:userID", ctrl.GetMe)
	authorized.Use(token.JwtAuthMiddleware()).PUT("/users/reset-password", ctrl.UpdatePassword)

	return ctrl
}

/**
 * 전화번호 인증 API
 * 요청받은 전화번호 검증 수행
 * @return : authnumber (6자리 난수)
 */
func (ctrl *Controller) SendSMS(c *gin.Context) {
	response := rest.NewApiResponse()

	var req dto.PostSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		for _, element := range err.(validator.ValidationErrors) {
			if element.ActualTag() == "required" {
				response.Error(&errorcode.MISSING_PARAMETERS, fmt.Sprintf("required: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			} else {
				response.Error(&errorcode.INVALID_PARAMETERS, fmt.Sprintf("tag: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}

	// 기존에 생성된 인증번호가 있는지 없는지 확인
	authnumber, _ := ctrl.usecase.GetAuthNumber()
	if authnumber == "" {
		authnumber, _ = ctrl.usecase.UpsertAuthNumber()
	}

	result := dto.PostSMSResponse{
		AuthNumber: authnumber,
	}

	response.Succeed("", result)
	c.JSON(http.StatusOK, response)
}

/**
 * 회원 가입 API
 * 요청받은 회원 정보 검증 수행
 * 요청받은 회원 정보가 기존 가입자인지 아닌지 확인 후 신규 가입 처리
 * @return : 가입 시 생성된 회원 정보 (w/ ID)
 */
func (ctrl *Controller) SignUp(c *gin.Context) {
	response := rest.NewApiResponse()

	var req dto.PostSignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		for _, element := range err.(validator.ValidationErrors) {
			if element.ActualTag() == "required" {
				response.Error(&errorcode.MISSING_PARAMETERS, fmt.Sprintf("required: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			} else {
				if len(fmt.Sprintf("%v", element.Value())) == 0 {
					break
				}
				response.Error(&errorcode.INVALID_PARAMETERS, fmt.Sprintf("tag: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}

	if err := ctrl.v.Struct(req); err != nil {
		response.Error(&errorcode.INVALID_PARAMETERS, err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	exists, _ := ctrl.usecase.GetOne(req.Email) // 이미 가입한 회원인지 확인
	if exists != nil {
		response.Error(&errorcode.AUTH_EMAIL_ALREADY_EXISTS, req.Email, "")
		c.JSON(errorcode.AUTH_EMAIL_ALREADY_EXISTS.HttpStatusCode, response)
		return
	}

	insertedID, err := ctrl.usecase.SaveOne(&req)
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	found, err := ctrl.usecase.GetOneByID(insertedID)
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	response.Succeed("", found)
	c.JSON(http.StatusOK, response)
}

/**
 * 회원 로그인 API
 * 요청받은 회원 정보 검증 수행
 * (이메일, 비밀번호) 혹은 (전화번호, 비밀번호) 로 로그인 가능하도록 구현
 * @return : 가입 시 생성된 회원 정보 (w/ ID, JWT)
 */
func (ctrl *Controller) SignIn(c *gin.Context) {
	response := rest.NewApiResponse()

	var req dto.PostSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		for _, element := range err.(validator.ValidationErrors) {
			if element.ActualTag() == "required" {
				response.Error(&errorcode.MISSING_PARAMETERS, fmt.Sprintf("required: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			} else {
				if len(fmt.Sprintf("%v", element.Value())) == 0 {
					break
				}
				response.Error(&errorcode.INVALID_PARAMETERS, fmt.Sprintf("tag: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}

	if err := ctrl.v.Struct(req); err != nil {
		response.Error(&errorcode.INVALID_PARAMETERS, err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 && len(strings.TrimSpace(req.Phone)) == 0 {
		response.Error(&errorcode.BAD_REQUEST, "", nil)
		c.JSON(errorcode.BAD_REQUEST.HttpStatusCode, response)
		return
	}

	var identifier string
	if len(strings.TrimSpace(req.Email)) == 0 {
		identifier = req.Phone
	} else if len(strings.TrimSpace(req.Phone)) == 0 {
		identifier = req.Email
	} else {
		identifier = req.Email
	}

	found, err := ctrl.usecase.GetOne(identifier, req.Password)
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	/**
	 * 로그인 성공 후 인증번호 갱신하는 이유
	 * 비밀번호 수정 시 토큰이 만료된 경우 토큰도 갱신해야하고 이어서 본인 전화번호 인증 절차가 필요함
	 * 그런데 인증번호가 갱신되지 않은 상태에서 토큰만 갱신한채로 비밀번호 수정 API 가 성공하는 경우의 수가 발생할 여지가 있음
	 * 즉, 토큰 갱신 후에 기존에 사용했던 인증번호로 본인 전화번호 인증 절차를 누락한채 비밀번호 수정을 진행할 수 있음
	 */
	_, err = ctrl.usecase.UpsertAuthNumber()
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	response.Succeed("", found)
	c.JSON(http.StatusOK, response)
}

/**
 * 회원 정보 조회 API
 * JWT 검증 과정 후 회원 정보 조회
 * @return : 가입 시 생성된 회원 정보 (w/ ID)
 */
func (ctrl *Controller) GetMe(c *gin.Context) {
	response := rest.NewApiResponse()

	userID := c.Param("userID")
	found, err := ctrl.usecase.GetOneByID(userID)
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	response.Succeed("", found)
	c.JSON(http.StatusOK, response)
}

/**
 * 비밀번호 수정 API
 * JWT 및 요청받은 정보에 대한 검증
 * JWT 만료된 경우
 * 		- 1) 회원 로그인 API 를 호출하여 신규 토큰 획득
 * 		- 2) 이어서 전화번호 인증 API 호출하여 신규 인증번호 획득
 * 		- 3) 이어서 새로 획득한 인증번호를 요청모델에 담아서 비밀번호 수정 API 호출
 * 기존 비밀번호로 회원 정보 조회에 성공한 후 요청받은 신규 비밀번호로 비밀번호 변경
 */
func (ctrl *Controller) UpdatePassword(c *gin.Context) {
	response := rest.NewApiResponse()

	var req dto.PutPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		for _, element := range err.(validator.ValidationErrors) {
			if element.ActualTag() == "required" {
				response.Error(&errorcode.MISSING_PARAMETERS, fmt.Sprintf("required: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			} else {
				response.Error(&errorcode.INVALID_PARAMETERS, fmt.Sprintf("tag: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}

	if err := ctrl.v.Struct(req); err != nil {
		response.Error(&errorcode.INVALID_PARAMETERS, err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// "기존 비밀번호"와 "신규 비밀번호"가 동일한지 확인
	if req.Password == req.NewPassword || req.Password == req.Confirmation {
		response.Error(&errorcode.BAD_REQUEST, "same as the existing password", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 사용자가 입력한 "신규 비밀번호"와 "비밀번호 확인"이 동일한지 확인
	if req.NewPassword != req.Confirmation {
		response.Error(&errorcode.BAD_REQUEST, "password mismatch", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	found, err := ctrl.usecase.GetOne(req.Email, req.Password)
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	_, err = ctrl.usecase.UpdatePassword(req.AuthNumber, found.Id, req.NewPassword)
	if err != nil {
		response.Error(err.CodeDesc, err.Message, err.Data)
		c.JSON(err.CodeDesc.HttpStatusCode, response)
		return
	}

	response.Succeed("", nil)
	c.JSON(http.StatusOK, response)
}
