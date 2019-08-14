package config

import (
	"go-open-registry/internal/log"
	"go-open-registry/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
)

// Auth struct with credentials
type Auth struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// AppConfig type struct
type AppConfig struct {
	App struct {
		Port             int    `json:"port"`
		CargoAPIURL      string `json:"cargo_api_url"`
		CargoDownloadURL string `json:"cargo_download_url"`
	}
	Repo struct {
		URL      string          `json:"url"`
		Path     string          `json:"path"`
		Instance *git.Repository `json:"instance"`
		Auth     Auth
	}
	Storage struct {
		Type     storage.Type `json:"type"`
		Instance storage.GenericStorage
	}
	LocalStorage       storage.LocalStorage
	ArtifactoryStorage storage.ArtifactoryStorage
	S3Storage          storage.S3Storage
	DB                 struct {
		URI     string        `json:"uri"`
		Timeout time.Duration `json:"timeout"`
		Client  *mongo.Client `json:"client"`
	}
}

// New Initialize new application config
func New() *AppConfig {
	appConfig := AppConfig{}

	viper.AutomaticEnv()

	// PORT
	viper.SetDefault("port", 8000)
	appConfig.App.Port = viper.GetInt("port")

	// CARGO_API_URL
	viper.SetDefault("cargo_api_url", "http://localhost:8000")
	appConfig.App.CargoAPIURL = viper.GetString("cargo_api_url")
	appConfig.App.CargoDownloadURL = viper.GetString("cargo_api_url") + "/api/v1/crates"

	// GIT_REPO_URL
	viper.SetDefault("git_repo_url", "")
	appConfig.Repo.URL = viper.GetString("git_repo_url")

	// GIT_REPO_PATH
	viper.SetDefault("git_repo_path", "tmpGit")
	appConfig.Repo.Path = viper.GetString("git_repo_path")

	// GIT_REPO_USERNAME
	viper.SetDefault("git_repo_username", "")
	appConfig.Repo.Auth.Name = viper.GetString("git_repo_username")

	// GIT_REPO_PASSWORD
	viper.SetDefault("git_repo_password", "")
	appConfig.Repo.Auth.Password = viper.GetString("git_repo_password")

	// GIT_REPO_EMAIL
	viper.SetDefault("git_repo_email", "")
	appConfig.Repo.Auth.Email = viper.GetString("git_repo_email")

	// MONGODB_URI
	viper.SetDefault("mongodb_uri", "mongodb://localhost:27017")
	appConfig.DB.URI = viper.GetString("mongodb_uri")

	// MONGO_CONNECTION_TIMEOUT
	viper.SetDefault("mongo_connection_timeout", 5)
	appConfig.DB.Timeout = viper.GetDuration("mongo_connection_timeout")

	// STORAGE_TYPE
	viper.SetDefault("storage_type", "local")

	// LOCAL_STORAGE_PATH
	viper.SetDefault("local_storage_path", "upload")
	appConfig.LocalStorage.Path = viper.GetString("local_storage_path")

	// ARTIFACTORY_URL
	viper.SetDefault("artifactory_url", "")
	appConfig.ArtifactoryStorage.URL = viper.GetString("artifactory_url")

	// ARTIFACTORY_LOGIN
	viper.SetDefault("artifactory_login", "")
	appConfig.ArtifactoryStorage.Login = viper.GetString("artifactory_login")

	// ARTIFACTORY_PASSWORD
	viper.SetDefault("artifactory_password", "")
	appConfig.ArtifactoryStorage.Password = viper.GetString("artifactory_password")

	// ARTIFACTORY_REPO_NAME
	viper.SetDefault("artifactory_repo_name", "")
	appConfig.ArtifactoryStorage.RepoName = viper.GetString("artifactory_repo_name")

	appConfig.ArtifactoryStorage.Path = appConfig.ArtifactoryStorage.URL + "/" + appConfig.ArtifactoryStorage.RepoName

	// AWS_ACCESS_KEY_ID
	viper.SetDefault("aws_access_key_id", "")
	appConfig.S3Storage.AccessKeyID = viper.GetString("aws_access_key_id")

	// AWS_SECRET_ACCESS_KEY
	viper.SetDefault("aws_secret_access_key", "")
	appConfig.S3Storage.SecretAccessKey = viper.GetString("aws_secret_access_key")

	// AWS_DEFAULT_REGION
	viper.SetDefault("aws_default_region", "")
	appConfig.S3Storage.DefaultRegion = viper.GetString("aws_default_region")

	// AWS_S3_BUCKET_NAME
	viper.SetDefault("aws_s3_bucket_name", "")
	appConfig.S3Storage.BucketName = viper.GetString("aws_s3_bucket_name")

	//AWS_S3_ENDPOINT
	viper.SetDefault("aws_s3_endpoint", "")
	appConfig.S3Storage.Endpoint = viper.GetString("aws_s3_endpoint")
	if appConfig.S3Storage.Endpoint == "s3" {
		appConfig.S3Storage.Endpoint = "s3." + appConfig.S3Storage.DefaultRegion + ".amazonaws.com"
	}

	// AWS_S3_USE_SSL
	viper.SetDefault("aws_s3_use_ssl", "")
	appConfig.S3Storage.UseSSL = viper.GetBool("aws_s3_use_ssl")

	switch viper.GetString("storage_type") {
	case "local":
		appConfig.Storage.Type = storage.Local
		appConfig.Storage.Instance = &appConfig.LocalStorage
		log.Info("Using local storage")
	case "s3":
		appConfig.Storage.Type = storage.S3
		appConfig.Storage.Instance = &appConfig.S3Storage
		log.Info("Using S3 storage")
	case "artifactory":
		appConfig.Storage.Type = storage.Artifactory
		appConfig.Storage.Instance = &appConfig.ArtifactoryStorage
		log.Info("Using artifactory storage")

	default:
		log.FatalWithFields("Storage config can be set one of: 'local', 's3', 'artifactory'",
			log.Fields{"storage": viper.GetString("storage_type")})
	}

	return &appConfig
}
