package state

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

// samples to share between all tests
var testDir = "/tmp/"

// a valid config
var validFile = testDir + "valid.json"
var validFileContents = `{
	"deployments": [
		{
			"name": "foo",
			"namespace": "bar",
			"replicas": 1
		}
	]
}`
var validConfig = DeploymentReplicaConfigs{
	Configs: []DeploymentReplicaConfig{
		{
			Name:      "foo",
			Namespace: "bar",
			Replicas:  1,
		},
	},
}

// an invalid config file
var invalidFile = testDir + "invalid.json"
var invalidFileContents = `{
	"foo": "bar"
}`

// mock a filesystem using https://github.com/spf13/afero
func mockFilesystem() error {
	appFS := afero.NewOsFs()

	// create test files
	if err := appFS.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("Error creating mock directory: %w", err)
	}
	if err := afero.WriteFile(appFS, validFile, []byte(validFileContents), 0644); err != nil {
		return fmt.Errorf("Error creating mock file: %w", err)
	}

	if err := afero.WriteFile(appFS, invalidFile, []byte(invalidFileContents), 0644); err != nil {
		return fmt.Errorf("Error creating mock file: %w", err)
	}

	return nil
}

// tests opening an existing file (or creating a new one if not exists)
func TestNew(t *testing.T) {
	if err := mockFilesystem(); err != nil {
		t.Fatal(err)
	}

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    DeploymentReplicaConfigs
		wantErr bool
	}{
		{
			name:    "Valid file",
			args:    args{filePath: validFile},
			want:    validConfig,
			wantErr: false,
		},
		{
			name:    "No file",
			args:    args{filePath: testDir + "not-exists.json"},
			want:    DeploymentReplicaConfigs{},
			wantErr: false, // a new file is created
		},
		{
			name:    "Invalid file",
			args:    args{filePath: invalidFile},
			want:    DeploymentReplicaConfigs{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

// tests reading a file
func TestRead(t *testing.T) {
	if err := mockFilesystem(); err != nil {
		t.Fatal(err)
	}

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Valid file",
			args:    args{filePath: validFile},
			want:    []byte(validFileContents),
			wantErr: false,
		},
		{
			name:    "No file",
			args:    args{filePath: ""},
			wantErr: true,
		},
		{
			name:    "Invalid file",
			args:    args{filePath: invalidFile},
			want:    []byte(invalidFileContents),
			wantErr: false, // read doesn't parse the file, so there would be no errors
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

// tests updating a file
func TestPut(t *testing.T) {
	if err := mockFilesystem(); err != nil {
		t.Fatal(err)
	}

	type args struct {
		data     DeploymentReplicaConfigs
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid file",
			args: args{
				filePath: validFile,
				data:     DeploymentReplicaConfigs{},
			},
			wantErr: false,
		},
		{
			name: "No file",
			args: args{
				filePath: "",
				data:     DeploymentReplicaConfigs{},
			},
			wantErr: true,
		},
		{
			name: "Invalid file",
			args: args{
				filePath: invalidFile,
				data:     DeploymentReplicaConfigs{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Put(tt.args.data, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
