package storage

import "github.com/sirupsen/logrus" //nolint:depguard

// S3Storage  type struct
type S3Storage struct {
	Path string
}

// PutFile implementation
func (s *S3Storage) PutFile(packageName, packageVersion string, content []byte) error {
	logrus.Info("Put a file to S3 storage")
	return nil
}

// GetFile implementation
func (s *S3Storage) GetFile(packageName, packageVersion string) ([]byte, error) {
	logrus.Info("Get a file from S3 storage")
	return []byte{}, nil
}
