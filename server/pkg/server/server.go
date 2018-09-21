package server

import (
	"memefy/server/pkg/config"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// RunServer starts the server
func RunServer(cancelCtx context.Context, ready chan bool, config *config.Config) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.Use(
		handlers.RecoveryHandler(
			handlers.RecoveryLogger(log.StandardLogger()),
			handlers.PrintRecoveryStack(true),
		),
	)
	router.Use(accessLoggingMiddleware)

	adminHandler := &AdminHandler{}

	//info endpoints
	router.Methods("GET").Path("/").Name("self").Handler(adminHandler.IndexHandler(router))
	router.Methods("GET").Path("/info").Name("info").Handler(adminHandler.AdminInfoHandler())
	router.Methods("GET").Path("/health").Name("health").Handler(adminHandler.HealthCheckHandler())

	//app endpoints
	router.Methods("POST").Path("/").Name("POST new meme").
		Handler(basicAuthMiddleware(config.Security)(adminHandler.PostMemeHandler()))

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Server.Port), Handler: router}

	go func() {
		log.Infof("Starting server at :%d", config.Server.Port)

		select {
		case ready <- true:
		default:
		}

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	select {
	case <-stop:
	case <-cancelCtx.Done():
	}

	log.Info("Shutting down the server...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Error(err)
	} else {
		log.Info("Server gracefully stopped")
	}
}
