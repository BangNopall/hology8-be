package firebase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	filepath "path/filepath"
	"strings"
	"sync"

	cstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	fstorage "firebase.google.com/go/v4/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"

	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/log"
)

var (
	extFileOpt = []string{"jpg", "jpeg", "png", "pdf"}
	metadata   = []string{"image/jpeg", "image/jpeg", "image/png", "application/pdf"}
)

type FirebaseStorage interface {
	Upload(relativePath string, file *multipart.FileHeader) (string, error)
	Update(relativePath string, file *multipart.FileHeader, oldUrl string) (string, error)
	Delete(relativePath string) error
}

type firebaseClient struct {
	client *fstorage.Client
}

func NewFirebaseStorage() FirebaseStorage {
	configStore := &firebase.Config{
		StorageBucket: env.AppEnv.FireBucket,
	}
	opt := option.WithCredentialsFile(env.AppEnv.FireCredPath)

	app, err := firebase.NewApp(context.Background(), configStore, opt)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][NewFirebaseStorage] failed to initialize firebase")
	}

	client, err := app.Storage(context.Background())

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][NewFirebaseStorage] failed to get firebase storage")
	}

	return &firebaseClient{client}
}

func (f *firebaseClient) Upload(relativePath string, file *multipart.FileHeader) (string, error) {

	ctx := context.Background()

	bucket, err := f.client.DefaultBucket()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][UploadFile] failed to upload file")
		return "", err
	}

	src, err := file.Open()

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][UploadFile] failed to open file")
		return "", err
	}
	defer src.Close()

	uuid := uuid.New().String()
	relativePath = strings.ReplaceAll(relativePath, "+", " ")

	obj := bucket.Object(relativePath)
	wr := obj.NewWriter(ctx)

	wr.ObjectAttrs.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": uuid,
	}

	for idx, ext := range extFileOpt {
		if strings.Contains(filepath.Ext(file.Filename), ext) {
			wr.ObjectAttrs.ContentType = metadata[idx]
		}
	}

	if _, err := io.Copy(wr, src); err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][UploadFile] failed to open file")
		return "", err
	}

	if err := wr.Close(); err != nil {
		return "", err
	}

	url, err := createDownloadUrl(ctx, obj)

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][UploadFile] failed to get download url")
		return "", err
	}

	return url, nil
}

func (f *firebaseClient) Update(relativePath string, file *multipart.FileHeader, oldUrl string) (string, error) {
	var (
		linkChan = make(chan string, 1)
		errChans = make(chan error, 2)
		wg       sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := f.Delete(oldUrl)

		if err != nil {
			errChans <- err
		}
	}()

	go func() {
		defer wg.Done()
		link, err := f.Upload(relativePath, file)

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
			}, "[FIREBASE][UpdateFile] failed to update file")

			return "", err
		}
	}

	link, ok := <-linkChan

	if !ok {
		err := errors.New("link is empty")

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][UpdateFile] failed to update file")

		return "", err
	}

	return link, nil
}

func (f *firebaseClient) Delete(oldUrl string) error {

	parsedUrl, err := url.Parse(oldUrl)

	if err != nil {
		return err
	}

	sanitized := strings.ReplaceAll(parsedUrl.Path, "%2F", "/")
	sanitized = sanitized[strings.Index(sanitized, "o/")+2:]
	sanitized = strings.ReplaceAll(sanitized, "+", " ")

	bucket, err := f.client.DefaultBucket()

	if err != nil {
		return err
	}

	object := bucket.Object(sanitized)

	if err := object.Delete(context.TODO()); err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[FIREBASE][DeleteFile] failed to delete file with file : "+sanitized)
		return err
	}

	return nil
}

func createDownloadUrl(ctx context.Context, obj *cstorage.ObjectHandle) (string, error) {
	attrs, err := obj.Attrs(ctx)

	if err != nil {
		return "", err
	}

	token, ok := attrs.Metadata["firebaseStorageDownloadTokens"]

	if !ok {
		return "", fmt.Errorf("file token not found")
	}

	parsed := url.QueryEscape(obj.ObjectName())

	return fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		obj.BucketName(),
		parsed,
		token,
	), nil
}
