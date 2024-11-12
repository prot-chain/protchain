package main

import (
	"log"
	"os"
	"os/signal"
	"protchain/internal/api/rest"
	"protchain/internal/config"
	"protchain/internal/dep"
	"protchain/internal/logic"
	"protchain/internal/service/queue"
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	allowConnectionsAfterShutdown = time.Second * 8
)

func main() {
	appConfig := config.LoadConfig()
	appDep := dep.New(appConfig)
	appLogic := logic.New(appDep)
	logic.InitWorkers(appLogic, appConfig)
	logic.InitCron(appLogic)

	notifyConnCloseCh := queue.RmqConn.NotifyClose(make(chan *amqp091.Error))

	go func() {
		for notifyConnCloseCh != nil {
			select {
			case err, ok := <-notifyConnCloseCh:

				if ok {

					log.Printf("worker connection closed due to: %v. Attempting to re-initilize workers", err)
					logic.InitWorkers(appLogic, appConfig)

				}
			}
		}
	}()

	restApi := rest.API{
		Config: appConfig,
		Dep:    appDep,
		Logic:  appLogic,
	}

	go func() {
		log.Fatal(restApi.Serve())
	}()

	// graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	logger, _ := zap.NewProduction()
	logger.Sugar().Infof("Request to shutdown server. Doing nothing for %v", allowConnectionsAfterShutdown)
	waitTimer := time.NewTimer(allowConnectionsAfterShutdown)
	<-waitTimer.C

	logger.Info("Shutting down server...")
}
