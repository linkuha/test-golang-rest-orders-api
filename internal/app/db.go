package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/linkuha/test-golang-rest-orders-api/config"
	"github.com/rs/zerolog/log"
)

func newDB(cfg *config.EnvParams) (*sql.DB, error) {
	var connStr string
	if cfg.DatabaseURL != "" {
		connStr = cfg.DatabaseURL
	} else {
		connStr = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPassword, cfg.PgDB)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Ping database...")
	if err = db.Ping(); err != nil {
		return nil, err
	}
	q := db.QueryRow("SELECT VERSION()")
	var ver string
	err = q.Scan(&ver)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("DB version: %s", ver)

	return db, nil
}
