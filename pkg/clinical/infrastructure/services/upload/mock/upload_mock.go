package mock

import (
	"context"
	"io"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeUpload is a mock of the fake upload
type FakeUpload struct {
	MockUploadMediaFn func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.MediaOutPut, error)
}

// NewFakeUploadMock initializes a new instance of upload mock
func NewFakeUploadMock() *FakeUpload {
	return &FakeUpload{
		MockUploadMediaFn: func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.MediaOutPut, error) {
			return &dto.MediaOutPut{
				URL: "https://google.com",
			}, nil
		},
	}
}

// UploadMedia is a mock implementation of uploading media to GCS
func (u *FakeUpload) UploadMedia(ctx context.Context, name string, file io.Reader, contentType string) (*dto.MediaOutPut, error) {
	return u.MockUploadMediaFn(ctx, name, file, contentType)
}
