package api

import (
	"net/http"
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"

	"github.com/gin-gonic/gin"
	"github.com/kkodecaffeine/go-common/errorcode"
	"github.com/kkodecaffeine/go-common/rest"
	"gopkg.in/validator.v2"
)

type Controller struct {
	usecase user.Usecase
}

// NewController returns new controller instance
func NewController(e *gin.Engine, uc user.Usecase) Controller {
	ctrl := Controller{uc}

	v1 := e.Group("/v1")
	v1.POST("/auth/sign-up", ctrl.SignUp)

	return ctrl
}

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
