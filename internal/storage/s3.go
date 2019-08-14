package storage

import (
	"bytes"
	"github.com/minio/minio-go/v6"
	"go-open-registry/internal/log"
)

// S3Storage  type struct
type S3Storage struct {
	Path            string `json:"path"`
	Endpoint        string `json:"endpoint"`
	DefaultRegion   string `json:"default_region"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"-"`
	BucketName      string `json:"bucket_name"`
	UseSSL          bool   `json:"use_ssl"`
}

// PutFile implementation
func (s *S3Storage) PutFile(packageName, packageVersion string, content []byte) (err error) {
	log.Info("Put a file to S3 storage")

	// Initialize minio client object.
	minioClient, err := minio.New(s.Endpoint, s.AccessKeyID, s.SecretAccessKey, s.UseSSL)
	if err != nil {
		log.ErrorWithFields("Error init mini client", log.Fields{"err": err})
		return err
	}

	reader := bytes.NewReader(content)
	objectName := packageName + "-" + packageVersion + ".crate"
	// Upload the zip file with FPutObject
	_, err = minioClient.PutObject(s.BucketName, objectName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		log.ErrorWithFields("Error put object", log.Fields{"err": err})
		return err
	}

	return err
}

// GetFile implementation
func (s *S3Storage) GetFile(packageName, packageVersion string) (bytes []byte, err error) {
	log.Info("Get a file from S3 storage")

	minioClient, err := minio.New(s.Endpoint, s.AccessKeyID, s.SecretAccessKey, s.UseSSL)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	objectName := packageName + "-" + packageVersion + ".crate"
	object, err := minioClient.GetObject(s.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	_, err = object.Read(bytes)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return bytes, err
}
