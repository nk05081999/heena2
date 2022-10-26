package dbrepo

import (
	"database/sql"

	"github.com/nk05081999/Reservation87/internal/config"
	"github.com/nk05081999/Reservation87/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
