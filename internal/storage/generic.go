package storage

import "github.com/sirupsen/logrus" //nolint:depguard

// Response struct to return from storage requests
type Response struct {
	message string
}

// GenericStorage interface, for method implementation
type GenericStorage interface {
	PutFile(packageName, packageVersion string, content []byte) (Response, error)
	GetFile(filename string) ([]byte, error)
}

// New storage by Type
func New(p Type) GenericStorage {
	switch p {
	case Local:
		return &LocalStorage{}
	case S3:
		return &S3Storage{}
	case Artifactory:
		return &ArtifactoryStorage{}
	default:
		logrus.WithField("storage", p.String()).Fatal("No storage defined")
	}
	return nil
}

// Type storage Enum
type Type int

const (
	// Local Storage type
	Local Type = iota
	// S3 storage type
	S3
	// Artifactory storage type
	Artifactory
)

// Storage type to String
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
