package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go-open-registry/internal/log"
)

// S3Storage  type struct
type S3Storage struct {
	Path       string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
}

// PutFile implementation
func (s *S3Storage) PutFile(packageName, packageVersion string, content []byte) error {
	log.Info("Put a file to S3 storage")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	},
	)

	if err != nil {
		log.ErrorWithFields("Error setup new s3 session", log.Fields{"error": err})
	}

	s3svc := s3.New(sess)

	uploader := s3manager.NewUploaderWithClient(s3svc)

	reader := bytes.NewReader(content)

	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Path),
		Key:    aws.String(packageName + "-" + packageVersion + ".crate"),
		Body:   reader,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info(output.Location)

	return nil
}

// GetFile implementation
func (s *S3Storage) GetFile(packageName, packageVersion string) ([]byte, error) {
	log.Info("Get a file from S3 storage")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	},
	)
	if err != nil {
		return nil, err
	}

	s3svc := s3.New(sess)

	downloader := s3manager.NewDownloaderWithClient(s3svc)

	params := &s3.GetObjectInput{
		Bucket: aws.String(s.Path),
		Key:    aws.String(packageName + "-" + packageVersion + ".crate"),
	}

	buf := aws.NewWriteAtBuffer([]byte{})

	if _, err := downloader.Download(buf, params); err != nil {
		return nil, err

	}

	return buf.Bytes(), nil
}
