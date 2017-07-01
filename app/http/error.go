package http

import (
	codes "net/http"
)

type ResErrorJson struct{
	Msg string `json:"error"`
}

type Error struct{
	Code int
	Msg string
}

func (e Error) Error() string{
	return e.Msg
}

func NewNotFoundError(msg string) Error {
	return Error{codes.StatusNotFound, msg}
}

func NewInternalError(msg string) Error {
	return Error{codes.StatusInternalServerError, msg}
}

func NewBadRequestError(msg string) Error {
	return Error{codes.StatusBadRequest, msg}
}

func ResError(err error) (int, ResErrorJson){
	er := err.(Error)
	return er.Code, ResErrorJson{Msg: er.Msg}
}

/*
func NewNotFoundOrInternalError(err string) (int, ResError){
	if strings.Contains(err, NOT_FOUND_SUBSTRING){
		return codes.StatusNotFound, ResError{
			Msg: err,
		}
	}

	return NewInternalError(err)
}

func NewBadRequestError1(err string) (int, ResError){
	return codes.StatusBadRequest, ResError{
		Msg: err,
	}
}

func NewInternalError1(err string) (int, ResError){
	return codes.StatusInternalServerError, ResError{
		Msg: err,
	}
}*/
