package k8s

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// test creating a watcher
func Test_watcher_Watch(t *testing.T) {
	type fields struct {
		client    kubernetes.Interface
		namespace string
	}
	type args struct {
		options metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create watcher",
			fields: fields{
				client:    clientSet,
				namespace: "",
			},
			args:    args{metav1.ListOptions{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &watcher{
				client:    tt.fields.client,
				namespace: tt.fields.namespace,
			}
			_, err := s.Watch(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("watcher.Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
