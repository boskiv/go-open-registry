package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-open-registry/internal/config"
	"net/http"
)



// NewCrateHandler to serve cargo publish command
func GetCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		// /api/v1/crates/bo-helper/0.1.2/download
		name := c.Param("name")
		version := c.Param("version")


		logrus.WithFields(logrus.Fields{
			"name": name,
			"version": version,
		}).Info("Got request")

		crateFile, err := appConfig.Storage.Instance.GetFile(name, version)

		if err != nil {
			logrus.Error(err)
		}
		h := sha256.New()
		h.Write(crateFile)
		cksum := hex.EncodeToString(h.Sum(nil))
		logrus.WithField("cksum", cksum).Info("Set cksum")
		crateFileReader := bytes.NewReader(crateFile)
		contentLength := int64(len(crateFile))
		contentType := "Content-Type: multipart/form-data; boundary=something"
		filename := name + "-" + version + ".crate"
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="`+ filename +`"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, crateFileReader, extraHeaders)
	}
}