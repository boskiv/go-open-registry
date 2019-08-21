package gitregistry

import (
	"encoding/json"
	"fmt"
	"go-open-registry/internal/config"
	"go-open-registry/internal/helpers"
	"go-open-registry/internal/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type cargoConfig struct {
	DownloadURL string `json:"dl"`
	APIUrl      string `json:"api"`
}

// New instance of git repository
func New(appConfig *config.AppConfig) *git.Repository {
	log.InfoWithFields("Init repo started", log.Fields{
		"repo": appConfig.Repo.URL,
	})
	repo, err := git.PlainOpen(appConfig.Repo.Path)
	if err != nil {
		log.ErrorWithFields("Error while open repo", log.Fields{
			"err": err,
		})
	}
	if repo == nil {
		log.Info("Repo folder does not exist, make clone")
		repo, err = git.PlainClone(appConfig.Repo.Path, false, &git.CloneOptions{
			URL:  appConfig.Repo.URL,
			Auth: &http.BasicAuth{Username: appConfig.Repo.Auth.Name, Password: appConfig.Repo.Auth.Password},
		})
		if err != nil {
			log.FatalWithFields("Error while clone repo", log.Fields{
				"err": err,
			})
		}
	}

	return repo
}

// RegistryAdd a git wrapper
// * create path in git repo,
// * create file with content,
// * commit it to repo
// * push commit to git remote origin
func RegistryAdd(
	appConfig *config.AppConfig,
	packageName string,
	packageVersion string,
	content []byte) error {
	log.InfoWithFields("RegistryAdd called", log.Fields{
		"package": packageName,
		"version": packageVersion,
		"size":    len(content),
	})

	// crate folder structure
	result, err := makePath(appConfig, packageName)
	if err != nil {
		log.ErrorWithFields("Error make path", log.Fields{
			"err": err,
		})
		return err
	}
	log.InfoWithFields("Folder created", log.Fields{
		"folder": result,
	})

	// create a file in path
	result, err = createFile(result, content)
	if err != nil {
		log.ErrorWithFields("Error create file", log.Fields{
			"err": err,
		})
		return err
	}
	log.InfoWithFields("File created", log.Fields{
		"file": result,
	})

	// Commit file to git
	result, err = commitFile(appConfig, packageName, packageVersion)
	if err != nil {
		log.ErrorWithFields("Error commit to repo", log.Fields{
			"err": err,
		})
		return err
	}
	log.InfoWithFields("File committed", log.Fields{
		"commit": result,
	})

	// push repo to origin
	result, err = pushRegistryRepo(appConfig)
	if err != nil {
		log.ErrorWithFields("Error pushing to repo", log.Fields{
			"err": err,
		})
		return err
	}
	log.InfoWithFields("Changes pushed", log.Fields{
		"push": result,
	})
	return nil
}

// InitConfig recommit config file every time on start
func InitConfig(appConfig *config.AppConfig) (err error) {
	configFileName := "config.json"
	pathToConfig := path.Join(appConfig.Repo.Path, configFileName)
	fileContent, err := os.OpenFile(pathToConfig, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		// file not found
		log.Error(err)
	}

	// read json
	var jsonConfig cargoConfig

	var bytes []byte
	bytes, err = ioutil.ReadFile(fileContent.Name())
	if err != nil {
		log.Error(err)
	}
	defer fileContent.Close()

	err = json.Unmarshal(bytes, &jsonConfig)
	if err != nil {
		log.ErrorWithFields("Bad json file,rewriting", log.Fields{"err": err})
		jsonConfig = cargoConfig{
			DownloadURL: "",
			APIUrl:      "",
		}

	} else {
		if jsonConfig.APIUrl == appConfig.App.CargoAPIURL &&
			jsonConfig.DownloadURL == appConfig.App.CargoDownloadURL {
			return err
		}
	}

	err = os.Truncate(fileContent.Name(), 0)
	if err != nil {
		log.Fatal(err)
	}

	jsonConfig.APIUrl = appConfig.App.CargoAPIURL
	jsonConfig.DownloadURL = appConfig.App.CargoDownloadURL

	bytes, err = json.Marshal(jsonConfig)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = fileContent.Write(bytes)
	if err != nil {
		log.Error(err)
		return err
	}
	defer fileContent.Close()

	w, err := appConfig.Repo.Instance.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(configFileName)
	if err != nil {
		log.ErrorWithFields("Error adding to local repo", log.Fields{"error": err})
		return err
	}

	log.InfoWithFields("File added to local repo", log.Fields{"file": configFileName})

	commitMsg := fmt.Sprintf("Commit config %s", configFileName)
	log.InfoWithFields("Commit config to repo", log.Fields{
		"message": commitMsg,
	})
	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  appConfig.Repo.Auth.Name,
			Email: appConfig.Repo.Auth.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = appConfig.Repo.Instance.CommitObject(commit)
	if err != nil {
		return err
	}
	_, err = pushRegistryRepo(appConfig)
	if err != nil {
		return err

	}
	return err
}

