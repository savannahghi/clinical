package clinical

import (
	"reflect"
	"testing"
	"time"
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
