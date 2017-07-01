package config

import "github.com/caarlos0/env"

type Config struct {
	Db                    DbConfig
	MinPlayerPointsAmount int `env:"MIN_PLAYER_POINTS_AMOUNT" envDefault:"50"`
	Port                  int `env:"PORT" envDefault:"8080"`
}

type DbConfig struct {
	DbName string `env:"DB_NAME"`
	DbUser string `env:"DB_USER"`
	DbPswd string `env:"DB_PASSWORD"`
	DbHost string `env:"DB_HOST"`
	DbPort string `env:"DB_PORT"`
}

func GetConfig() *Config {
	cnf := &Config{}

	if err := env.Parse(cnf); err != nil {
		panic(err)
	}

	if err := env.Parse(&cnf.Db); err != nil {
		panic(err)
	}

	return cnf
}