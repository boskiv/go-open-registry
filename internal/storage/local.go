package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"go-open-registry/internal/log"
	"os"
	"path"
	"strings"
) //nolint:depguard

// LocalStorage type struct
type LocalStorage struct {
	Path string
}

// PutFile implementation
func (l *LocalStorage) PutFile(packageName, packageVersion string, content []byte) error {
	log.Info("Put a file to local storage")
	log.InfoWithFields("Got package upload request",log.Fields{
		"package": packageName,
		"version": packageVersion,
	})

	var resultPath []string
	resultPath = append(resultPath, l.Path)
	resultPath = append(resultPath, packageName)
	resultPath = append(resultPath, packageVersion)
	resultPath = append(resultPath, packageName+"-"+packageVersion+".crate")
	resultPathString := strings.Join(resultPath, string(os.PathSeparator))

	crateDir, crateFile := path.Split(resultPathString)
	log.InfoWithFields("Got path",log.Fields{
		"directory": crateDir,
		"file":      crateFile,
	})
	// create dir tree
	err := os.MkdirAll(crateDir, os.ModePerm)
	if err != nil {
		log.ErrorWithFields("Error mkdir", log.Fields{
			"err": err,
		})
		return err
	}

	// write file

	f, err := os.Create(resultPathString)

	if err != nil {
		log.ErrorWithFields("Error create file", log.Fields{
			"err": err,
		})
		return err
	}

	defer f.Close()

	h := sha256.New()
	_, err = h.Write(content)
	if err != nil {
		log.ErrorWithFields("Error make hash", log.Fields{
			"err": err,
		})
	}
	cksum := hex.EncodeToString(h.Sum(nil))
	log.InfoWithFields("Content cksum from PutFile", log.Fields{
		"cksum": cksum,
	})
	if _, err := f.Write(content); err != nil {

		log.ErrorWithFields("Error write file", log.Fields{
			"err": err,
		})
		return err
	}

	return nil
}

// GetFile implementation
func (l *LocalStorage) GetFile(packageName, packageVersion string) ([]byte, error) {
	filename := packageName + "-" + packageVersion + ".crate"
	filenamePath := path.Join(l.Path, packageName, packageVersion, filename)
	log.InfoWithFields("Get a file from local storage", log.Fields{
		"path": filenamePath,
	})
	crateFile, err := os.Open(filenamePath)
	if err != nil {
		log.ErrorWithFields("Error open file path", log.Fields{
			"err": err,
		})
		return nil, err
	}
	defer crateFile.Close()
	stats, statsErr := crateFile.Stat()
	if statsErr != nil {
		return nil, statsErr
	}
	log.InfoWithFields("", log.Fields{
		"stats": stats,
	})
	var size = stats.Size()
	crateFileBytes := make([]byte, size)

	buffer := bufio.NewReader(crateFile)
	_, err = buffer.Read(crateFileBytes)

	return crateFileBytes, err
}
