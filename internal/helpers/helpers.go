package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus" //nolint:depguard
)

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// FatalIfError should be used to naively panics if an error is not nil.
func FatalIfError(err error) {
	if err == nil {
		return
	}
	logrus.Fatal(err)
}

// CheckSum SHA256 of []byte return as string
func CheckSum(content []byte) string {
	h := sha256.New()
	_, err := h.Write(content)
	if err != nil {
		logrus.Error(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// MakeCratePath slice with package name parts
// https://doc.rust-lang.org/cargo/reference/registries.html#index-format
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

// FullCratePath Return full path with parts from MakeCratePath and package name
func FullCratePath() {
	//withUploadDir := append([]string{localStorage.path}, paths...)
	//_ = os.MkdirAll(strings.Join(withUploadDir, string(os.PathSeparator)), os.ModePerm)
	//withPackageName := append(withUploadDir, packageName)
}
