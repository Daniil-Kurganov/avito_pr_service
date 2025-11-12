package conf

import (
	"log/slog"
	"os"
)

const ServerHTTPServeSocket = "127.0.0.1:8080"

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug.Level()}))
}
