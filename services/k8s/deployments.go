package k8s

import (
	"context"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	toolsWatch "k8s.io/client-go/tools/watch"
)

// WatchDeployments returns a retry watcher that watches for changes to Deployments
func WatchDeployments(ctx context.Context, client kubernetes.Interface, namespace string) (*toolsWatch.RetryWatcher, error) {
	watcher := &watcher{client, namespace}

	// RetryWatcher will make sure that in case the underlying watcher is closed (e.g. due to API timeout or etcd timeout)
	// it will get restarted from the last point without the consumer even knowing about it.
	retryWatcher, err := toolsWatch.NewRetryWatcher("1", watcher)

	if err != nil {
		return nil, err
	}

	return retryWatcher, nil
}

// GetScale returns a deployment's current replicas
func GetScale(ctx context.Context, client kubernetes.Interface, namespace, name string) (*autoscalingv1.Scale, error) {
	res, err := client.AppsV1().Deployments(namespace).GetScale(ctx, name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateScale updates a deployment's current replicas
func UpdateScale(ctx context.Context, client kubernetes.Interface, namespace, name string, replicas int) (*autoscalingv1.Scale, error) {
	scale, err := GetScale(ctx, client, namespace, name)
	if err != nil {
		return nil, err
	}

	scale.Spec.Replicas = int32(replicas)

	res, err := client.AppsV1().Deployments(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})

	if err != nil {
		return nil, err
	}

	return res, nil
}
