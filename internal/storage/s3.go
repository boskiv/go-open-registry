package storage

import (
	"go-open-registry/internal/log"
)

// S3Storage  type struct
type S3Storage struct {
	Path string
}

// PutFile implementation
func (s *S3Storage) PutFile(packageName, packageVersion string, content []byte) error {
	log.Info("Put a file to S3 storage")
	return nil
}

// GetFile implementation
func (s *S3Storage) GetFile(packageName, packageVersion string) ([]byte, error) {
	log.Info("Get a file from S3 storage")
	return []byte{}, nil
}
