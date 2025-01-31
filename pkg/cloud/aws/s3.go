package aws

import (
	"context"
	"io"
	"stori/pkg/cloud"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var _ cloud.BucketClient = S3{}

type S3 struct {
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3(session *session.Session, timeout time.Duration) S3 {
	s3manager.NewUploader(session)
	return S3{
		timeout:    timeout,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

func (s S3) Create(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return err
	}

	if err := s.client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return err
	}

	return nil
}

func (s S3) UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (s S3) DownloadObject(ctx context.Context, bucket, fileName string, body io.WriterAt) error {
	if _, err := s.downloader.DownloadWithContext(ctx, body, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return err
	}

	return nil
}

func (s S3) DeleteObject(ctx context.Context, bucket, fileName string) error {
	if _, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return err
	}

	if err := s.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return err
	}

	return nil
}

func (s S3) FetchObject(ctx context.Context, bucket, fileName string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (s S3) Delete(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.client.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return err
	}

	return nil
}
