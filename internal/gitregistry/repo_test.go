package gitregistry

import (
	"go-open-registry/internal/config"
	"go-open-registry/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"testing"
	"time"
)

func Test_createFile(t *testing.T) {
	type args struct {
		resultPathString string
		content          []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{"Simple", args{"testFile", []byte{}}, "testFile", false},
		{"err", args{"/testFile", []byte{}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := createFile(tt.args.resultPathString, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("createFile() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
	_ = os.Remove("testFile")

}

func Test_makePath(t *testing.T) {
	type args struct {
		appConfig   *config.AppConfig
		packageName string
	}

	_appConfig := &config.AppConfig{
		App: struct {
			Port int `json:"port"`
		}{8000},
		Repo: struct {
			URL      string          `json:"url"`
			Path     string          `json:"path"`
			Instance *git.Repository `json:"instance"`
			Bot      config.RepoBot
		}{"", "tmpGit", nil, config.RepoBot{
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
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{"First", args{
			appConfig:   _appConfig,
			packageName: "nohup",
		}, "tmpGit/no/hu/nohup", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := makePath(tt.args.appConfig, tt.args.packageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("makePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("makePath() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
