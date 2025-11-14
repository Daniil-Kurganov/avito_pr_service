package main

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/db"
	"avito_pr_service/src/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db.Connect()
	listener := http.StartHTTPServer()
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	acceptedSignal := <-osSignals
	close(osSignals)
	conf.Logger.Info(fmt.Sprintf("Accepted OS %s signal", acceptedSignal.String()))
	if err := listener.Close(); err != nil {
		conf.Logger.Error(fmt.Sprintf("Error on closing HTTP-server: %v", err))
	}
	if err := db.CloseConnection(); err != nil {
		conf.Logger.Error(fmt.Sprintf("Error on closing DB connection: %v", err))
	}
}
