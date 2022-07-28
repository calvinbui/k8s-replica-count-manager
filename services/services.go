package services

import (
	"github.com/calvinbui/teleport-sre-challenge/config"
	"github.com/calvinbui/teleport-sre-challenge/services/k8s"
	"github.com/calvinbui/teleport-sre-challenge/services/state"

	"k8s.io/client-go/kubernetes"
)

type Services struct {
	K8s      *kubernetes.Clientset
	DpConfig state.DeploymentReplicaConfigs
}

func New(conf *config.Config) (Services, error) {
	k8s, err := k8s.New()
	if err != nil {
		return Services{}, err
	}

	dpConfig, err := state.New(conf.FilePath)
	if err != nil {
		return Services{}, err
	}

	svc := Services{
		K8s:      k8s,
		DpConfig: dpConfig,
	}

	return svc, nil
}
