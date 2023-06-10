package slack

import "testing"

func TestExtractFirstURLFromUrlsConcatByPipe(t *testing.T) {
	t.Parallel()

	type args struct {
		urls string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "|でURLが区切られている",
			args: args{
				urls: "https://example.com|https://example.com",
			},
			want: "https://example.com",
		},
		{
			name: "|でURLが区切られていない",
			args: args{
				urls: "https://example.com",
			},
			want: "https://example.com",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ExtractFirstURLFromUrlsConcatByPipe(tt.args.urls); got != tt.want {
				t.Errorf("ExtractFirstURLFromUrlsConcatByPipe() = %v, want %v", got, tt.want)
			}
		})
	}
}
