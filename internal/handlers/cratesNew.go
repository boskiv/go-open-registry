package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go-open-registry/internal/config"
	"go-open-registry/internal/gitregistry"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/log"
	"go-open-registry/internal/parser"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// NewCrateHandler to serve cargo publish command
func NewCrateHandler(appConfig *config.AppConfig) func(c *gin.Context) {

	return func(c *gin.Context) {
		// Read the Body content
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			jsonFile, crateFile, err := parser.ReadBinary(c.Request.Body)
			cksum := helpers.CheckSHA256Sum(crateFile)
			if err != nil {
				log.ErrorWithFields("Error while parser.ReadBinary repo", log.Fields{
					"err": err,
				})
			}
			var crateJSON parser.CrateJSON
			err = json.Unmarshal(jsonFile, &crateJSON)
			if err != nil {
				log.ErrorWithFields("Error while json.Unmarshal repo", log.Fields{
					"err": err,
				})
			}
			crateJSON.Cksum = cksum
			log.InfoWithFields("Content cksum", log.Fields{
				"cksum": cksum,
			})

			jsonFileWithCksum, err := json.Marshal(crateJSON)
			if err != nil {
				log.ErrorWithFields("Error 400 throw",log.Fields{"error": err})
				c.JSON(400, gin.H{
					"error": err,
				})
				return
			}

			err = addDBVersion(appConfig, crateJSON)
			if err != nil {
				log.ErrorWithFields("Error 400 throw",log.Fields{"error": err})
				c.JSON(400, gin.H{
					"error": err,
				})
				return
			}

			err = registryAdd(appConfig, crateJSON, jsonFileWithCksum)
			if err != nil {
				log.ErrorWithFields("Error 400 throw",log.Fields{"error": err})
				c.JSON(400, gin.H{
					"error": err,
				})
				return
			}

			err = storagePut(appConfig, crateJSON, crateFile)
			if err != nil {
				log.ErrorWithFields("Error 400 throw",log.Fields{"error": err})
				c.JSON(400, gin.H{
					"error": err,
				})
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

func storagePut(appConfig *config.AppConfig, crateJSON parser.CrateJSON, crateFile []byte) (err error) {
	err = appConfig.Storage.Instance.PutFile(crateJSON.Name, crateJSON.Vers, crateFile)
	if err != nil {
		err = rollBackDBVersion(appConfig, crateJSON)
		return err
	}
	return err
}

func registryAdd(appConfig *config.AppConfig, crateJSON parser.CrateJSON, jsonFile []byte) (err error) {
	err = gitregistry.RegistryAdd(appConfig, crateJSON.Name, crateJSON.Vers, jsonFile)
	log.ErrorWithFields("Error while add file to registry", log.Fields{
		"err": err,
	})
	if err != nil {
		err = rollBackDBVersion(appConfig, crateJSON)
		return err
	}
	return err
}

func addDBVersion(appConfig *config.AppConfig, crateJSON parser.CrateJSON) (err error) {
	// Validate version
	collection := appConfig.DB.Client.Database("crates").Collection("packages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
	if err != nil {
		return err
	}
	if res != nil {
		id := res.InsertedID
		log.InfoWithFields("Package version added to mongo", log.Fields{"id": id})
	}
	return nil
}

func rollBackDBVersion(appConfig *config.AppConfig, crateJSON parser.CrateJSON) (err error) {
	// Validate version
	log.InfoWithFields("Rolling back record", log.Fields{
		"package": crateJSON.Name,
		"version": crateJSON.Vers,
	})

	collection := appConfig.DB.Client.Database("crates").Collection("packages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"name": crateJSON.Name, "version": crateJSON.Vers})
	if err != nil {
		log.ErrorWithFields("rollBackDBVersion: Error while deleting key from db", log.Fields{
			"err": err,
		})
		return err
	}
	if res != nil {
		err = errors.New("rollBackDBVersion: previous step failed")
		count := res.DeletedCount
		log.InfoWithFields("Deleted from database", log.Fields{"count": count})
	}
	return err
}
