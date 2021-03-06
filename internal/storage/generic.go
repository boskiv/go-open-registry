package storage

// GenericStorage interface, for method implementation
type GenericStorage interface {
	PutFile(packageName, packageVersion string, content []byte) error
	GetFile(packageName, packageVersion string) ([]byte, error)
}

//// New storage by Type
//func New(p Type, path string, login, password string) GenericStorage {
//	switch p {
//	case Local:
//		return &LocalStorage{Path: path}
//	case S3:
//		return &S3Storage{Path: path, AccessKey: login, SecretKey: password}
//	case Artifactory:
//		log.InfoWithFields("Path", log.Fields{"path": path})
//		return &ArtifactoryStorage{Path: path, Login: login, Password: password}
//
//	default:
//		return nil
//	}
//
//}

// Type storage Enum
type Type int

const (
	// Local Storage type
	Local Type = iota
	// S3 storage type
	S3
	// Artifactory storage type
	Artifactory
	// Unknown storage type
	Unknown
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
