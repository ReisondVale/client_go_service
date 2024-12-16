package db

import (
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
	"project/internal/config"
)

func Connect() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", config.GetDatabaseURL())
}
