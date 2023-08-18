package pkg

import "testing"

func TestValid2(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "vaild2",
			args: args{filename: "../cmd/openapi-cli/testdata/hello-java-sec-api-docs.json"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid2(tt.args.filename); got != tt.want {
				t.Errorf("Valid2() = %v, want %v", got, tt.want)
			}
		})
	}
}
