package http

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	errorResponse struct {
		Error error `json:"error"`
	}
)

func addTeam(gctx *gin.Context) {
	conf.Logger.Debug(fmt.Sprintf("%s: add new team request", conf.LogHeaders.HTTPServer))
	var team usecase.Team
	if err := gctx.ShouldBindBodyWithJSON(&team); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on binding requst body: %v", err))
		gctx.Status(http.StatusInternalServerError)
	}
	// TODO: call logic method and response creation
	if err := team.Add(); err != nil {
		// TEAM EXIST validation
	}
	gctx.JSON(http.StatusCreated, team)
}

func getTeam(gctx *gin.Context) {}

func setActiveUser(gctx *gin.Context) {}

func getUsersRewie(gctx *gin.Context) {}

func createPullRequest(gctx *gin.Context) {}

func mergePullRequest(gctx *gin.Context) {}

func reassignPullRequest(gctx *gin.Context) {}
