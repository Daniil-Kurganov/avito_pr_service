package db

import (
	"avito_pr_service/src/conf"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TeamNameLenght = 50

var Connection *pgxpool.Pool

func Connect() {
	connectionURL := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable",
		conf.PSQLUser, conf.PSQLPassword, conf.PSQLDBName)
	var err error
	if Connection, err = pgxpool.New(context.Background(), connectionURL); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on connection to DB: %v", conf.LogHeaders.PSQL, err))
		os.Exit(conf.OSExitCode.InvalidFunction)
	}
	conf.Logger.Debug(fmt.Sprintf("%s: Successfully connected to DB", conf.LogHeaders.PSQL))
}
