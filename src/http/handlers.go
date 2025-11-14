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
		gctx.Status(http.StatusInternalServerError)
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
		gctx.JSON(http.StatusInternalServerError, nil)
	}
	gctx.JSON(http.StatusCreated, team)
}

func getTeam(gctx *gin.Context) {}

func setActiveUser(gctx *gin.Context) {}

func getUsersRewie(gctx *gin.Context) {}

func createPullRequest(gctx *gin.Context) {}

func mergePullRequest(gctx *gin.Context) {}

func reassignPullRequest(gctx *gin.Context) {}
