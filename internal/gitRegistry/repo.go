package gitRegistry

import (
	"github.com/sirupsen/logrus"
	"go-open-registry/internal/config"
	"go-open-registry/internal/helpers"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path"
	"strings"
)

func New(url string) *git.Repository {
	logrus.WithFields(logrus.Fields{
		"repo": url,
	}).Info("Init repo started")
	repo, err := git.PlainOpen("tmpGit")
	if repo == nil {
		logrus.Info("Repo folder does not exist, make clone")
		repo, err = git.PlainClone("tmpGit", false, &git.CloneOptions{
			URL: url,
		})
	}

	helpers.CheckIfError(err)
	return repo
}

func HeadRepo(repo *git.Repository) {
	result, err := repo.Head()
	helpers.CheckIfError(err)
	helpers.Info("%s: %s",result.Name(), result.Hash())
}

func CommitCrateJson(appConfig *config.AppConfig,packageName string, content string) {
	r := appConfig.Repo.Instance
	logrus.Info(r)
	var fullJsonCratePath []string
	fullJsonCratePath = append(fullJsonCratePath, appConfig.Repo.Path)
	crateJsonPath := strings.Join(fullJsonCratePath,string(os.PathSeparator))

	//paths := helpers.MakeCratePath(packageName)

	crateDir, crateFile := path.Split(crateJsonPath)
	logrus.WithFields(logrus.Fields{
		"directory": crateDir,
		"file": crateFile,
	}).Info("Got path")
	// create dir tree
	err := os.MkdirAll(crateDir,os.ModePerm)
	helpers.CheckIfError(err)

	// write file
	//f, err := os.OpenFile(crateJsonPath,
	//	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Println(err)
	//}
	//defer f.Close()
	//if _, err := f.WriteString("text to append\n"); err != nil {
	//	log.Println(err)
	//}
	//
	//w, err := r.Worktree()
	//helpers.CheckIfError(err)
	//
	//// ... we need a file to commit so let's create a new file inside of the
	//// worktree of the project using the go standard library.
	//helpers.Info("echo \"hello world!\" > example-gitRegistry-file")
	//filename := filepath.Join(directory, "example-gitRegistry-file")
	//err = ioutil.WriteFile(filename, []byte(content), 0644)
	//helpers.CheckIfError(err)
	//
	//// Adds the new file to the staging area.
	//helpers.Info("gitRegistry add example-gitRegistry-file")
	//_, err = w.Add("example-gitRegistry-file")
	//helpers.CheckIfError(err)
	//
	//// We can verify the current status of the worktree using the method Status.
	//helpers.Info("gitRegistry status --porcelain")
	//status, err := w.Status()
	//helpers.CheckIfError(err)
	//
	//fmt.Println(status)
	//
	//// Commits the current staging area to the repository, with the new file
	//// just created. We should provide the object.Signature of Author of the
	//// commit.
	//helpers.Info("gitRegistry commit -m \"example go-gitRegistry commit\"")
	//commit, err := w.Commit("example go-gitRegistry commit", &git.CommitOptions{
	//	Author: &object.Signature{
	//		Name:  "John Doe",
	//		Email: "john@doe.org",
	//		When:  time.Now(),
	//	},
	//})
	//
	//helpers.CheckIfError(err)
	//
	//// Prints the current HEAD to verify that all worked well.
	//helpers.Info("gitRegistry show -s")
	//obj, err := r.CommitObject(commit)
	//helpers.CheckIfError(err)
	//
	//fmt.Println(obj)
}
