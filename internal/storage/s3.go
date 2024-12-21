package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	bucketName string
	client     *s3.Client
}

func (c *S3Client) Upload(fileName string, fileData []byte) (string, error) {
	_, err := c.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(fileData),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.bucketName, c.client.Options().Region, fileName)
	return fileURL, nil
}

func (c *S3Client) Download(fileName string) (io.ReadCloser, error) {
	output, err := c.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return output.Body, nil
}
