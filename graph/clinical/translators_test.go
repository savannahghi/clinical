package clinical

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
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

func TestRelationshipTypeDisplay(t *testing.T) {
	type args struct {
		val RelationshipType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "unknown - random stuff",
			args: args{
				val: RelationshipType("some random stuff"),
			},
			want: "Unknown",
		},
		{
			name: "unknown - real unknown enum",
			args: args{
				val: RelationshipTypeU,
			},
			want: "Unknown",
		},
		{
			name: "state agency",
			args: args{
				val: RelationshipTypeS,
			},
			want: "State Agency",
		},
		{
			name: "other",
			args: args{
				val: RelationshipTypeO,
			},
			want: "Other",
		},
		{
			name: "next of kin",
			args: args{
				val: RelationshipTypeN,
			},
			want: "Next-of-Kin",
		},
		{
			name: "insurance company",
			args: args{
				val: RelationshipTypeI,
			},
			want: "Insurance Company",
		},
		{
			name: "federal agency",
			args: args{
				val: RelationshipTypeF,
			},
			want: "Federal Agency",
		},
		{
			name: "employer",
			args: args{
				val: RelationshipTypeE,
			},
			want: "Employer",
		},
		{
			name: "emergency contact",
			args: args{
				val: RelationshipTypeC,
			},
			want: "Emergency Contact",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RelationshipTypeDisplay(tt.args.val); got != tt.want {
				t.Errorf("RelationshipTypeDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaritalStatusDisplay(t *testing.T) {
	type args struct {
		val MaritalStatus
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "unknown - random stuff",
			args: args{
				val: MaritalStatus("some random stuff"),
			},
			want: "unknown",
		},
		{
			name: "unknown - real unknown enum",
			args: args{
				val: MaritalStatusUnk,
			},
			want: "unknown",
		},
		{
			name: "annulled",
			args: args{
				val: MaritalStatusA,
			},
			want: "Annulled",
		},
		{
			name: "divorced",
			args: args{
				val: MaritalStatusD,
			},
			want: "Divorced",
		},
		{
			name: "interlocutory",
			args: args{
				val: MaritalStatusI,
			},
			want: "Interlocutory",
		},
		{
			name: "Legally Separated",
			args: args{
				val: MaritalStatusL,
			},
			want: "Legally Separated",
		},
		{
			name: "Married",
			args: args{
				val: MaritalStatusM,
			},
			want: "Married",
		},
		{
			name: "Polygamous",
			args: args{
				val: MaritalStatusP,
			},
			want: "Polygamous",
		},
		{
			name: "Never Married",
			args: args{
				val: MaritalStatusS,
			},
			want: "Never Married",
		},
		{
			name: "Domestic Partner",
			args: args{
				val: MaritalStatusT,
			},
			want: "Domestic Partner",
		},
		{
			name: "Unmarried",
			args: args{
				val: MaritalStatusU,
			},
			want: "unmarried",
		},
		{
			name: "Widowed",
			args: args{
				val: MaritalStatusW,
			},
			want: "Widowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaritalStatusDisplay(tt.args.val); got != tt.want {
				t.Errorf("MaritalStatusDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhysicalPostalAddressesToCombinedFHIRAddress(t *testing.T) {
	type args struct {
		physical []*PhysicalAddress
		postal   []*PostalAddress
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "nil inputs",
			args: args{
				physical: nil,
				postal:   nil,
			},
			wantNil: true,
		},
		{
			name: "non nil physical",
			args: args{
				physical: []*PhysicalAddress{
					{
						MapsCode:        gofakeit.Address().Zip,
						PhysicalAddress: gofakeit.Address().Address,
					},
				},
				postal: nil,
			},
			wantNil: false,
		},
		{
			name: "non nil postal",
			args: args{
				postal: []*PostalAddress{
					{
						PostalAddress: gofakeit.Address().Address,
						PostalCode:    gofakeit.Address().City,
					},
				},
				physical: nil,
			},
			wantNil: false,
		},
		{
			name: "both non nil",
			args: args{
				postal: []*PostalAddress{
					{
						PostalAddress: gofakeit.Address().Address,
						PostalCode:    gofakeit.Address().City,
					},
				},
				physical: []*PhysicalAddress{
					{
						MapsCode:        gofakeit.Address().Zip,
						PhysicalAddress: gofakeit.Address().Address,
					},
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PhysicalPostalAddressesToCombinedFHIRAddress(tt.args.physical, tt.args.postal)
			if !tt.wantNil && got == nil {
				t.Errorf("unexpected nil result")
				return
			}
		})
	}
}
