package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go-open-registry/internal/config"
	"go-open-registry/internal/log"
	"net/http"
)

// GetCrateHandler to serve cargo publish command
func GetCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		// /api/v1/crates/bo-helper/0.1.2/download
		name := c.Param("name")
		version := c.Param("version")

		log.InfoWithFields("Got request",log.Fields{
			"name":    name,
			"version": version,
		})

		crateFile, err := appConfig.Storage.Instance.GetFile(name, version)

		if err != nil {
			log.ErrorWithFields("Error getting file from storage", log.Fields{"err": err})
		}
		h := sha256.New()
		_, err = h.Write(crateFile)
		if err != nil {
			log.ErrorWithFields("Error while writing file",log.Fields{"error": err})
			c.JSON(400, gin.H{
				"error": err,
			})
			return
		}

		cksum := hex.EncodeToString(h.Sum(nil))
		log.InfoWithFields("Cksum get", log.Fields{"cksum": cksum})
		crateFileReader := bytes.NewReader(crateFile)
		contentLength := int64(len(crateFile))
		contentType := "Content-Type: multipart/form-data; boundary=something"
		filename := name + "-" + version + ".crate"
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="` + filename + `"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, crateFileReader, extraHeaders)
	}
}
