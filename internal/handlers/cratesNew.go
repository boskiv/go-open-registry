package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/parser"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// NewCrateHandler to serve cargo publish command
func NewCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {

	return func(c *gin.Context) {
		gitregistry.HeadRepo(appConfig.Repo.Instance)
		// Read the Body content
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			jsonFile, crateFile, err := parser.ReadBinary(c.Request.Body)
			h := sha256.New()
			h.Write(crateFile)
			cksum := hex.EncodeToString(h.Sum(nil))
			helpers.FatalIfError(err)
			logrus.Debug(jsonFile)
			var crateJSON parser.CrateJSON
			err = json.Unmarshal(jsonFile, &crateJSON)
			helpers.FatalIfError(err)
			crateJSON.Cksum = cksum
			logrus.WithField("cksum", cksum).Info("Set cksum")

			jsonFileWithCksum, err := json.Marshal(crateJSON)

			// Validate version
			collection := appConfig.DB.Client.Database("crates").Collection("packages")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			res, err := collection.InsertOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
			if err != nil {
				logrus.WithField("error", err).Error("Error 400")
				c.JSON(400, gin.H{
					"error": err,
				})
				return

			}
			if res != nil {
				id := res.InsertedID
				logrus.WithField("id", id).Info("Package version added to mongo")
			}

			commitError := gitregistry.CommitCrateJSON(appConfig, crateJSON.Name, crateJSON.Vers, jsonFileWithCksum)
			if commitError != nil {
				logrus.WithField("commitError", commitError).Info("Error while commit to git")
				res, err := collection.DeleteOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
				if err != nil {
					c.JSON(400, gin.H{
						"error": err,
					})
					return
				}
				if res != nil {
					count := res.DeletedCount
					logrus.WithField("count", count).Info("Deleted from database")
				}

			}
			storageError := appConfig.Storage.Instance.PutFile(crateJSON.Name, crateJSON.Vers, crateFile)
			if storageError != nil {
				logrus.WithField("commitError", commitError).Info("Error while commit to git")
				res, err := collection.DeleteOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
				if err != nil {
					c.JSON(400, gin.H{
						"error": storageError,
					})
					return
				}
				if res != nil {
					count := res.DeletedCount
					logrus.WithField("count", count).Info("Deleted from database")
				}

			}
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
			"warnings": resp,
		})
	}
}
