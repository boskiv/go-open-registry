package git

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-open-registry/internal/helpers"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func InitRepo() (*git.Repository, error) {
	logrus.WithFields(logrus.Fields{
		"repo": viper.Get("repo"),
	}).Info("Init repo started")

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: viper.GetString("repo"),
	})
	if err != nil {
		return repo, err
	}
	return repo, err
}

func HeadRepo(repo *git.Repository) {
	_, err := repo.Head()
	helpers.CheckIfError(err)
	helpers.Info("Done")
	//CommitRepo("hello world")
}

func CommitRepo(content string) {
	helpers.CheckArgs("<directory>")
	directory := os.Args[1]

	// Opens an already existing repository.
	r, err := InitRepo()
	helpers.CheckIfError(err)

	w, err := r.Worktree()
	helpers.CheckIfError(err)

	// ... we need a file to commit so let's create a new file inside of the
	// worktree of the project using the go standard library.
	helpers.Info("echo \"hello world!\" > example-git-file")
	filename := filepath.Join(directory, "example-git-file")
	err = ioutil.WriteFile(filename, []byte(content), 0644)
	helpers.CheckIfError(err)

	// Adds the new file to the staging area.
	helpers.Info("git add example-git-file")
	_, err = w.Add("example-git-file")
	helpers.CheckIfError(err)

	// We can verify the current status of the worktree using the method Status.
	helpers.Info("git status --porcelain")
	status, err := w.Status()
	helpers.CheckIfError(err)

	fmt.Println(status)

	// Commits the current staging area to the repository, with the new file
	// just created. We should provide the object.Signature of Author of the
	// commit.
	helpers.Info("git commit -m \"example go-git commit\"")
	commit, err := w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	helpers.CheckIfError(err)

	// Prints the current HEAD to verify that all worked well.
	helpers.Info("git show -s")
	obj, err := r.CommitObject(commit)
	helpers.CheckIfError(err)

	fmt.Println(obj)
}
