package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-open-registry/internal/storage"
)

type AppConfig struct {
	repo      string
	uploadDir string
	Storage   storage.GenericStorage
}

func InitConfig() *AppConfig {
	appConfig := AppConfig{}

	viper.AutomaticEnv()

	viper.SetDefault("repo", "")
	appConfig.repo = viper.GetString("repo")

	viper.SetDefault("upload_dir", "upload")
	appConfig.uploadDir = viper.GetString("upload_dir")

	viper.SetDefault("storage", "local")

	switch viper.GetString("storage") {
	case "local":
		appConfig.Storage = &storage.LocalStorage{}

		logrus.Info("Using local storage")
	case "s3":

		logrus.Info("Using S3 storage")
	case "artifactory":
		logrus.Info("Using artifactory storage")
	default:
		logrus.WithField("storage", viper.GetString("storage")).
			Fatal("Storage config can be set one of: 'local', 's3', 'artifactory'")
	}
	appConfig.Storage.NewStorage()

	return &appConfig
}
