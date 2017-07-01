package daos

import (
	"gopkg.in/pg.v5"
	"fmt"
	"errors"
	"github.com/social-tournament-service/app/config"
)

type daos struct{
	Player *PlayerDao
	Tournament *TournamentDao
	Default *DefaultDao
}

func Init(cnf *config.Config) (*daos, error){
	db, err := initDB(cnf.Db)
	if err != nil{
		return nil, err
	}
	
	defaultDao := &DefaultDao{db}
	if err := defaultDao.createDBSchema(); err != nil{
		panic(err)
	}
	
	return &daos{
		Player: &PlayerDao{db},
		Tournament: &TournamentDao{db, cnf.MinPlayerPointsAmount},
		Default: defaultDao,
	}, nil
}

func initDB(cnf config.DbConfig) (*pg.DB, error) {	
	db := pg.Connect(&pg.Options{
		Database: cnf.DbName,
		Addr: fmt.Sprintf("%s:%s", cnf.DbHost, cnf.DbPort),
		User: cnf.DbUser,
		Password: cnf.DbPswd,
	})
	if db == nil{
		return nil, errors.New("Could not establish database connection")
	}
	
	return db, nil
}