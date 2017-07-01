package services

import (
	"github.com/social-tournament-service/app/daos"
	"github.com/social-tournament-service/app/models"
	errors "github.com/social-tournament-service/app/http"
)

type PlayerService struct{
	playerDao *daos.PlayerDao
}

func NewPlayerService(playerDao *daos.PlayerDao) *PlayerService{
	return &PlayerService{
		playerDao: playerDao,
	}
}

func (s *PlayerService) Take(playerId string, points int) error {
	player := &models.Player{
		Id: playerId,
	}
	
	if err := s.playerDao.Get(player); err != nil{
		return err
	}
	
	if player.Points <= points{
		return errors.NewBadRequestError("Points can't be taken from player's because then his balance will be zero")
	}
	
	player.Points -= points
	
	return s.playerDao.UpdatePoints(player)
}

func (s *PlayerService) Fund(playerId string, points int) error {
	player := &models.Player{
		Id: playerId,
		Points: points,
	}
	
	if points <= 0 && !s.playerDao.Exists(player){
		return errors.NewBadRequestError("Player with zero points can't be created")
	}
	
	return s.playerDao.Upsert(player)
}

func (s *PlayerService) Balance(playerId string) (models.PlayerBalanceDto, error) {
	player := &models.Player{
		Id: playerId,
	}
	
	if err := s.playerDao.GetBalance(player); err != nil{
		return models.PlayerBalanceDto{}, err
	}
	
	return models.PlayerBalanceDto{
		Id: player.Id,
		Points: player.Points,
	}, nil
}