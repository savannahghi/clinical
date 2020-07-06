package clinical

import "testing"

func Test_trimString(t *testing.T) {
	type args struct {
		inp       string
		maxLength int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "short string",
			args: args{
				inp:       "a short string",
				maxLength: 20,
			},
			want: "a short string",
		},
		{
			name: "a long string",
			args: args{
				inp:       "this string is longer than the indicated max length",
				maxLength: 20,
			},
			want: "this string is lo...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimString(tt.args.inp, tt.args.maxLength); got != tt.want {
				t.Errorf("trimString() = %v, want %v", got, tt.want)
			}
		})
	}
}