// takes config
// push repo changes to origin
// return updated last commit hash or error
func pushRegistryRepo(appConfig *config.AppConfig) (result string, err error) {
	err = appConfig.Repo.Instance.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: appConfig.Repo.Auth.Name,
			Password: appConfig.Repo.Auth.Password,
		},
	})
	if err != nil {
		log.ErrorWithFields("Error pushing to repo", log.Fields{
			"err": err,
		})
		return result, err
	}

	ref, err := appConfig.Repo.Instance.Head()
	if err != nil {
		return result, err
	}
	result = ref.Hash().String()
	return result, err
}

// commit file to gir repository
func commitFile(appConfig *config.AppConfig, packageName, packageVersion string) (result string, err error) {
	folderStructure := helpers.MakeCratePath(packageName)
	var resultPath []string
	resultPath = append(resultPath, folderStructure...)
	resultPath = append(resultPath, packageName)
	resultPathString := strings.Join(resultPath, string(os.PathSeparator))
	w, err := appConfig.Repo.Instance.Worktree()
	if err != nil {
		return result, err
	}
	_, err = w.Add(resultPathString)
	if err != nil {
		log.ErrorWithFields("Error adding to local repo", log.Fields{"error": err})
		return "", err
	}

	log.InfoWithFields("File added to local repo", log.Fields{"file": resultPathString})

	commitMsg := fmt.Sprintf("Commit package %s version %s",
		packageName, packageVersion)
	log.InfoWithFields("Commit file to repo", log.Fields{
		"message": commitMsg,
	})
	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  appConfig.Repo.Auth.Name,
			Email: appConfig.Repo.Auth.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return result, err
	}

	obj, err := appConfig.Repo.Instance.CommitObject(commit)
	if err != nil {
		return result, err
	}
	result = obj.Hash.String()
	return result, err
}

// append text to file by provided path and []byte content
// if file does not exist it will be created
// line will end with \n character
func createFile(resultPathString string, content []byte) (result string, err error) {
	f, err := os.OpenFile(resultPathString,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(string(content) + "\n")
	//removedGitFolder := strings.Split(resultPathString, string(os.PathSeparator))

	return resultPathString, err
}

// make path in git repository
// following https://doc.rust-lang.org/cargo/reference/registries.html#index-format
func makePath(appConfig *config.AppConfig, packageName string) (result string, err error) {
	folderStructure := helpers.MakeCratePath(packageName)
	var resultPath []string
	resultPath = append(resultPath, appConfig.Repo.Path)
	resultPath = append(resultPath, folderStructure...)
	resultPath = append(resultPath, packageName)
	resultPathString := strings.Join(resultPath, string(os.PathSeparator))
	crateDir, crateFile := path.Split(resultPathString)
	log.InfoWithFields("Got path", log.Fields{
		"directory": crateDir,
		"file":      crateFile,
	})
	// create dir tree
	err = os.MkdirAll(crateDir, os.ModePerm)
	return resultPathString, err
}
