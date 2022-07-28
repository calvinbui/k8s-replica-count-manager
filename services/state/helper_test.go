package state

import (
	"reflect"
	"testing"
)

// tests finding a config
func TestDeploymentReplicaConfigs_Find(t *testing.T) {
	testName := "foo"
	testNamespace := "bar"
	testConfig := DeploymentReplicaConfig{
		Name:      testName,
		Namespace: testNamespace,
		Replicas:  1,
	}

	type fields struct {
		Configs []DeploymentReplicaConfig
	}
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *DeploymentReplicaConfig
		want1  int
	}{
		{
			name:   "Found",
			fields: fields{Configs: []DeploymentReplicaConfig{testConfig}},
			args:   args{name: testName, namespace: testNamespace},
			want:   &testConfig,
			want1:  0,
		},
		{
			name:   "Not found",
			fields: fields{Configs: []DeploymentReplicaConfig{testConfig}},
			args:   args{name: testName, namespace: "baz"},
			want:   nil,
			want1:  -1,
		},
		{
			name:   "Empty config",
			fields: fields{Configs: []DeploymentReplicaConfig{}},
			args:   args{name: testName, namespace: testNamespace},
			want:   nil,
			want1:  -1,
		},
		{
			name:   "Empty fields",
			fields: fields{Configs: []DeploymentReplicaConfig{testConfig}},
			args:   args{name: "", namespace: ""},
			want:   nil,
			want1:  -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := &DeploymentReplicaConfigs{
				Configs: tt.fields.Configs,
			}
			got, got1 := dp.Find(tt.args.name, tt.args.namespace)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeploymentReplicaConfigs.Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeploymentReplicaConfigs.Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
