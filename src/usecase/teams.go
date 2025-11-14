package usecase

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/db"
	"context"
	"errors"
	"fmt"
)

type (
	TeamMember struct {
		UserId   string `json:"user_id" binding:"required"`
		Username string `json:"username" binding:"required"`
		IsActive bool   `json:"is_active" binding:"required"`
	}
	Team struct {
		TeamName string       `json:"team_name" binding:"required"`
		Members  []TeamMember `json:"members" binding:"required"`
	}

	customError error
)

var ErrorTeamDuplication customError = errors.New("ОШИБКА: повторяющееся значение ключа нарушает ограничение уникальности \"teams_name_key\" (SQLSTATE 23505)")

func (tm *TeamMember) add(teamId int64) (err error) {
	if _, err = db.Connection.Exec(context.Background(),
		"insert into users (user_id, username, team_id, is_active) values ($1, $2, $3, $4)",
		tm.UserId, tm.Username, teamId, tm.IsActive); err != nil {
		err = fmt.Errorf("error on inserting new user: %w", err)
	}
	return
}

func (t *Team) Add() (err error) {
	var teamId int64
	if err = db.Connection.QueryRow(context.Background(), "insert into teams (name) values ($1) returning id", t.TeamName).Scan(&teamId); err != nil {
		return
	}
	for _, currentMember := range t.Members {
		if err = currentMember.add(teamId); err != nil {
			conf.Logger.Error(fmt.Sprintf("%s: %v", conf.LogHeaders.Usecase, err))
			break
		}
	}
	if err != nil {
		if _, err = db.Connection.Exec(context.Background(), "delete from teams where id=$1", teamId); err != nil {
			conf.Logger.Error(fmt.Sprintf("%s: error on removing raw team by id = %d", conf.LogHeaders.Usecase, teamId))
		}
		return
	}
	conf.Logger.Info(fmt.Sprintf("%s: new team successfully added", conf.LogHeaders.Usecase))
	return
}
