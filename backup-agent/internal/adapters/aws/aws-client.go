package aws

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	// TODO config
	ydxResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           "https://storage.yandexcloud.net",
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})
)

func NewS3BackupStorage(bucket string, backupObjectKey string) (*S3BackupStorage, error) {
    httpClient := &http.Client{Transport: &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(ydxResolver),
        config.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	return &S3BackupStorage{s3.NewFromConfig(cfg), bucket, backupObjectKey}, nil
}

type S3BackupStorage struct {
	client *s3.Client

	bucket string
	key    string
}

func (s S3BackupStorage) Get(version string) ([]byte, error) {
	o, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:    &s.bucket,
		Key:       &s.key,
		VersionId: &version,
	})
	if err != nil {
		return nil, err
	}
	return io.ReadAll(o.Body)
}

func (s S3BackupStorage) GetLatest() ([]byte, error) {
	o, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
	})
	if err != nil {
		return nil, err
	}
	return io.ReadAll(o.Body)
}

func (s S3BackupStorage) Put(obj []byte) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    &s.key,
        Body: bytes.NewReader(obj),
	})
	return err
}
