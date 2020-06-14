package env

import (
	"os"
	"testing"
	"time"
)

const envName = "env_name"

func TestInt(t *testing.T) {
	type args struct {
		env      string
		fallback int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "no env",
			args: args{
				env:      "",
				fallback: 25,
			},
			want: 25,
		},
		{
			name: "invalid int",
			args: args{
				env:      "invalid",
				fallback: 25,
			},
			want: 25,
		},
		{
			name: "valid int",
			args: args{
				env:      "33",
				fallback: 25,
			},
			want: 33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(envName, tt.args.env)
			if err != nil {
				t.Error("wasn't able to set an env variable")
			}
			if got := Int(envName, tt.args.fallback); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		env      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no env",
			args: args{
				env:      "",
				fallback: "fallback",
			},
			want: "fallback",
		},
		{
			name: "valid string",
			args: args{
				env:      "some string",
				fallback: "fallback",
			},
			want: "some string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(envName, tt.args.env)
			if err != nil {
				t.Error("wasn't able to set an env variable")
			}
			if got := String(envName, tt.args.fallback); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_duration(t *testing.T) {
	type args struct {
		env      string
		fallback time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "no env",
			args: args{
				env:      "",
				fallback: time.Duration(25),
			},
			want: time.Duration(25),
		},
		{
			name: "invalid duration",
			args: args{
				env:      "52da2",
				fallback: time.Duration(25),
			},
			want: time.Duration(25),
		},
		{
			name: "valid duration",
			args: args{
				env:      "522s",
				fallback: time.Duration(25),
			},
			want: 522 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(envName, tt.args.env)
			if err != nil {
				t.Error("wasn't able to set an env variable")
			}
			if got := Duration(envName, tt.args.fallback); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	type args struct {
		env      string
		fallback bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no env",
			args: args{
				env:      "",
				fallback: false,
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				env:      "true",
				fallback: false,
			},
			want: true,
		},
		{
			name: "false",
			args: args{
				env:      "false",
				fallback: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(envName, tt.args.env)
			if err != nil {
				t.Error("wasn't able to set an env variable")
			}
			if got := Bool(envName, tt.args.fallback); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}
