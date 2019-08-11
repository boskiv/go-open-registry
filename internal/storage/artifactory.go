package storage

import "github.com/sirupsen/logrus" //nolint:depguard

// ArtifactoryStorage struct
type ArtifactoryStorage struct {
	Path string
}

// PutFile implementation
func (a ArtifactoryStorage) PutFile(packageName, packageVersion string, content []byte) error {
	logrus.Info("Put a file to artifactory storage")
	return nil
}

// GetFile implementation
func (a ArtifactoryStorage) GetFile(packageName, packageVersion string) ([]byte, error) {
	logrus.Info("Get a file from artifactory storage")
	return []byte{}, nil
}
