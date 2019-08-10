package storage

import "github.com/sirupsen/logrus" //nolint:depguard

// LocalStorage type struct
type LocalStorage struct {
	Path string
}

// PutFile implementation
func (l *LocalStorage) PutFile(packageName, packageVersion string, content []byte) (Response, error) {
	logrus.Info("Put a file to local storage")
	return Response{message: "Ok"}, nil
}

// GetFile implementation
func (l *LocalStorage) GetFile(filename string) ([]byte, error) {
	logrus.Info("Get a file from local storage")
	return []byte{}, nil
}
