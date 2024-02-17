package dto

import (
	"bytes"
	"strconv"
	"testing"
)

func TestConsentState_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     SegmentationCategory
		wantW string
	}{
		{
			name:  "valid type s",
			c:     SegmentationCategoryLowRisk,
			wantW: strconv.Quote("CERVICAL_CANCER_LOW_RISK"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SegmentationCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSegmentationCategory_String(t *testing.T) {
	tests := []struct {
		name string
		e    SegmentationCategory
		want string
	}{
		{
			name: "CERVICAL_CANCER_TIPS",
			e:    SegmentationCategoryNoRisk,
			want: "CERVICAL_CANCER_TIPS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("SegmentationCategory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegmentationCategory_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    SegmentationCategory
		want bool
	}{
		{
			name: "valid type",
			e:    SegmentationCategoryHighRiskNegative,
			want: true,
		},
		{
			name: "invalid type",
			e:    SegmentationCategory("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("SegmentationCategory.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegmentationCategory_UnmarshalGQL(t *testing.T) {
	value := SegmentationCategoryHighRiskNegative
	invalid := SegmentationCategory("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *SegmentationCategory
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "CERVICAL_CANCER_HIGH_RISK",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SegmentationCategory.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSegmentationCategory_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     SegmentationCategory
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     SegmentationCategoryHighRiskNegative,
			b:     w,
			wantW: strconv.Quote("CERVICAL_CANCER_HIGH_RISK"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SegmentationCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
