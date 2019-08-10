package helpers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}
	logrus.Fatal(err)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func MakeCratePath(packageName string) []string {
	var path []string

	switch len(packageName) {
	case 1:
		path = append(path, "1")
	case 2:
		path = append(path, "2")
	case 3:
		path = append(path, "3", packageName[0:1])
	default:
		path = append(path, packageName[0:2], packageName[2:4])
	}

	return path
}

func FullCratePath() {
	//withUploadDir := append([]string{localStorage.path}, paths...)
	//_ = os.MkdirAll(strings.Join(withUploadDir, string(os.PathSeparator)), os.ModePerm)
	//withPackageName := append(withUploadDir, packageName)
}
