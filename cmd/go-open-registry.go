package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/handlers"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus" //nolint:depguard
	ginlogrus "github.com/toorop/gin-logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB(appConfig *config.AppConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), appConfig.DB.Timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.DB.URI))
	helpers.FatalIfError(err)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.WithField("mongo", appConfig.DB.URI).Fatal("Mongo connection failed after timeout")
	} else {
		logrus.WithField("mongo", appConfig.DB.URI).Info("Mongo connected")
		appConfig.DB.Client = client

		result, err := client.Database("crates").Collection("packages").Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.M{
				"name":    1,
				"version": 1,
			},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			println(err)
			logrus.WithField("result", err).Info("Index already exist")
		}
		logrus.WithField("index", result).Info("Index created")
	}
}

func main() {
	appConfig := config.New()
	appRepo := gitregistry.New(appConfig)
	appConfig.Repo.Instance = appRepo
	appStorage := storage.New(appConfig.Storage.Type, appConfig.Storage.Path)
	appConfig.Storage.Instance = appStorage

	initDB(appConfig)

	logger := logrus.New()

	if gin.Mode() != gin.ReleaseMode {
		logrus.Info("Config: ")
		jsonOutput, err := json.MarshalIndent(appConfig, "", "  ")
		if err != nil {
			logger.Fatal(err)
		}
		fmt.Println(string(jsonOutput))
	}

	engine := gin.New()

	engine.Use(ginlogrus.Logger(logger), gin.Recovery())

	engine.PUT("/api/v1/crates/new", handlers.NewCrateHandler(appConfig))
	engine.GET("/api/v1/crates/:name/:version/*download", handlers.GetCrateHandler(appConfig))
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
