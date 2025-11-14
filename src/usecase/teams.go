package usecase

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

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
	// if err = t.validate(); err != nil {
	// 	return err
	// }
	var result pgconn.CommandTag
	if result, err = db.Connection.Exec(context.Background(), "insert into teams (name) values ($1)", t.TeamName); err != nil {
		err = fmt.Errorf("error on add new team: %w", err)
		return
	}
	conf.Logger.Debug(fmt.Sprintf("%s: number added rows: %d", conf.LogHeaders.Usecase, result.RowsAffected()))
	return
}
