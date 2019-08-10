package storage

import "github.com/sirupsen/logrus"

type S3Storage struct {
	Path string
}

func (s *S3Storage) PutFile(packageName string, packageVersion string, content []byte) (Response, error) {
	logrus.Info("Put a file to S3 storage")
	return Response{message:"Ok"}, nil
}

func (s *S3Storage) GetFile(filename string) ([]byte, error) {
	logrus.Info("Get a file from S3 storage")
	return []byte{}, nil
}

