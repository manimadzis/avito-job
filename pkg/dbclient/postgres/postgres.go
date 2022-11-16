package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func New(config Config) (conn *sqlx.DB, err error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err = sqlx.Connect("postgres", connString)
	if err != nil {
		return
	}

	if err = conn.Ping(); err != nil {
		return
	}

	return
}
