package config

import (
	"go-open-registry/internal/log"
	"go-open-registry/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

// AppConfig type struct
type AppConfig struct {
	App struct {
		Port int `json:"port"`
	}
	Repo struct {
		URL      string          `json:"url"`
		Path     string          `json:"path"`
		Instance *git.Repository `json:"instance"`
		Bot      struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
	}
	Storage struct {
		Type     storage.Type           `json:"type"`
		Path     string                 `json:"path"`
		Instance storage.GenericStorage `json:"instance"`
	}
	DB struct {
		URI     string        `json:"uri"`
		Timeout time.Duration `json:"timeout"`
		Client  *mongo.Client `json:"client"`
	}
}

// New Initialize new application config
func New() *AppConfig {
	appConfig := AppConfig{}

	viper.AutomaticEnv()

	viper.SetDefault("port", 8000)
	appConfig.App.Port = viper.GetInt("port")

	viper.SetDefault("mongodb_uri", "mongodb://localhost:27017")
	appConfig.DB.URI = viper.GetString("mongodb_uri")

	viper.SetDefault("mongo_connection_timeout", 5)
	appConfig.DB.Timeout = viper.GetDuration("mongo_connection_timeout")

	viper.SetDefault("git_repo_url", "")
	appConfig.Repo.URL = viper.GetString("git_repo_url")

	viper.SetDefault("git_repo_username", "")
	appConfig.Repo.Bot.Name = viper.GetString("git_repo_username")

	viper.SetDefault("git_repo_email", "")
	appConfig.Repo.Bot.Email = viper.GetString("git_repo_email")

	viper.SetDefault("git_repo_password", "")
	appConfig.Repo.Bot.Password = viper.GetString("git_repo_password")

	viper.SetDefault("git_repo_path", "tmpGit")
	appConfig.Repo.Path = viper.GetString("git_repo_path")

	viper.SetDefault("storage_path", "upload")
	appConfig.Storage.Path = viper.GetString("storage_path")

	viper.SetDefault("storage_type", "local")

	switch viper.GetString("storage_type") {
	case "local":
		appConfig.Storage.Type = storage.Local
		log.Info("Using local storage")
	case "s3":
		appConfig.Storage.Type = storage.S3
		log.Info("Using S3 storage")
	case "artifactory":
		appConfig.Storage.Type = storage.Artifactory
		log.Info("Using artifactory storage")
	default:
		log.FatalWithFields("Storage config can be set one of: 'local', 's3', 'artifactory'",
			log.Fields{"storage": viper.GetString("storage")})
	}

	return &appConfig
}
