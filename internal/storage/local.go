package storage

type LocalStorage struct {
	Path string
}

func (l *LocalStorage) PutFile(packageName string, packageVersion string, content []byte) (Response, error) {
	panic("implement me")
	return Response{message:"Ok"}, nil
}

func (l *LocalStorage) GetFile(filename string) ([]byte, error) {
	panic("implement me")
	return []byte{}, nil
}



