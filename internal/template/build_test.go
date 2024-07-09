package template

import (
	"os"
	"testing"
)

func TestBuildBicepTemplate(t *testing.T) {
	type args struct {
		bicepFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid_file",
			args: args{
				bicepFile: "./testdata/basic.bicep",
			},
			wantErr: false,
		},
		{
			name: "invalid_file",
			args: args{
				bicepFile: "./testdata/invalid.bicep",
			},
			wantErr: true,
		},
		{
			name: "non_existent_file",
			args: args{
				bicepFile: "./testdata/non_existent.bicep",
			},
			wantErr: true,
		},
		{
			name: "invalid_file_extension",
			args: args{
				bicepFile: "./testdata/main.md",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := BuildBicepTemplate(tt.args.bicepFile)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BuildBicepTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if err := os.Remove(path); err == nil {
					t.Fatalf("BuildBicepTemplate() did create file %s", path)
				}
			} else {
				if err := os.Remove(path); err != nil {
					t.Fatalf("BuildBicepTemplate() did not create file %s", path)
				}
			}
		})
	}
}

func Test_commandExists(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ls",
			args: args{
				cmd: "ls",
			},
			want: true,
		},
		{
			name: "non_existent_command",
			args: args{
				cmd: "some-non-existent-command",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandExists(tt.args.cmd); got != tt.want {
				t.Fatalf("commandExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
