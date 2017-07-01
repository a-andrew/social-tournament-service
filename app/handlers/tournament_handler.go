package handlers

import (
	"github.com/social-tournament-service/app/services"
	"github.com/social-tournament-service/app/http"
	"strconv"
	codes "net/http"
	"github.com/social-tournament-service/app/models"
)

type TournamentHandler struct{
	service *services.TournamentService
}

func InitTournamentHandler(router *http.Router, service *services.TournamentService) {
	h := TournamentHandler{service}
	
	router.Get("/announceTournament", h.announce, "tournamentId", "deposit")
	router.Get("/joinTournament", h.join, "tournamentId", "playerId")
	router.Post("/resultTournament", h.result)
}

func (h *TournamentHandler) announce(ctx *http.Context) (int, interface{}){
	params := ctx.GetUrlParams()
	tournamentId, depositStr := params.Get("tournamentId"), params.Get("deposit")

	deposit, err := strconv.Atoi(depositStr)
	if err != nil{
		return http.NewBadRequestError("Could not convert `deposit` to int")
	}

	if deposit < 0{
		return http.NewBadRequestError("`deposit` must be greater than or equal zero")
	}
	
	if err := h.service.Announce(tournamentId, deposit); err != nil{
		return http.NewInternalError(err.Error())
	}

	return codes.StatusOK, struct{}{}
}

func (h *TournamentHandler) join(ctx *http.Context) (int, interface{}){
	params := ctx.GetUrlParams()
	
	backerIds := []string{}
	if b, ok := params["backerId"]; ok{
		backerIds = b
	}
	
	if err := h.service.Join(params.Get("tournamentId"), params.Get("playerId"), backerIds); err != nil{
		return http.NewInternalError(err.Error())
	}

	return codes.StatusOK, struct{}{}
}

func (h *TournamentHandler) result(ctx *http.Context) (int, interface{}){
	body := &models.TournamentResultJson{}
	ctx.ParseBody(body)
	
	if err := h.service.Result(body); err != nil{
		return http.NewInternalError(err.Error())
	}

	return codes.StatusOK, body
}