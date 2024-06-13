package storage

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 implements the Storer interface for Amazon S3 Bucket
type S3 struct {
	BasePath string
}

func NewS3Storage(basePath string) *S3 {
	return &S3{BasePath: basePath}
}

func (s *S3) Save(batchDir, filename string, data []string) error {
	key := filepath.Join(s.BasePath, batchDir, filename)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	svc := s3.New(sess)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(data[0])),

		// Add the metadata
		Metadata: map[string]*string{
			"key": aws.String(key),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3) Read(key string) ([]string, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	svc := s3.New(sess)

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)

	return []string{buf.String()}, nil
}

func (s *S3) Delete(key string) error {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))

	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	return nil

}

func (s *S3) List(ctx context.Context, delim string, timeout int, returnFunc ListReturnFunc) ([]string, error) {
	return []string{}, nil
}

func (s *S3) ReadBatch(readItems []ReadItem) (map[string][]string, error) {
	// Implement the ReadBatch method
	return nil, nil
}
