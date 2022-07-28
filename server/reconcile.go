package server

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/calvinbui/teleport-sre-challenge/services/k8s"
	"github.com/calvinbui/teleport-sre-challenge/services/logger"
)

func (s *server) WatchAndReconcileDeployments(ctx context.Context) error {
	logger.Debug("Reconciler: Creating watcher")

	// watch changes to Deployment resources
	watcher, err := k8s.WatchDeployments(ctx, s.services.K8s, s.config.WatchNamespace)

	if err != nil {
		return fmt.Errorf("Error watching Deployments: %w", err)
	}

	logger.Debug("Reconciler: Watching deployment events")
	for event := range watcher.ResultChan() {
		item, ok := event.Object.(*appsv1.Deployment)
		// the message returned may not be a deployment resource
		// this can be due to timeouts or problems with the Kubernetes API
		if !ok {
			return fmt.Errorf("Error casting event to Deployment")
		}

		// react when a deployment is added or modified
		switch event.Type {
		case watch.Added, watch.Modified:
			if config, idx := s.services.DpConfig.Find(item.Name, item.Namespace); idx != -1 {
				// check if replica count matches configured replica count and act accordingly
				if *item.Spec.Replicas != int32(config.Replicas) {
					logger.Info(fmt.Sprintf("%s/%s: %s. Replicas count %v does not match desired %v. Reconciling...", config.Namespace, config.Name, string(event.Type), *item.Spec.Replicas, config.Replicas))
					_, err := k8s.UpdateScale(ctx, s.services.K8s, config.Namespace, config.Name, config.Replicas)
					if err != nil {
						logger.Error(fmt.Sprintf("Error updating %s/%s replicas to %v", config.Name, config.Namespace, config.Replicas), err)
						continue
					}
					logger.Info(fmt.Sprintf("%s/%s: Reconciled replica count to %v", config.Namespace, config.Name, config.Replicas))
				}
			}
		}
	}

	return nil
}
