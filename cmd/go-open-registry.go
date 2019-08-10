package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitRegistry"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/parser"
	"go-open-registry/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)



func NewCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {

	return func(c *gin.Context) {
		gitRegistry.HeadRepo(appConfig.Repo.Instance)
		// Read the Body content
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			jsonFile, crateFile, err := parser.ReadBinary(c.Request.Body)
			helpers.CheckIfError(err)
			fmt.Printf("%s", jsonFile)
			var crateJson parser.CrateJson
			err = json.Unmarshal(jsonFile, &crateJson)
			helpers.CheckIfError(err)

			gitRegistry.CommitCrateJson(appConfig, crateJson.Name, crateJson.Vers)
			_, _ = appConfig.Storage.Instance.PutFile(crateJson.Name, crateJson.Vers, crateFile)
		}

		resp := map[string][]string{
			// Array of strings of categories that are invalid and ignored.
			"invalid_categories": {},
			// Array of strings of badge names that are invalid and ignored.
			"invalid_badges": {},
			// Array of strings of arbitrary warnings to display to the user.
			"other": {},
		}
		c.JSON(200, gin.H{
			// Optional object of warnings to display to the user.
			"warnings": resp,
		})
	}
}

func main() {
	appConfig := config.New()
	appRepo := gitRegistry.New(appConfig.Repo.Url)
	appConfig.Repo.Instance = appRepo
	appStorage := storage.New(appConfig.Storage.Type)
	appConfig.Storage.Instance = appStorage
	log := logrus.New()

	engine := gin.New()

	engine.Use(ginlogrus.Logger(log), gin.Recovery())

	engine.PUT("/api/v1/crates/new", NewCrateHandler(appConfig))

	logrus.WithField("port", appConfig.App.Port).Info("Staring server on port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.App.Port),
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
