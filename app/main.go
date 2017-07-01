package main

import (
	d "github.com/social-tournament-service/app/daos"
	h "github.com/social-tournament-service/app/http"
	"github.com/social-tournament-service/app/handlers"
	"github.com/social-tournament-service/app/services"
	c "github.com/social-tournament-service/app/config"
	"log"
)

func main(){
	config := c.GetConfig()
	
	router := h.NewRouter()
	daos, err := d.Init(config)
	if err != nil{
		panic(err)
	}
	
	handlers.InitPlayerHandler(router, services.NewPlayerService(daos.Player))
	handlers.InitTournamentHandler(router, services.NewTournamentService(daos.Tournament, daos.Player))
	handlers.InitDefaultHandler(router, services.NewDefaultService(daos.Default))
	
	log.Fatal(router.Serve(config.Port))
}