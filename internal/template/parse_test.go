/*
Package template provides functions to build and parse Bicep and the corresponding ARM templates.
*/
package template

import (
	"reflect"
	"testing"

	"github.com/christosgalano/bicep-docs/internal/types"
)

func TestParseTemplates(t *testing.T) {
	type args struct {
		bicepFile string
		armFile   string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Template
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTemplates(tt.args.bicepFile, tt.args.armFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTemplates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBicepTemplate(t *testing.T) {
	type args struct {
		bicepFile string
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Module
		want1   []types.Resource
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseBicepTemplate(tt.args.bicepFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBicepTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBicepTemplate() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseBicepTemplate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
