package user

import (
	"signupin-api/internal/app/api/dto"

	"github.com/kkodecaffeine/go-common/core/database/mongo/errortype"
	"github.com/kkodecaffeine/go-common/errorcode"
	"github.com/kkodecaffeine/go-common/rest"
	"github.com/kkodecaffeine/go-common/utils"
)

// UseCase interface definition
type Usecase interface {
	SaveOne(req *dto.PostSignUpRequest) (string, *rest.CustomError)

	// GET
	GetOne(email string, password ...string) (*dto.GetUserWithTokenResponse, *rest.CustomError)
	GetOneByID(ID string) (*dto.GetUserResponse, *rest.CustomError)

	// UPDATE
	UpdatePassword(authnumber, ID, newpassword string) (*dto.GetUserResponse, *rest.CustomError)
}

type usecase struct {
	repo Repository
}

func (u *usecase) SaveOne(req *dto.PostSignUpRequest) (string, *rest.CustomError) {
	user := newUser(req)

	if user == nil {
		return "", &rest.CustomError{CodeDesc: &errorcode.BAD_REQUEST, Message: "auth number mismatch"}
	}

	insertedID, err := u.repo.SaveOne(user)
	if err != nil {
		if errortype.IsDecodeError(err) {
			return "", &rest.CustomError{CodeDesc: &errorcode.FAILED_DB_PROCESSING, Message: err.Error()}
		} else if errortype.IsNotFoundErr(err) {
			return "", &rest.CustomError{CodeDesc: &errorcode.NOT_FOUND_ERROR, Message: err.Error()}
		} else {
			return "", &rest.CustomError{CodeDesc: &errorcode.FAILED_INTERNAL_ERROR, Message: err.Error()}
		}
	}
	return insertedID, nil
}

func (u *usecase) GetOne(email string, password ...string) (*dto.GetUserWithTokenResponse, *rest.CustomError) {
	var response *dto.GetUserWithTokenResponse
	var err error

	if len(password) == 0 {
		response, err = u.repo.GetOne(email)
	} else {
		response, err = u.repo.GetOne(email, password[0])
	}

	if err != nil {
		if errortype.IsDecodeError(err) {
			return response, &rest.CustomError{CodeDesc: &errorcode.FAILED_DB_PROCESSING, Message: err.Error()}
		} else if errortype.IsNotFoundErr(err) {
			return response, &rest.CustomError{CodeDesc: &errorcode.NOT_FOUND_ERROR, Message: err.Error()}
		} else {
			return response, &rest.CustomError{CodeDesc: &errorcode.FAILED_INTERNAL_ERROR, Message: err.Error()}
		}
	}

	if response == nil {
		return response, &rest.CustomError{CodeDesc: &errorcode.NOT_FOUND_ERROR, Message: ""}
	} else {
		token, err := utils.GenerateToken(response.Id)

		if err != nil {
			return response, &rest.CustomError{CodeDesc: &errorcode.ACCESS_DENIED, Message: err.Error()}
		}
		response.AccessToken = token
	}

	return response, nil
}

func (u *usecase) GetOneByID(ID string) (*dto.GetUserResponse, *rest.CustomError) {
	response, err := u.repo.GetOneByID(ID)
	if err != nil {
		if errortype.IsDecodeError(err) {
			return response, &rest.CustomError{CodeDesc: &errorcode.FAILED_DB_PROCESSING, Message: err.Error()}
		} else if errortype.IsNotFoundErr(err) {
			return response, &rest.CustomError{CodeDesc: &errorcode.NOT_FOUND_ERROR, Message: err.Error()}
		} else {
			return response, &rest.CustomError{CodeDesc: &errorcode.FAILED_INTERNAL_ERROR, Message: err.Error()}
		}
	}

	if response == nil {
		return response, &rest.CustomError{CodeDesc: &errorcode.NOT_FOUND_ERROR, Message: ""}
	} else {
		if err != nil {
			return response, &rest.CustomError{CodeDesc: &errorcode.ACCESS_DENIED, Message: err.Error()}
		}
	}

	return response, nil
}

func (u *usecase) UpdatePassword(authnumber, ID, newpassword string) (*dto.GetUserResponse, *rest.CustomError) {
	if !compareAuthNumber(authnumber) {
		return nil, &rest.CustomError{CodeDesc: &errorcode.BAD_REQUEST, Message: "auth number mismatch"}
	}

	objectID, _ := utils.MapToObjectID(ID)

	response, err := u.repo.UpdatePassword(objectID, newpassword)
	if err != nil {
		if errortype.IsDecodeError(err) {
			return response, &rest.CustomError{CodeDesc: &errorcode.FAILED_DB_PROCESSING, Message: err.Error()}
		} else if errortype.IsNotFoundErr(err) {
			return response, &rest.CustomError{CodeDesc: &errorcode.NOT_FOUND_ERROR, Message: err.Error()}
		} else {
			return response, &rest.CustomError{CodeDesc: &errorcode.FAILED_INTERNAL_ERROR, Message: err.Error()}
		}
	}
	return response, nil
}

// NewUsecase returns new Usecase implementation
func NewUsecase(userRepo Repository) Usecase {
	return &usecase{repo: userRepo}
}

var _ Usecase = &usecase{}
