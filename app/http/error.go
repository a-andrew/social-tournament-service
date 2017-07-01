package http

import (
	"strings"
	codes "net/http"
)

const NOT_FOUND_SUBSTRING = "not found"

type ResError struct{
	Msg string `json:"error"`
}

func (e *ResError) Error() string{
	return e.Msg
}

func NewNotFoundOrInternalError(err string) (int, ResError){
	if strings.Contains(err, NOT_FOUND_SUBSTRING){
		return codes.StatusNotFound, ResError{
			Msg: err,
		}
	}

	return NewInternalError(err)
}

func NewBadRequestError(err string) (int, ResError){
	return codes.StatusBadRequest, ResError{
		Msg: err,
	}
}

func NewInternalError(err string) (int, ResError){
	return codes.StatusInternalServerError, ResError{
		Msg: err,
	}
}