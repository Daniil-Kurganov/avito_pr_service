package usecase

import "fmt"

type (
	TeamMember struct {
		UserId   string `json:"user_id"`
		Username string `json:"username"`
		IsActive bool   `json:"is_active"`
	}
	Team struct {
		TeamName string       `json:"team_name" binding:"required"`
		Members  []TeamMember `json:"members"`
	}
)

func (t *Team) validate() error {
	if t.TeamName == "" {
		return fmt.Errorf(`invalid "team_name" field: must be non-empty`)
	}
	if len(t.Members) == 0 {
		return fmt.Errorf(`''members'' filed must be non-empty`)
	}
	return nil
}

func (t *Team) Add() (err error) {
	if err = t.validate(); err != nil {
		return err
	}
	// DB work
	return
}
