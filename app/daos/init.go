package daos

import (
	"github.com/caarlos0/env"
	"gopkg.in/pg.v5"
	"fmt"
	"errors"
)

type dbConfig struct{
	DbName string `env:"DB_NAME"`
	DbUser string `env:"DB_USER"`
	DbPswd string `env:"DB_PASSWORD"`
	DbHost string `env:"DB_HOST"`
	DbPort string `env:"DB_PORT"`
}

type daos struct{
	Player *PlayerDao
	Tournament *TournamentDao
	Default *DefaultDao
}

func Init() (*daos, error){
	db, err := initDB()
	if err != nil{
		return nil, err
	}
	
	defaultDao := &DefaultDao{db}
	if err := defaultDao.createDBSchema(); err != nil{
		panic(err)
	}
	
	return &daos{
		Player: &PlayerDao{db},
		Tournament: &TournamentDao{db},
		Default: defaultDao,
	}, nil
}

func initDB() (*pg.DB, error) {
	dbCfg := dbConfig{}
	if err := env.Parse(&dbCfg); err != nil{
		return nil, err
	}
	
	db := pg.Connect(&pg.Options{
		Database: dbCfg.DbName,
		Addr: fmt.Sprintf("%s:%s", dbCfg.DbHost, dbCfg.DbPort),
		User: dbCfg.DbUser,
		Password: dbCfg.DbPswd,
	})
	if db == nil{
		return nil, errors.New("Could not establish database connection")
	}
	
	return db, nil
}