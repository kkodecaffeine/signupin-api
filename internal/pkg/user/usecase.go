package user

import (
	"signupin-api/internal/app/api/dto"

	"github.com/kkodecaffeine/go-common/core/database/mongo/errortype"
	"github.com/kkodecaffeine/go-common/errorcode"
	"github.com/kkodecaffeine/go-common/rest"
)

// UseCase interface definition
type Usecase interface {
	SaveOne(request *dto.PostSignUpRequest) (string, error)

	// GET
	GetOne(email string) *rest.ApiResponse
	GetOneByID(ID string) *rest.ApiResponse
}

type usecase struct {
	repo Repository
}

func (u *usecase) SaveOne(req *dto.PostSignUpRequest) (string, error) {
	user := newUser(req)

	return u.repo.SaveOne(user)
}

func (u *usecase) GetOne(email string) *rest.ApiResponse {
	response := rest.NewApiResponse()

	result, err := u.repo.GetOne(email)
	if err != nil {
		if errortype.IsDecodeError(err) {
			return response.Error(&errorcode.FAILED_DB_PROCESSING, err.Error(), nil)
		} else if errortype.IsNotFoundErr(err) {
			return response.Error(&errorcode.NOT_FOUND_ERROR, err.Error(), nil)
		} else {
			return response.Error(&errorcode.FAILED_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	if result == nil {
		response.Code = errorcode.NOT_FOUND_ERROR.Code
	} else {
		response.Succeed("", result)
	}

	return response
}

func (u *usecase) GetOneByID(ID string) *rest.ApiResponse {
	response := rest.NewApiResponse()

	result, err := u.repo.GetOneByID(ID)
	if err != nil {
		if errortype.IsDecodeError(err) {
			return response.Error(&errorcode.FAILED_DB_PROCESSING, err.Error(), nil)
		} else if errortype.IsNotFoundErr(err) {
			return response.Error(&errorcode.NOT_FOUND_ERROR, err.Error(), nil)
		} else {
			return response.Error(&errorcode.FAILED_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	if result == nil {
		response.Code = errorcode.NOT_FOUND_ERROR.Code
	} else {
		response.Succeed("", result)
	}

	return response
}

// NewUsecase returns new Usecase implementation
func NewUsecase() Usecase {
	return &usecase{}
}

var _ Usecase = &usecase{}
