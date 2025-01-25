package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	// "protchain/internal/api/rest"
	"protchain/internal/config"
	"protchain/internal/dep"
	"protchain/internal/logging"
	"protchain/internal/logic"
	"protchain/internal/restapi"
	"syscall"
	"time"
)

const (
	allowConnectionsAfterShutdown = time.Second * 8
)

func main() {
	appConfig := config.LoadConfig()
	appDep := dep.New(appConfig)
	_ = logic.New(appDep)

	restApi := restapi.API{
		Config: appConfig,
		Deps:   appDep,
	}

	logging.New(appConfig)
	go func() {
		log.Fatal(restApi.Serve())
	}()

	// graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	slog.Info(fmt.Sprintf("Request to shutdown server. Doing nothing for %v", allowConnectionsAfterShutdown))
	waitTimer := time.NewTimer(allowConnectionsAfterShutdown)
	<-waitTimer.C

	slog.Info("Shutting down server...")
}
