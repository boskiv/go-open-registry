package storage

import (
	"go-open-registry/internal/log"
)

// ArtifactoryStorage struct
type ArtifactoryStorage struct {
	Path string
}

// PutFile implementation
func (a ArtifactoryStorage) PutFile(packageName, packageVersion string, content []byte) error {
	log.Info("Put a file to artifactory storage")
	return nil
}

// GetFile implementation
func (a ArtifactoryStorage) GetFile(packageName, packageVersion string) ([]byte, error) {
	log.Info("Get a file from artifactory storage")
	return []byte{}, nil
}
