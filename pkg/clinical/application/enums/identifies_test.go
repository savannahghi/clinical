package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestIDDocumentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    IDDocumentType
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    IDDocumentTypeCCC,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    IDDocumentType("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("IDDocumentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDDocumentType_String(t *testing.T) {
	tests := []struct {
		name string
		c    IDDocumentType
		want string
	}{
		{
			name: "Happy Case",
			c:    IDDocumentTypeCCC,
			want: IDDocumentTypeCCC.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("IDDocumentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDDocumentType_UnmarshalGQL(t *testing.T) {
	validValue := IDDocumentTypeCCC
	invalidType := IDDocumentType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *IDDocumentType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: IDDocumentTypeCCC.String(),
			},
			c:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			c:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(float)",
			args: args{
				v: 45.1,
			},
			c:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("IDDocumentType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIDDocumentType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     IDDocumentType
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     IDDocumentTypeCCC,
			wantW: strconv.Quote("ccc_number"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("IDDocumentType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
