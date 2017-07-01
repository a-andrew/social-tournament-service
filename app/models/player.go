package models

type Player struct{
	Id     string `sql:"id,pk" json:"playerId"`
	Points int `sql:"points" json:"balance"`
}

type PlayerBalanceDto struct{
	Id     string `json:"playerId"`
	Points int `json:"balance"`
}