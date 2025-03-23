package dbrepo

import (
	"database/sql"

	"github.com/aparkinlot/Bookings/internal/config"
	"github.com/aparkinlot/Bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
