package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/handlers"
	"go-open-registry/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus" //nolint:depguard
	ginlogrus "github.com/toorop/gin-logrus"
)

func main() {
	appConfig := config.New()
	appRepo := gitregistry.New(appConfig.Repo.URL)
	appConfig.Repo.Instance = appRepo
	appStorage := storage.New(appConfig.Storage.Type)
	appConfig.Storage.Instance = appStorage
	log := logrus.New()

	if gin.Mode() != gin.ReleaseMode {
		logrus.Info("Config: ")
		jsonOutput, err := json.MarshalIndent(appConfig, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonOutput))
	}

	engine := gin.New()

	engine.Use(ginlogrus.Logger(log), gin.Recovery())

	engine.PUT("/api/v1/crates/new", handlers.NewCrateHandler(appConfig))

	logrus.WithField("port", appConfig.App.Port).Info("Staring server on port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.App.Port),
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithField("error", err).Fatal("listen: ")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithField("error", err).Fatal("Server Shutdown: ", err)
	}

	logrus.Info("Server exiting")
}
