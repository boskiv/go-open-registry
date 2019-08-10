package gitregistry

import (
	"fmt"
	"go-open-registry/internal/config"
	"go-open-registry/internal/helpers"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus" //nolint:depguard
	"gopkg.in/src-d/go-git.v4"
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
			URL: appConfig.Repo.URL,
			Auth: &http.BasicAuth{Username: appConfig.Repo.Bot.Name , Password: appConfig.Repo.Bot.Password},
		})
	}



	helpers.FatalIfError(err)
	return repo
}

// HeadRepo information of repository
// for example current branch Name and last Commit Hash
func HeadRepo(repo *git.Repository) {
	result, err := repo.Head()
	helpers.FatalIfError(err)
	helpers.Info("%s: %s", result.Name(), result.Hash())
}

// CommitCrateJSON information form crate file to git registry
// if file exist, information will be append to it
// It also manage directory structure followed by
// https://doc.rust-lang.org/cargo/reference/registries.html#index-format
func CommitCrateJSON(appConfig *config.AppConfig, packageName string, packageVersion string, content []byte) error {
	logrus.WithFields(logrus.Fields{
		"package": packageName,
		"version": packageVersion,
		"size":    len(content),
	}).Info("Commit function called")
	r := appConfig.Repo.Instance
	logrus.Info(r)
	// Get slice of directories to append to git registry root
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
	err := os.MkdirAll(crateDir, os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// write file

	f, err := os.OpenFile(resultPathString,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		logrus.Error(err)
		return err
	}

	defer f.Close()
	if _, err := f.WriteString(string(content) + "\n"); err != nil {

		logrus.Error(err)
		return err
	}
	logrus.Info("Getting git work tree")
	w, err := r.Worktree()
	if err != nil {
		logrus.Error(err)
		return err
	}

	addPackageName := append(folderStructure, packageName)
	commitPath := strings.Join(addPackageName, string(os.PathSeparator))
	logrus.WithField("path",commitPath).Info("Add file to stage")
	_, err = w.Add(commitPath)
	if err != nil {
		logrus.Error(err)
		return err
	}


	logrus.Info("Commit file to repo")
	commit, err := w.Commit(fmt.Sprintf("Commit package %s version %s",packageName, packageVersion), &git.CommitOptions{
		Author: &object.Signature{
			Name:  appConfig.Repo.Bot.Name,
			Email: appConfig.Repo.Bot.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Info("Getting new head")
	obj, err := r.CommitObject(commit)
	if err != nil {
		logrus.Error(err)
		return err
	}
	//
	logrus.Info(obj)



	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: appConfig.Repo.Bot.Name,
			Password: appConfig.Repo.Bot.Password,
		},
	})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return err
}
