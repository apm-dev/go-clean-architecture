package di

import (
	"github.com/apm-dev/go-clean-architecture/core/configs"
	"github.com/go-pg/pg/v10"
)

var (
	dbPostgres *pg.DB
)

func providePostgresDB() *pg.DB {
	if dbPostgres == nil {
		dbPostgres = pg.Connect(&pg.Options{
			Addr:     ":" + configs.Env("POSTGRES_PORT"),
			User:     configs.Env("POSTGRES_USER"),
			Password: configs.Env("POSTGRES_PASSWORD"),
			Database: configs.Env("POSTGRES_DB"),
		})
	}
	return dbPostgres
}
