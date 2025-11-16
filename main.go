package main

import (
	"avito_pr_service/src/db"
	"avito_pr_service/src/http"
)

func main() {
	db.Connect()
	http.StartHTTPServer()
	defer db.Connection.Close()
}
