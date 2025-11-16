package http

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/usecase"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	errorBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	errorResponse struct {
		Error errorBody `json:"error"`
	}
)

func addTeam(gctx *gin.Context) {
	conf.Logger.Debug(fmt.Sprintf("%s: add new team request", conf.LogHeaders.HTTPServer))
	var team usecase.Team
	if err := gctx.ShouldBindJSON(&team); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on binding requst body: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := team.Add(gctx); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on adding new team: %v", conf.LogHeaders.HTTPServer, err))
		if strings.Contains(err.Error(), usecase.ErrorTeamDuplication.Error()) {
			gctx.JSON(http.StatusBadRequest, errorResponse{Error: errorBody{
				Code:    "TEAM_EXIST",
				Message: "team_name already exist",
			}})
			return
		}
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gctx.JSON(http.StatusCreated, team)
}

func getTeam(gctx *gin.Context) {
	conf.Logger.Debug(fmt.Sprintf("%s: get team request", conf.LogHeaders.HTTPServer))
	var teamName string
	var ok bool
	if teamName, ok = gctx.GetQuery("team_name"); !ok {
		err := errors.New("''team_name'' is required query parameter")
		conf.Logger.Error(fmt.Sprintf("%s: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	team := usecase.Team{TeamName: teamName}
	if err := team.Get(gctx); err != nil {
		conf.Logger.Error("%s: error on getting team data: %v", conf.LogHeaders.HTTPServer, err)
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if len(team.Members) == 0 {
		gctx.JSON(http.StatusNotFound, errorResponse{Error: errorBody{
			Code:    "NOT_FOUND",
			Message: "resource not found",
		}})
		return
	}
	conf.Logger.Debug(fmt.Sprintf("%s: current team - %v", conf.LogHeaders.HTTPServer, team))
	gctx.JSON(http.StatusOK, team)
}

func setActiveUser(gctx *gin.Context) {
	type response struct {
		usecase.TeamMember
		TeamName string `json:"team_name"`
	}

	conf.Logger.Debug(fmt.Sprintf("%s: set user active request", conf.LogHeaders.HTTPServer))
	var user usecase.TeamMember
	var err error
	if err = gctx.ShouldBindJSON(&user); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on binding requst body: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fullUserData := response{TeamMember: user}
	if fullUserData.TeamName, err = user.SetActive(gctx); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on set active user: %v", conf.LogHeaders.HTTPServer, err))
		if strings.Contains(err.Error(), usecase.ErrorNotFound.Error()) {
			gctx.JSON(http.StatusNotFound, errorResponse{Error: errorBody{
				Code:    "NOT_FOUND",
				Message: "resource not found",
			}})
			return
		}
		return
	}
	gctx.JSON(http.StatusOK, fullUserData)
}

func getUsersRewie(gctx *gin.Context) {
	type response struct {
		UserId       string                     `json:"user_id"`
		PullRequests []usecase.ShortPullRequest `json:"pull_requests"`
	}

	conf.Logger.Debug(fmt.Sprintf("%s: get user's PR requets", conf.LogHeaders.HTTPServer))
	var userId string
	var ok bool
	if userId, ok = gctx.GetQuery("user_id"); !ok {
		err := errors.New("''user_id'' is required query parameter")
		conf.Logger.Error(fmt.Sprintf("%s: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := usecase.TeamMember{UserId: userId}
	userPRData := response{UserId: userId}
	var err error
	if userPRData.PullRequests, err = user.GetRewiew(gctx); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	gctx.JSON(http.StatusOK, userPRData)
}

func createPullRequest(gctx *gin.Context) {
	type response struct {
		PR usecase.PullRequest `json:"pr"`
	}

	conf.Logger.Debug(fmt.Sprintf("%s: create PR request", conf.LogHeaders.HTTPServer))
	var pr usecase.PullRequest
	var err error
	if err = gctx.ShouldBindJSON(&pr); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on binding requst body: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err = pr.Create(gctx); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on creation PR: %v", conf.LogHeaders.HTTPServer, err))
		if errors.Is(err, usecase.ErrorAuthorTeamNotFound) || strings.Contains(err.Error(), usecase.ErrorInvalidAuthor.Error()) {
			gctx.JSON(http.StatusNotFound, errorResponse{Error: errorBody{
				Code:    "NOT_FOUND",
				Message: "resource not found",
			}})
			return
		}
		if strings.Contains(err.Error(), usecase.ErrorPRDuplication.Error()) {
			gctx.JSON(http.StatusConflict, errorResponse{Error: errorBody{
				Code:    "PR_EXISTS",
				Message: "PR id already exists",
			}})
			return
		}
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp := response{PR: pr}
	gctx.JSON(http.StatusOK, resp)
}

func mergePullRequest(gctx *gin.Context) {
	type response struct {
		PR       usecase.PullRequest
		MergedAt string `json:"merged_at"`
	}

	conf.Logger.Debug(fmt.Sprintf("%s: merge PR request", conf.LogHeaders.HTTPServer))
	var pr usecase.PullRequest
	var err error
	if err = gctx.ShouldBindJSON(&pr); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on binding requst body: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var transactionTime time.Time
	if transactionTime, err = pr.Merge(gctx); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on merging PR: %v", conf.LogHeaders.HTTPServer, err))
		if strings.Contains(err.Error(), usecase.ErrorPRNotFound.Error()) {
			gctx.JSON(http.StatusNotFound, errorResponse{Error: errorBody{
				Code:    "NOT_FOUND",
				Message: "resource not found",
			}})
			return
		}
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gctx.JSON(http.StatusOK, response{
		PR:       pr,
		MergedAt: transactionTime.Format("2006-01-02T15:04:05Z"),
	})
}

func reassignPullRequest(gctx *gin.Context) {
	type (
		request struct {
			PullRequestId string `json:"pull_request_id"`
			OldReviewerId string `json:"old_reviewer_id"`
		}
		response struct {
			PR         usecase.PullRequest `json:""`
			ReplacedBy string              `json:"replaced_by"`
		}
	)

	conf.Logger.Debug(fmt.Sprintf("%s: reassign PR request", conf.LogHeaders.HTTPServer))
	var err error
	var req request
	if err = gctx.ShouldBindJSON(&req); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on binding requst body: %v", conf.LogHeaders.HTTPServer, err))
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pr := usecase.PullRequest{PullRequestId: req.PullRequestId}
	var newReviewerId string
	if newReviewerId, err = pr.Reassign(gctx, req.OldReviewerId); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on reasign PR: %v", conf.LogHeaders.HTTPServer, err))
		if strings.Contains(err.Error(), fmt.Sprintf("error on selecting new reviewer: %s", usecase.ErrorNotFound.Error())) {
			gctx.JSON(http.StatusConflict, errorResponse{Error: errorBody{
				Code:    "NO_CANDIDATE",
				Message: "no active replacement candidate in team",
			}})
			return
		}
		if strings.Contains(err.Error(), usecase.ErrorNotFound.Error()) {
			gctx.JSON(http.StatusNotFound, errorResponse{Error: errorBody{
				Code:    "NOT_FOUND",
				Message: "resource not found",
			}})
			return
		}
		if strings.Contains(err.Error(), fmt.Sprintf("error on reading PR data: %s", usecase.ErrorPRAuthorNotFound.Error())) {
			gctx.JSON(http.StatusConflict, errorResponse{Error: errorBody{
				Code:    "NOT_ASSIGNED",
				Message: "reviewer is not assigned to this PR",
			}})
			return
		}
		if errors.Is(err, usecase.ErrorPRReassignMerge) {
			gctx.JSON(http.StatusConflict, errorResponse{Error: errorBody{
				Code:    "PR_MERGED",
				Message: "cannot reassign on merged PR",
			}})
			return
		}
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gctx.JSON(http.StatusOK, response{
		PR:         pr,
		ReplacedBy: newReviewerId,
	})
}
