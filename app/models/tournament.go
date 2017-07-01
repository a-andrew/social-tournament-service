package models

type Tournament struct{
	Id string `sql:"id,pk"`
	Deposit int `sql:"deposit"`
	Finished bool `sql:"finished"`
}