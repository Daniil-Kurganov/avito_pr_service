package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// newTeam := usecase.Team{ // team/add
	// 	TeamName: "building",
	// 	Members: []usecase.TeamMember{
	// 		{
	// 			UserId:   "u3",
	// 			Username: "Sofia",
	// 			IsActive: true,
	// 		},
	// 		{
	// 			UserId:   "u4",
	// 			Username: "Mark",
	// 			IsActive: true,
	// 		},
	// 	},
	// }

	userActiveUpdation := struct { // users/setIsActive
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}{
		UserID:   "u5",
		IsActive: true,
	}

	var newData []byte
	var err error
	if newData, err = json.Marshal(userActiveUpdation); err != nil {
		log.Fatalf("Error on marshaling data: %v", err)
	}
	var request *http.Request
	if request, err = http.NewRequest("POST", "http://127.0.0.1:8080/users/setIsActive", bytes.NewBuffer(newData)); err != nil {
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
