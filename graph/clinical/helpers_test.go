package clinical

import (
	"reflect"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/base"
)

func TestParseDate(t *testing.T) {

	rfcTime, _ := time.Parse(time.RFC3339, "2020-09-24T18:02:38.661033Z")
	monthNumberShortFormTime, _ := time.Parse(StringTimeParseMonthNumberLayout, "2018-01-01")
	monthNameShortFormTime, _ := time.Parse(StringTimeParseMonthNameLayout, "2018-Jan-01")

	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "RFC 3339 Time string",
			args: args{
				date: "2020-09-24T18:02:38.661033Z",
			},
			want: rfcTime,
		},
		{
			name: "Month number short form Time string",
			args: args{
				date: "2018-01-01",
			},
			want: monthNumberShortFormTime,
		},
		{
			name: "Month name short form Time string",
			args: args{
				date: "2018-Jan-01",
			},
			want: monthNameShortFormTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDate(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyOTP(t *testing.T) {
	config, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("unable to load dependencies from YAML: %s", err)
		return
	}

	otpClient, err := base.SetupISCclient(*config, OtpService)
	if err != nil {
		t.Errorf("unable to set up engagement ISC client: %v", err)
		return
	}

	validPhone := "+254723002959"
	validOTP, err := RequestOTP(validPhone, otpClient)
	if err != nil {
		t.Errorf("unable to generate OTP: %v", err)
		return
	}

	type args struct {
		msisdn string
		otp    string
	}
	tests := []struct {
		name         string
		args         args
		wantVerified bool
		wantErr      bool
	}{
		{
			name: "valid case",
			args: args{
				msisdn: validPhone,
				otp:    validOTP,
			},
			wantVerified: true,
			wantErr:      false,
		},
		{
			name: "invalid case",
			args: args{
				msisdn: validPhone,
				otp:    "not a valid OTP",
			},
			wantVerified: false,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verified, _, err := VerifyOTP(tt.args.msisdn, tt.args.otp, otpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if verified != tt.wantVerified {
				t.Errorf("VerifyOTP() got = %v, want %v", verified, tt.wantVerified)
				return
			}
		})
	}
}

func TestRequestOTP(t *testing.T) {
	config, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("unable to load dependencies from YAML: %s", err)
		return
	}

	otpClient, err := base.SetupISCclient(*config, OtpService)
	if err != nil {
		t.Errorf("unable to set up engagement ISC client: %v", err)
		return
	}

	type args struct {
		msisdn string
	}
	tests := []struct {
		name      string
		args      args
		wantBlank bool
		wantErr   bool
	}{
		{
			name: "valid phone number - international format",
			args: args{
				msisdn: "+254723002959",
			},
			wantBlank: false,
			wantErr:   false,
		},
		{
			name: "valid phone number - local format",
			args: args{
				msisdn: "0723002959",
			},
			wantBlank: false,
			wantErr:   false,
		},
		{
			name: "invalid phone number",
			args: args{
				msisdn: "not a real number",
			},
			wantBlank: true,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := RequestOTP(tt.args.msisdn, otpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if otp == "" && !tt.wantBlank {
				t.Errorf("got a blank OTP when expecting non blank")
				return
			}
		})
	}
}
