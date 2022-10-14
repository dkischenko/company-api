package app

import (
	"context"
	"flag"
	"fmt"
	"githib.com/dkischenko/company-api/configs"
	xm_logger "githib.com/dkischenko/company-api/pkg/logger"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunServer(router *mux.Router, logger *xm_logger.Logger, config *configs.Config) {
	logger.Entry.Info("start application")
	logger.Entry.Info("listen TCP")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.AppHost, config.AppPort))

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Entry.Infof("server listening address %s:%s", config.AppHost, config.AppPort)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("shutting down")
	os.Exit(0)
}
