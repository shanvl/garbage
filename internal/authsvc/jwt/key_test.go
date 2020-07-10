package jwt

import (
	"testing"
)

func TestPrivateKeyFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing file",
			args: args{
				path: "./keys_test/test.rsa",
			},
			wantErr: false,
		},
		{
			name: "unknown file",
			args: args{
				path: "./no_such_dir/test.rsa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got == nil {
				t.Errorf("PrivateKeyFromFile() error == nil, key == nil")
			}

		})
	}
}

func TestPublicKeyFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing file",
			args: args{
				path: "./keys_test/test.rsa.pub",
			},
			wantErr: false,
		},
		{
			name: "unknown file",
			args: args{
				path: "./no_such_dir/test.rsa",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublicKeyFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got == nil {
				t.Errorf("PublicKeyFromFile() error == nil, key == nil")
			}

		})
	}
}
