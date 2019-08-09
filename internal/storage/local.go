package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-open-registry/internal/helpers"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type LocalStorage struct {
	path string
}

func (localStorage *LocalStorage) NewStorage() {
	localStorage.NewLocalStorage()
}

func (localStorage *LocalStorage) NewLocalStorage() *LocalStorage {
	uploadDir := viper.GetString("upload_dir")
	if len(uploadDir) == 0 {
		logrus.Fatal("No local path in config")
	}

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Created local upload dir")
	}

	return &LocalStorage{path: uploadDir}
}

func (localStorage *LocalStorage) PutFile(packageName string, content []byte) (Response, error) {
	logrus.WithFields(logrus.Fields{
		"file":   packageName,
		"folder": viper.GetString("upload_dir"),
	}).Info("Uploading file to local folder")

	paths := helpers.MakeCratePath(packageName)
	withUploadDir := append([]string{viper.GetString("upload_dir")}, paths...)
	_ = os.MkdirAll(strings.Join(withUploadDir, string(os.PathSeparator)), os.ModePerm)
	withPackageName := append(withUploadDir, packageName)
	err := ioutil.WriteFile(path.Join(withPackageName...), content, 0644)
	return Response{message: "Ok"}, err
}

func (localStorage *LocalStorage) GetFile(file string) ([]byte, error) {
	panic("implement me")
}
