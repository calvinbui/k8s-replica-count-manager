package k8s

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// There is an issue with testing fake subresources that create requests of different object types than their parent resource, causing a panic when testing
// I've written the tests but commented them out for now.
//
// More:
// - https://github.com/kubernetes/kubernetes/issues/80354
// - https://github.com/kubernetes/kubernetes/issues/71598
// - https://github.com/kubernetes/kubernetes/pull/71599
// - https://github.com/kubernetes/client-go/issues/734
// - https://github.com/kubernetes/kubernetes/pull/97937
// - https://github.com/kubernetes/client-go/issues/1077

// samples to share between all tests
var clientSet = fake.NewSimpleClientset()
var name = "foo"
var namespace = "bar"
var currentReplicas = int32(10)
var labels = map[string]string{"foo": "bar"}
var dp = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: &currentReplicas,
		Selector: &metav1.LabelSelector{MatchLabels: labels},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  name,
						Image: name,
					},
				},
			},
		},
	},
}

// names must be unique
func generateFakeDeployment(name string) error {
	dp.ObjectMeta.Name = dp.ObjectMeta.Name + "-" + name

	if _, err := clientSet.AppsV1().Deployments(namespace).Create(context.TODO(), dp, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("Error generating fake deployment: %s", err)
	}

	return nil
}

func TestGetScale(t *testing.T) {
	if err := generateFakeDeployment("get-scale"); err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx       context.Context
		client    kubernetes.Interface
		namespace string
		name      string
	}
	tests := []struct {
		name    string
		args    args
		want    *autoscalingv1.Scale
		wantErr bool
	}{
		// {
		// 	name: "Get Deployment",
		// 	args: args{
		// 		ctx:       context.TODO(),
		// 		client:    clientSet,
		// 		namespace: namespace,
		// 		name:      name,
		// 	},
		// 	want: &autoscalingv1.Scale{
		// 		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		// 		Spec:       autoscalingv1.ScaleSpec{Replicas: currentReplicas},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "Deployment does not exist",
		// 	args: args{
		// 		ctx:       context.TODO(),
		// 		client:    clientSet,
		// 		namespace: "",
		// 		name:      "",
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetScale(tt.args.ctx, tt.args.client, tt.args.namespace, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetScale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetScale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateScale(t *testing.T) {
	if err := generateFakeDeployment("update-scale"); err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx       context.Context
		client    kubernetes.Interface
		namespace string
		name      string
		replicas  int
	}
	tests := []struct {
		name    string
		args    args
		want    *autoscalingv1.Scale
		wantErr bool
	}{
		// {
		// 	name: "Scale up deployment",
		// 	args: args{
		// 		ctx:       context.TODO(),
		// 		client:    clientSet,
		// 		namespace: namespace,
		// 		name:      name,
		// 		replicas:  10,
		// 	},
		// 	want: &autoscalingv1.Scale{
		// 		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		// 		Spec:       autoscalingv1.ScaleSpec{Replicas: 10},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "Scale up deployment",
		// 	args: args{
		// 		ctx:       context.TODO(),
		// 		client:    clientSet,
		// 		namespace: namespace,
		// 		name:      name,
		// 		replicas:  0,
		// 	},
		// 	want: &autoscalingv1.Scale{
		// 		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: labels},
		// 		Spec:       autoscalingv1.ScaleSpec{Replicas: 0},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "Deployment does not exist",
		// 	args: args{
		// 		ctx:       context.TODO(),
		// 		client:    clientSet,
		// 		namespace: namespace,
		// 		name:      name,
		// 		replicas:  5,
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateScale(tt.args.ctx, tt.args.client, tt.args.namespace, tt.args.name, tt.args.replicas)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateScale() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateScale() = %v, want %v", got, tt.want)
			}
		})
	}
}
