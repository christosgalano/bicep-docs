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
			name: "valid bicep file",
			args: args{
				bicepFile: "./testdata/main.bicep",
			},
			wantErr: false,
		},
		{
			name: "invalid bicep file",
			args: args{
				bicepFile: "./testdata/invalid.bicep",
			},
			wantErr: true,
		},
		{
			name: "non existent bicep file",
			args: args{
				bicepFile: "./testdata/non_existent.bicep",
			},
			wantErr: true,
		},
		{
			name: "invalid file extension",
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
