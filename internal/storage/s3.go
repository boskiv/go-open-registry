package storage

import "github.com/sirupsen/logrus" //nolint:depguard

// S3Storage  type struct
type S3Storage struct {
	Path string
}

// PutFile implementation
func (s *S3Storage) PutFile(packageName, packageVersion string, content []byte) (Response, error) {
	logrus.Info("Put a file to S3 storage")
	return Response{message: "Ok"}, nil
}

// GetFile implementation
func (s *S3Storage) GetFile(filename string) ([]byte, error) {
	logrus.Info("Get a file from S3 storage")
	return []byte{}, nil
}
