package usecase

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/db"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type (
	PullRequest struct {
		PullRequestId     string   `json:"pull_request_id"`
		PullRequestName   string   `json:"pull_request_name"`
		AuthorId          string   `json:"author_id"`
		Status            string   `json:"status"`
		AssignedReviewers []string `json:"assigned_reviewers"`
	}
	ShortPullRequest struct {
		PullRequestId   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorId        string `json:"author_id"`
		Status          string `json:"status"`
	}
)

const mergedStatus = "MERGED"

var (
	ErrorAuthorTeamNotFound = errors.New("author or his teammates not found")
	ErrorPRDuplication      = errors.New("ОШИБКА: повторяющееся значение ключа нарушает ограничение уникальности \"pull_requests_pkey\" (SQLSTATE 23505)")
	ErrorPRDidntMerged      = errors.New("cannot scan NULL into *time.Time")
	ErrorPRNotFound         = errors.New("no rows in result set")
	ErrorPRAuthorNotFound   = errors.New("can't scan into dest[3] (col: array_position): cannot scan NULL into *int64")
	ErrorPRReassignMerge    = errors.New("PR has been merged")
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
		err = fmt.Errorf("%w", ErrorAuthorTeamNotFound)
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

func (pr *PullRequest) Merge() (transactionTime time.Time, err error) {
	var row pgx.Row
	row = db.Connection.QueryRow(context.Background(),
		"select pr_name, author_id, assigned_reviewers, merged_at, status from pull_requests where pr_id = $1", pr.PullRequestId)
	if err = row.Scan(&pr.PullRequestName, &pr.AuthorId, &pr.AssignedReviewers, &transactionTime, &pr.Status); err != nil {
		if !errors.Is(err, ErrorPRDidntMerged) {
			err = fmt.Errorf("error on getting PR data: %w", err)
			return
		}
	}
	if pr.Status == mergedStatus {
		return
	}
	row = db.Connection.QueryRow(context.Background(),
		"update pull_requests set status = $1, merged_at = NOW() where pr_id = $2 returning merged_at",
		mergedStatus, pr.PullRequestId)
	if err = row.Scan(&transactionTime); err != nil {
		err = fmt.Errorf("error on updation PR data: %w", err)
		return
	}
	pr.Status = mergedStatus
	return
}

func (pr *PullRequest) Reassign(oldReviewerId string) (newReviewerId string, err error) {
	var row pgx.Row
	row = db.Connection.QueryRow(context.Background(),
		`select author_id, status, team_id, array_position(assigned_reviewers, $1)
		from pull_requests join users on pull_requests.author_id = users.user_id
		where pr_id = $2`,
		oldReviewerId, pr.PullRequestId)
	var teamId, oldReviewerIndex int64
	if err = row.Scan(&pr.AuthorId, &pr.Status, &teamId, &oldReviewerIndex); err != nil {
		err = fmt.Errorf("error on reading PR data: %w", err)
		return
	}
	if pr.Status == mergedStatus {
		err = ErrorPRReassignMerge
		return
	}
	row = db.Connection.QueryRow(context.Background(),
		`select user_id
		from users
		where is_active = true and team_id = $1 and
		user_id not in (select unnest(assigned_reviewers) from pull_requests where pr_id = $2) and user_id != $3 and user_id != $4
		order by random()
		limit 1`,
		teamId, pr.PullRequestId, oldReviewerId, pr.AuthorId)
	if err = row.Scan(&newReviewerId); err != nil {
		err = fmt.Errorf("error on selecting new reviewer: %w", err)
		return
	}
	row = db.Connection.QueryRow(context.Background(),
		"update pull_requests set assigned_reviewers[$1] = $2 where pr_id = $3 returning pr_name, assigned_reviewers",
		oldReviewerIndex, newReviewerId, pr.PullRequestId)
	if err = row.Scan(&pr.PullRequestName, &pr.AssignedReviewers); err != nil {
		err = fmt.Errorf("error on updation PR data: %w", err)
		return
	}
	return
}
