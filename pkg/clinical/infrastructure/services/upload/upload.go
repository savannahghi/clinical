package upload

import (
	"context"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/serverutils"
	"google.golang.org/api/option"
)

// ServiceUpload holds the upload service methods
type ServiceUpload interface {
	UploadMedia(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error)
}

// ServiceUploadImpl represents upload service implementations
type ServiceUploadImpl struct {
	Client storage.Client
}

// NewServiceUpload returns new instance of upload service
func NewServiceUpload(ctx context.Context) *ServiceUploadImpl {
	credentials := serverutils.MustGetEnvVar("GOOGLE_APPLICATION_CREDENTIALS")

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentials))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	return &ServiceUploadImpl{
		Client: *client,
	}
}

// UploadMedia uploads media to GCS
func (u *ServiceUploadImpl) UploadMedia(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
	bucketName := serverutils.MustGetEnvVar("CLINICAL_BUCKET_NAME")

	// Set the content type based on the request header
	object := u.Client.Bucket(bucketName).Object(name)

	// Set up the resumable upload
	wc := object.NewWriter(ctx)
	wc.ContentType = contentType
	wc.ChunkSize = 256 * 1024 // 256 KB chunk size

	// Write the file to Google Cloud Storage
	if _, err := io.Copy(wc, file); err != nil {
		wc.Close()
		return nil, err
	}

	// Close the writer to flush the data to Google Cloud Storage
	timeoutCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	if err := wc.Close(); err != nil {
		return nil, err
	}

	// Generate a signed URL for the uploaded file
	url, err := object.Attrs(timeoutCtx)
	if err != nil {
		return nil, err
	}

	output := &dto.Media{
		URL:         url.MediaLink,
		Name:        url.Name,
		ContentType: url.ContentType,
	}

	return output, nil
}
