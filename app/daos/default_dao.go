package daos

import (
	"gopkg.in/pg.v5"
	"fmt"
)

const DB_SCHEMA = "public"

const DB_PLAYERS_TABLE_SQL = `
	CREATE TABLE players(
	    id varchar(32) not null primary key,
	    points int4
	)`

const DB_TOURNAMENTS_TABLE_SQL = `
	CREATE TABLE tournaments(
	    id int4 not null primary key,
	    deposit int4,
	    finished boolean default false not null
	)`

const DB_TOURNAMENTPLAYERS_TABLE_SQL = `
	CREATE TABLE tournamentplayers(
	    playerid varchar(32) not null,
	    tournamentid int4 not null,
	    backerIds varchar(32) [],
	    CONSTRAINT "playerid_tournamentid" UNIQUE( "playerid", "tournamentid" )
	)`

type DefaultDao struct {
	db *pg.DB
}

func (d *DefaultDao) createDBSchema() error{
	var needMigration bool
	
	_, err := d.db.Query(&needMigration, `SELECT count(*) != 3 as exists FROM pg_class WHERE relkind = 'r' and relname = 'players' or relname = 'tournaments' or relname = 'tournamentplayers'`)
	if err != nil{
		return err
	}
	
	if needMigration{
		return d.Reset()
	}
	
	return nil
}

func (d *DefaultDao) Reset() error{
	return d.db.RunInTransaction(func(txn *pg.Tx) error{
		if _, err := txn.Exec(fmt.Sprintf("drop schema if exists %s cascade; create schema %s", DB_SCHEMA, DB_SCHEMA)); err != nil{
			return err
		}

		if _, err := txn.Exec(fmt.Sprintf("%s;%s;%s", DB_PLAYERS_TABLE_SQL, DB_TOURNAMENTS_TABLE_SQL, DB_TOURNAMENTPLAYERS_TABLE_SQL)); err != nil{
			return err
		}
		
		return nil
	})
}