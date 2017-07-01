package models

type TournamentResultJson struct{
	TournamentId string `json:"tournamentId"`
	Winners []TournamentWinnersJson `json:"winners"`
}

type TournamentWinnersJson struct{
	PlayerId string `json:"playerId"`
	Prize int `json:"prize"`
}