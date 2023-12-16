package service

import "testing"

func Test_formatURL(t *testing.T) {
	type args struct {
		targetURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "formatURL",
			args: args{
				targetURL: "https://example.com",
			},
			want:    "example.com:443",
			wantErr: false,
		},
		{
			name: "formatURL",
			args: args{
				targetURL: "example.com",
			},
			want:    "example.com:443",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatURL(tt.args.targetURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
