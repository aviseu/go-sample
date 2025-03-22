package infrastructure

import "github.com/jmoiron/sqlx"

type Config struct {
	DSN string `default:"postgres://api:pwd@localhost:5433/todo?sslmode=disable"`
}

func SetupDatabase(cfg Config) (*sqlx.DB, error) {
	return sqlx.Open("postgres", cfg.DSN)
}
