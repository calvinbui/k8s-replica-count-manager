package state

import (
	"encoding/json"
	"fmt"
	"os"

	"sync"

	"github.com/calvinbui/teleport-sre-challenge/services/logger"
)

var fileMutex sync.Mutex

// New creates a new state config
func New(filePath string) (DeploymentReplicaConfigs, error) {
	file, err := Read(filePath)
	if err != nil {
		// creates a new state file if it does not exist
		if os.IsNotExist(err) {
			logger.Info(fmt.Sprintf("%s does not exist. Creating empty config.", filePath))
			return DeploymentReplicaConfigs{}, nil
		}
		return DeploymentReplicaConfigs{}, fmt.Errorf("Error reading %s: %w", filePath, err)
	}

	// load the file
	var data DeploymentReplicaConfigs
	if err = json.Unmarshal(file, &data); err != nil {
		return DeploymentReplicaConfigs{}, fmt.Errorf("Error unmarshaling JSON from %s: %w", filePath, err)
	}

	// check if the file is empty or doesn't contain the right keys
	if data.Configs == nil {
		return DeploymentReplicaConfigs{}, fmt.Errorf("Required field 'deployments' is missing")
	}

	return data, nil
}

// opens the state file and returns its contents
func Read(filePath string) ([]byte, error) {
	logger.Info("Reading " + filePath)
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// writes to the state file
func Put(data DeploymentReplicaConfigs, filePath string) error {
	logger.Info(fmt.Sprintf("%s: Marshalling JSON", filePath))
	file, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%s: Writing to local file", filePath))

	// mutex to lock the file if two requests come to it at the same time
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if err = os.WriteFile(filePath, file, 0644); err != nil {
		return err
	}

	return nil
}

// not a part of the challenge
// func Delete() {}
