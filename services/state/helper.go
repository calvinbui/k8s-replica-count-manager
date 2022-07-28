package state

import "golang.org/x/exp/slices"

// returns the index of the config in the slice.
// returns -1 if not found
func (dp *DeploymentReplicaConfigs) Find(name, namespace string) (*DeploymentReplicaConfig, int) {
	idx := slices.IndexFunc(dp.Configs, func(c DeploymentReplicaConfig) bool {
		return c.Name == name && c.Namespace == namespace
	})

	if idx == -1 {
		return nil, idx
	}

	return &dp.Configs[idx], idx
}
