package config

import (
	"go-open-registry/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/src-d/go-git.v4"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_appConfig := AppConfig{
		App: struct {
			Port int `json:"port"`
		}{8000},
		Repo: struct {
			URL      string          `json:"url"`
			Path     string          `json:"path"`
			Instance *git.Repository `json:"instance"`
			Bot      RepoBot
		}{"", "tmpGit", nil, RepoBot{
			Name:     "",
			Email:    "",
			Password: "",
		}},
		Storage: struct {
			Type     storage.Type           `json:"type"`
			Path     string                 `json:"path"`
			Login    string                 `json:"username"`
			Password string                 `json:"-"`
			Instance storage.GenericStorage `json:"instance"`
		}{storage.Local, "upload", "", "", nil},
		DB: struct {
			URI     string        `json:"uri"`
			Timeout time.Duration `json:"timeout"`
			Client  *mongo.Client `json:"client"`
		}{"mongodb://localhost:27017", 5, nil},
	}

	tests := []struct {
		name string
		want *AppConfig
	}{
		{"DefaultConfig", &_appConfig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
