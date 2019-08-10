package handlers

import (
	"encoding/json"
	"fmt"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/parser"

	"github.com/gin-gonic/gin"
)

// NewCrateHandler to serve cargo publish command
func NewCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {

	return func(c *gin.Context) {
		gitregistry.HeadRepo(appConfig.Repo.Instance)
		// Read the Body content
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			jsonFile, crateFile, err := parser.ReadBinary(c.Request.Body)
			helpers.CheckIfError(err)
			fmt.Printf("%s", jsonFile)
			var crateJSON parser.CrateJSON
			err = json.Unmarshal(jsonFile, &crateJSON)
			helpers.CheckIfError(err)

			gitregistry.CommitCrateJSON(appConfig, crateJSON.Name, crateJSON.Vers)
			_, _ = appConfig.Storage.Instance.PutFile(crateJSON.Name, crateJSON.Vers, crateFile)
		}

		resp := map[string][]string{
			// Array of strings of categories that are invalid and ignored.
			"invalid_categories": {},
			// Array of strings of badge names that are invalid and ignored.
			"invalid_badges": {},
			// Array of strings of arbitrary warnings to display to the user.
			"other": {},
		}
		c.JSON(200, gin.H{
			// Optional object of warnings to display to the user.
			"warnings": resp,
		})
	}
}
