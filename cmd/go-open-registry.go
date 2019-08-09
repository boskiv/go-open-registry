package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"go-open-registry/internal/config"
	"go-open-registry/internal/git"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/parser"
)

var appConfig *config.AppConfig

type CrateDependency struct {
	Name            string   `json:"name"`
	Req             string   `json:"req"`
	Features        []string `json:"features"`
	Optional        bool     `json:"optional"`
	DefaultFeatures bool     `json:"default_features"`
	Target          string   `json:"target"`
	Kind            string   `json:"kind"`
	Registry        string   `json:"registry"`
	Package         string   `json:"package"`
}

type CrateJson struct {
	Name     string            `json:"name"`
	Vers     string            `json:"vers"`
	Deps     []CrateDependency `json:"deps"`
	Cksum    string            `json:"cksum"`
	Features interface{}       `json:"features"`
	Yanked   bool              `json:"yanked"`
	Links    string            `json:"links"`
}

func NewCrateHandler(c *gin.Context) {
	// Read the Body content
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		jsonFile, crateFile, err := parser.ReadBinary(c.Request.Body)
		helpers.CheckIfError(err)
		fmt.Printf("%s", jsonFile)
		var crateJson CrateJson
		err = json.Unmarshal(jsonFile, &crateJson)
		helpers.CheckIfError(err)
		_, _ = appConfig.Storage.PutFile(crateJson.Name, jsonFile)
		_, _ = appConfig.Storage.PutFile(crateJson.Name+"-"+crateJson.Vers+".crate", crateFile)
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

func main() {
	appConfig = config.InitConfig()

	log := logrus.New()

	repo, err := git.InitRepo()
	helpers.CheckIfError(err)
	git.HeadRepo(repo)

	r := gin.New()

	r.Use(ginlogrus.Logger(log), gin.Recovery())

	r.PUT("/api/v1/crates/new", NewCrateHandler)
	_ = r.Run()
}
