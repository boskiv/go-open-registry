package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

	h := sha256.New()
	_, err = h.Write(content)
	if err != nil {
		log.ErrorWithFields("Error while write hash from file", log.Fields{"error": err})
		return err
	}

	cksum := hex.EncodeToString(h.Sum(nil))
	log.InfoWithFields("Cksum from s3 put", log.Fields{"cksum": cksum})

	reader := bytes.NewReader(content)

	objectName := packageName + "-" + packageVersion + ".crate"
	// Upload the zip file with FPutObject
	response, err := minioClient.PutObject(s.BucketName, objectName, reader, reader.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.ErrorWithFields("Error put object", log.Fields{"err": err})
		return err
	}
	log.InfoWithFields("Response", log.Fields{"response": response})

	return err
}

// GetFile implementation
func (s *S3Storage) GetFile(packageName, packageVersion string) (content []byte, err error) {
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

	stat, err := object.Stat()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.InfoWithFields("S3 Stats:", log.Fields{"stat": stat})

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(object)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	content = buf.Bytes()

	h := sha256.New()
	_, err = h.Write(content)
	if err != nil {
		log.ErrorWithFields("Error while write hash from file", log.Fields{"error": err})
		return nil, err
	}

	cksum := hex.EncodeToString(h.Sum(nil))
	log.InfoWithFields("Cksum from s3 get", log.Fields{"cksum": cksum})

	return content, err
}
