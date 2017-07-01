package daos

import (
	"gopkg.in/pg.v5"
	"github.com/social-tournament-service/app/models"
	errors "github.com/social-tournament-service/app/http"
	"strings"
	"fmt"
)

type TournamentDao struct {
	db *pg.DB
	MinPlayerPointsAmount int
}

func (d *TournamentDao) Create(tournament *models.Tournament) error{
	err := d.db.Insert(tournament)
	if err != nil{
		if strings.Contains(err.Error(), "duplicate key value"){
			return errors.NewBadRequestError(fmt.Sprintf("Tournament with ID `%s` already exists", tournament.Id))
		}
		
		return errors.NewInternalError(err.Error())
	}
	
	return nil
}

func (d *TournamentDao) Get(tournament *models.Tournament, columns ...string) error{
	model := d.db.Model(tournament)

	if len(columns) > 0{
		model.Column(columns...)
	}

	if err := model.Where("id = ?id").Select(); err != nil{
		if err.Error() == "pg: no rows in result set"{
			return errors.NewNotFoundError(fmt.Sprintf("Tournament with ID `%s` not found", tournament.Id))
		}

		return errors.NewInternalError(err.Error())
	}

	return nil
}

func (d *TournamentDao) GetDeposit(tournament *models.Tournament) error{
	return d.Get(tournament, "deposit")
}

func (d *TournamentDao) GetData(tournamentPlayer *models.TournamentPlayer, columns ...string) error{
	model := d.db.Model(tournamentPlayer)

	if len(columns) > 0{
		model.Column(columns...)
	}

	if err := model.Where("tournamentid = ?tournamentid and playerid = ?playerid").Select(); err != nil{
		if err.Error() == "pg: no rows in result set"{
			return errors.NewNotFoundError(fmt.Sprintf("Tournament with ID `%s` is not started", tournamentPlayer.TournamentId))
		}

		return errors.NewInternalError(err.Error())
	}

	return nil
}

func (d *TournamentDao) GetBackers(tournamentPlayer *models.TournamentPlayer) error{
	return d.GetData(tournamentPlayer, "backerids")
}

func (d *TournamentDao) Exists(tournament *models.Tournament) bool{
	if err := d.db.Select(tournament); err != nil{
		return false
	}
	
	return true
}

func (d *TournamentDao) ExistsAndNotFinished(tournament *models.Tournament) bool{
	if err := d.db.Model(tournament).
		Where("id = ?id and finished is not true").
		Select(); err != nil{
		return false
	}
	
	return true
}

func (d *TournamentDao) IsStarted(tournamentPlayer *models.TournamentPlayer) bool{
	if err := d.db.Model(tournamentPlayer).
		Where("tournamentid = ?tournamentid").
		Limit(1).
		Select(); err != nil{
		return false
	}
	
	return true
}

func (d *TournamentDao) IsJoinedAsPlayer(tournamentPlayer *models.TournamentPlayer) bool{
	if err := d.db.Model(tournamentPlayer).
		Where("tournamentid = ?tournamentid and playerid = ?playerid").
		Select(); err != nil{
		return false
	}
	
	return true
}

func (d *TournamentDao) IsJoinedAsBacker(tournamentPlayer *models.TournamentPlayer) bool{
	if err := d.db.Model(tournamentPlayer).
		Where("tournamentid = ?tournamentid and ?playerid = any(tournament_player.backerids)").
		Select(); err != nil{
		return false
	}
	
	return true
}

func (d *TournamentDao) CreateTournamentPlayer(tournamentPlayer *models.TournamentPlayer) error{
	err := d.db.Insert(tournamentPlayer)
	if err != nil{
		return errors.NewInternalError(err.Error())
	}

	return nil
}

func (d *TournamentDao) CreateTournamentPlayerTxn(txn *pg.Tx, tournamentPlayer *models.TournamentPlayer) error{
	if err := txn.Insert(tournamentPlayer); err != nil{
		return errors.NewInternalError(err.Error())
	}
	
	return nil
}

func (d *TournamentDao) FinishTxn(txn *pg.Tx, tournament *models.Tournament) error{
	res, err := txn.Model(tournament).Column("finished").Update()
	if err != nil {
		return errors.NewInternalError(err.Error())
	}

	if res.RowsAffected() == 0{
		return errors.NewNotFoundError(fmt.Sprintf("Tournament with ID `%s` not found", tournament.Id))
	}
	
	return nil
}

func (d *TournamentDao) Transaction(do func(txn *pg.Tx) error) error{
	return d.db.RunInTransaction(do)
}