package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type watcher struct {
	client    kubernetes.Interface
	namespace string
}

// function for the watcher client to watch deployments
func (s *watcher) Watch(options metav1.ListOptions) (watch.Interface, error) {
	return s.client.AppsV1().Deployments(s.namespace).Watch(context.Background(), metav1.ListOptions{})
}
