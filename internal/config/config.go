package config

import (
	"go-open-registry/internal/storage"

	"github.com/sirupsen/logrus" //nolint:depguard
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
	}
	Storage struct {
		Type     storage.Type           `json:"type"`
		Path     string                 `json:"path"`
		Instance storage.GenericStorage `json:"instance"`
	}
}

// New Initialize new application config
func New() *AppConfig {
	appConfig := AppConfig{}

	viper.AutomaticEnv()

	viper.SetDefault("port", 8000)
	appConfig.App.Port = viper.GetInt("port")

	viper.SetDefault("git_repo_url", "")
	appConfig.Repo.URL = viper.GetString("git_repo_url")

	viper.SetDefault("git_repo_path", "./tmp")
	appConfig.Repo.Path = viper.GetString("git_repo_path")

	viper.SetDefault("storage_path", "./upload")
	appConfig.Storage.Path = viper.GetString("storage_path")

	viper.SetDefault("storage_type", "local")

	switch viper.GetString("storage_type") {
	case "local":
		appConfig.Storage.Type = storage.Local
		logrus.Info("Using local storage")
	case "s3":
		appConfig.Storage.Type = storage.S3
		logrus.Info("Using S3 storage")
	case "artifactory":
		appConfig.Storage.Type = storage.Artifactory
		logrus.Info("Using artifactory storage")
	default:
		logrus.WithField("storage", viper.GetString("storage")).
			Fatal("Storage config can be set one of: 'local', 's3', 'artifactory'")
	}

	return &appConfig
}
