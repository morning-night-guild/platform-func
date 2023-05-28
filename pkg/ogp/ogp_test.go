package ogp

import (
	"context"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		url string
	}

	tests := []struct {
		name string
		args args
		want OGP
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				url: "",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Create(tt.args.ctx, tt.args.url)
			t.Logf("got: %+v, err: %+v", got, err)
		})
	}
}
