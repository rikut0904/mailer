package aws

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rikut0904/mailer-backend/internal/domain/entity"
	"github.com/rikut0904/mailer-backend/internal/domain/repository"
)

type s3Client struct {
	client     *s3.Client
	bucketName string
}

func NewS3ClientFromDomain(domain *entity.S3Domain) (repository.MailStorageRepository, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(domain.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			domain.AccessKeyID,
			domain.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	return &s3Client{
		client:     client,
		bucketName: domain.Bucket,
	}, nil
}

func (s *s3Client) ListKeys(prefix string, continuationToken *string, maxKeys int) ([]string, *string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucketName),
		MaxKeys: aws.Int32(int32(maxKeys)),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}
	if continuationToken != nil {
		input.ContinuationToken = continuationToken
	}

	result, err := s.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, nil, err
	}

	keys := make([]string, 0, len(result.Contents))
	for _, obj := range result.Contents {
		keys = append(keys, *obj.Key)
	}

	var nextToken *string
	if result.IsTruncated != nil && *result.IsTruncated {
		nextToken = result.NextContinuationToken
	}

	return keys, nextToken, nil
}

func (s *s3Client) GetObject(key string) ([]byte, error) {
	result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *s3Client) DeleteObject(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}
