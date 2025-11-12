package http

import (
	"avito_pr_service/src/conf"
	"fmt"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	reuse "github.com/libp2p/go-reuseport"
)

func StartHTTPServer() *net.Listener {
	router := gin.Default()
	team := router.Group("team")
	{
		team.POST("add", addTeam)
		team.GET("get", getTeam)
	}
	users := router.Group("users")
	{
		users.POST("setIsActive", setActiveUser)
		users.GET("getReview", getUsersRewie)
	}
	pullRequests := router.Group("pullRequest")
	{
		pullRequests.POST("create", createPullRequest)
		pullRequests.POST("merge", mergePullRequest)
		pullRequests.POST("reassign", reassignPullRequest)
	}
	var listener net.Listener
	var err error
	if listener, err = reuse.Listen("tcp", conf.ServerHTTPServeSocket); err != nil {
		conf.Logger.Error(fmt.Sprintf("Error on creating listener: %s", err))
		os.Exit(1)
	}
	if err = router.RunListener(listener); err != nil {
		conf.Logger.Error(fmt.Sprintf("Error on starting HTTP-server: %s", err))
		os.Exit(1)
	}
	conf.Logger.Info("HTTP-server has been started")
	return &listener
}
