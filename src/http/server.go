package http

import (
	"avito_pr_service/src/conf"
	"avito_pr_service/src/db"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	reuse "github.com/libp2p/go-reuseport"
)

func StartHTTPServer() {
	var listener net.Listener

	exit := func() {
		if err := listener.Close(); err != nil {
			conf.Logger.Error(fmt.Sprintf("Error on closing HTTP-server: %v", err))
		}
		db.Connection.Close()

	}

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
	var err error
	if listener, err = reuse.Listen("tcp", conf.ServerHTTPServeSocket); err != nil {
		conf.Logger.Error(fmt.Sprintf("Error on creating listener: %s", err))
		os.Exit(1)
	}
	var adminToken, userToken string
	if adminToken, userToken, err = generateTokens(); err != nil {
		conf.Logger.Error(fmt.Sprintf("%s: error on getting tokens: %v", conf.LogHeaders.HTTPServer, err))
		exit()
	}
	conf.Logger.Info(fmt.Sprintf("%s: admin token - %s, user token - %s", conf.LogHeaders.HTTPServer, adminToken, userToken))
	go func() {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
		acceptedSignal := <-osSignals
		close(osSignals)
		conf.Logger.Info(fmt.Sprintf("Accepted OS %s signal", acceptedSignal.String()))
		exit()
	}()
	if err = router.RunListener(listener); err != nil {
		conf.Logger.Warn(fmt.Sprintf("Error on listening HTTP-server: %s", err))
		os.Exit(1)
	}
}
