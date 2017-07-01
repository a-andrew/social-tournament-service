package handlers

import (
	"github.com/social-tournament-service/app/http"
	"github.com/social-tournament-service/app/services"
	codes "net/http"
)

type DefaultHandler struct{
	service *services.DefaultService
}

func InitDefaultHandler(router *http.Router, service *services.DefaultService) {
	handler := DefaultHandler{service}
	
	router.Get("/reset", handler.reset)
}

func (h *DefaultHandler) reset(ctx *http.Context) (int, interface{}){
	if err := h.service.Reset(); err != nil{
		return codes.StatusInternalServerError, err.Error()
	}
	
	return codes.StatusOK, struct{}{}
}