package state

type DeploymentReplicaConfigs struct {
	Configs []DeploymentReplicaConfig `json:"deployments"`
}

type DeploymentReplicaConfig struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int    `json:"replicas"`
}
