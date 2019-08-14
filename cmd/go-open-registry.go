package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v6"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/handlers"
	"go-open-registry/internal/log"
	"go-open-registry/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus" //nolint:depguard
	ginlogrus "github.com/toorop/gin-logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB(appConfig *config.AppConfig) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), appConfig.DB.Timeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.DB.URI))
	if err != nil {
		log.FatalWithFields("Mongo connection failed after timeout", log.Fields{"err": err})
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.FatalWithFields("Mongo ping failed after timeout", log.Fields{"mongo": appConfig.DB.URI})
		return err
	}

	log.InfoWithFields("Mongo connected", log.Fields{"mongo": appConfig.DB.URI})
	appConfig.DB.Client = client

	result, err := client.Database("crates").Collection("packages").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"name":    1,
			"version": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.InfoWithFields("Index already exist", log.Fields{"result": err})
		return nil // Todo: handle duplicate index
	}
	log.InfoWithFields("Index created", log.Fields{"index": result})
	return err

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found. Searching config data in environment variables")
	}

	appConfig := config.New()
	appRepo := gitregistry.New(appConfig)
	appConfig.Repo.Instance = appRepo
	//appStorage := storage.New(appConfig.Storage.Type, appConfig.Storage.Path, appConfig.Storage.Login, appConfig.Storage.Password)
	//appConfig.Storage.Instance = appStorage

	err = initDB(appConfig)
	if err != nil {
		log.ErrorWithFields("Error from initDB", log.Fields{
			"err": err,
		})
	}

	err = initGit(appConfig)
	if err != nil {
		log.ErrorWithFields("Error from initGit", log.Fields{
			"err": err,
		})
	}

	if appConfig.Storage.Type == storage.S3 {
		err = initS3Storage(appConfig)
		if err != nil {
			log.ErrorWithFields("Error from initS3Storage", log.Fields{
				"err": err,
			})
		}
	}

	if gin.Mode() != gin.ReleaseMode {
		log.SetLogLevel(logrus.DebugLevel)
		logrus.Info("Config: ")
		jsonOutput, err := json.MarshalIndent(appConfig, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonOutput))
	} else {
		// gin release mode
		log.SetLogLevel(logrus.InfoLevel)
	}

	engine := gin.New()

	engine.Use(ginlogrus.Logger(log.Logger), gin.Recovery())

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

func initGit(appConfig *config.AppConfig) (err error) {
	err = gitregistry.InitConfig(appConfig)
	if err != nil {
		return err
	}
	return err
}

func initS3Storage(appConfig *config.AppConfig) (err error) {
	log.Info("Init Client")
	// Initialize minio client object.
	minioClient, err := minio.New(appConfig.S3Storage.Endpoint, appConfig.S3Storage.AccessKeyID, appConfig.S3Storage.SecretAccessKey, appConfig.S3Storage.UseSSL)
	if err != nil {
		log.Error(err)
		return err
	}

	err = minioClient.MakeBucket(appConfig.S3Storage.BucketName, appConfig.S3Storage.DefaultRegion)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(appConfig.S3Storage.BucketName)
		if errBucketExists == nil && exists {
			log.InfoWithFields("We already own", log.Fields{
				"bucket": appConfig.S3Storage.BucketName,
			})
		} else {
			log.Fatal(err)
		}
	} else {
		log.InfoWithFields("Successfully created", log.Fields{
			"bucket": appConfig.S3Storage.BucketName,
		})
	}
	return err
}
