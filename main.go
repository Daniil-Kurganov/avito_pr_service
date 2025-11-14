package main

import (
	"avito_pr_service/src/db"
	"avito_pr_service/src/usecase"
	"log"
)

func main() {
	t := usecase.Team{TeamName: "trisigma"}
	err := t.Add()
	log.Print(err)
	db.CloseConnection()
}
