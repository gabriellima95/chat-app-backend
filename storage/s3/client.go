package s3

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

const AttachmentsBucket = "messenger.chat-attachments"

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// used in the examples.
// It contains S3Client, an Amazon S3 service client that is used to perform bucket
// and object actions.
type S3Client struct {
	s3Client  *s3.Client
	presigner *s3.PresignClient
}

func NewS3Client() (S3Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		return S3Client{}, err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	presigner := s3.NewPresignClient(s3Client)
	return S3Client{s3Client: s3Client, presigner: presigner}, nil
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (c S3Client) UploadBase64(filename string, content string, chatID string) (string, error) {
	fileBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		log.Printf("error decoding base64: %s", err.Error())
		return "", fmt.Errorf("error decoding base64: %w", err)
	}

	filepath := chatID + "/" + filename
	_, err = c.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(AttachmentsBucket),
		Key:    aws.String(filepath),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			filename, AttachmentsBucket, filename, err)
	}

	return filepath, nil
}

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
// type Presigner struct {
// 	PresignClient *s3.PresignClient
// }

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (c S3Client) GeneratePresignedURL(
	objectKey string, lifetimeSecs int64) (string, error) {
	request, err := c.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(AttachmentsBucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			AttachmentsBucket, objectKey, err)
	}
	return request.URL, err
}

func (c S3Client) Upload(filename string, content []byte) (string, error) {
	return "", nil
}
