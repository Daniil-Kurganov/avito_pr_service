package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// newTeam := usecase.Team{ // team/add
	// 	TeamName: "games",
	// 	Members: []usecase.TeamMember{
	// 		{

	// 			UserId:   "u1",
	// 			Username: "Max",
	// 			IsActive: true,
	// 		},
	// 		{

	// 			UserId:   "u2",
	// 			Username: "Victor",
	// 			IsActive: true,
	// 		},
	// 		{

	// 			UserId:   "u3",
	// 			Username: "Wiliam",
	// 			IsActive: true,
	// 		},
	// 	},
	// }

	userActiveUpdation := struct { // users/setIsActive
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}{
		UserID:   "u2",
		IsActive: false,
	}

	// prCreation := struct { // pullRequest/create
	// 	PullRequestId   string `json:"pull_request_id"`
	// 	PullRequestName string `json:"pull_request_name"`
	// 	AuthorId        string `json:"author_id"`
	// }{
	// 	PullRequestId:   "pr-6",
	// 	PullRequestName: "New app1",
	// 	AuthorId:        "u123",
	// }

	// prMerge := struct { // pullRequest/merge
	// 	PullRequestId string `json:"pull_request_id"`
	// }{
	// 	PullRequestId: "pr-54",
	// }

	// prReassing := struct { // pullRequest/reassign
	// 	PullRequestId string `json:"pull_request_id"`
	// 	OldReviewerId string `json:"old_reviewer_id"`
	// }{
	// 	PullRequestId: "pr-1",
	// 	OldReviewerId: "u2",
	// }

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
	bodyBytes, _ := io.ReadAll(response.Body)
	log.Printf("Response body: %s", string(bodyBytes))
	response.Body.Close()
}
