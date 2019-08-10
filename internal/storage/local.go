package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"github.com/sirupsen/logrus"
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
	logrus.Info("Put a file to local storage")
	logrus.WithFields(logrus.Fields{
		"package": packageName,
		"version": packageVersion,
	}).Info("Got package upload request")


	var resultPath []string
	resultPath = append(resultPath, l.Path)
	resultPath = append(resultPath, packageName)
	resultPath = append(resultPath, packageVersion)
	resultPath = append(resultPath, packageName + "-" + packageVersion + ".crate")
	resultPathString := strings.Join(resultPath, string(os.PathSeparator))

	crateDir, crateFile := path.Split(resultPathString)
	logrus.WithFields(logrus.Fields{
		"directory": crateDir,
		"file":      crateFile,
	}).Info("Got path")
	// create dir tree
	err := os.MkdirAll(crateDir, os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// write file

	f, err := os.Create(resultPathString)

	if err != nil {
		logrus.Error(err)
		return err
	}

	defer f.Close()

	h := sha256.New()
	h.Write(content)
	cksum := hex.EncodeToString(h.Sum(nil))
	logrus.WithField("cksum", cksum).Info("Content cksum from PutFile")

	if _, err := f.Write(content); err != nil {

		logrus.Error(err)
		return err
	}

	return nil
}

// GetFile implementation
func (l *LocalStorage) GetFile(packageName, packageVersion string ) ([]byte, error) {
	logrus.Info("Get a file from local storage")
	filename := packageName + "-" + packageVersion + ".crate"
	filenamePath := path.Join(l.Path, packageName, packageVersion, filename)
	crateFile, err := os.Open(filenamePath)
	if err != nil {
		logrus.Error(err)
	}
	defer crateFile.Close()
	stats, statsErr := crateFile.Stat()
	if statsErr != nil {
		return nil, statsErr
	}
	logrus.Info(stats)
	var size = stats.Size()
	crateFileBytes := make([]byte, size)

	bufr := bufio.NewReader(crateFile)
	_,err = bufr.Read(crateFileBytes)

	return crateFileBytes, nil
}
