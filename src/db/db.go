package db

import (
	"avito_pr_service/src/conf"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

const TeamNameLenght = 50

var Connection *pgx.Conn

func Connect() {
	connectionURL := fmt.Sprintf("postgres://%s:%s@localhost:5433/%s",
		conf.PSQLUser, conf.PSQLPassword, conf.PSQLDBName)
	var err error
	if Connection, err = pgx.Connect(context.Background(), connectionURL); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on connection to DB: %v", conf.LogHeaders.PSQL, err))
		os.Exit(conf.OSExitCode.InvalidFunction)
	}
	conf.Logger.Debug(fmt.Sprintf("%s: Successfully connected to DB", conf.LogHeaders.PSQL))
}

func CloseConnection() (err error) {
	err = Connection.Close(context.Background())
	return
}
