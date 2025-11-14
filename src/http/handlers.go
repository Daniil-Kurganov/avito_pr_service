package http

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/usecase"
	"errors"
	"fmt"
	"net/http"

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
	if err := team.Add(); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on adding new team: %v", conf.LogHeaders.HTTPServer, err))
		if errors.As(err, &usecase.ErrorTeamDuplication) {
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
	if err := team.Get(); err != nil {
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
	if fullUserData.TeamName, err = user.SetActive(); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on set active user: %v", conf.LogHeaders.HTTPServer, err))
		if errors.As(err, &usecase.ErrorNotFound) {
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

func getUsersRewie(gctx *gin.Context) {}

func createPullRequest(gctx *gin.Context) {}

func mergePullRequest(gctx *gin.Context) {}

func reassignPullRequest(gctx *gin.Context) {}
