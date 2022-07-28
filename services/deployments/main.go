package deployments

import (
	"context"
	"fmt"

	"github.com/calvinbui/teleport-sre-challenge/config"
	pb "github.com/calvinbui/teleport-sre-challenge/proto/gen"
	"github.com/calvinbui/teleport-sre-challenge/services"
	"github.com/calvinbui/teleport-sre-challenge/services/k8s"
	"github.com/calvinbui/teleport-sre-challenge/services/logger"
	"github.com/calvinbui/teleport-sre-challenge/services/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"k8s.io/apimachinery/pkg/api/errors"
)

type deploymentsService struct {
	services *services.Services
	config   *config.Config
	pb.UnimplementedDeploymentsServiceServer
}

func New(svc *services.Services, conf *config.Config) *deploymentsService {
	return &deploymentsService{
		services: svc,
		config:   conf,
	}
}

// ListDeployment is the gRPC endpoint to list all configurations
func (s *deploymentsService) ListDeployments(ctx context.Context, req *pb.ListDeploymentsRequest) (*pb.DeploymentsReplicaConfigs, error) {
	logger.Debug(fmt.Sprintf("ListDeployments: Request for namespace '%s'", req.Namespace))

	var resp pb.DeploymentsReplicaConfigs

	// this could be more efficient so we don't need to scan the entire file
	for _, d := range s.services.DpConfig.Configs {
		if d.Namespace == req.Namespace {
			resp.DeploymentsReplicaConfig = append(resp.DeploymentsReplicaConfig, &pb.DeploymentsReplicaConfig{
				Name:      d.Name,
				Namespace: d.Namespace,
				Replicas:  int32(d.Replicas),
			})
		}
	}

	logger.Debug(fmt.Sprintf("ListDeployments: Responding to request with %+v", &resp))
	return &resp, nil
}

// GetDeployment is the gRPC endpoint to get a specific configuration
// this could be cleaner
func (s *deploymentsService) GetDeployment(ctx context.Context, req *pb.GetDeploymentRequest) (*pb.GetDeploymentResponse, error) {
	logger.Debug(fmt.Sprintf("GetDeployment: Request for %s/%s", req.Name, req.Namespace))

	// keep track of the current replica count
	var currentReplicas int32

	// this could be more efficient so we don't need to scan the entire file
	if config, idx := s.services.DpConfig.Find(req.Name, req.Namespace); idx != -1 {
		scale, err := k8s.GetScale(ctx, s.services.K8s, config.Namespace, config.Name)
		if err != nil {
			if errors.IsNotFound(err) {
				currentReplicas = 0
			} else {
				return nil, status.Error(codes.NotFound, fmt.Sprintf("Deployment %s/%s not found", req.Namespace, req.Name))
			}
		} else {
			currentReplicas = scale.Spec.Replicas
		}

		resp := pb.GetDeploymentResponse{
			Name:            req.Name,
			Namespace:       req.Namespace,
			DesiredReplicas: int32(config.Replicas),
			CurrentReplicas: currentReplicas,
		}

		logger.Debug(fmt.Sprintf("GetDeployment: Responding to request with %+v", &resp))
		return &resp, nil
	}

	return nil, status.Error(codes.NotFound, fmt.Sprintf("Deployment %s/%s not found", req.Namespace, req.Name))
}

// SetDeployment is the gRPC endpoint to set the replica count for a deployment
func (s *deploymentsService) SetDeployment(ctx context.Context, req *pb.DeploymentsReplicaConfig) (*pb.SetDeploymentResponse, error) {
	logger.Debug(fmt.Sprintf("SetDeployment: Request to set %s/%s to %v replicas", req.Name, req.Namespace, req.Replicas))

	resp := pb.SetDeploymentResponse{
		Name:        req.Name,
		Namespace:   req.Namespace,
		NewReplicas: req.Replicas,
	}

	logger.Info(fmt.Sprintf("%s/%s: Setting replicas to %v", req.Name, req.Namespace, req.Replicas))
	// this could be more efficient so we don't need to scan the entire file
	if config, idx := s.services.DpConfig.Find(req.Name, req.Namespace); idx != -1 {
		logger.Info(fmt.Sprintf("%s/%s: Existing config found with %v replicas", req.Name, req.Namespace, config.Replicas))

		// 0 if doesn't exist
		resp.OldReplicas = int32(config.Replicas)

		// set to new replica value
		s.services.DpConfig.Configs[idx].Replicas = int(req.Replicas)
	} else {
		logger.Info(fmt.Sprintf("%s/%s: No existing config found. Creating...", req.Name, req.Namespace))
		s.services.DpConfig.Configs = append(s.services.DpConfig.Configs, state.DeploymentReplicaConfig{
			Name:      req.Name,
			Namespace: req.Namespace,
			Replicas:  int(req.Replicas),
		})
	}

	logger.Debug(fmt.Sprintf("SetDeployment: Updating %s with %+v", s.config.FilePath, s.services.DpConfig))

	// write the changes to persist state
	if err := state.Put(s.services.DpConfig, s.config.FilePath); err != nil {
		return nil, err
	}

	// update the deployment's replica count
	if _, err := k8s.UpdateScale(ctx, s.services.K8s, req.Namespace, req.Name, int(req.Replicas)); err != nil {
		return &resp, nil
	}

	logger.Debug(fmt.Sprintf("SetDeployment: Responding to request with %+v", &resp))
	return &resp, nil
}
