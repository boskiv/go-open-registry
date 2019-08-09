package storage

type Response struct {
	message string
}

type GenericStorage interface {
	NewStorage()
	PutFile(filename string, content []byte) (Response, error)
	GetFile(filename string) ([]byte, error)
}
