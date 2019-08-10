package storage

import "github.com/sirupsen/logrus"

type Response struct {
	message string
}

type GenericStorage interface {
	PutFile(packageName string, packageVersion string, content []byte) (Response, error)
	GetFile(filename string) ([]byte, error)
}


func New(p Type) GenericStorage {
	switch p {
	case Local:
		return &LocalStorage{}
	case S3:
		return &S3Storage{}
	case Artifactory:
		return &ArtifactoryStorage{}
	default:
		logrus.Fatal("No storage %s defined", p.String())
	}
	return nil
}

type Type int
const (
	Local Type = iota
	S3
	Artifactory
)

func (name Type) String() string {
	names := [...]string{
		"Local",
		"S3",
		"Artifactory",
		}
	if name < Local || name > Artifactory {
		return "Unknown"
	}
	return names[name]
}