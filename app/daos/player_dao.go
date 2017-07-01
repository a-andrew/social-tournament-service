package daos

import (
	"gopkg.in/pg.v5"
	"github.com/social-tournament-service/app/models"
	"errors"
	"fmt"
)

//TODO: move it out into config
const MIN_PLAYER_POINTS_AMOUNT = 0

type PlayerDao struct {
	db *pg.DB
}

func (d *PlayerDao) Get(player *models.Player, columns ...string) error{
	model := d.db.Model(player)
	
	if len(columns) > 0{
		model.Column(columns...)
	}
	
	if err := model.Where("id = ?id").Select(); err != nil{
		if err.Error() == "pg: no rows in result set"{
			return errors.New(fmt.Sprintf("Player with ID `%s` not found", player.Id))
		}
		
		return err
	}
	
	return nil
}

func (d *PlayerDao) GetBalance(player *models.Player) error{
	return d.Get(player, "points")
}

func (d *PlayerDao) Upsert(player *models.Player) error{
	if _, err := d.db.Model(player).
		OnConflict("(id) DO UPDATE").
		Set("points = player.points + ?points").
		Insert(); err != nil{
		return err
	}

	return nil
}

func (d *PlayerDao) Update(player *models.Player, columns ...string) error{
	model := d.db.Model(player)
	
	if len(columns) > 0{
		model.Column(columns...)
	}

	res, err := model.Update()
	if err != nil {
		return err
	}
	
	if res.RowsAffected() == 0{
		return errors.New(fmt.Sprintf("Player with ID `%s` not found", player.Id))
	}

	return nil
}

func (d *PlayerDao) UpdateTxn(txn *pg.Tx, player *models.Player, columns ...string) error{
	model := txn.Model(player)
	
	if len(columns) > 0{
		model.Column(columns...)
	}

	res, err := model.Update()
	if err != nil {
		return err
	}
	
	if res.RowsAffected() == 0{
		return errors.New(fmt.Sprintf("Player with ID `%s` not found", player.Id))
	}

	return nil
}

func (d *PlayerDao) UpdatePoints(player *models.Player) error{
	return d.Update(player, "points")
}

func (d *PlayerDao) UpdatePointsTxn(txn *pg.Tx, player *models.Player) error{
	return d.UpdateTxn(txn, player, "points")
}

func (d *PlayerDao) AddPointsTxn(txn *pg.Tx, player *models.Player) error{
	res, err := txn.Model(player).Set("points = player.points + ?points").Update()
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0{
		return errors.New(fmt.Sprintf("Player with ID `%s` not found", player.Id))
	}
	
	return nil
}

func (d *PlayerDao) Exists(player *models.Player) bool{
	if err := d.db.Select(player); err != nil{
		return false
	}

	return true
}