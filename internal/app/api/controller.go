package api

import (
	"fmt"
	"net/http"
	"os"
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"

	"gopkg.in/validator.v2"

	"github.com/gin-gonic/gin"
	v10 "github.com/go-playground/validator/v10"
	"github.com/kkodecaffeine/go-common/errorcode"
	"github.com/kkodecaffeine/go-common/rest"
)

type Controller struct {
	usecase user.Usecase
}

// NewController returns new controller instance
func NewController(e *gin.Engine, uc user.Usecase) Controller {
	ctrl := Controller{uc}

	v1 := e.Group("/v1")
	v1.POST("/auth/sms", ctrl.SendSMS)
	v1.POST("/auth/sign-up", ctrl.SignUp)

	return ctrl
}

// 전화 번호 인증 API
func (ctrl *Controller) SendSMS(c *gin.Context) {
	response := rest.NewApiResponse()

	var req dto.PostSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		for _, element := range err.(v10.ValidationErrors) {
			if element.ActualTag() == "required" {
				response.Error(&errorcode.MISSING_PARAMETERS, fmt.Sprintf("required: %s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			} else {
				response.Error(&errorcode.INVALID_PARAMETERS, fmt.Sprintf("%s", element.Field()), nil)
				c.JSON(http.StatusBadRequest, response)
				return
			}
		}
	}

	var res dto.PostSMSResponse
	res.AuthNumber = os.Getenv("AUTH_NUMBER")

	response.Succeed("", res)
	c.JSON(http.StatusOK, response)
}

// 회원 가입 API
func (ctrl *Controller) SignUp(c *gin.Context) {
	response := rest.NewApiResponse()

	var req dto.PostSignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(&errorcode.MISSING_PARAMETERS, err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(&errorcode.INVALID_PARAMETERS, err.Error(), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	found := ctrl.usecase.GetOne(req.Email)
	if found.Code == errorcode.SUCCESS.Code {
		response.Error(&errorcode.AUTH_EMAIL_ALREADY_EXISTS, "", nil)
		c.JSON(errorcode.AUTH_EMAIL_ALREADY_EXISTS.HttpStatusCode, response)
		return
	}

	insertedID, _ := ctrl.usecase.SaveOne(&req)
	found = ctrl.usecase.GetOneByID(insertedID)

	if found.Code == errorcode.SUCCESS.Code {
		response.Error(&errorcode.AUTH_EMAIL_ALREADY_EXISTS, "", nil)
		c.JSON(errorcode.AUTH_EMAIL_ALREADY_EXISTS.HttpStatusCode, response)
		return
	}

	response.Created("", found)
	c.JSON(http.StatusOK, response)
}
