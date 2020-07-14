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

func TestKeysFromFiles(t *testing.T) {
	type args struct {
		prPath  string
		pubPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing file",
			args: args{
				prPath:  "./keys_test/test.rsa",
				pubPath: "./keys_test/test.rsa.pub",
			},
			wantErr: false,
		},
		{
			name: "unknown private file",
			args: args{
				prPath:  "./no_such_dir/test.rsa",
				pubPath: "./keys_test/test.rsa.pub",
			},
			wantErr: true,
		},
		{
			name: "unknown public file",
			args: args{
				prPath:  "./keys_test/test.rsa",
				pubPath: "./no_such_dir/test.rsa.pub",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prKey, pubKey, err := KeysFromFiles(tt.args.prPath, tt.args.pubPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublicKeyFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && (prKey == nil || pubKey == nil) {
				t.Errorf("PublicKeyFromFile() error == nil, prKey == %v, pubKey == %v", prKey, pubKey)
			}

		})
	}
}
