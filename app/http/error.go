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
