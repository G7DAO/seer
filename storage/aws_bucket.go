package storage

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ethereum/go-ethereum/log"
)

// S3 implements the Storer interface for Amazon S3

type S3 struct {
	// Add S3 specific fields
}

func NewS3Storage() *S3 {
	// Implement the NewS3 function
	log.Warn("AWS bucket support not implemented yet")
	return &S3{}
}

func (s *S3) Save(key string, data []string) error {
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

func (s *S3) ReadBatch(readItems []ReadItem) (map[string][]string, error) {
	// Implement the ReadBatch method
	return nil, nil
}
