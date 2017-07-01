package handlers

import (
	"github.com/social-tournament-service/app/http"
	"github.com/social-tournament-service/app/services"
	"strconv"
	codes "net/http"
)

type PlayerHandler struct{
	service *services.PlayerService
}

func InitPlayerHandler(router *http.Router, service *services.PlayerService) {
	h := PlayerHandler{service}

	router.Get("/take", h.take, "playerId", "points")
	router.Get("/fund", h.fund, "playerId", "points")
	router.Get("/balance", h.balance, "playerId")
}

func (h *PlayerHandler) take(ctx *http.Context) (int, interface{}){
	params := ctx.GetUrlParams()
	playerId, pointsStr := params.Get("playerId"), params.Get("points")
	
	points, err := strconv.Atoi(pointsStr)
	if err != nil{
		return http.ResError(http.NewBadRequestError("Could not convert `points` to int"))
	}
	
	if err := h.service.Take(playerId, points); err != nil{
		return http.ResError(err)
	}

	return codes.StatusOK, struct{}{}
}

func (h *PlayerHandler) fund(ctx *http.Context) (int, interface{}){
	params := ctx.GetUrlParams()
	playerId, pointsStr := params.Get("playerId"), params.Get("points")
	
	points, err := strconv.Atoi(pointsStr)
	if err != nil{
		return http.ResError(http.NewBadRequestError("Could not convert `points` to int"))
	}
	
	if err := h.service.Fund(playerId, points); err != nil{
		return http.ResError(err)
	}
	
	return codes.StatusOK, struct{}{}
}

func (h *PlayerHandler) balance(ctx *http.Context) (int, interface{}){
	balance, err := h.service.Balance(ctx.GetUrlParams().Get("playerId"))
	if err != nil{
		return http.ResError(err)
	}
	
	return codes.StatusOK, balance
}