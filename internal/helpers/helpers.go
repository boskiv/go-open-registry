package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"go-open-registry/internal/log"
)

// CheckSHA256Sum SHA256 of []byte return as string
func CheckSHA256Sum(content []byte) string {
	h := sha256.New()
	_, err := h.Write(content)
	if err != nil {
		log.ErrorWithFields("", log.Fields{
			"err": err,
		})
	}
	return hex.EncodeToString(h.Sum(nil))
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
