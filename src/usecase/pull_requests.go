package usecase

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/db"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PullRequest struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorId          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

var (
	ErrorAuthorTeamNotFound customError = errors.New("author or his teammates not found")
	ErrorPRDuplication      customError = errors.New("ОШИБКА: повторяющееся значение ключа нарушает ограничение уникальности \"pull_requests_pkey\" (SQLSTATE 23505)")
)

func (pr *PullRequest) Create() (err error) {
	var rows pgx.Rows
	if rows, err = db.Connection.Query(context.Background(),
		`select user_id 
		from users 
		where user_id != $1 and is_active = true and 
		team_id = (select team_id from users where user_id = $2)`,
		pr.AuthorId, pr.AuthorId); err != nil {
		err = fmt.Errorf("error on getting author's teammaets: %w", err)
		return
	}
	defer rows.Close()
	teammatesArray := `{`
	teammetesCounter := 0
	for rows.Next() && teammetesCounter < conf.PRReviewersMax {
		var currentUserId string
		if err = rows.Scan(&currentUserId); err != nil {
			err = fmt.Errorf("error on parsing teammates data: %w", err)
			return
		}
		teammatesArray = fmt.Sprintf(`%s"%s",`, teammatesArray, currentUserId)
		teammetesCounter++
	}
	if teammetesCounter == 0 {
		err = ErrorAuthorTeamNotFound
		return
	}
	teammatesArray = fmt.Sprintf("%s}", teammatesArray[:len(teammatesArray)-1])
	if _, err = db.Connection.Exec(context.Background(),
		`insert into pull_requests (pr_id, pr_name, author_id, assigned_reviewers) values ($1, $2, $3, $4)`,
		pr.PullRequestId, pr.PullRequestName, pr.AuthorId, teammatesArray); err != nil {
		err = fmt.Errorf("error on inserting new PR data: %w", err)
		return
	}
	return
}
