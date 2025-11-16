package conf

import (
	"log/slog"
	"os"
)

const (
	ServerHTTPServeSocket = "0.0.0.0:8080"

	PRReviewersMax = 2
)

var (
	Logger *slog.Logger

	OSExitCode = struct {
		InvalidFunction int
		InvalidHandle   int
	}{
		InvalidFunction: 1,
		InvalidHandle:   6,
	}
	LogHeaders = struct {
		HTTPServer string
		PSQL       string
		Usecase    string
	}{
		HTTPServer: "[HTTP server]",
		PSQL:       "[PostgreSQL]",
		Usecase:    "[Usecase]",
	}

	PSQLUser, PSQLPassword, PSQLDBName = os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASSWORD"), os.Getenv("PSQL_DB_NAME")
)

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug.Level()}))
}
