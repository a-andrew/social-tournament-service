package models

type TournamentPlayer struct{
	TableName struct{} `sql:"tournamentplayers"`
	TournamentId string `sql:"tournamentid"`
	PlayerId string `sql:"playerid"`
	BackerIds []string `sql:"backerids" pg:",array"`
}