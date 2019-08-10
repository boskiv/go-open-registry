package storage

import "github.com/sirupsen/logrus"

type ArtifactoryStorage struct {
	Path string
}

func (a *ArtifactoryStorage) New(p Type) {
	panic("implement me")
}

func (a *ArtifactoryStorage) NewStorage() {
	logrus.Info("Init new Artifactory storage")
}

func (a *ArtifactoryStorage) PutFile(packageName string, packageVersion string, content []byte) (Response, error) {
	logrus.Info("Put a file to Artifactory storage")
	return Response{message:"Ok"}, nil
}

func (a *ArtifactoryStorage) GetFile(filename string) ([]byte, error) {
	logrus.Info("Get a file from Artifactory storage")
	return []byte{}, nil
}

