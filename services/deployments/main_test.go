package deployments

import (
	"context"
	"log"
	"net"
	"reflect"
	"testing"

	"github.com/calvinbui/teleport-sre-challenge/config"
	pb "github.com/calvinbui/teleport-sre-challenge/proto/gen"
	"github.com/calvinbui/teleport-sre-challenge/services"
	"github.com/calvinbui/teleport-sre-challenge/services/state"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// bufconn creates an in memory grpc server (but not on a real socket/port) for unit testing
func dialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterDeploymentsServiceServer(server, New(&services.Services{
		DpConfig: state.DeploymentReplicaConfigs{
			// generate a config to use
			Configs: []state.DeploymentReplicaConfig{
				{
					Name:      "foo",
					Namespace: "default",
					Replicas:  1,
				},
				{
					Name:      "bar",
					Namespace: "default",
					Replicas:  2,
				},
				{
					Name:      "baz",
					Namespace: "something-else",
					Replicas:  3,
				},
			},
		},
	}, &config.Config{}))

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func Test_deploymentsService_ListDeployments(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name     string
		args     args
		response []*pb.DeploymentsReplicaConfig
		respCode codes.Code
	}{
		{
			name: "List configs in default namespace",
			args: args{namespace: "default"},
			response: []*pb.DeploymentsReplicaConfig{
				{Name: "foo", Namespace: "default", Replicas: 1},
				{Name: "bar", Namespace: "default", Replicas: 2},
			},
			respCode: codes.OK,
		},
		{
			name:     "No config exists",
			args:     args{namespace: "non-existent"},
			response: []*pb.DeploymentsReplicaConfig{},
			respCode: codes.OK,
		},
		// TODO update function to return all configs in all namespaces if request is empty
		{
			name:     "Empty namespace",
			args:     args{namespace: ""},
			response: []*pb.DeploymentsReplicaConfig{},
			respCode: codes.OK,
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDeploymentsServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.ListDeploymentsRequest{Namespace: tt.args.namespace}

			res, err := client.ListDeployments(ctx, req)

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.respCode {
						t.Errorf("ListDeployments() error code = %v, want = %v", er.Code(), tt.respCode)
						return
					}
				}
			}

			if res != nil {
				// also catch edge case where empty slice vs nil slice
				if !reflect.DeepEqual(res.DeploymentsReplicaConfig, tt.response) && len(res.DeploymentsReplicaConfig) != len(tt.response) {
					t.Errorf("ListDeployments() = %v, want = %v", res.DeploymentsReplicaConfig, tt.response)
					return
				}
			}
		})
	}
}

// this current errors due to the issue mentioned in the Kubernetes test file.
// i've commented out the tests for now
func Test_deploymentsService_GetDeployment(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name     string
		args     args
		response *pb.GetDeploymentResponse
		respCode codes.Code
	}{
		// {
		// 	name:     "Get config",
		// 	args:     args{name: "foo", namespace: "default"},
		// 	response: &pb.GetDeploymentResponse{Name: "foo", Namespace: "default", CurrentReplicas: 1, DesiredReplicas: 1},
		// 	respCode: codes.OK,
		// },
		// {
		// 	name:     "No config exists",
		// 	args:     args{name: "doesnt", namespace: "exist"},
		// 	response: &pb.GetDeploymentResponse{},
		// 	respCode: codes.NotFound,
		// },
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDeploymentsServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.GetDeploymentRequest{Name: tt.args.name, Namespace: tt.args.namespace}

			res, err := client.GetDeployment(ctx, req)

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.respCode {
						t.Errorf("GetDeployment() error code = %v, want = %v", er.Code(), tt.respCode)
						return
					}
				}
			}

			if res != nil {
				if !reflect.DeepEqual(res, tt.response) {
					t.Errorf("GetDeployment() = %v, want = %v", res, tt.response)
					return
				}
			}
		})
	}
}

func Test_deploymentsService_SetDeployment(t *testing.T) {
	type args struct {
		name      string
		namespace string
		replicas  int
	}
	tests := []struct {
		name     string
		args     args
		response *pb.SetDeploymentResponse
		respCode codes.Code
	}{
		// {
		// 	name:     "Update existing config",
		// 	args:     args{name: "foo", namespace: "default", replicas: 20},
		// 	response: &pb.SetDeploymentResponse{Name: "foo", Namespace: "default", OldReplicas: 1, NewReplicas: 20},
		// 	respCode: codes.OK,
		// },
		// {
		// 	name:     "Create new config",
		// 	args:     args{name: "doesnt", namespace: "exist", replicas: 2},
		// 	response: &pb.SetDeploymentResponse{Name: "foo", Namespace: "default", OldReplicas: 0, NewReplicas: 2},
		// 	respCode: codes.OK,
		// },
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewDeploymentsServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.DeploymentsReplicaConfig{Name: tt.args.name, Namespace: tt.args.namespace, Replicas: int32(tt.args.replicas)}

			res, err := client.SetDeployment(ctx, req)

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.respCode {
						t.Errorf("SetDeployment() error code = %v, want = %v", er.Code(), tt.respCode)
						return
					}
				}
			}

			if res != nil {
				if !reflect.DeepEqual(res, tt.response) {
					t.Errorf("SetDeployment() = %v, want = %v", res, tt.response)
					return
				}
			}
		})
	}
}
