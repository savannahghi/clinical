package clinical

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
)

func TestPhotosToAttachments(t *testing.T) {
	srv := NewService()
	bs, err := ioutil.ReadFile("testdata/photo.jpg")
	if err != nil {
		t.Fatalf("unable to read test photo %s: ", err)
	}
	var photoBase64 = base64.StdEncoding.EncodeToString(bs)

	type args struct {
		ctx        context.Context
		photos     []*PhotoInput
		engagement *base.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: context.Background(),
				photos: []*PhotoInput{
					{
						PhotoContentType: base.ContentTypePng,
						PhotoFilename:    ksuid.New().String(),
						PhotoBase64data:  photoBase64,
					},
				},
				engagement: srv.engagement,
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PhotosToAttachments(tt.args.ctx, tt.args.photos, tt.args.engagement)
			if (err != nil) != tt.wantErr {
				t.Errorf("PhotosToAttachments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantNil {
				t.Errorf("got nil attachments, expected non nil")
				return
			}
		})
	}
}
