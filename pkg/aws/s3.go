package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BangNopall/hology8-be/internal/infra/env"
	"github.com/BangNopall/hology8-be/pkg/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	allowedExts  = []string{"jpg", "jpeg", "png", "pdf"}
	allowedMimes = []string{"image/jpeg", "image/jpeg", "image/png", "application/pdf"}
)

type CloudStorage interface {
	Upload(relativePath string, file *multipart.FileHeader) (string, error)
	Update(relativePath string, file *multipart.FileHeader, oldUrl string) (string, error)
	Delete(oldUrl string) error
}

type s3Client struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Storage() CloudStorage {
	if env.AppEnv.AwsBucket == "" || env.AppEnv.AwsRegion == "" {
		log.Fatal(log.LogInfo{}, "[S3][NewS3Storage] missing AWS_BUCKET or AWS_REGION in env")
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.AppEnv.AwsRegion),
		config.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
		config.WithResponseChecksumValidation(aws.ResponseChecksumValidationWhenRequired),
	)
	if err != nil {
		log.Fatal(log.LogInfo{
			"error":  err.Error(),
			"region": env.AppEnv.AwsRegion,
		}, "[S3][NewS3Storage] failed to load AWS config")
	}

	return &s3Client{
		client: s3.NewFromConfig(cfg),
		bucket: env.AppEnv.AwsBucket,
		region: env.AppEnv.AwsRegion,
	}
}

func (s *s3Client) Upload(relativePath string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	src, err := file.Open()
	if err != nil {
		log.Error(log.LogInfo{
			"error":    err.Error(),
			"filename": file.Filename,
		}, "[S3][UploadFile] failed to open multipart file")
		return "", err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		log.Error(log.LogInfo{
			"error":    err.Error(),
			"filename": file.Filename,
		}, "[S3][UploadFile] failed to copy file to buffer")
		return "", err
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(file.Filename)), ".")
	contentType := ""
	for i, e := range allowedExts {
		if e == ext {
			contentType = allowedMimes[i]
			break
		}
	}
	if contentType == "" {
		err := errors.New("unsupported file type")
		log.Error(log.LogInfo{
			"filename": file.Filename,
			"ext":      ext,
			"error":    err.Error(),
		}, "[S3][UploadFile] unsupported file extension")
		return "", err
	}

	key := strings.ReplaceAll(relativePath, "+", " ")

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
		Metadata: map[string]string{
			"uuid": uuid.New().String(),
		},
	})
	if err != nil {
		log.Error(log.LogInfo{
			"error":  err.Error(),
			"key":    key,
			"bucket": s.bucket,
		}, "[S3][UploadFile] failed to upload file to S3")
		return "", err
	}

	url := buildPublicObjectURL(s.bucket, s.region, key)

	log.Info(log.LogInfo{
		"url":      url,
		"key":      key,
		"filename": file.Filename,
	}, "[S3][UploadFile] upload successful")

	return url, nil
}

func (s *s3Client) Delete(oldUrl string) error {
	ctx := context.Background()

	key, err := extractKeyFromURL(oldUrl)
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
			"url":   oldUrl,
		}, "[S3][DeleteFile] failed to extract key from URL")
		return err
	}

	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error(log.LogInfo{
			"error":  err.Error(),
			"bucket": s.bucket,
			"key":    key,
		}, "[S3][DeleteFile] failed to delete object from S3")
		return err
	}

	log.Info(log.LogInfo{
		"key":    key,
		"bucket": s.bucket,
	}, "[S3][DeleteFile] delete successful")

	return nil
}

func (s *s3Client) Update(relativePath string, file *multipart.FileHeader, oldUrl string) (string, error) {
	var (
		linkChan = make(chan string, 1)
		errChans = make(chan error, 2)
		wg       sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.Delete(oldUrl); err != nil {
			errChans <- err
		}
	}()

	go func() {
		defer wg.Done()
		link, err := s.Upload(relativePath, file)
		if err != nil {
			errChans <- err
		} else {
			linkChan <- link
		}
	}()

	go func() {
		wg.Wait()
		close(linkChan)
		close(errChans)
	}()

	for err := range errChans {
		if err != nil {
			log.Error(log.LogInfo{
				"error": err.Error(),
				"path":  relativePath,
			}, "[S3][UpdateFile] failed to update file")
			return "", err
		}
	}

	link, ok := <-linkChan
	if !ok {
		err := errors.New("link is empty")
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[S3][UpdateFile] empty result channel")
		return "", err
	}
	return link, nil
}

func extractKeyFromURL(fileURL string) (string, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(u.Path, "/"), nil
}

func buildPublicObjectURL(bucket, region, key string) string {
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3.%s.amazonaws.com", bucket, region),
		Path:   "/" + key,
	}
	return u.String()
}
