package main

import (
	"avito_pr_service/src/usecase"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	newTeam := usecase.Team{
		TeamName: "payment",
		Members: []usecase.TeamMember{
			{
				UserId:   "u1",
				Username: "Aleksey",
				IsActive: true,
			},
			{
				UserId:   "u2",
				Username: "Julia",
				IsActive: true,
			},
		},
	}
	var newTeamJSONData []byte
	var err error
	if newTeamJSONData, err = json.Marshal(newTeam); err != nil {
		log.Fatalf("Error on marshaling team data: %v", err)
	}
	var request *http.Request
	if request, err = http.NewRequest("POST", "http://127.0.0.1:8080/team/add", bytes.NewBuffer(newTeamJSONData)); err != nil {
		log.Fatalf("Error on creation request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	var response *http.Response
	if response, err = client.Do(request); err != nil {
		log.Fatalf("Error on sending request: %v", err)
	}
	log.Printf("Response: status - %s", response.Status)
	response.Body.Close()
}
