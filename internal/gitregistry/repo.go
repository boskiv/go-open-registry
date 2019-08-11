package gitregistry

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-open-registry/internal/config"
	"go-open-registry/internal/helpers"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path"
	"strings"
	"time"
)

// New instance of git repository
func New(appConfig *config.AppConfig) *git.Repository {
	logrus.WithFields(logrus.Fields{
		"repo": appConfig.Repo.URL,
	}).Info("Init repo started")
	repo, err := git.PlainOpen(appConfig.Repo.Path)
	if repo == nil {
		logrus.Info("Repo folder does not exist, make clone")
		repo, err = git.PlainClone(appConfig.Repo.Path, false, &git.CloneOptions{
			URL:  appConfig.Repo.URL,
			Auth: &http.BasicAuth{Username: appConfig.Repo.Bot.Name, Password: appConfig.Repo.Bot.Password},
		})
	}

	helpers.FatalIfError(err)
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
	logrus.WithFields(logrus.Fields{
		"package": packageName,
		"version": packageVersion,
		"size":    len(content),
	}).Info("Commit function called")

	// crate folder structure
	result, err := makePath(appConfig, packageName)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"folder": result,
	}).Info("Folder created")

	// create a file in path
	result, err = createFile(result, content)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"file": result,
	}).Info("File created")

	// Commit file to git
	result, err = commitFile(appConfig, packageName, packageVersion)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"commit": result,
	}).Info("File committed")

	// push repo to origin
	result, err = pushRegistryRepo(appConfig)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"push": result,
	}).Info("Changes pushed")
	return nil
}

// takes config
// push repo changes to origin
// return updated last commit hash or error
func pushRegistryRepo(appConfig *config.AppConfig) (result string, err error) {
	err = appConfig.Repo.Instance.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: appConfig.Repo.Bot.Name,
			Password: appConfig.Repo.Bot.Password,
		},
	})
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
		logrus.WithField("error", err).Error("Error adding to local repo")
		return "", err
	}

	logrus.WithField("file", resultPathString).Info("File added to local repo")

	commitMsg := fmt.Sprintf("Commit package %s version %s",
		packageName, packageVersion)
	logrus.Info("Commit file to repo")
	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  appConfig.Repo.Bot.Name,
			Email: appConfig.Repo.Bot.Email,
			When:  time.Now(),
		},
	})

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

	return resultPathString, nil
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
	logrus.WithFields(logrus.Fields{
		"directory": crateDir,
		"file":      crateFile,
	}).Info("Got path")
	// create dir tree
	err = os.MkdirAll(crateDir, os.ModePerm)
	return resultPathString, err
}
