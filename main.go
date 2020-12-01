package main

import (
	"context"
	"flag"
	"fmt"
	"go-gin-gorm-mysql/internal/core/config"
	"go-gin-gorm-mysql/internal/core/database"
	"go-gin-gorm-mysql/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "go-gin-gorm-mysql/docs"

	"github.com/sirupsen/logrus"
)

// @title Learning API
// @version 1.0
// @description Learning API Description

// @host go-mysql-api-xapa2c3k2q-as.a.run.app
// @BasePath /api
func main() {
	//=======================================================
	// Read enironment for server config
	environment := flag.String("environment", "dev", "set working environment")
	configs := flag.String("config", "configs", "set configs path, default as: 'configs'")
	port := flag.String("port", "8080", "set port for start service")

	flag.Parse()

	//=======================================================

	// Init configuration
	if err := config.InitConfig(*configs, *environment); err != nil {
		panic(err)
	}
	logrus.Infof("Initial 'Configuration'. %+v", config.CF)
	//=======================================================

	// Init return result
	if err := config.InitReturnResult("configs"); err != nil {
		panic(err)
	}
	//=======================================================

	// Init database mysql
	if err := database.InitConnection(config.CF); err != nil {
		panic(err)
	}

	ro := handler.Routes{}
	handler := ro.Init(config.CF, config.RR)

	srv := &http.Server{
		Addr:              fmt.Sprint(":", *port),
		Handler:           handler,
		ReadTimeout:       config.CF.HTTPServer.ReadTimeout,
		WriteTimeout:      config.CF.HTTPServer.WriteTimeout,
		ReadHeaderTimeout: config.CF.HTTPServer.ReadHeaderTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	cancel()

	srvCtx, srvCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer srvCancel()
	logrus.Infof("shutting down http server...")
	if err := srv.Shutdown(srvCtx); err != nil {
		logrus.Panicln("http server shutdown with error:", err)
	}
}
