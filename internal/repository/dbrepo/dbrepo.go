package dbrepo

import (
	"database/sql"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//NewPostgresRepo will return postresdb
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

//NewTestingRepo will return testingdb
func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
