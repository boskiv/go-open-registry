package storage

import (
	"bytes"
	"fmt"
	"go-open-registry/internal/log"
	"io/ioutil"
	"net/http"
)

// ArtifactoryStorage struct
type ArtifactoryStorage struct {
	Path     string `json:"path"`
	Login    string `json:"login"`
	Password string `json:"-"`
	URL      string `json:"url"`
	RepoName string `json:"repo_name"`
}

// PutFile implementation
func (a ArtifactoryStorage) PutFile(packageName, packageVersion string, content []byte) error {
	log.Info("Put a file to artifactory storage")
	uri := a.Path + "/" + packageName + "-" + packageVersion + ".crate"
	reader := bytes.NewReader(content)
	req, err := http.NewRequest("PUT", uri, reader)
	if err != nil {
		log.ErrorWithFields("Error new request", log.Fields{"err": err})
		return err
	}

	req.SetBasicAuth(a.Login, a.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.ErrorWithFields("Error new request", log.Fields{"err": err})
		return err
	}
	defer resp.Body.Close()
	log.DebugWithFields("Response", log.Fields{"response": resp})
	switch resp.StatusCode {
	case http.StatusCreated:
		log.Info("File uploaded")
		return nil
	case http.StatusForbidden:
		err = fmt.Errorf("respone code:%d status: %s", resp.StatusCode, resp.Status)
		log.ErrorWithFields("File already exist, but no rights to rewrite it or you don't have write rights to repo at all", log.Fields{"err": err})
		return err
	case http.StatusUnauthorized:
		err = fmt.Errorf("respone code:%d status: %s", resp.StatusCode, resp.Status)
		return err
	default:
		err = fmt.Errorf("respone code:%d status: %s", resp.StatusCode, resp.Status)
		log.ErrorWithFields("Error new request", log.Fields{"err": err})
		return err
	}
}

// GetFile implementation
func (a ArtifactoryStorage) GetFile(packageName, packageVersion string) ([]byte, error) {
	log.Info("Get a file from artifactory storage")
	uri := a.Path + "/" + packageName + "-" + packageVersion + ".crate"
	client := &http.Client{}
	resp, err := client.Get(uri)
	if err != nil {
		log.ErrorWithFields("Error new request", log.Fields{"err": err})
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("respone code:%d status: %s", resp.StatusCode, resp.Status)
		log.ErrorWithFields("Error new request", log.Fields{"err": err})
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorWithFields("Error read bytes from body", log.Fields{"err": err})
	}

	return bodyBytes, nil
}
