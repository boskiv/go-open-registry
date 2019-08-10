package handlers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// NewCrateHandler to serve cargo publish command
func NewCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {

	return func(c *gin.Context) {
		gitregistry.HeadRepo(appConfig.Repo.Instance)
		// Read the Body content
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			jsonFile, crateFile, err := parser.ReadBinary(c.Request.Body)
			cksum := helpers.CheckSum(crateFile)
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
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

			// Todo: Refactor it smart
			done, commitError := registryCommit(ctx, appConfig, crateJSON, jsonFileWithCksum, collection, c)
			if done {
				return
			}
			if storagePut(ctx, appConfig, crateJSON, crateFile, collection, c, commitError) {
				return
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

func storagePut(ctx context.Context, appConfig *config.AppConfig, crateJSON parser.CrateJSON, crateFile []byte, collection *mongo.Collection, c *gin.Context, commitError error) bool {
	storageError := appConfig.Storage.Instance.PutFile(crateJSON.Name, crateJSON.Vers, crateFile)
	if storageError != nil {
		logrus.WithField("commitError", commitError).Info("Error while commit to git")
		res, err := collection.DeleteOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
		if err != nil {
			c.JSON(400, gin.H{
				"error": storageError,
			})
			return true
		}
		if res != nil {
			count := res.DeletedCount
			logrus.WithField("count", count).Info("Deleted from database")
		}

	}
	return false
}

func registryCommit(ctx context.Context, appConfig *config.AppConfig, crateJSON parser.CrateJSON, jsonFileWithCksum []byte, collection *mongo.Collection, c *gin.Context) (bool, error) {
	commitError := gitregistry.CommitCrateJSON(appConfig, crateJSON.Name, crateJSON.Vers, jsonFileWithCksum)
	if commitError != nil {
		logrus.WithField("commitError", commitError).Info("Error while commit to git")
		res, err := collection.DeleteOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
		if err != nil {
			c.JSON(400, gin.H{
				"error": err,
			})
			return true, nil
		}
		if res != nil {
			count := res.DeletedCount
			logrus.WithField("count", count).Info("Deleted from database")
		}

	}
	return false, commitError
}
